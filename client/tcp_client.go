package client

import (
	"context"
	"fmt"
	"net"
	"time"
)

const (
	defaultDialTimeout  = 10 * time.Second
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
)

type Client struct {
	conn         net.Conn
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewClient(dialTimeout, readTimeout, writeTimeout time.Duration) *Client {
	if dialTimeout == 0 {
		dialTimeout = defaultDialTimeout
	}
	if readTimeout == 0 {
		readTimeout = defaultReadTimeout
	}
	if writeTimeout == 0 {
		writeTimeout = defaultWriteTimeout
	}

	return &Client{
		dialTimeout:  dialTimeout,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (c *Client) Connect(ctx context.Context, host, port string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	
	dialer := net.Dialer{Timeout: c.dialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to server at %s: %v", address, err)
	}
	
	c.conn = conn
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) Write(data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("connection not established")
	}

	if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %v", err)
	}

	_, err := c.conn.Write(data)
	return err
}

func (c *Client) Read(buffer []byte) (int, error) {
	if c.conn == nil {
		return 0, fmt.Errorf("connection not established")
	}

	if err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
		return 0, fmt.Errorf("failed to set read deadline: %v", err)
	}

	return c.conn.Read(buffer)
}
