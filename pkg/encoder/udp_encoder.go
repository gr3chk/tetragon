// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

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

// UDPEncoder implements EventEncoder interface for sending events over UDP
type UDPEncoder struct {
	addr     *net.UDPAddr
	mu       sync.RWMutex
	closed   int32
	jsonOpts protojson.MarshalOptions
	connPool sync.Pool
	poolSize int
}

// NewUDPEncoder creates a new UDP encoder that sends events to the specified address and port
func NewUDPEncoder(address string, port int, bufferSize int) (*UDPEncoder, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address %s:%d: %w", address, port, err)
	}

	encoder := &UDPEncoder{
		addr:     addr,
		poolSize: 10, // Connection pool size
		jsonOpts: protojson.MarshalOptions{
			UseProtoNames: true, // Maintain backward compatibility with snake_case
		},
	}

	// Initialize connection pool with buffer size configuration
	encoder.connPool.New = func() interface{} {
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return nil
		}

		// Set socket buffer size if specified
		if bufferSize > 0 {
			if err := conn.SetWriteBuffer(bufferSize); err != nil {
				logger.GetLogger().Warn("Failed to set UDP write buffer size",
					"size", bufferSize,
					logfields.Error, err)
			}
		}

		return conn
	}

	return encoder, nil
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

	// Get connection from pool
	connObj := u.connPool.Get()
	if connObj == nil {
		// Fallback to creating new connection if pool is empty
		conn, err := net.DialUDP("udp", nil, u.addr)
		if err != nil {
			logger.GetLogger().Warn("Failed to create UDP connection",
				"address", u.addr.String(),
				logfields.Error, err)
			return err
		}
		defer conn.Close()
		_, err = conn.Write(data)
		return err
	}

	conn := connObj.(*net.UDPConn)
	defer u.connPool.Put(conn)

	// Send the data over UDP
	_, err = conn.Write(data)
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

// Write implements io.Writer interface for compatibility with existing exporter
func (u *UDPEncoder) Write(p []byte) (n int, err error) {
	if atomic.LoadInt32(&u.closed) == 1 {
		return 0, fmt.Errorf("UDP encoder is closed")
	}

	// Get connection from pool
	connObj := u.connPool.Get()
	if connObj == nil {
		// Fallback to creating new connection if pool is empty
		conn, err := net.DialUDP("udp", nil, u.addr)
		if err != nil {
			logger.GetLogger().Warn("Failed to create UDP connection",
				"address", u.addr.String(),
				logfields.Error, err)
			return 0, err
		}
		defer conn.Close()
		return conn.Write(p)
	}

	conn := connObj.(*net.UDPConn)
	defer u.connPool.Put(conn)

	return conn.Write(p)
}

// GetRemoteAddr returns the remote UDP address
func (u *UDPEncoder) GetRemoteAddr() string {
	return u.addr.String()
}
