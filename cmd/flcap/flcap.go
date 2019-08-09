package main

import (
	"time"

	f "github.com/CSUNetSec/flowride"
	"fmt"
	"os"
	"github.com/intel-go/nff-go/common"
	"github.com/intel-go/nff-go/flow"
	"github.com/intel-go/nff-go/packet"
)

const (
	TCP = int(0)
	UDP
	ICMP
)

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
			(*res)[i+20] = uint8(9)
		}
		if (*mask)[i+22] {
			(*res)[i+22] = uint8(8)
		}
		if (*mask)[i+24] {
			(*res)[i+24] = uint8(7)
		}
		if (*mask)[i+26] {
			(*res)[i+26] = uint8(6)
		}
		if (*mask)[i+28] {
			(*res)[i+28] = uint8(5)
		}
		if (*mask)[i+30] {
			(*res)[i+30] = uint8(4)
		}
	}
}

func genHandle(ind int) func([]*packet.Packet, *[32]bool, flow.UserContext) {
	return func(pkt []*packet.Packet, mask *[32]bool, ctx flow.UserContext) {
		for i := 0; i < 32; i++ {
			if (*mask)[i] {
				bufpos := (*privCnt)[ind] * 12 // index times 2 IPs + 2 ports
				if bufpos >= bufsize-13 {      //one pair of IPs +  ports and a safe.
					//fmt.Println("q: ", i, " reseting buffer!")
					(*privCnt)[ind] = 0
					bufpos = 0 // ahaha. edw itan.
					time.Now()
				}
				pktIPv4, _, _ := pkt[i].ParseAllKnownL3()
				if pktIPv4 != nil {
					//pktTCP, pktUDP, pktICMP := pkt[i].ParseAllKnownL4ForIPv4()
					pkt[i].ParseAllKnownL4ForIPv4()
					(*privBufs)[ind][bufpos] = byte(pktIPv4.SrcAddr)
					(*privBufs)[ind][bufpos+1] = byte(pktIPv4.SrcAddr >> 8)
					(*privBufs)[ind][bufpos+2] = byte(pktIPv4.SrcAddr >> 16)
					(*privBufs)[ind][bufpos+3] = byte(pktIPv4.SrcAddr >> 24)
					(*privBufs)[ind][bufpos+4] = byte(pktIPv4.DstAddr)
					(*privBufs)[ind][bufpos+5] = byte(pktIPv4.DstAddr >> 8)
					(*privBufs)[ind][bufpos+6] = byte(pktIPv4.DstAddr >> 16)
					(*privBufs)[ind][bufpos+7] = byte(pktIPv4.DstAddr >> 24)
					//copy((*privBufs)[ind][bufpos], []byte{srcAddr, dstAddr, pktIPv4.SrcPort, pktIPv4.DstPort})
					(*privCnt)[ind] = (*privCnt)[ind] + 1
				}
			}
		}
	}
}

const (
	bufsize = 1 << 15
)

var (
	privCnt         *[10]uint64
	privBufs        *[10][bufsize]byte
	flconf		f.FlConf
)

func main() {
	if len(os.Args) < 2 {
		f.LogFatal("please provde config file")
	}
	flconf, err := f.ConfigFromFileName(os.Args[1])
	f.CheckLogFatal(err)
	fmt.Printf("config is :%v\n", flconf)
	
	config := flow.Config{
		CPUList: flconf.FlCapConf.CpuList,
		DPDKArgs: flconf.FlCapConf.DPDKArgs,
		LogType: common.Debug,
	}
	f.CheckLogFatal(flow.SystemInit(&config))
	inFlow, err := flow.SetReceiver(uint16(flconf.FlCapConf.DpdkInPort))
	f.CheckLogFatal(err)
	sflows, _ := flow.SetVectorSplitter(inFlow, vsplit, 10, nil)
	privCnt = &[10]uint64{}
	privBufs = &[10][bufsize]byte{}
	for i := range sflows {
		flow.SetVectorHandler(sflows[i], genHandle(i), nil)
		flow.SetSender(sflows[i], uint16(flconf.FlCapConf.DpdkOutPort))
	}
	f.CheckLogFatal(flow.SystemStart())
}
