package main

import (
	"os"
	"io"
	"net"
)

const (
	PEER_TYPE_SERVER = iota
	PEER_TYPE_CLIENT
)

type PeerType uint32

type Peer struct {
	Logger   *Logger
	Type PeerType
	Addr *net.TCPAddr
	Conn io.ReadWriteCloser
	Listener *net.TCPListener
	SentBytes     uint64
	ReceivedBytes uint64
}

func NewPeer(addr *net.TCPAddr, peerType PeerType, logger *Logger) *Peer {
	return &Peer{
		Addr: addr,
		Type: peerType,
		Logger: logger,
	}
}

func (p *Peer) Setup() {
	if p.Type == PEER_TYPE_SERVER {
		listener, err := net.ListenTCP("tcp", p.Addr)
		if err != nil {
			p.Logger.Log("Failed to open local port to listen: %s\n", err)
			os.Exit(1)
		}
		p.Listener = listener
	}
}

func (p *Peer) DialOrAccept() (*net.TCPConn, error) {
	if p.Type == PEER_TYPE_SERVER {
		conn, err := p.Listener.AcceptTCP()
		if err != nil {
			p.Logger.Log("Failed to accept connection '%s'\n", err)
			return nil, err
		}
		p.Conn = conn

		return conn, nil
	}

	conn, err := net.DialTCP("tcp", nil, p.Addr)
	if err != nil {
		p.Logger.Log("Failed to accept connection '%s'\n", err)
		return nil, err
	}
	p.Conn = conn

	return conn, nil
} 
