// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package encoder

import (
	"net"
	"testing"

	"github.com/cilium/tetragon/api/v1/tetragon"
)

func BenchmarkUDPEncoder_Encode(b *testing.B) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port)
	if err != nil {
		b.Fatal(err)
	}
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

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if err := encoder.Encode(event); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUDPEncoder_EncodeLargeEvent(b *testing.B) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port)
	if err != nil {
		b.Fatal(err)
	}
	defer encoder.Close()

	// Create a large test event (simulating ~1500 bytes)
	largeArgs := ""
	for i := 0; i < 100; i++ {
		largeArgs += "very-long-argument-that-makes-the-event-larger "
	}

	event := &tetragon.GetEventsResponse{
		Event: &tetragon.GetEventsResponse_ProcessExec{
			ProcessExec: &tetragon.ProcessExec{
				Process: &tetragon.Process{
					Binary:    "/bin/very-long-binary-name-that-contributes-to-size",
					Arguments: largeArgs,
				},
			},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if err := encoder.Encode(event); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUDPEncoder_EncodeVeryLargeEvent(b *testing.B) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port)
	if err != nil {
		b.Fatal(err)
	}
	defer encoder.Close()

	// Create a very large test event (simulating ~9000 bytes)
	veryLargeArgs := ""
	for i := 0; i < 500; i++ {
		veryLargeArgs += "extremely-long-argument-that-makes-the-event-very-large "
	}

	event := &tetragon.GetEventsResponse{
		Event: &tetragon.GetEventsResponse_ProcessExec{
			ProcessExec: &tetragon.ProcessExec{
				Process: &tetragon.Process{
					Binary:    "/bin/extremely-long-binary-name-that-contributes-significantly-to-size",
					Arguments: veryLargeArgs,
				},
			},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if err := encoder.Encode(event); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUDPEncoder_Write(b *testing.B) {
	// Start a test UDP server
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close()

	serverAddr := conn.LocalAddr().(*net.UDPAddr)

	// Create UDP encoder
	encoder, err := NewUDPEncoder(serverAddr.IP.String(), serverAddr.Port)
	if err != nil {
		b.Fatal(err)
	}
	defer encoder.Close()

	// Test data
	testData := []byte("test message\n")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if _, err := encoder.Write(testData); err != nil {
			b.Fatal(err)
		}
	}
}
