package client

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync"
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

type SessionManager struct {
	mu         sync.Mutex
	sessionIDs map[string]string
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessionIDs: make(map[string]string),
	}
}

func (sm *SessionManager) GenerateSessionID() string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Generate a unique session ID
	sessionID := make([]byte, 16)
	rand.Read(sessionID)

	// Convert to hex string for readability and uniqueness
	return fmt.Sprintf("%x", sessionID)
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

func (c *Client) WriteMessageWithHeader(sessionID string, messageType string, payload []byte) error {
	// Calculate total message length (header + payload)
	totalLength := len(payload) + 19

	// Create header
	header := make([]byte, 19)
	
	// Convert sessionID to bytes
	sessionBytes := []byte(sessionID)
	
	// Ensure sessionID is exactly 16 bytes
	if len(sessionBytes) > 16 {
		sessionBytes = sessionBytes[:16]
	} else if len(sessionBytes) < 16 {
		paddedSessionID := make([]byte, 16)
		copy(paddedSessionID, sessionBytes)
		sessionBytes = paddedSessionID
	}

	// First 16 bytes: Session ID
	copy(header[0:16], sessionBytes)
	
	// Next 3 bytes: total message length (big-endian)
	binary.BigEndian.PutUint32(header[16:19], uint32(totalLength))

	// Combine header and payload
	fullMessage := append(header, payload...)

	// Write the full message
	_, err := c.conn.Write(fullMessage)
	return err
}

func (c *Client) ReadMessageWithHeader() (sessionID string, payload []byte, err error) {
	// Read header
	header := make([]byte, 19)
	_, err = c.conn.Read(header)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read header: %v", err)
	}

	// Extract session ID (first 16 bytes)
	sessionID = fmt.Sprintf("%x", header[0:16])

	// Extract total message length (last 3 bytes)
	totalLength := binary.BigEndian.Uint32(header[16:19])

	// Read payload
	payload = make([]byte, totalLength-19)
	_, err = c.conn.Read(payload)
	if err != nil {
		return sessionID, nil, fmt.Errorf("failed to read payload: %v", err)
	}

	return sessionID, payload, nil
}
