package flowride

import (
	"github.com/intel-go/nff-go/types"
	"time"
)

const (
	CAPBUFSIZE = 1 << 20
)

type PktCap struct {
	srcIP    types.IPv4Address
	dstIP    types.IPv4Address
	sport    uint16
	dport    uint16
	tcpFlags types.TCPFlags
	len      uint16
}

func PktCapFromBytes(ip1, ip2 types.IPv4Address, p1, p2 uint16, fl types.TCPFlags, pl uint16) PktCap {
	return PktCap{
		srcIP:    ip1,
		dstIP:    ip2,
		sport:    p1,
		dport:    p2,
		tcpFlags: fl,
		len:      pl,
	}
}

type PktBuf struct {
	Buf    []PktCap  //a slice of packet headers that will be used to create flows
	Ind    uint64    //the index where the handling goroutines will write the next cap
	Tstamp time.Time // the export time of this buffer
}

func newPktBuf() PktBuf {
	return PktBuf{
		Buf: make([]PktCap, CAPBUFSIZE),
		Ind: 0,
	}
}

// Flowride context is a way for the goroutines started by NFF to
// communucate with the handle functions and provide spaces to
// write packet captures
type FlowrideContext struct {
	BufChan  chan *PktBuf
	termChan chan struct{}
	IndChan  chan uint64
}

func (f FlowrideContext) Copy() interface{} {
	return NewFlowrideContext()
}

func (f FlowrideContext) Delete() {
	LogInfo("closing Flowride Context")
	f.termChan <- struct{}{}
}

func NewFlowrideContext() FlowrideContext {
	LogInfo("starting new Flowride Context")
	ret := FlowrideContext{
		BufChan:  make(chan *PktBuf),
		termChan: make(chan struct{}),
		IndChan:  make(chan uint64),
	}
	go BufMgr(ret.BufChan, ret.termChan, ret.IndChan)
	return ret
}

func BufMgr(bc chan *PktBuf, tc chan struct{}, ic chan uint64) {
	var (
		curbuf, tmpbuf PktBuf // current is turned to tmp before it is passed for processing
	)
	curbuf = newPktBuf()
	for {
		select {
		case <-tc:
			break
		case bc <- &curbuf:
			nind := <-ic
			curbuf.Ind += uint64(nind)
			// That would mean that on the next call it might overflow
			// because a vector call can write a max of 32 packets at once.
			if curbuf.Ind >= CAPBUFSIZE-33 {
				tmpbuf = curbuf
				go ProcessBuf(&tmpbuf)
				curbuf = newPktBuf()
			}
		}
	}
	LogInfo("BufMgr terminating")
	close(bc)
}

func ProcessBuf(a *PktBuf) {
	a.Tstamp = time.Now()
}
