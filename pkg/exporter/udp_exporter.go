// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package exporter

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc/metadata"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"github.com/cilium/tetragon/pkg/encoder"
	"github.com/cilium/tetragon/pkg/logger"
	"github.com/cilium/tetragon/pkg/logger/logfields"
	"github.com/cilium/tetragon/pkg/ratelimit"
	"github.com/cilium/tetragon/pkg/server"
)

// UDPExporter implements server.Listener interface for UDP output
type UDPExporter struct {
	ctx         context.Context
	request     *tetragon.GetEventsRequest
	server      *server.Server
	encoder     *encoder.UDPEncoder
	rateLimiter *ratelimit.RateLimiter
	mu          sync.Mutex
	closed      bool
}

// NewUDPExporter creates a new UDP exporter
func NewUDPExporter(
	ctx context.Context,
	request *tetragon.GetEventsRequest,
	server *server.Server,
	udpEncoder *encoder.UDPEncoder,
	rateLimiter *ratelimit.RateLimiter,
) *UDPExporter {
	return &UDPExporter{
		ctx:         ctx,
		request:     request,
		server:      server,
		encoder:     udpEncoder,
		rateLimiter: rateLimiter,
	}
}

// Start starts the UDP exporter
func (e *UDPExporter) Start() error {
	var readyWG sync.WaitGroup
	var exporterStartErr error
	readyWG.Add(1)
	go func() {
		if err := e.server.GetEventsWG(e.request, e, e.encoder, &readyWG); err != nil {
			exporterStartErr = fmt.Errorf("error starting UDP exporter: %w", err)
		}
	}()
	readyWG.Wait()
	return exporterStartErr
}

// SendMetadataEvent sends a metadata event over UDP
func (e *UDPExporter) SendMetadataEvent(hostname string, udpDestination string, udpBufferSize int) error {
	metadataEvent := NewMetadataEvent(hostname, udpDestination, udpBufferSize)
	event := metadataEvent.ToGetEventsResponse()

	// Send the metadata event directly through the encoder
	if err := e.encoder.Encode(event); err != nil {
		logger.GetLogger().Warn("Failed to encode metadata event for UDP", logfields.Error, err)
		return err
	}

	logger.GetLogger().Info("Metadata event sent over UDP",
		"event", "agent_init",
		"hostname", hostname,
		"udp_destination", udpDestination)

	return nil
}

// Send implements server.Listener.Send
func (e *UDPExporter) Send(event *tetragon.GetEventsResponse) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return fmt.Errorf("UDP exporter is closed")
	}

	if e.rateLimiter != nil && !e.rateLimiter.Allow() {
		e.rateLimiter.Drop()
		rateLimitDropped.Inc()
		return nil
	}

	if err := e.encoder.Encode(event); err != nil {
		logger.GetLogger().Warn("Failed to encode event for UDP", logfields.Error, err)
		return err
	}

	eventsExportedTotal.Inc()
	eventsExportTimestamp.Set(float64(event.GetTime().GetSeconds()))
	return nil
}

// Close closes the UDP exporter
func (e *UDPExporter) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil
	}

	e.closed = true
	return e.encoder.Close()
}

// SetHeader implements server.Listener.SetHeader
func (e *UDPExporter) SetHeader(metadata.MD) error {
	return nil
}

// SendHeader implements server.Listener.SendHeader
func (e *UDPExporter) SendHeader(metadata.MD) error {
	return nil
}

// SetTrailer implements server.Listener.SetTrailer
func (e *UDPExporter) SetTrailer(metadata.MD) {
}

// Context implements server.Listener.Context
func (e *UDPExporter) Context() context.Context {
	return e.ctx
}

// SendMsg implements server.Listener.SendMsg
func (e *UDPExporter) SendMsg(_ interface{}) error {
	return nil
}

// RecvMsg implements server.Listener.RecvMsg
func (e *UDPExporter) RecvMsg(_ interface{}) error {
	return nil
}
