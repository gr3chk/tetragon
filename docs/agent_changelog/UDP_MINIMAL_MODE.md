# UDP Minimal Mode Implementation

## Overview

This document describes the implementation of UDP Minimal Mode in Tetragon, which automatically disables unnecessary services and ports when UDP output is enabled. This provides a truly minimal operation mode focused solely on event generation and UDP export.

## What Gets Disabled

When `--udp-output-enabled` is set to `true`, the following services are automatically disabled:

### 1. Health Server (Port 6789)
- **Default**: `:6789` (enabled)
- **UDP Mode**: Disabled automatically
- **Reason**: Not needed for UDP-only operation
- **Log Message**: `"Health server disabled for UDP-only operation"`

### 2. gRPC Server (Port 54321)
- **Default**: `localhost:54321` (enabled)
- **UDP Mode**: Disabled automatically unless explicitly enabled with `--grpc-enabled`
- **Reason**: Core gRPC server not needed for UDP export
- **Log Message**: `"gRPC server disabled for UDP-only operation"`

### 3. Gops Server (Port 8118)
- **Default**: `localhost:8118` (enabled)
- **UDP Mode**: Disabled automatically
- **Reason**: Debugging server not needed for production UDP export
- **Log Message**: `"Gops server disabled for UDP-only operation"`

### 4. Metrics Server (Port 2112)
- **Default**: `:2112` (disabled by default)
- **UDP Mode**: Disabled if previously enabled
- **Reason**: Metrics collection not needed for UDP-only operation
- **Log Message**: `"Metrics server disabled for UDP-only operation"`

### 5. Pprof Server (Port 6060)
- **Default**: Disabled by default
- **UDP Mode**: Disabled if previously enabled
- **Reason**: Profiling server not needed for UDP-only operation
- **Log Message**: `"Pprof server disabled for UDP-only operation"`

### 6. Kubernetes API Access
- **Default**: `false` (disabled by default)
- **UDP Mode**: Disabled if previously enabled
- **Reason**: Pod association not needed for UDP-only operation
- **Log Message**: `"Kubernetes API access disabled for UDP-only operation"`

### 7. Policy Filtering
- **Default**: `false` (disabled by default)
- **UDP Mode**: Disabled if previously enabled
- **Reason**: Policy filtering not needed for UDP-only operation
- **Log Message**: `"Policy filtering disabled for UDP-only operation"`

### 8. Pod Info & Tracing Policy CRD
- **Default**: `false` (disabled by default)
- **UDP Mode**: Disabled if previously enabled
- **Reason**: Kubernetes CRD operations not needed for UDP-only operation
- **Log Message**: `"Pod info and tracing policy CRD disabled for UDP-only operation"`

### 9. CRI (Container Runtime Interface)
- **Default**: `false` (disabled by default)
- **UDP Mode**: Disabled if previously enabled
- **Reason**: Container runtime integration not needed for UDP-only operation
- **Log Message**: `"CRI disabled for UDP-only operation"`

## Implementation Details

### Configuration Logic

The automatic service disabling is implemented in `pkg/option/flags.go` in the `ReadAndSetFlags()` function. The logic runs after all configuration values are read from viper, allowing it to override any previously set values.

### Conditional Disabling

Services are only disabled if:
1. UDP output is enabled (`Config.UDPOutputEnabled == true`)
2. The service was previously enabled (either by default or user configuration)

### Logging

Each disabled service logs an informational message explaining why it was disabled, making it clear to users what's happening.

## Usage Examples

### Basic UDP-Only Mode
```bash
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514
```

**Result**: All unnecessary services automatically disabled, only UDP export active.

### UDP + Explicit gRPC
```bash
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514 --grpc-enabled
```

**Result**: UDP export active, gRPC server active, other services disabled.

### UDP + Custom Health Server
```bash
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514 --health-server-address=:9999
```

**Result**: UDP export active, custom health server on port 9999, other services disabled.

## Benefits

### 1. Security
- **Reduced Attack Surface**: Fewer open ports and services
- **Minimal Permissions**: No unnecessary API access
- **Isolated Operation**: Focused solely on event generation and export

### 2. Performance
- **Lower Resource Usage**: No unnecessary background services
- **Reduced Memory Footprint**: Fewer active components
- **Focused Processing**: All resources dedicated to event handling

### 3. Simplicity
- **Single Purpose**: Clear, focused operation
- **Easy Deployment**: Minimal configuration required
- **Predictable Behavior**: Consistent service state

### 4. Production Ready
- **Firewall Friendly**: Only necessary UDP port open
- **Load Balancer Compatible**: No health check endpoints to manage
- **Monitoring Simplified**: Single export stream to monitor

## Configuration Override

If you need to re-enable any specific service while keeping UDP output enabled, you can explicitly set the configuration:

```bash
# Re-enable health server on custom port
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514 --health-server-address=:9999

# Re-enable gRPC server
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514 --grpc-enabled

# Re-enable metrics server
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514 --metrics-server=:2112
```

## Troubleshooting

### Service Still Running
If a service is still running when you expect it to be disabled:

1. **Check Configuration**: Verify `--udp-output-enabled` is set to `true`
2. **Check Logs**: Look for the disable messages in startup logs
3. **Explicit Override**: Check if the service was explicitly enabled with a custom flag
4. **Configuration Order**: Ensure UDP output flag comes before other service flags

### Missing Log Messages
If you don't see the disable messages:

1. **Log Level**: Ensure log level is set to `info` or lower
2. **Configuration**: Verify UDP output is actually enabled
3. **Service State**: Check if the service was already disabled

## Future Enhancements

### 1. Configurable Service Disabling
- Allow users to specify which services to keep enabled
- Provide a whitelist/blacklist approach for service management

### 2. Service Dependencies
- Automatically disable dependent services
- Provide warnings for incompatible service combinations

### 3. Health Check Alternatives
- Provide UDP-based health checking
- Implement lightweight status reporting via UDP

### 4. Metrics Export
- Allow metrics to be exported via UDP instead of HTTP
- Provide structured metrics in UDP packets

## Summary

UDP Minimal Mode provides a truly minimal operation mode for Tetragon, automatically disabling all unnecessary services when UDP output is enabled. This creates a focused, secure, and performant deployment that's ideal for production environments where only event generation and UDP export are required.

The implementation is transparent, logging all changes for user awareness, and allows explicit overrides when specific services are needed. This makes it easy to deploy Tetragon in minimal configurations while maintaining flexibility for custom requirements. 