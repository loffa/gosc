package gosc

import (
	"bufio"
	"bytes"
	"errors"
	"net"
)

type transportUDP struct {
	conn       net.PacketConn
	bufferSize int
}

// NewUDPTransport returns the default UDP Transport for clients.
func NewUDPTransport(address string, bufferSize int) (Transport, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}
	return &transportUDP{
		conn:       conn.(net.PacketConn),
		bufferSize: bufferSize,
	}, nil
}

// NewUDPListen returns the default UDP Transport for servers.
func NewUDPListen(address string, bufferSize int) (Transport, error) {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return nil, err
	}
	return &transportUDP{
		conn:       conn,
		bufferSize: bufferSize,
	}, nil
}

// Send uses buffering to send a complete Package on the UDP socket.
func (t *transportUDP) Send(pack Package, addr net.Addr) error {
	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)

	err := writePackage(pack, w)
	if err != nil {
		return err
	}
	_ = w.Flush()
	_, err = t.conn.WriteTo(buf.Bytes(), addr)
	if errors.Is(err, net.ErrWriteToConnected) {
		conn := t.conn.(*net.UDPConn)
		_, err = conn.Write(buf.Bytes())
	}
	if err != nil {
		return err
	}
	return nil
}

// Receive reads the UDP socket buffered and returns a Package when found.
func (t *transportUDP) Receive() (pack Package, from net.Addr, err error) {
	buf := make([]byte, t.bufferSize)
	_, from, err = t.conn.ReadFrom(buf)
	if err != nil {
		return nil, from, err
	}
	r := bufio.NewReaderSize(bytes.NewReader(buf), len(buf))
	pack, err = readPackage(r)
	return pack, from, err
}
