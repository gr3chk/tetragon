// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package exporter

import (
	"os"
	"time"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"github.com/cilium/tetragon/pkg/version"
	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// MetadataEvent represents the agent initialization metadata
type MetadataEvent struct {
	Timestamp       time.Time `json:"@timestamp"`
	Event           string    `json:"event"`
	TetragonVersion string    `json:"tetragon_version"`
	BuildCommit     string    `json:"build_commit,omitempty"`
	BuildDate       string    `json:"build_date,omitempty"`
	Hostname        string    `json:"hostname"`
	OS              string    `json:"os"`
	KernelVersion   string    `json:"kernel_version"`
	PID             int       `json:"pid"`
	UDPDestination  string    `json:"udp_destination"`
	UDPBufferSize   int       `json:"udp_buffer_size"`
	Uptime          string    `json:"uptime"`
}

// NewMetadataEvent creates a new metadata event for agent initialization
func NewMetadataEvent(hostname string, udpDestination string, udpBufferSize int) *MetadataEvent {
	// Get build information
	buildInfo := version.ReadBuildInfo()

	// Get kernel version
	kernelVersion := "unknown"
	var uname unix.Utsname
	if err := unix.Uname(&uname); err == nil {
		kernelVersion = unix.ByteSliceToString(uname.Release[:])
	}

	return &MetadataEvent{
		Timestamp:       time.Now().UTC(),
		Event:           "agent_init",
		TetragonVersion: version.Version,
		BuildCommit:     buildInfo.Commit,
		BuildDate:       buildInfo.Time,
		Hostname:        hostname,
		OS:              "linux", // We'll make this configurable later
		KernelVersion:   kernelVersion,
		PID:             os.Getpid(),
		UDPDestination:  udpDestination,
		UDPBufferSize:   udpBufferSize,
		Uptime:          "initialized at 0",
	}
}

// ToGetEventsResponse converts the metadata event to a Tetragon GetEventsResponse
func (m *MetadataEvent) ToGetEventsResponse() *tetragon.GetEventsResponse {
	// Create a custom event that represents metadata
	// We'll use a generic event structure since metadata doesn't fit standard Tetragon events
	return &tetragon.GetEventsResponse{
		Event: &tetragon.GetEventsResponse_ProcessExec{
			ProcessExec: &tetragon.ProcessExec{
				Process: &tetragon.Process{
					Binary:    "tetragon_metadata",
					Arguments: "agent_init",
					Pid:       &wrapperspb.UInt32Value{Value: uint32(m.PID)},
				},
			},
		},
		// Add metadata as annotations or custom fields
		// Note: This is a simplified approach - in a real implementation,
		// you might want to create a custom event type for metadata
	}
}
