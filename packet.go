package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	VerLen    = 2
	AskLen    = 4
	BodyLen   = 4
	HeaderLen = VerLen + AskLen + BodyLen

	VerOffset     = 0
	AskOffset     = VerLen
	BodyLenOffset = VerLen + AskLen
)

type Packet struct {
	Ver   uint16 `json:"ver"`
	AskID uint32 `json:"ask_id"`
	Body  []byte `json:"body"`
}

func (p *Packet) ReadTCP(rr *bufio.Reader) (err error) {
	buf := make([]byte, HeaderLen) //todo
	if _, err = io.ReadFull(rr, buf); err != nil {
		return
	}
	p.Ver = binary.BigEndian.Uint16(buf[VerOffset:AskOffset])
	p.AskID = binary.BigEndian.Uint32(buf[AskOffset:BodyLenOffset])
	rawLen := binary.BigEndian.Uint32(buf[BodyLenOffset:])
	body := make([]byte, rawLen) //todo
	if _, err = io.ReadFull(rr, body); err != nil {
		return
	}
	p.Body = body

	return nil
}

func (p *Packet) WriteTCP(wr *bufio.Writer) (err error) {
	rawLen := len(p.Body)
	buf := make([]byte, HeaderLen+rawLen)
	binary.BigEndian.PutUint16(buf[:AskOffset], p.Ver)
	binary.BigEndian.PutUint32(buf[AskOffset:BodyLenOffset], p.AskID)
	binary.BigEndian.PutUint32(buf[BodyLenOffset:HeaderLen], uint32(rawLen))
	//buf[HeaderLen:] = p.Body
	copy(buf[HeaderLen:], p.Body)
	if _, err = wr.Write(buf); err != nil {
		return
	}
	if err = wr.Flush(); err != nil {
		return
	}

	return nil
}

func (p *Packet) ToString() string {
	return fmt.Sprintf("\n-------- packet --------\nver: %d\naskId: %d\nbody: %s\n-----------------------", p.Ver, p.AskID, p.Body)
}
