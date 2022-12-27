package main

import "encoding/binary"

// Package : wcl package proto
type Package struct {
	datas [][]byte
}

// Parse : parse message
func (p *Package) Parse(msg []byte) bool {
	pckLen := len(msg)
	if pckLen < 4 {
		return false
	}

	lenBuff := msg[0:4]
	lenSize := p.parseInt(lenBuff)

	if lenSize+4 != pckLen {
		return false
	}

	datas := msg[4 : lenSize-4]
	offset := 0

	for offset < lenSize {
		dataSize := p.parseInt(datas[offset : offset+4])
		data := datas[offset+4 : offset+4+dataSize]
		offset += dataSize + 4

		p.datas = append(p.datas, data)
	}
	return true
}

func (p *Package) parseInt(buff []byte) int {
	return int(binary.LittleEndian.Uint32(buff))
}

// LogPck : log package
type LogPck struct {
	*Package

	FName string
	Data  []byte
}

// ParseMsg : parse wcl message
func (p *LogPck) ParseMsg(msg []byte) bool {
	ret := p.Parse(msg)
	p.FName = string(p.datas[0])
	p.Data = p.datas[1]

	return ret
}

// GetName : get file name
func (p *LogPck) GetName() string {
	return p.FName
}

// GetData : get write data
func (p *LogPck) GetData() []byte {
	return p.Data
}

// NewLogPck : create new packet
func NewLogPck() *LogPck {
	return &LogPck{
		&Package{
			datas: make([][]byte, 0),
		},
		"",
		nil,
	}
}
