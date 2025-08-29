---
title: "UDP Output"
description: "Configure Tetragon to send events and logs over UDP using fire-and-forget packet transmission"
---

# UDP Output

Tetragon can send events and logs over UDP to a configurable destination using **fire-and-forget packet transmission**. This feature is useful for integrating with log aggregation systems, SIEM platforms, or custom monitoring solutions that accept UDP input.

## Key Features

- **No listener required** - UDP packets are sent even if the destination port is closed or unreachable
- **Fire-and-forget** - True connectionless UDP transmission using WriteToUDP
- **High performance** - Efficient packet transmission without connection overhead
- **Network resilient** - Survives network interruptions and destination unavailability
- **Metadata export** - Automatically exports agent initialization metadata on startup

## How It Works

Tetragon uses unbound UDP sockets with `WriteToUDP()` to send packets directly to the destination address and port. This approach:

- **Eliminates connection requirements** - No need for a listener on the destination
- **Provides true UDP behavior** - Packets are transmitted regardless of destination state
- **Improves reliability** - No connection establishment failures or "connection refused" errors
- **Enables packet dumping** - Perfect for scenarios where you want to send logs even without a service listening
- **Minimal operation mode** - Automatically disables unnecessary services (health server, gRPC, metrics, etc.) when UDP output is enabled 

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

Events are sent as JSON over UDP, with each event on a separate line. The format is identical to the JSON export format used by other exporters.

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

## Technical Implementation

### WriteToUDP Architecture

Tetragon's UDP output uses unbound UDP sockets with `WriteToUDP()` for optimal performance:

- **Unbound Sockets**: Created with `net.ListenUDP("udp", ":0")` for maximum flexibility
- **Direct Packet Transmission**: Uses `conn.WriteToUDP(data, destAddr)` for fire-and-forget sending
- **Connection Pooling**: Efficient reuse of UDP sockets for better performance
- **No Connection State**: Eliminates connection establishment failures

### Socket Management

- **Socket Pool**: Maintains a pool of reusable UDP sockets
- **Automatic Fallback**: Creates new sockets when pool is empty
- **Buffer Optimization**: Configurable socket buffer sizes for performance tuning
- **Resource Cleanup**: Automatic socket cleanup and memory management

## Performance Considerations

- **Fire-and-forget**: No acknowledgment expected, maximum throughput
- **Network Resilience**: Packets sent regardless of destination availability
- **Buffer Tuning**: UDP buffer size significantly impacts performance
- **Rate Limiting**: Use `--export-rate-limit` for high-volume deployments
- **Memory Usage**: Efficient connection pooling minimizes resource overhead

## Use Cases

### Packet Dumping (No Listener Required)

Perfect for scenarios where you want to send logs even without a service listening:

```bash
# Send logs to port 514 even if no syslog server is running
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514
```

### Network Monitoring

Send events to network monitoring tools or packet analyzers:

```bash
# Send to network monitoring system
tetragon --udp-output-enabled --udp-output-address=monitor.example.com --udp-output-port=9000
```

### Log Aggregation

Integrate with log aggregation systems:

```bash
# Send to Fluentd
tetragon --udp-output-enabled --udp-output-address=fluentd.example.com --udp-output-port=12201
```

### Agent Metadata Export

On startup, Tetragon automatically exports a metadata event containing:

- **@timestamp**: ISO 8601 UTC timestamp of agent startup
- **event**: "agent_init" identifier
- **tetragon_version**: Current Tetragon version
- **build_commit**: Git commit hash (if available)
- **build_date**: Build timestamp (if available)
- **hostname**: System hostname
- **os**: Operating system identifier
- **kernel_version**: Linux kernel version
- **pid**: Process ID of the running agent
- **udp_destination**: Configured UDP destination (host:port)
- **udp_buffer_size**: Configured UDP buffer size
- **uptime**: Initialized at 0

This metadata event is sent as the first UDP packet, providing essential context for all subsequent events.

### Minimal Operation Mode

When UDP output is enabled, Tetragon automatically enters minimal operation mode:

- **Health server disabled**: Port 6789 automatically closed
- **gRPC server disabled**: Port 54321 automatically closed (unless explicitly enabled)
- **Gops server disabled**: Port 8118 automatically closed
- **Metrics server disabled**: Port 2112 automatically closed
- **Other services disabled**: Kubernetes API, policy filtering, CRI, etc.

This creates a focused deployment with minimal attack surface and maximum performance for UDP export scenarios:

```bash
# Minimal mode - only UDP export active
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514

# Minimal mode with custom health server
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514 --health-server-address=:9999
```

## Troubleshooting

### No More "Connection Refused" Errors

With the new WriteToUDP approach, you should no longer see connection-related errors. UDP packets will be sent successfully even when:

- Destination port is closed
- No service is listening
- Network is unreachable
- Firewall blocks traffic

### Check UDP Output Status

Look for UDP exporter startup messages in the Tetragon logs:

```
Starting UDP exporter address=192.168.1.100:514
```

### Verify Packet Transmission

Use network monitoring tools to verify UDP packets are being sent:

```bash
# Monitor UDP traffic to destination
sudo tcpdump -i any udp and host 192.168.1.100 and port 514
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