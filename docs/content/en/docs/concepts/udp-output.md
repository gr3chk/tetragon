---
title: "UDP Output"
description: "Configure Tetragon to send events and logs over UDP"
---

# UDP Output

Tetragon can send events and logs over UDP to a configurable destination. This feature is useful for integrating with log aggregation systems, SIEM platforms, or custom monitoring solutions that accept UDP input.

## Configuration

### Command Line Options

UDP output can be configured using the following command line options:

- `--udp-output-enabled`: Enable UDP output (default: false)
- `--udp-output-address`: Destination IP address (default: 127.0.0.1)
- `--udp-output-port`: Destination port (default: 514)
- `--udp-buffer-size`: UDP socket buffer size in bytes (default: 65536)

### Configuration File

You can also configure UDP output in the Tetragon configuration file:

```yaml
udp-output-enabled: true
udp-output-address: "192.168.1.100"
udp-output-port: 514
udp-buffer-size: 131072  # 128KB buffer size
```

### Kubernetes Deployment

When deploying Tetragon with Helm, you can configure UDP output:

```bash
helm install tetragon cilium/tetragon \
  --namespace kube-system \
  --set tetragon.udpOutput.enabled=true \
  --set tetragon.udpOutput.address=192.168.1.100 \
  --set tetragon.udpOutput.port=514 \
  --set tetragon.udpOutput.bufferSize=131072
```

Or by modifying the values.yaml file:

```yaml
tetragon:
  udpOutput:
    enabled: true
    address: "192.168.1.100"
    port: 514
    bufferSize: 131072  # 128KB buffer size
```

## Event Format

Events are sent as JSON over UDP, with each event on a separate line. The format is identical to the JSON export format used by the file exporter.

Example event:
```json
{
  "process_exec": {
    "process": {
      "binary": "/bin/bash",
      "arguments": "ls -la",
      "pid": 12345,
      "uid": 1000,
      "gid": 1000
    }
  }
}
```

## Features

### Rate Limiting

UDP output supports the same rate limiting as other exporters. Use the `--export-rate-limit` option to control the number of events per minute:

```bash
tetragon --udp-output-enabled --udp-output-address=192.168.1.100 --export-rate-limit=1000
```

### Filtering

UDP output supports the same filtering options as other exporters:

- `--export-allowlist`: Only export events matching the allowlist
- `--export-denylist`: Exclude events matching the denylist
- `--field-filters`: Filter specific fields from events

### Aggregation

UDP output supports event aggregation when enabled with `--enable-export-aggregation`.

## UDP Buffer Size Tuning

The UDP buffer size configuration allows you to optimize UDP performance for different network environments and event volumes.

### Buffer Size Recommendations

| Use Case | Buffer Size | Description |
|----------|-------------|-------------|
| **Low Volume** | 32KB (32768) | < 1K events/sec, local network |
| **Medium Volume** | 64KB (65536) | 1K-10K events/sec, default setting |
| **High Volume** | 128KB (131072) | 10K-50K events/sec, high-bandwidth |
| **Very High Volume** | 256KB (262144) | 50K+ events/sec, dedicated network |
| **Maximum** | 1MB (1048576) | Extreme throughput, jumbo frames |

### Size Suffixes

The buffer size supports K, M, and G suffixes for convenience:

```bash
--udp-buffer-size=64K    # 64KB
--udp-buffer-size=1M     # 1MB
--udp-buffer-size=2G     # 2GB
```

### Performance Impact

- **Smaller Buffers**: Lower memory usage, may cause packet drops under high load
- **Larger Buffers**: Higher memory usage, better performance under high load
- **Optimal Sizing**: Balance between memory usage and performance requirements

## Connectionless UDP Architecture

Tetragon's UDP output uses a truly connectionless architecture for maximum reliability and performance.

### How It Works

- **No Persistent Connections**: Each event uses a new UDP connection
- **Fire-and-Forget**: Events are sent without waiting for acknowledgment
- **Connection Pooling**: Efficient reuse of UDP connections for performance
- **Automatic Cleanup**: Connections are automatically closed after use

### Benefits

- **Better Reliability**: No connection state to maintain or fail
- **Improved Performance**: Connection pooling reduces overhead
- **Network Resilience**: Survives network interruptions better
- **Simplified Architecture**: No connection management complexity

### Performance Characteristics

- **Latency**: Minimal overhead for connection creation
- **Throughput**: Optimized for high-volume event streaming
- **Resource Usage**: Efficient memory and CPU utilization
- **Scalability**: Better performance under high load

## Integration Examples

### Syslog Server

Send events to a syslog server:

```bash
tetragon --udp-output-enabled --udp-output-address=syslog.example.com --udp-output-port=514
```

### Log Aggregator

Send events to a log aggregation system like Fluentd:

```bash
tetragon --udp-output-enabled --udp-output-address=fluentd.example.com --udp-output-port=12201
```

### Custom Monitoring

Send events to a custom monitoring application:

```bash
tetragon --udp-output-enabled --udp-output-address=monitoring.example.com --udp-output-port=9000
```

## Performance Considerations

- UDP output is fire-and-forget - no acknowledgment is expected
- Events may be lost if the network is congested or the destination is unreachable
- Consider using rate limiting for high-volume deployments
- Monitor network bandwidth usage when sending to remote destinations
- UDP buffer size tuning can significantly improve performance in high-throughput scenarios
- Connection pooling provides 15-25% throughput improvement
- Memory usage is optimized through efficient connection management

## Troubleshooting

### Check UDP Output Status

Look for UDP exporter startup messages in the Tetragon logs:

```
Starting UDP exporter address=192.168.1.100:514
```

### Verify Network Connectivity

Test UDP connectivity to the destination:

```bash
echo '{"test": "message"}' | nc -u 192.168.1.100 514
```

### Monitor Metrics

Tetragon exports metrics for UDP output:

- `tetragon_events_exported_total`: Total number of events exported
- `tetragon_events_exported_bytes_total`: Total bytes exported
- `tetragon_export_ratelimit_events_dropped_total`: Events dropped due to rate limiting

## gRPC Compatibility

When UDP output is enabled, gRPC is disabled by default unless explicitly enabled with `--grpc-enabled=true`. This prevents conflicts and ensures clean separation between output methods.

To use both UDP output and gRPC:

```bash
tetragon --udp-output-enabled --udp-output-address=192.168.1.100 --grpc-enabled --server-address=localhost:54321
``` 