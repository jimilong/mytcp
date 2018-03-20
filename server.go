package main

import (
	"log"
	"net"
)

type Server struct {
	Addr            string
	ReceiveChanSize int
	SendChanSize    int
	SocketBufSize   int //todo conn-buf bufio-buf ??
}

func NewServer(addr string, rchan int, schan int, sbuf int) *Server {
	return &Server{
		Addr:            addr,
		ReceiveChanSize: rchan,
		SendChanSize:    schan,
		SocketBufSize:   sbuf,
	}
}

func (s *Server) Start() {
	var (
		tcpAddr  *net.TCPAddr
		listener *net.TCPListener
		conn     *net.TCPConn
		err      error
	)
	if tcpAddr, err = net.ResolveTCPAddr("tcp", s.Addr); err != nil {
		log.Printf("net.ResolveTCPAddr(\"tcp\", \"%s\") error(%v) \n", s.Addr, err)
		return
	}
	if listener, err = net.ListenTCP("tcp", tcpAddr); err != nil {
		log.Printf("net.ListenTCP(\"tcp\", \"%s\") error(%v) \n", tcpAddr, err)
		return
	}
	log.Printf("start tcp listen: \"%s\" \n", s.Addr)

	for {
		if conn, err = listener.AcceptTCP(); err != nil {
			log.Printf("listener.Accept(\"%s\") error(%v)", listener.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(true); err != nil {
			log.Printf("conn.SetKeepAlive() error(%v)", err)
			return
		}
		if err = conn.SetReadBuffer(1024); err != nil {
			log.Printf("conn.SetReadBuffer() error(%v)", err)
			return
		}
		if err = conn.SetWriteBuffer(1024); err != nil {
			log.Printf("conn.SetWriteBuffer() error(%v)", err)
			return
		}

		go handleConn(s, conn)
	}
}

func (s *Server) Close() {

}
