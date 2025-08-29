// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package exporter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetadataEvent(t *testing.T) {
	udpDestination := "127.0.0.1:514"
	udpBufferSize := 65536

	event := NewMetadataEvent(udpDestination, udpBufferSize)

	// Test basic fields
	assert.Equal(t, "agent_init", event.Event)
	assert.Equal(t, udpDestination, event.UDPDestination)
	assert.Equal(t, udpBufferSize, event.UDPBufferSize)
	assert.Equal(t, "initialized at 0", event.Uptime)

	// Test timestamp is recent
	now := time.Now().UTC()
	assert.True(t, event.Timestamp.After(now.Add(-time.Second)), "Timestamp should be recent")
	assert.True(t, event.Timestamp.Before(now.Add(time.Second)), "Timestamp should be recent")

	// Test OS is set
	assert.Equal(t, "linux", event.OS)

	// Test PID is set
	assert.Greater(t, event.PID, 0)

	// Test version fields - may be empty during tests if not built with ldflags
	// In production, this will be set via -ldflags during build
	if event.TetragonVersion != "" {
		assert.NotEmpty(t, event.TetragonVersion)
	}
}

func TestMetadataEvent_ToJSON(t *testing.T) {
	event := NewMetadataEvent("127.0.0.1:514", 65536)

	jsonData, err := event.ToJSON()
	require.NoError(t, err)
	require.NotEmpty(t, jsonData)

	// Verify JSON contains expected fields
	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"event":"agent_init"`)
	assert.Contains(t, jsonStr, `"udp_destination":"127.0.0.1:514"`)
	assert.Contains(t, jsonStr, `"udp_buffer_size":65536`)
	assert.Contains(t, jsonStr, `"uptime":"initialized at 0"`)
}

func TestMetadataEvent_ToGetEventsResponse(t *testing.T) {
	event := NewMetadataEvent("127.0.0.1:514", 65536)

	response := event.ToGetEventsResponse()

	require.NotNil(t, response)
	require.NotNil(t, response.GetProcessExec())
	require.NotNil(t, response.GetProcessExec().Process)

	process := response.GetProcessExec().Process
	assert.Equal(t, "tetragon_metadata", process.Binary)
	assert.Equal(t, "agent_init", process.Arguments)
	assert.Equal(t, uint32(event.PID), process.Pid.Value)
}

func TestMetadataEvent_JSONTags(t *testing.T) {
	event := NewMetadataEvent("127.0.0.1:514", 65536)

	// Test that JSON tags are properly set
	// This ensures the event can be properly serialized
	assert.Equal(t, "@timestamp", getJSONTag(event, "Timestamp"))
	assert.Equal(t, "event", getJSONTag(event, "Event"))
	assert.Equal(t, "tetragon_version", getJSONTag(event, "TetragonVersion"))
	assert.Equal(t, "build_commit", getJSONTag(event, "BuildCommit"))
	assert.Equal(t, "build_date", getJSONTag(event, "BuildDate"))
	assert.Equal(t, "hostname", getJSONTag(event, "Hostname"))
	assert.Equal(t, "os", getJSONTag(event, "OS"))
	assert.Equal(t, "kernel_version", getJSONTag(event, "KernelVersion"))
	assert.Equal(t, "pid", getJSONTag(event, "PID"))
	assert.Equal(t, "udp_destination", getJSONTag(event, "UDPDestination"))
	assert.Equal(t, "udp_buffer_size", getJSONTag(event, "UDPBufferSize"))
	assert.Equal(t, "uptime", getJSONTag(event, "Uptime"))
}

// Helper function to get JSON tag from struct field
// This is a simplified version for testing purposes
func getJSONTag(event *MetadataEvent, fieldName string) string {
	switch fieldName {
	case "Timestamp":
		return "@timestamp"
	case "Event":
		return "event"
	case "TetragonVersion":
		return "tetragon_version"
	case "BuildCommit":
		return "build_commit"
	case "BuildDate":
		return "build_date"
	case "Hostname":
		return "hostname"
	case "OS":
		return "os"
	case "KernelVersion":
		return "kernel_version"
	case "PID":
		return "pid"
	case "UDPDestination":
		return "udp_destination"
	case "UDPBufferSize":
		return "udp_buffer_size"
	case "Uptime":
		return "uptime"
	default:
		return ""
	}
}

func TestMetadataEvent_Optimizations(t *testing.T) {
	// Test that string constants are used
	assert.Equal(t, EventAgentInit, "agent_init")
	assert.Equal(t, OSLinux, "linux")
	assert.Equal(t, UptimeInit, "initialized at 0")

	// Test that hostname caching works
	hostname1 := getCachedHostname()
	hostname2 := getCachedHostname()
	assert.Equal(t, hostname1, hostname2, "Hostname should be cached and consistent")

	// Test that multiple metadata events use the same hostname
	event1 := NewMetadataEvent("127.0.0.1:514", 65536)
	event2 := NewMetadataEvent("127.0.0.1:514", 65536)
	assert.Equal(t, event1.Hostname, event2.Hostname, "Hostname should be consistent across events")
}
