package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	errConnClosing  = errors.New("use of closed network connection")
	errSendBlocking = errors.New("write packet was blocking")
)

type tcpConn struct {
	conn        *net.TCPConn
	reader      *bufio.Reader
	writer      *bufio.Writer
	receiveChan chan *Packet
	sendChan    chan *Packet
	closeChan   chan struct{}
	closeOnce   sync.Once
}

func handleConn(server *Server, conn *net.TCPConn) {
	log.Printf("start tcp serve \"%s\" with \"%s\"",
		conn.LocalAddr().String(), conn.RemoteAddr().String())
	tcp := &tcpConn{
		conn:        conn,
		reader:      bufio.NewReader(conn),
		writer:      bufio.NewWriter(conn),
		receiveChan: make(chan *Packet, server.ReceiveChanSize),
		sendChan:    make(chan *Packet, server.SendChanSize),
		closeChan:   make(chan struct{}),
	}

	go tcp.readLoop()
	go tcp.handleLoop()
	go tcp.writeLoop()
}

func (tc *tcpConn) AsyncSend(p *Packet) error {
	select {
	case tc.sendChan <- p:
		return nil
	case <-tc.closeChan:
		return errConnClosing
	default:
		return errSendBlocking
	}
}

func (tc *tcpConn) readLoop() {
	defer tc.close()

	var (
		packet *Packet
		err    error
	)
	for {
		select {
		case <-tc.closeChan:
			return
		default:
		}
		packet = &Packet{} // todo packetPool
		if err = packet.ReadTCP(tc.reader); err != nil {
			log.Printf("read tcp error(%v)\n", err)
			break
		}
		tc.receiveChan <- packet
	}
}

func (tc *tcpConn) handleLoop() {
	defer tc.close()

	for {
		select {
		case <-tc.closeChan:
			return
		case p, ok := <-tc.receiveChan:
			if !ok {
				return
			}
			fmt.Println(p.ToString())
			//todo dispatch
		}
	}
}

func (tc *tcpConn) writeLoop() {
	defer tc.close()

	var err error
	for {
		select {
		case <-tc.closeChan:
			return
		case p, ok := <-tc.sendChan:
			if !ok {
				return
			}
			if err = p.WriteTCP(tc.writer); err != nil {
				log.Printf("write tcp error (%v) \n", err)
				return
			}
		}
	}
}

func (tc *tcpConn) close() {
	tc.closeOnce.Do(func() {
		close(tc.closeChan)
		close(tc.receiveChan)
		close(tc.sendChan)
		tc.conn.Close()
		log.Printf("stop tcp serve \"%s\" with \"%s\"",
			tc.conn.LocalAddr().String(), tc.conn.RemoteAddr().String())
	})
}
