package main

import (
	"fmt"
	f "github.com/CSUNetSec/flowride"
	"github.com/intel-go/nff-go/common"
	"github.com/intel-go/nff-go/flow"
	"github.com/intel-go/nff-go/packet"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
)

const (
	TCP = int(0)
	UDP
	ICMP
)

// splits the packets from a 32 vector input to 16 different flows. so 2 packets per flow)
func vsplit(pkt []*packet.Packet, mask *[32]bool, res *[32]uint8, ctx flow.UserContext) {
	for i := uint8(0); i < 2; i++ {
		if (*mask)[i] {
			(*res)[i] = uint8(0)
		}
		if (*mask)[i+2] {
			(*res)[i+2] = uint8(1)
		}
		if (*mask)[i+4] {
			(*res)[i+4] = uint8(2)
		}
		if (*mask)[i+6] {
			(*res)[i+6] = uint8(3)
		}
		if (*mask)[i+8] {
			(*res)[i+8] = uint8(4)
		}
		if (*mask)[i+10] {
			(*res)[i+10] = uint8(5)
		}
		if (*mask)[i+12] {
			(*res)[i+12] = uint8(6)
		}
		if (*mask)[i+14] {
			(*res)[i+14] = uint8(7)
		}
		if (*mask)[i+16] {
			(*res)[i+16] = uint8(8)
		}
		if (*mask)[i+18] {
			(*res)[i+18] = uint8(9)
		}
		if (*mask)[i+20] {
			(*res)[i+20] = uint8(10)
		}
		if (*mask)[i+22] {
			(*res)[i+22] = uint8(11)
		}
		if (*mask)[i+24] {
			(*res)[i+24] = uint8(12)
		}
		if (*mask)[i+26] {
			(*res)[i+26] = uint8(13)
		}
		if (*mask)[i+28] {
			(*res)[i+28] = uint8(14)
		}
		if (*mask)[i+30] {
			(*res)[i+30] = uint8(15)
		}
	}
}

func pktHandler(pkt []*packet.Packet, mask *[32]bool, ctx flow.UserContext) {
	fctx := ctx.(f.FlowrideContext)
	pktBuf := <-fctx.BufChan
	processed := uint64(0)
	startind := pktBuf.Ind
	var p f.PktCap
	for i := 0; i < 32; i++ {
		if (*mask)[i] {
			pktIPv4, _, _ := pkt[i].ParseAllKnownL3()
			if pktIPv4 != nil {
				pktTCP, pktUDP, pktICMP := pkt[i].ParseAllKnownL4ForIPv4()
				if pktTCP != nil {
					p = f.PktCapFromBytes(pktIPv4.SrcAddr, pktIPv4.DstAddr, pktTCP.SrcPort, pktTCP.DstPort, 0, 0)
				} else if pktUDP != nil {
					p = f.PktCapFromBytes(pktIPv4.SrcAddr, pktIPv4.DstAddr, pktUDP.SrcPort, pktUDP.DstPort, 0, 0)
				} else if pktICMP != nil {
					p = f.PktCapFromBytes(pktIPv4.SrcAddr, pktIPv4.DstAddr, 0, 0, 0, 0)
				}
				if startind >= f.CAPBUFSIZE-32 {
					f.LogFatal(fmt.Sprintf("ind is:%d processed:%d\n", pktBuf.Ind, processed))
				}
				pktBuf.Buf[startind+processed] = p
				//		(*privCnt)[ppair][ind] = (*privCnt)[ppair][ind] + 1
				processed++
			}
		}
	}
	fctx.IndChan <- processed
}

var (
	flconf f.FlConf
)

func main() {
	if len(os.Args) < 2 {
		f.LogFatal("please provde config file")
	}
	flconf, err := f.ConfigFromFileName(os.Args[1])
	f.CheckLogFatal(err)
	fmt.Printf("config is :%v\n", flconf)
	conf := flconf.FlCapConf
	if conf.Profiler {
		fmt.Printf("starting profiler on localhost:6161")
		go func() {
			log.Println(http.ListenAndServe("localhost:6161", nil))
		}()
	}
	dpdkwlist := conf.DPDKArgs
	for _, v := range append(conf.DpdkInPorts, conf.DpdkOutPorts...) {
		dpdkwlist = append(dpdkwlist, "-w "+v)
	}
	if len(conf.DpdkInPorts) != len(conf.DpdkOutPorts) {
		f.LogFatal("currently you should have the same numer of in and out ports")
	}
	config := flow.Config{
		CPUList:  flconf.FlCapConf.CpuList,
		DPDKArgs: dpdkwlist,
		LogType:  common.Debug,
	}
	fmt.Printf("\n DPDPKArgs: %v\n", dpdkwlist)
	f.CheckLogFatal(flow.SystemInit(&config))
	fctx := f.NewFlowrideContext()

	for i := 0; i < len(conf.DpdkInPorts); i++ {
		fmt.Println("trying to get ", conf.DpdkInPorts[i])
		inp, err := flow.GetPortByName(strings.Trim(conf.DpdkInPorts[i], " "))
		fmt.Println("got ", inp, " err ", err)
		//f.CheckLogFatal(err)
		fmt.Println("trying to get ", conf.DpdkOutPorts[i])
		outp, err := flow.GetPortByName(strings.Trim(conf.DpdkOutPorts[i], " "))
		fmt.Println("got ", outp)
		//f.CheckLogFatal(err)
		// XXX HACK wtf is happening here
		inp = 0
		outp = 1
	}

	inFlow, err := flow.SetReceiver(0)
	f.CheckLogFatal(err)
	err = flow.SetVectorHandler(inFlow, pktHandler, fctx)
	f.CheckLogFatal(err)
	err = flow.SetSender(inFlow, 1)
	f.CheckLogFatal(err)

	inFlow1, err := flow.SetReceiver(2)
	f.CheckLogFatal(err)
	err = flow.SetVectorHandler(inFlow1, pktHandler, fctx)
	f.CheckLogFatal(err)
	err = flow.SetSender(inFlow1, 3)
	f.CheckLogFatal(err)
	f.CheckLogFatal(flow.SystemStart())
}
