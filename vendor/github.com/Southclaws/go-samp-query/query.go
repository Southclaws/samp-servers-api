package sampquery

import (
	"bytes"
	"context"
	"encoding/binary"
	"math/rand"
	"net"

	"github.com/pkg/errors"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/htmlindex"
)

// Server contains all the information retreived from the server query API.
type Server struct {
	Address    string
	Hostname   string
	Players    int
	MaxPlayers int
	Gamemode   string
	Language   string
	Password   bool
	Rules      map[string]string
}

// QueryType represents a query method from the SA:MP set: i, r, c, d, x, p
type QueryType uint8

const (
	// Info is the 'i' packet type
	Info QueryType = 'i'
	// Rules is the 'r' packet type
	Rules QueryType = 'r'
	// Players is the 'c' packet type
	Players QueryType = 'c'
	// Ping is the 'p' packet type
	Ping QueryType = 'p'
)

// Query stores state for masterlist queries
type Query struct {
	addr *net.UDPAddr
	Data Server
}

// GetServerInfo wraps a set of legacy queries and returns a new Server object with the available
// fields populated.
func GetServerInfo(ctx context.Context, host string, attemptDecode bool) (server Server, err error) {
	query, err := NewLegacyQuery(host)
	if err != nil {
		return
	}
	defer func() {
		query.Close()
	}()

	err = query.GetPing(ctx)
	if err != nil {
		return
	}

	server, err = query.GetInfo(ctx, attemptDecode)
	if err != nil {
		return
	}
	server.Address = host

	server.Rules, err = query.GetRules(ctx)
	if err != nil {
		return
	}

	return
}

// NewLegacyQuery creates a new legacy query handler for a server
func NewLegacyQuery(host string) (query *Query, err error) {
	query = new(Query)

	query.addr, err = net.ResolveUDPAddr("udp", host)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve host")
	}

	return query, nil
}

// Close closes a legacy query manager's connection
func (query *Query) Close() error {
	return nil
}

// SendQuery writes a SA:MP format query with the specified opcode, returns the raw response bytes
func (query *Query) SendQuery(ctx context.Context, opcode QueryType) (response []byte, err error) {
	request := new(bytes.Buffer)

	port := [2]byte{
		byte(query.addr.Port & 0xFF),
		byte((query.addr.Port >> 8) & 0xFF),
	}

	if err = binary.Write(request, binary.LittleEndian, []byte("SAMP")); err != nil {
		return
	}
	if err = binary.Write(request, binary.LittleEndian, query.addr.IP.To4()); err != nil {
		return
	}
	if err = binary.Write(request, binary.LittleEndian, port[0]); err != nil {
		return
	}
	if err = binary.Write(request, binary.LittleEndian, port[1]); err != nil {
		return
	}
	if err = binary.Write(request, binary.LittleEndian, opcode); err != nil {
		return
	}
	if opcode == Ping {
		p := make([]byte, 4)
		_, err = rand.Read(p)
		if err != nil {
			return
		}
		if err = binary.Write(request, binary.LittleEndian, p); err != nil {
			return
		}
	}

	conn, err := openConnection(query.addr)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Write(request.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "failed to write")
	}

	type resultData struct {
		data  []byte
		bytes int
		err   error
	}
	waitResult := make(chan resultData)

	go func() {
		response := make([]byte, 2048)
		n, errInner := conn.Read(response)
		if errInner != nil {
			waitResult <- resultData{err: errors.Wrap(errInner, "failed to read response")}
		}
		if n > cap(response) {
			waitResult <- resultData{err: errors.New("read response over buffer capacity")}
		}
		waitResult <- resultData{data: response, bytes: n}
	}()

	var result resultData
	select {
	case <-ctx.Done():
		return nil, errors.New("socket read timed out")

	case result = <-waitResult:
		break
	}

	if result.err != nil {
		return nil, result.err
	}

	return result.data[:result.bytes], nil
}

// GetPing sends and receives a packet to measure ping
func (query *Query) GetPing(ctx context.Context) (err error) {
	_, err = query.SendQuery(ctx, Ping)
	if err != nil {
		return err
	}

	return
}

// GetInfo returns the core server info for displaying on the browser list.
func (query *Query) GetInfo(ctx context.Context, attemptDecode bool) (server Server, err error) {
	response, err := query.SendQuery(ctx, Info)
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

	hostnameRaw := response[ptr : ptr+hostnameLen]
	ptr += hostnameLen

	gamemodeLen := int(binary.LittleEndian.Uint16(response[ptr : ptr+4]))
	ptr += 4

	gamemodeRaw := response[ptr : ptr+gamemodeLen]
	ptr += gamemodeLen

	languageLen := int(binary.LittleEndian.Uint16(response[ptr : ptr+4]))
	ptr += 4

	languageRaw := response[ptr : ptr+languageLen]
	ptr += languageLen

	guessHelper := bytes.Join([][]byte{
		hostnameRaw,
		gamemodeRaw,
		languageRaw,
	}, []byte(" "))

	if attemptDecode {
		server.Gamemode = attemptDecodeANSI(gamemodeRaw, guessHelper)
		server.Hostname = attemptDecodeANSI(hostnameRaw, guessHelper)
	} else {
		server.Gamemode = string(gamemodeRaw)
		server.Hostname = string(hostnameRaw)
	}

	if languageLen > 0 && attemptDecode {
		server.Language = attemptDecodeANSI(languageRaw, guessHelper)
	} else {
		server.Language = "-"
	}
	return
}

// GetRules returns a map of rule properties from a server. The legacy query uses established keys
// such as "Map" and "Version"
func (query *Query) GetRules(ctx context.Context) (rules map[string]string, err error) {
	response, err := query.SendQuery(ctx, Rules)
	if err != nil {
		return rules, err
	}

	rules = make(map[string]string)

	var (
		key string
		val string
		len int
	)

	ptr := 11
	amount := binary.LittleEndian.Uint16(response[ptr : ptr+2])
	ptr += 2

	for i := uint16(0); i < amount; i++ {
		len = int(response[ptr])
		ptr++

		key = string(response[ptr : ptr+len])
		ptr += len

		len = int(response[ptr])
		ptr++

		val = string(response[ptr : ptr+len])
		ptr += len

		rules[key] = val
	}

	return
}

// GetPlayers simply returns a slice of strings, score is rather arbitrary so it's omitted.
func (query *Query) GetPlayers(ctx context.Context) (players []string, err error) {
	response, err := query.SendQuery(ctx, Players)
	if err != nil {
		return
	}

	var (
		count  uint16
		length int
	)

	ptr := 11
	count = binary.LittleEndian.Uint16(response[ptr : ptr+2])
	ptr += 2

	players = make([]string, count)

	for i := uint16(0); i < count; i++ {
		length = int(response[ptr])
		ptr++

		players[i] = string(response[ptr : ptr+length])
		ptr += length
		ptr += 4 // score, unused
	}

	return players, nil
}

func openConnection(addr *net.UDPAddr) (conn *net.UDPConn, err error) {
	conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial")
	}
	return
}

func attemptDecodeANSI(input []byte, extra []byte) (result string) {
	result = string(input)
	detector, err := chardet.NewTextDetector().DetectBest(extra)
	if err != nil {
		return
	}
	e, err := htmlindex.Get(detector.Charset)
	if err != nil {
		return
	}
	dec := e.NewDecoder()
	decoded, err := dec.Bytes(input)
	if err != nil {
		return
	}
	return string(decoded)
}
