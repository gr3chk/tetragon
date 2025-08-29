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

	// Cached metadata for performance optimization
	cachedMetadata []byte
	metadataOnce   sync.Once
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

// initCachedMetadata initializes the cached metadata once for performance optimization
func (e *UDPExporter) initCachedMetadata(udpDestination string, udpBufferSize int) {
	e.metadataOnce.Do(func() {
		metadataEvent := NewMetadataEvent(udpDestination, udpBufferSize)
		if jsonData, err := metadataEvent.ToJSON(); err == nil {
			e.cachedMetadata = jsonData
		} else {
			logger.GetLogger().Warn("Failed to cache metadata event", logfields.Error, err)
		}
	})
}

// SendMetadataEvent sends a metadata event over UDP
func (e *UDPExporter) SendMetadataEvent(udpDestination string, udpBufferSize int) error {
	// Initialize cached metadata once
	e.initCachedMetadata(udpDestination, udpBufferSize)

	// Use cached metadata if available
	if e.cachedMetadata != nil {
		if err := e.encoder.WriteRaw(e.cachedMetadata); err != nil {
			logger.GetLogger().Warn("Failed to send cached metadata event over UDP", logfields.Error, err)
			return err
		}

		logger.GetLogger().Info("Cached metadata event sent over UDP",
			"event", "agent_init",
			"hostname", getCachedHostname(),
			"udp_destination", udpDestination)

		return nil
	}

	// Fallback to dynamic creation if caching failed
	metadataEvent := NewMetadataEvent(udpDestination, udpBufferSize)
	jsonData, err := metadataEvent.ToJSON()
	if err != nil {
		logger.GetLogger().Warn("Failed to marshal metadata event to JSON", logfields.Error, err)
		return err
	}

	if err := e.encoder.WriteRaw(jsonData); err != nil {
		logger.GetLogger().Warn("Failed to send metadata event over UDP", logfields.Error, err)
		return err
	}

	logger.GetLogger().Info("Metadata event sent over UDP",
		"event", "agent_init",
		"hostname", getCachedHostname(),
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
