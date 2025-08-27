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

### Configuration File

You can also configure UDP output in the Tetragon configuration file:

```yaml
udp-output-enabled: true
udp-output-address: "192.168.1.100"
udp-output-port: 514
```

### Kubernetes Deployment

When deploying Tetragon with Helm, you can configure UDP output:

```bash
helm install tetragon cilium/tetragon \
  --namespace kube-system \
  --set tetragon.udpOutput.enabled=true \
  --set tetragon.udpOutput.address=192.168.1.100 \
  --set tetragon.udpOutput.port=514
```

Or by modifying the values.yaml file:

```yaml
tetragon:
  udpOutput:
    enabled: true
    address: "192.168.1.100"
    port: 514
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