// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

// Package encoder provides UDP event encoding functionality.
// Tests verify UDP packet transmission using WriteToUDP on unbound sockets.
package encoder

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cilium/tetragon/api/v1/tetragon"
)

func TestNewUDPEncoder(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	// Get the actual address the server is listening on
	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)
	defer encoder.Close()

	assert.NotNil(t, encoder)
	assert.Equal(t, serverAddr.String(), encoder.GetRemoteAddr())
}

func TestUDPEncoder_Encode(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)
	defer encoder.Close()

	// Create a test event
	event := &tetragon.GetEventsResponse{
		Event: &tetragon.GetEventsResponse_ProcessExec{
			ProcessExec: &tetragon.ProcessExec{
				Process: &tetragon.Process{
					Binary:    "/bin/test",
					Arguments: "test arg",
				},
			},
		},
	}

	// Encode the event
	err = encoder.Encode(event)
	require.NoError(t, err)

	// Read the data from the server
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	require.NoError(t, err)

	// Verify the received data
	receivedData := buffer[:n]
	assert.Contains(t, string(receivedData), "/bin/test")
	assert.Contains(t, string(receivedData), "test")
	assert.Contains(t, string(receivedData), "arg")
}

func TestUDPEncoder_EncodeWithoutListener(t *testing.T) {
	// Test that UDP packets can be sent even without a listener
	// This demonstrates the fire-and-forget behavior of WriteToUDP

	// Use a port that's unlikely to have a listener (high port number)
	testPort := 65535

	// Create UDP encoder targeting a port with no listener
	encoder, err := NewUDPEncoder("127.0.0.1", testPort, 65536)
	require.NoError(t, err)
	defer encoder.Close()

	// Create a test event
	event := &tetragon.GetEventsResponse{
		Event: &tetragon.GetEventsResponse_ProcessExec{
			ProcessExec: &tetragon.ProcessExec{
				Process: &tetragon.Process{
					Binary:    "/bin/test",
					Arguments: "test arg",
				},
			},
		},
	}

	// Encode the event - this should succeed even without a listener
	// because WriteToUDP doesn't require a connection
	err = encoder.Encode(event)
	require.NoError(t, err, "UDP packet should be sent successfully even without a listener")

	// Verify the encoder reports the correct remote address
	assert.Equal(t, "127.0.0.1:65535", encoder.GetRemoteAddr())
}

func TestUDPEncoder_InvalidEvent(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)
	defer encoder.Close()

	// Try to encode an invalid event
	err = encoder.Encode("invalid event")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidEvent, err)
}

func TestUDPEncoder_Write(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)
	defer encoder.Close()

	// Test data
	testData := []byte("test message\n")

	// Write data using WriteToUDP internally (no listener required)
	n, err := encoder.Write(testData)
	require.NoError(t, err)
	assert.Equal(t, len(testData), n)

	// Read the data from the server
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 1024)
	n, _, err = conn.ReadFromUDP(buffer)
	require.NoError(t, err)

	// Verify the received data
	receivedData := buffer[:n]
	assert.Equal(t, testData, receivedData)
}

func TestUDPEncoder_WriteClosed(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)

	// Close the encoder
	err = encoder.Close()
	require.NoError(t, err)

	// Try to write after closing
	testData := []byte("test message\n")
	_, err = encoder.Write(testData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UDP encoder is closed")
}

func TestUDPEncoder_Closed(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)

	// Close the encoder
	err = encoder.Close()
	require.NoError(t, err)

	// Try to encode after closing
	event := &tetragon.GetEventsResponse{
		Event: &tetragon.GetEventsResponse_ProcessExec{
			ProcessExec: &tetragon.ProcessExec{
				Process: &tetragon.Process{
					Binary: "/bin/test",
				},
			},
		},
	}

	err = encoder.Encode(event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UDP encoder is closed")
}

func TestUDPEncoder_InvalidAddress(t *testing.T) {
	// Try to create encoder with invalid address
	encoder, err := NewUDPEncoder("invalid-address", 12345, 65536)
	assert.Error(t, err)
	assert.Nil(t, encoder)
	assert.Contains(t, err.Error(), "failed to resolve UDP address")
}

func TestUDPEncoder_SinglePacketPerEvent(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)
	defer encoder.Close()

	// Create a test event with a very long binary path to test size limits
	longPath := strings.Repeat("/very/long/path/", 1000) + "binary"
	event := &tetragon.GetEventsResponse{
		Event: &tetragon.GetEventsResponse_ProcessExec{
			ProcessExec: &tetragon.ProcessExec{
				Process: &tetragon.Process{
					Binary:    longPath,
					Arguments: "test arg",
				},
			},
		},
	}

	// Encode the event
	err = encoder.Encode(event)
	require.NoError(t, err)

	// Read the data from the server
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 65536) // Large enough to capture any data
	n, _, err := conn.ReadFromUDP(buffer)
	require.NoError(t, err)

	// Verify the received data fits in a single UDP packet
	receivedData := buffer[:n]
	assert.LessOrEqual(t, len(receivedData), MaxUDPSize, "Event should fit in single UDP packet")
	assert.Contains(t, string(receivedData), "binary")
	assert.True(t, strings.HasSuffix(string(receivedData), "\n"), "Event should end with newline")
}

func TestUDPEncoder_BufferSizeConfiguration(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Test with custom buffer size
	customBufferSize := 131072 // 128KB
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, customBufferSize)
	require.NoError(t, err)
	defer encoder.Close()

	// Create a test event
	event := &tetragon.GetEventsResponse{
		Event: &tetragon.GetEventsResponse_ProcessExec{
			ProcessExec: &tetragon.ProcessExec{
				Process: &tetragon.Process{
					Binary:    "/bin/test",
					Arguments: "test arg",
				},
			},
		},
	}

	// Encode the event
	err = encoder.Encode(event)
	require.NoError(t, err)

	// Read the data from the server
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	require.NoError(t, err)

	// Verify the received data
	receivedData := buffer[:n]
	assert.Contains(t, string(receivedData), "/bin/test")
	assert.Contains(t, string(receivedData), "test")
	assert.Contains(t, string(receivedData), "arg")
}

func TestUDPEncoder_MinimalModeCompatibility(t *testing.T) {
	// Test that UDP encoder works correctly in minimal mode
	// This test verifies that the encoder can function without
	// requiring any additional services or listeners

	// Create UDP encoder with minimal configuration
	encoder, err := NewUDPEncoder("127.0.0.1", 65535, 65536)
	require.NoError(t, err)
	defer encoder.Close()

	// Create a test event
	event := &tetragon.GetEventsResponse{
		Event: &tetragon.GetEventsResponse_ProcessExec{
			ProcessExec: &tetragon.ProcessExec{
				Process: &tetragon.Process{
					Binary:    "/bin/minimal",
					Arguments: "minimal arg",
				},
			},
		},
	}

	// Encode the event (should succeed even without listener)
	err = encoder.Encode(event)
	require.NoError(t, err, "UDP encoder should work in minimal mode without requiring listeners")

	// Verify encoder state
	assert.Equal(t, "127.0.0.1:65535", encoder.GetRemoteAddr())
	assert.True(t, encoder.IsMinimalMode(), "Encoder should indicate minimal mode operation")
}
