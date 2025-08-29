// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package exporter

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"github.com/cilium/tetragon/pkg/encoder"
	"github.com/cilium/tetragon/pkg/ratelimit"
	"github.com/cilium/tetragon/pkg/server"
)

func TestNewUDPExporter(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	udpEncoder, err := encoder.NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)
	defer udpEncoder.Close()

	// Create a mock server
	ctx := context.Background()
	req := &tetragon.GetEventsRequest{}
	mockServer := &server.Server{}

	// Create UDP exporter
	exporter := NewUDPExporter(ctx, req, mockServer, udpEncoder, nil)

	assert.NotNil(t, exporter)
	assert.Equal(t, ctx, exporter.Context())
}

func TestUDPExporter_Send(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	udpEncoder, err := encoder.NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)
	defer udpEncoder.Close()

	// Create a mock server
	ctx := context.Background()
	req := &tetragon.GetEventsRequest{}
	mockServer := &server.Server{}

	// Create UDP exporter
	exporter := NewUDPExporter(ctx, req, mockServer, udpEncoder, nil)

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

	// Send the event
	err = exporter.Send(event)
	require.NoError(t, err)

	// Read the data from the server
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	require.NoError(t, err)

	// Verify the received data
	receivedData := buffer[:n]
	assert.Contains(t, string(receivedData), "/bin/test")
	assert.Contains(t, string(receivedData), "test arg")
}

func TestUDPExporter_WithRateLimit(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	udpEncoder, err := encoder.NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)
	defer udpEncoder.Close()

	// Create a rate limiter (1 event per second)
	ctx := context.Background()
	rateLimiter := ratelimit.NewRateLimiter(ctx, time.Second, 1, udpEncoder)

	// Create a mock server
	req := &tetragon.GetEventsRequest{}
	mockServer := &server.Server{}

	// Create UDP exporter with rate limiter
	exporter := NewUDPExporter(ctx, req, mockServer, udpEncoder, rateLimiter)

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

	// Send the first event (should succeed)
	err = exporter.Send(event)
	require.NoError(t, err)

	// Send the second event immediately (should be rate limited)
	err = exporter.Send(event)
	require.NoError(t, err) // Rate limiter drops events silently

	// Read only one event from the server (the second should be dropped)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	require.NoError(t, err)

	// Verify only one event was received
	receivedData := buffer[:n]
	assert.Contains(t, string(receivedData), "/bin/test")
}

func TestUDPExporter_Close(t *testing.T) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	conn, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	udpEncoder, err := encoder.NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port, 65536)
	require.NoError(t, err)

	// Create a mock server
	ctx := context.Background()
	req := &tetragon.GetEventsRequest{}
	mockServer := &server.Server{}

	// Create UDP exporter
	exporter := NewUDPExporter(ctx, req, mockServer, udpEncoder, nil)

	// Close the exporter
	err = exporter.Close()
	require.NoError(t, err)

	// Try to send an event after closing
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

	err = exporter.Send(event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UDP exporter is closed")
}
