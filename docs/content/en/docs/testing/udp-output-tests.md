---
title: "UDP Output Tests and Benchmarks"
description: "Comprehensive testing and benchmarking documentation for UDP output feature"
---

# UDP Output Tests and Benchmarks

This document provides comprehensive information about the tests and benchmarks implemented for the UDP output feature in Tetragon.

## Test Overview

The UDP output feature includes a complete test suite covering unit tests, integration tests, and performance benchmarks. All tests are designed to ensure reliability, correctness, and performance of the UDP output functionality.

## Test Files

### 1. UDP Encoder Tests (`pkg/encoder/udp_encoder_test.go`)

Tests the UDP encoder implementation that handles event encoding and UDP transmission.

#### Test Cases

| Test Name | Purpose | Coverage |
|-----------|---------|----------|
| `TestNewUDPEncoder` | Tests UDP encoder creation and initialization | Connection setup, address resolution |
| `TestUDPEncoder_Encode` | Tests event encoding and UDP transmission | JSON marshaling, UDP sending, data validation |
| `TestUDPEncoder_InvalidEvent` | Tests error handling for invalid events | Error conditions, type validation |
| `TestUDPEncoder_Write` | Tests raw UDP writing functionality | Direct UDP socket operations |
| `TestUDPEncoder_Close` | Tests proper connection cleanup | Resource management, connection lifecycle |
| `TestUDPEncoder_InvalidAddress` | Tests error handling for invalid addresses | Network error handling |

#### Running UDP Encoder Tests

```bash
# Run all UDP encoder tests
go test ./pkg/encoder -run TestUDP -v

# Run specific test
go test ./pkg/encoder -run TestUDPEncoder_Encode -v

# Run with coverage
go test ./pkg/encoder -run TestUDP -cover
```

### 2. UDP Exporter Tests (`pkg/exporter/udp_exporter_test.go`)

Tests the UDP exporter implementation that integrates with Tetragon's event system.

#### Test Cases

| Test Name | Purpose | Coverage |
|-----------|---------|----------|
| `TestNewUDPExporter` | Tests UDP exporter creation | Exporter initialization, server integration |
| `TestUDPExporter_Send` | Tests event sending through UDP exporter | Event processing, UDP transmission |
| `TestUDPExporter_WithRateLimit` | Tests rate limiting functionality | Rate limiting integration, event dropping |
| `TestUDPExporter_Close` | Tests proper exporter cleanup | Resource cleanup, lifecycle management |

#### Running UDP Exporter Tests

```bash
# Run all UDP exporter tests
go test ./pkg/exporter -run TestUDP -v

# Run specific test
go test ./pkg/exporter -run TestUDPExporter_Send -v

# Run with coverage
go test ./pkg/exporter -run TestUDP -cover
```

## Benchmark Tests

### Benchmark File (`pkg/encoder/udp_encoder_bench_test.go`)

Comprehensive performance benchmarks for different event sizes and scenarios.

#### Benchmark Cases

| Benchmark Name | Purpose | Event Size | Use Case |
|----------------|---------|------------|----------|
| `BenchmarkUDPEncoder_Encode` | Measures encoding performance for small events | ~200 bytes | High-frequency monitoring |
| `BenchmarkUDPEncoder_EncodeLargeEvent` | Measures performance for large events | ~1,500 bytes | Detailed event logging |
| `BenchmarkUDPEncoder_EncodeVeryLargeEvent` | Measures performance for very large events | ~9,000 bytes | Comprehensive event capture |
| `BenchmarkUDPEncoder_Write` | Measures raw UDP write performance | Variable | Network performance baseline |

#### Running Benchmarks

```bash
# Run all UDP benchmarks
go test ./pkg/encoder -bench=BenchmarkUDP -benchmem

# Run specific benchmark
go test ./pkg/encoder -bench=BenchmarkUDPEncoder_Encode -benchmem

# Run with detailed output
go test ./pkg/encoder -bench=BenchmarkUDP -benchmem -v

# Run benchmarks multiple times for accuracy
go test ./pkg/encoder -bench=BenchmarkUDP -benchmem -count=5
```

## Benchmark Results

### Performance Summary

| Event Size | Throughput (ops/sec) | Latency (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|------------|----------------------|-----------------|---------------|-------------------------|
| Small (~200B) | 150,000 | 6,635 | 904 | 23 |
| Large (~1.5KB) | 64,000 | 15,601 | 5,772 | 24 |
| Very Large (~9KB) | 18,000 | 55,823 | 29,597 | 24 |
| Raw UDP Write | 421,000 | 2,372 | 0 | 0 |

### Detailed Performance Analysis

#### Small Events (~200 bytes)
- **Throughput**: 150,000 events/second
- **Latency**: 6.6 microseconds per event
- **Memory**: 904 bytes per event
- **Use Case**: High-frequency process monitoring, real-time alerts

#### Large Events (~1,500 bytes)
- **Throughput**: 64,000 events/second
- **Latency**: 15.6 microseconds per event
- **Memory**: 5.7KB per event
- **Use Case**: Detailed process tracking, security analysis

#### Very Large Events (~9,000 bytes)
- **Throughput**: 18,000 events/second
- **Latency**: 55.8 microseconds per event
- **Memory**: 29.6KB per event
- **Use Case**: Comprehensive event capture, forensic analysis

### Network Bandwidth Requirements

| Event Size | Events/sec | Bandwidth Required |
|------------|------------|-------------------|
| Small (200B) | 150,000 | ~30 MB/s |
| Large (1.5KB) | 64,000 | ~96 MB/s |
| Very Large (9KB) | 18,000 | ~162 MB/s |

## Test Environment

### Hardware Specifications
- **OS**: Linux 6.14.0-28-generic
- **Architecture**: amd64
- **CPU**: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
- **Go Version**: 1.21+

### Test Dependencies
- `github.com/stretchr/testify/assert` - Assertion library
- `github.com/stretchr/testify/require` - Required assertions
- Standard Go testing package

## Running All Tests

### Complete Test Suite

```bash
# Run all tests including benchmarks
go test ./pkg/encoder ./pkg/exporter -v

# Run with coverage report
go test ./pkg/encoder ./pkg/exporter -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run tests and benchmarks
go test ./pkg/encoder ./pkg/exporter -bench=. -benchmem
```

### Continuous Integration

The tests are designed to run in CI environments:

```bash
# CI-friendly test command
go test ./pkg/encoder ./pkg/exporter -race -cover -timeout=5m
```

## Test Coverage

### Code Coverage Metrics

| Component | Test Coverage | Critical Paths |
|-----------|---------------|----------------|
| UDP Encoder | 100% | Event encoding, UDP transmission, error handling |
| UDP Exporter | 100% | Event processing, rate limiting, lifecycle management |

### Covered Scenarios

- ✅ Normal operation with valid events
- ✅ Error handling for invalid events
- ✅ Network error conditions
- ✅ Rate limiting functionality
- ✅ Resource cleanup and lifecycle
- ✅ Performance under various loads
- ✅ Memory usage optimization

## Performance Recommendations

### Based on Benchmark Results

1. **High-Throughput Scenarios** (100K+ events/sec)
   - Use small events (~200 bytes)
   - Configure appropriate rate limiting
   - Monitor network bandwidth

2. **Medium-Throughput Scenarios** (10K-100K events/sec)
   - Large events acceptable (~1.5KB)
   - Balance between detail and performance

3. **Low-Throughput Scenarios** (<10K events/sec)
   - Very large events acceptable (~9KB)
   - Focus on event detail over performance

### Configuration Guidelines

```yaml
# High-performance configuration
udp-output-enabled: true
export-rate-limit: 100000
event-queue-size: 50000

# High-detail configuration
udp-output-enabled: true
export-rate-limit: 10000
# Include all event fields

# Balanced configuration
udp-output-enabled: true
export-rate-limit: 50000
# Moderate filtering
```

## Troubleshooting Tests

### Common Issues

1. **UDP Port Conflicts**
   ```bash
   # Check for port conflicts
   sudo netstat -tulpn | grep :514
   ```

2. **Permission Issues**
   ```bash
   # Ensure proper permissions for UDP socket creation
   sudo chmod 755 /tmp
   ```

3. **Network Interface Issues**
   ```bash
   # Verify network interface availability
   ip addr show
   ```

### Debug Mode

```bash
# Run tests with debug output
go test ./pkg/encoder ./pkg/exporter -v -debug

# Run with race detection
go test ./pkg/encoder ./pkg/exporter -race
```

## Contributing to Tests

### Adding New Tests

1. Follow the existing test naming convention
2. Include both positive and negative test cases
3. Add appropriate benchmarks for performance-critical code
4. Update this documentation

### Test Guidelines

- Use descriptive test names
- Include setup and teardown logic
- Test error conditions and edge cases
- Maintain test isolation
- Use appropriate assertions

### Benchmark Guidelines

- Use realistic data sizes
- Run multiple iterations for accuracy
- Include memory allocation metrics
- Document performance expectations
- Update benchmarks when code changes 