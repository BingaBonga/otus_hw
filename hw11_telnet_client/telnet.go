package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClientImpl struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClientImpl{address: address, timeout: timeout, in: in, out: out}
}

func (t *telnetClientImpl) Receive() error {
	if _, err := io.Copy(t.out, t.Conn); err != nil {
		return fmt.Errorf("write message error: %w", err)
	}
	return nil
}

func (t *telnetClientImpl) Send() error {
	if _, err := io.Copy(t.Conn, t.in); err != nil {
		return fmt.Errorf("send message error: %w", err)
	}
	return nil
}

func (t *telnetClientImpl) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	t.Conn = conn
	return nil
}

func (t *telnetClientImpl) Close() error {
	return t.Conn.Close()
}
