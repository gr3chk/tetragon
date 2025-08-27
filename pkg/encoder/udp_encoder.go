// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package encoder

import (
	"fmt"
	"net"
	"sync"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"github.com/cilium/tetragon/pkg/logger"
	"github.com/cilium/tetragon/pkg/logger/logfields"

	"google.golang.org/protobuf/encoding/protojson"
)

// UDPEncoder implements EventEncoder interface for sending events over UDP
type UDPEncoder struct {
	conn     *net.UDPConn
	addr     *net.UDPAddr
	mu       sync.Mutex
	closed   bool
	jsonOpts protojson.MarshalOptions
}

// NewUDPEncoder creates a new UDP encoder that sends events to the specified address and port
func NewUDPEncoder(address string, port int) (*UDPEncoder, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address %s:%d: %w", address, port, err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection to %s:%d: %w", address, port, err)
	}

	return &UDPEncoder{
		conn: conn,
		addr: addr,
		jsonOpts: protojson.MarshalOptions{
			UseProtoNames: true, // Maintain backward compatibility with snake_case
		},
	}, nil
}

// Encode implements EventEncoder.Encode
func (u *UDPEncoder) Encode(v interface{}) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		return fmt.Errorf("UDP encoder is closed")
	}

	event, ok := v.(*tetragon.GetEventsResponse)
	if !ok {
		return ErrInvalidEvent
	}

	// Marshal the event to JSON
	data, err := u.jsonOpts.Marshal(event)
	if err != nil {
		logger.GetLogger().Warn("Failed to marshal event to JSON", logfields.Error, err)
		return err
	}

	// Add newline for proper log formatting
	data = append(data, '\n')

	// Send the data over UDP
	_, err = u.conn.Write(data)
	if err != nil {
		logger.GetLogger().Warn("Failed to send event over UDP",
			"address", u.addr.String(),
			logfields.Error, err)
		return err
	}

	return nil
}

// Close closes the UDP connection
func (u *UDPEncoder) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		return nil
	}

	u.closed = true
	return u.conn.Close()
}

// Write implements io.Writer interface for compatibility with existing exporter
func (u *UDPEncoder) Write(p []byte) (n int, err error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		return 0, fmt.Errorf("UDP encoder is closed")
	}

	return u.conn.Write(p)
}

// GetRemoteAddr returns the remote UDP address
func (u *UDPEncoder) GetRemoteAddr() string {
	return u.addr.String()
}
