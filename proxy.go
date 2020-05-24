package main

import (
	"io"
	"fmt"
)

type Proxy struct {
	Logger   *Logger
	DidError bool
	OnExit   chan bool
	Peer1    *Peer
	Peer2    *Peer
}

func NewProxy(peer1, peer2 *Peer, logger *Logger) *Proxy {
	return &Proxy{
		Peer1: peer1,
		Peer2: peer2,
		Logger: logger,
		OnExit: make(chan bool),
		DidError:  false,
	}
}

func (p *Proxy) Start() {

	defer p.Peer1.Conn.Close()
	defer p.Peer2.Conn.Close()
	p.Logger.Log("Opened %s <<>> %s\n", p.Peer1.Addr.String(), p.Peer2.Addr.String())

	go p.pipe(p.Peer1.Conn, p.Peer2.Conn)
	go p.pipe(p.Peer2.Conn, p.Peer1.Conn)

	<-p.OnExit
	p.Logger.Log("Peer1 Closed (%d bytes sent, %d bytes recieved)\n", p.Peer1.SentBytes, p.Peer1.ReceivedBytes)
	p.Logger.Log("Peer2 Closed (%d bytes sent, %d bytes recieved)\n", p.Peer2.SentBytes, p.Peer2.ReceivedBytes)
}

func (p *Proxy) err(s string, err error) {
	if p.DidError {
		return
	}

	if err != io.EOF {
		p.Logger.Log(s, err)
	}

	p.OnExit <- true
	p.DidError = true
}

func (p *Proxy) pipe(src, dst io.ReadWriter) {
	isLeftPeer := src == p.Peer1.Conn

	var bytesNoFmt string
	if isLeftPeer {
		bytesNoFmt = "[P1 ==> P2] %d bytes | %s"
	} else {
		bytesNoFmt = "[P1 <== P2] %d bytes | %s"
	}

	var byteFormat string
	if *outputHex {
		byteFormat = "%x"
	} else {
		byteFormat = "%s"
	}

	buff := make([]byte, 0xffff)
	for {
		n, err := src.Read(buff)
		if err != nil {
			p.err("Read failed: %s\n", err)
			return
		}

		b := buff[:n]
		p.Logger.Log(bytesNoFmt, n, fmt.Sprintf(byteFormat, b))

		n, err = dst.Write(b)
		if err != nil {
			p.err("Write failed: %s\n", err)
			return
		}

		if isLeftPeer {
			p.Peer1.SentBytes += uint64(n)
			p.Peer2.ReceivedBytes += uint64(n)
		} else {
			p.Peer1.SentBytes += uint64(n)
			p.Peer2.ReceivedBytes += uint64(n)
		}
	}
}
