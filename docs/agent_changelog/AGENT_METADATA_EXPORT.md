# Agent Metadata Export

## Overview

Tetragon now automatically exports agent initialization metadata over UDP on startup when UDP output is enabled. This provides essential context for all subsequent events and enables better integration with log aggregation and monitoring systems.

## What Gets Exported

On agent startup, Tetragon sends a metadata event containing:

| Field | Description | Example |
|-------|-------------|---------|
| `@timestamp` | ISO 8601 UTC timestamp of agent startup | `2024-01-15T10:30:00Z` |
| `event` | Event identifier | `agent_init` |
| `tetragon_version` | Current Tetragon version | `v1.1.0` |
| `build_commit` | Git commit hash (if available) | `a1b2c3d4e5f6` |
| `build_date` | Build timestamp (if available) | `2024-01-15T09:00:00Z` |
| `hostname` | System hostname | `tetragon-node-1` |
| `os` | Operating system identifier | `linux` |
| `kernel_version` | Linux kernel version | `6.14.0-28-generic` |
| `pid` | Process ID of the running agent | `12345` |
| `udp_destination` | Configured UDP destination (host:port) | `127.0.0.1:514` |
| `udp_buffer_size` | Configured UDP buffer size | `65536` |
| `uptime` | Initialization status | `initialized at 0` |

## How It Works

1. **Startup Sequence**: During Tetragon initialization, metadata is collected
2. **UDP Exporter Start**: UDP exporter is initialized and started
3. **Metadata Export**: Metadata event is immediately sent as the first UDP packet
4. **Runtime Events**: All subsequent security events are exported normally

## Implementation Details

### Metadata Collection

Metadata is collected using:
- `os.Hostname()` for system hostname
- `version.ReadBuildInfo()` for build information
- `unix.Uname()` for kernel version
- `os.Getpid()` for process ID
- Configuration values for UDP settings

### Event Format

The metadata is exported as **raw JSON** over UDP, not as a Tetragon event structure. This ensures the metadata is easily parseable by log aggregation systems and monitoring tools.

The JSON format includes all the specified fields:
- `@timestamp`: ISO 8601 UTC timestamp
- `event`: "agent_init" identifier
- `tetragon_version`: Current version
- `build_commit` & `build_date`: Build information
- `hostname`: System hostname
- `os`: Operating system
- `kernel_version`: Linux kernel version
- `pid`: Process ID
- `udp_destination`: UDP target (host:port)
- `udp_buffer_size`: UDP buffer size
- `uptime`: "initialized at 0"

### Performance Optimizations

The metadata export system is designed for high performance and low resource usage:

#### Metadata Caching
- **One-time JSON marshaling**: Metadata is converted to JSON once and cached
- **Eliminates repeated processing**: Subsequent metadata exports use cached data
- **Reduces CPU usage**: No repeated JSON encoding operations
- **Memory efficient**: Single cached copy shared across all exports

#### String Constants
- **Static string optimization**: Common strings like "agent_init", "linux", "initialized at 0" are constants
- **Reduces allocations**: Prevents repeated string creation
- **Improves memory locality**: Constants are stored in read-only memory
- **Faster comparisons**: Direct constant comparisons vs. string allocations

#### Lazy Hostname Resolution
- **Single system call**: `os.Hostname()` is called only once during initialization
- **Cached resolution**: Hostname is stored and reused for all metadata events
- **Reduces kernel context switches**: Avoids repeated system call overhead
- **Consistent hostname**: Ensures all metadata events use the same hostname value

#### UDP Transmission Efficiency
- **Direct cached data transmission**: Cached JSON is sent without re-processing
- **Minimal memory copying**: Data flows directly from cache to UDP socket
- **Reduced latency**: Faster metadata export with cached data
- **Better throughput**: Higher UDP packet transmission rates

### Error Handling

- Metadata export failures don't prevent Tetragon startup
- Warnings are logged if metadata export fails
- UDP exporter continues to function normally

## Usage Examples

### Basic UDP Export with Metadata

```bash
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514
```

This will:
1. Start Tetragon with UDP output enabled
2. Automatically send metadata event to `127.0.0.1:514`
3. Export all runtime security events to the same destination

### Custom UDP Configuration

```bash
tetragon \
  --udp-output-enabled \
  --udp-output-address=192.168.1.100 \
  --udp-output-port=9000 \
  --udp-buffer-size=131072
```

### Integration with Log Aggregation

```bash
# Send to Fluentd
tetragon --udp-output-enabled --udp-output-address=fluentd.example.com --udp-output-port=12201

# Send to syslog
tetragon --udp-output-enabled --udp-output-address=syslog.example.com --udp-output-port=514

# Send to custom monitoring system
tetragon --udp-output-enabled --udp-output-address=monitor.example.com --udp-output-port=9000
```

## Benefits

### Operational Visibility

- **Startup Context**: Know exactly when and where Tetragon started
- **Version Tracking**: Track Tetragon versions across deployments
- **System Information**: Understand the environment where Tetragon is running
- **Configuration Validation**: Verify UDP settings are applied correctly

### Integration Benefits

- **Log Correlation**: Correlate Tetragon events with system logs
- **Monitoring Dashboards**: Create dashboards showing agent status
- **Alerting**: Set up alerts for agent restarts or version changes
- **Audit Trail**: Maintain complete audit trail of agent lifecycle

### Troubleshooting

- **Startup Issues**: Identify problems during agent initialization
- **Configuration Problems**: Verify UDP settings are correct
- **Version Mismatches**: Detect version inconsistencies across nodes
- **System Changes**: Track when agents restart or reconfigure

## Monitoring and Alerting

### Metadata Event Detection

Monitor for metadata events to track agent lifecycle:

```bash
# Monitor for agent_init events
grep "event.*agent_init" /var/log/tetragon.log

# Check UDP export status
grep "Metadata event sent over UDP" /var/log/tetragon.log
```

### Health Checks

Use metadata events for health monitoring:

- **Agent Restarts**: Detect when agents restart
- **Version Changes**: Alert on version updates
- **Configuration Changes**: Monitor UDP setting modifications
- **System Changes**: Track hostname or kernel changes

## Future Enhancements

### Planned Features

- **Runtime Metadata Updates**: Periodic metadata refresh during operation
- **Custom Metadata Fields**: User-configurable metadata fields
- **Metadata Compression**: Efficient encoding for high-volume deployments
- **Metadata Validation**: Schema validation for metadata events

### Integration Opportunities

- **Prometheus Metrics**: Export metadata as Prometheus labels
- **OpenTelemetry**: Integrate with OpenTelemetry tracing
- **Kubernetes Events**: Send metadata as Kubernetes events
- **Custom Exporters**: Extend metadata export to other output formats

## Summary

Agent metadata export provides essential operational visibility and integration capabilities for Tetragon deployments. By automatically exporting initialization metadata on startup, Tetragon enables better monitoring, troubleshooting, and integration with existing infrastructure.

The metadata event serves as a foundation for understanding the agent's context and enables more sophisticated operational workflows in production environments. 