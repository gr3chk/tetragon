// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package encoder

import (
	"net"
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
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port)
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
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port)
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

func TestUDPEncoder_InvalidEvent(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port)
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
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port)
	require.NoError(t, err)
	defer encoder.Close()

	// Write test data
	testData := []byte("test message\n")
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

func TestUDPEncoder_Close(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port)
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

	// Try to write after closing
	_, err = encoder.Write([]byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UDP encoder is closed")
}

func TestUDPEncoder_InvalidAddress(t *testing.T) {
	// Try to create UDP encoder with invalid address
	_, err := NewUDPEncoder("invalid-address", 12345)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to resolve UDP address")
}
