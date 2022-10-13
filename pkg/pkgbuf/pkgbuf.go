package pkgbuf

import (
	"encoding/binary"
	"github.com/songgao/packets/ethernet"
	"io"
)

type PkgBuff struct {
	buf   []byte
	end   int
	start int
	rd    io.Reader
}

func New(rd io.Reader) *PkgBuff {
	pkgBuff := PkgBuff{
		buf: make([]byte, 8192),
		rd:  rd,
	}
	return &pkgBuff
}

func (p *PkgBuff) ReadFrame() (ethernet.Frame, error) {
	needInput := false
	for {
		if needInput {
			needInput = false
			if p.end+2048 > 8192 {
				p.start = 0
				p.end = 0
			}
			n, err := p.rd.Read(p.buf[p.end : p.end+2048])
			if err != nil {
				return nil, err
			}
			p.end += n
		}

		isFound := false
		for p.start < p.end-3 {
			if p.buf[p.start] == byte(0xBC) {
				if p.buf[p.start+1] == byte(0xBC) {
					isFound = true
				}
				break
			}

			p.start++
		}

		if p.start != 0 && p.end != p.start {
			copy(p.buf, p.buf[p.start:p.end])
			p.end -= p.start
			p.start = 0
		}

		if !isFound {
			needInput = true
			continue
		}

		length := int(binary.BigEndian.Uint16(p.buf[2:4]))
		if p.end-p.start < 4+length {
			needInput = true
			continue
		}

		frame := make([]byte, length)
		copy(frame, p.buf[p.start+4:p.start+length+4])
		copy(p.buf[p.start:], p.buf[p.start+length+4:p.end])
		p.end -= length + 4

		return frame, nil
	}
}
