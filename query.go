// Some of the code in this module was from urShadow, it was adapted and modified 2017-07-01

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// LegacyQuery stores state for old-style masterlist queries
type LegacyQuery struct {
	addr *net.UDPAddr
	conn net.Conn
}

// GetServerLegacyInfo wraps a set of legacy queries and returns a new Server object with the
// available fields populated.
func GetServerLegacyInfo(host string) (server Server, err error) {
	lq, err := NewLegacyQuery(host)
	if err != nil {
		return server, err
	}

	server, err = lq.GetInfo()
	if err != nil {
		return server, err
	}

	// server.PlayerList, err = lq.GetPlayers()
	// if err != nil {
	// 	return server, err
	// }

	err = lq.Close()

	return server, err
}

// NewLegacyQuery creates a new legacy query handler for a server
func NewLegacyQuery(host string) (lq *LegacyQuery, err error) {
	lq = new(LegacyQuery)
	lq.addr, err = net.ResolveUDPAddr("udp", host)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve: %v", err)
	}

	lq.conn, err = net.DialUDP("udp", nil, lq.addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return lq, nil
}

// Close closes a legacy query manager's connection
func (lq *LegacyQuery) Close() error {
	return lq.conn.Close()
}

func (lq *LegacyQuery) sendQuery(id rune) ([]byte, error) {
	request := new(bytes.Buffer)
	response := make([]byte, 2048)

	port := [2]byte{
		byte(lq.addr.Port & 0xFF),
		byte((lq.addr.Port >> 8) & 0xFF),
	}

	lq.conn.SetDeadline(time.Now().Add(3000 * time.Millisecond))

	binary.Write(request, binary.LittleEndian, []byte("SAMP"))
	binary.Write(request, binary.LittleEndian, lq.addr.IP.To4())
	binary.Write(request, binary.LittleEndian, port[0])
	binary.Write(request, binary.LittleEndian, port[1])
	binary.Write(request, binary.LittleEndian, uint8(id))

	if id == 'p' {
		binary.Write(request, binary.LittleEndian, uint32(0))
	}

	lq.conn.Write(request.Bytes())

	n, err := lq.conn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if n > cap(response) {
		return nil, fmt.Errorf("read response over buffer capacity")
	}

	return response[:n], nil
}

// GetInfo returns the core server info for displaying on the browser list.
func (lq *LegacyQuery) GetInfo() (server Server, err error) {
	response, err := lq.sendQuery('i')
	if err != nil {
		return server, err
	}

	ptr := 11

	server.Password = (response[ptr] == 1)
	ptr++

	server.Players = int(binary.LittleEndian.Uint16(response[ptr : ptr+2]))
	ptr += 2

	server.MaxPlayers = int(binary.LittleEndian.Uint16(response[ptr : ptr+2]))
	ptr += 2

	hostnameLen := int(binary.LittleEndian.Uint16(response[ptr : ptr+4]))
	ptr += 4

	server.Hostname = string(response[ptr : ptr+hostnameLen])
	ptr += hostnameLen

	gamemodeLen := int(binary.LittleEndian.Uint16(response[ptr : ptr+4]))
	ptr += 4

	server.Gamemode = string(response[ptr : ptr+gamemodeLen])
	ptr += gamemodeLen

	languageLen := int(binary.LittleEndian.Uint16(response[ptr : ptr+4]))
	ptr += 4

	if languageLen > 0 {
		server.Language = string(response[ptr : ptr+languageLen])
		ptr += languageLen
	} else {
		server.Language = "-"
	}

	return
}

// GetRules returns a map of rule properties from a server. The legacy query uses established keys
// such as "Map" and "Version"
func (lq *LegacyQuery) GetRules() (rules map[string]string, err error) {
	response, err := lq.sendQuery('r')
	if err != nil {
		return rules, err
	}

	var (
		body                = bytes.NewBuffer(response)
		amount              uint16
		rulename, rulevalue string
		strlen              uint8
		strbuf              []byte
	)

	body.Next(11)
	binary.Read(body, binary.LittleEndian, &amount)

	for i := uint16(0); i < amount; i++ {
		binary.Read(body, binary.LittleEndian, &strlen)
		strbuf = make([]byte, strlen)
		binary.Read(body, binary.LittleEndian, &strbuf)
		rulename = string(strbuf)

		binary.Read(body, binary.LittleEndian, &strlen)
		strbuf = make([]byte, strlen)
		binary.Read(body, binary.LittleEndian, &strbuf)
		rulevalue = string(strbuf)

		rules[rulename] = rulevalue
	}

	return
}

// GetPlayers simply returns a slice of strings, score is rather arbitrary so it's omitted.
func (lq *LegacyQuery) GetPlayers() (players []string, err error) {
	response, err := lq.sendQuery('c')
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(response)

	var (
		count    uint16
		nickname string
		strlen   uint8
		strbuf   []byte
	)

	body.Next(11)
	binary.Read(body, binary.LittleEndian, &count)

	list := make([]string, count)

	for i := uint16(0); i < count; i++ {
		binary.Read(body, binary.LittleEndian, &strlen)
		strbuf = make([]byte, strlen)
		binary.Read(body, binary.LittleEndian, &strbuf)
		nickname = string(strbuf)

		list[i] = nickname
	}

	return list, nil
}
