package main

import (
	"os"
	"net"
	"flag"
	"log"
	"fmt"
)

const (
	MODE_L2R = iota
	MODE_L2L
	MODE_R2R
)

var (
	logger *log.Logger
	proxyId uint32

	localAddr1   = flag.String("l1", "", "local address")
	localAddr2   = flag.String("l2", ":9999", "second local address")
	remoteAddr1  = flag.String("r1", "", "remote address")
	remoteAddr2  = flag.String("r2", "", "second remote address")
	outputHex    = flag.Bool("h", false, "output hex")
)

func main() {
	flag.Parse()
	logger = log.New(os.Stdout, "", log.LstdFlags)
	
	var mode uint32
	if *localAddr1 == "" {
		mode = MODE_R2R
		logger.Printf("Proxying from %v to %v\n", *remoteAddr1, *remoteAddr2)
	} else if *localAddr1 != "" && *remoteAddr1 != "" {
		mode = MODE_L2R
		logger.Printf("Proxying from %v to %v\n", *localAddr1, *remoteAddr1)
	} else {
		mode = MODE_L2L
		logger.Printf("Proxying from %v to %v\n", *localAddr1, *localAddr2)
	}


	laddr1, err := net.ResolveTCPAddr("tcp", *localAddr1)
	if err != nil {
		logger.Printf("Failed to resolve local address: %s\n", err)
		os.Exit(1)
	}
	laddr2, err := net.ResolveTCPAddr("tcp", *localAddr2)
	if err != nil {
		logger.Printf("Failed to resolve local address: %s\n", err)
		os.Exit(1)
	}
	raddr1, err := net.ResolveTCPAddr("tcp", *remoteAddr1)
	if err != nil {
		logger.Printf("Failed to resolve remote address: %s\n", err)
		os.Exit(1)
	}
	raddr2, err := net.ResolveTCPAddr("tcp", *remoteAddr2)
	if err != nil {
		logger.Printf("Failed to resolve remote address2: %s\n", err)
		os.Exit(1)
	}

	peerLogger1 := NewLogger("[P1]")
	peerLogger2 := NewLogger("[P2]")
	var peer1, peer2 *Peer
	var async bool
	switch(mode) {
	case MODE_L2L:
		peer1 = NewPeer(laddr1, PEER_TYPE_SERVER, peerLogger1)
		peer2 = NewPeer(laddr2, PEER_TYPE_SERVER, peerLogger2)
		async = true
		break
	case MODE_R2R:
		peer1 = NewPeer(raddr1, PEER_TYPE_CLIENT, peerLogger1)
		peer2 = NewPeer(raddr2, PEER_TYPE_CLIENT, peerLogger2)
		async = false
		break
	case MODE_L2R:
		peer1 = NewPeer(laddr1, PEER_TYPE_SERVER, peerLogger1)
		peer2 = NewPeer(raddr1, PEER_TYPE_CLIENT, peerLogger2)
		async = true
		break
	}
	
	start(peer1, peer2, async)
}

func start(peer1, peer2 *Peer, async bool) {
	peer1.Setup()
	peer2.Setup()
	for {
		_, err := peer1.DialOrAccept()
		if err != nil {
			continue
		}

		_, err = peer2.DialOrAccept()
		if err != nil {
			continue
		}

		proxyId++
		l := NewLogger(fmt.Sprintf("[Proxy-%d] ", proxyId))
		p := NewProxy(peer1, peer2, l)

		if async {
			go p.Start()
		} else {
			p.Start()
		}
	}
}
