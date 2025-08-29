// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

// Package encoder provides UDP event encoding functionality.
// Note: UDP is connectionless by nature. We use unbound UDP sockets with WriteToUDP
// to send packets without requiring a listener on the destination port.
// This allows fire-and-forget UDP packet transmission.
package encoder

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"github.com/cilium/tetragon/pkg/logger"
	"github.com/cilium/tetragon/pkg/logger/logfields"

	"google.golang.org/protobuf/encoding/protojson"
)

const (
	// MaxUDPSize is the maximum size for a single UDP packet to avoid fragmentation
	// Standard UDP packet size limit is 65507 bytes (65535 - 20 IP header - 8 UDP header)
	MaxUDPSize = 65507
)

// UDPEncoder implements EventEncoder interface for sending events over UDP
// It uses unbound UDP sockets with WriteToUDP for fire-and-forget packet transmission.
type UDPEncoder struct {
	addr       *net.UDPAddr
	mu         sync.RWMutex
	closed     int32
	jsonOpts   protojson.MarshalOptions
	connPool   sync.Pool
	poolSize   int
	bufferSize int
}

// NewUDPEncoder creates a new UDP encoder that sends events to the specified address and port
func NewUDPEncoder(address string, port int, bufferSize int) (*UDPEncoder, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address %s:%d: %w", address, port, err)
	}

	encoder := &UDPEncoder{
		addr:     addr,
		poolSize: 10, // UDP socket pool size
		jsonOpts: protojson.MarshalOptions{
			UseProtoNames: true, // Maintain backward compatibility with snake_case
		},
		bufferSize: bufferSize,
	}

	// Initialize UDP socket pool with unbound sockets for WriteToUDP
	encoder.connPool.New = func() interface{} {
		// Create unbound UDP socket (bound to any available port)
		localAddr, err := net.ResolveUDPAddr("udp", ":0")
		if err != nil {
			logger.GetLogger().Debug("Failed to resolve local address for UDP socket",
				logfields.Error, err)
			return nil
		}

		conn, err := net.ListenUDP("udp", localAddr)
		if err != nil {
			logger.GetLogger().Debug("Failed to create unbound UDP socket for pool",
				logfields.Error, err)
			return nil
		}

		// Set socket buffer size if specified
		if bufferSize > 0 {
			if err := conn.SetWriteBuffer(bufferSize); err != nil {
				logger.GetLogger().Warn("Failed to set UDP socket write buffer size",
					"size", bufferSize,
					logfields.Error, err)
			}
		}

		return conn
	}

	return encoder, nil
}

// createUnboundUDPSocket creates an unbound UDP socket for WriteToUDP operations
func (u *UDPEncoder) createUnboundUDPSocket() (*net.UDPConn, error) {
	// Create unbound UDP socket (bound to any available port)
	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve local address: %w", err)
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create unbound UDP socket: %w", err)
	}

	// Set socket buffer size if specified
	if u.bufferSize > 0 {
		if err := conn.SetWriteBuffer(u.bufferSize); err != nil {
			logger.GetLogger().Warn("Failed to set UDP socket write buffer size",
				"size", u.bufferSize,
				logfields.Error, err)
		}
	}

	return conn, nil
}

// Encode implements EventEncoder.Encode
func (u *UDPEncoder) Encode(v interface{}) error {
	if atomic.LoadInt32(&u.closed) == 1 {
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

	// Ensure single-packet per event by checking size
	if len(data) > MaxUDPSize {
		logger.GetLogger().Warn("Event too large for single UDP packet, truncating",
			"size", len(data),
			"max_size", MaxUDPSize)
		// Truncate to fit in single packet, preserving newline
		data = data[:MaxUDPSize-1]
		data = append(data, '\n')
	}

	// Get UDP socket from pool
	connObj := u.connPool.Get()
	if connObj == nil {
		// Fallback: create new unbound UDP socket if pool is empty
		conn, err := u.createUnboundUDPSocket()
		if err != nil {
			logger.GetLogger().Warn("Failed to create unbound UDP socket",
				"address", u.addr.String(),
				logfields.Error, err)
			return err
		}
		defer conn.Close()
		_, err = conn.WriteToUDP(data, u.addr)
		return err
	}

	conn := connObj.(*net.UDPConn)
	defer u.connPool.Put(conn)

	// Send the data over UDP using WriteToUDP (no listener required)
	_, err = conn.WriteToUDP(data, u.addr)
	if err != nil {
		logger.GetLogger().Warn("Failed to send event over UDP",
			"address", u.addr.String(),
			logfields.Error, err)
		return err
	}

	return nil
}

// Close closes the UDP encoder
func (u *UDPEncoder) Close() error {
	atomic.StoreInt32(&u.closed, 1)
	return nil
}

// IsMinimalMode returns true if the encoder is operating in minimal mode
// Minimal mode means the encoder uses unbound sockets and WriteToUDP for
// fire-and-forget packet transmission without requiring listeners
func (u *UDPEncoder) IsMinimalMode() bool {
	return true // UDP encoder always operates in minimal mode
}

// Write implements io.Writer interface for compatibility with existing exporter
func (u *UDPEncoder) Write(p []byte) (n int, err error) {
	if atomic.LoadInt32(&u.closed) == 1 {
		return 0, fmt.Errorf("UDP encoder is closed")
	}

	// Ensure single-packet per write by checking size
	if len(p) > MaxUDPSize {
		logger.GetLogger().Warn("Data too large for single UDP packet, truncating",
			"size", len(p),
			"max_size", MaxUDPSize)
		p = p[:MaxUDPSize]
	}

	// Get UDP socket from pool
	connObj := u.connPool.Get()
	if connObj == nil {
		// Fallback: create new unbound UDP socket
		conn, err := u.createUnboundUDPSocket()
		if err != nil {
			logger.GetLogger().Warn("Failed to create unbound UDP socket",
				"address", u.addr.String(),
				logfields.Error, err)
			return 0, err
		}
		defer conn.Close()
		return conn.WriteToUDP(p, u.addr)
	}

	conn := connObj.(*net.UDPConn)
	defer u.connPool.Put(conn)

	return conn.WriteToUDP(p, u.addr)
}

// GetRemoteAddr returns the remote UDP address
func (u *UDPEncoder) GetRemoteAddr() string {
	return u.addr.String()
}
