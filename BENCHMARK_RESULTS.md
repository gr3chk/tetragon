# UDP Output Benchmark Results

## Quick Reference

This document provides the benchmark results for the UDP output feature. Use this for quick performance reference and capacity planning.

## Performance Summary Table

| Event Size | Throughput | Latency | Memory | Bandwidth | Use Case |
|------------|------------|---------|---------|-----------|----------|
| Small (200B) | 150K ops/sec | 6.6μs | 904B | 30 MB/s | High-frequency monitoring |
| Large (1.5KB) | 64K ops/sec | 15.6μs | 5.7KB | 96 MB/s | Detailed logging |
| Very Large (9KB) | 18K ops/sec | 55.8μs | 29.6KB | 162 MB/s | Comprehensive capture |
| Raw UDP Write | 421K ops/sec | 2.4μs | 0B | N/A | Network baseline |

## Detailed Benchmark Results

### Benchmark Environment
- **OS**: Linux 6.14.0-28-generic
- **Architecture**: amd64
- **CPU**: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
- **Go Version**: 1.21+

### Small Events (~200 bytes)
```
BenchmarkUDPEncoder_Encode-8    193659    6635 ns/op    904 B/op    23 allocs/op
```
- **Throughput**: 150,000 events/second
- **Latency**: 6.6 microseconds per event
- **Memory**: 904 bytes per event
- **Allocations**: 23 allocations per event
- **Bandwidth**: ~30 MB/s at full throughput

### Large Events (~1,500 bytes)
```
BenchmarkUDPEncoder_EncodeLargeEvent-8    78613    15601 ns/op    5772 B/op    24 allocs/op
```
- **Throughput**: 64,000 events/second
- **Latency**: 15.6 microseconds per event
- **Memory**: 5.7KB per event
- **Allocations**: 24 allocations per event
- **Bandwidth**: ~96 MB/s at full throughput

### Very Large Events (~9,000 bytes)
```
BenchmarkUDPEncoder_EncodeVeryLargeEvent-8    21811    55823 ns/op    29597 B/op    24 allocs/op
```
- **Throughput**: 18,000 events/second
- **Latency**: 55.8 microseconds per event
- **Memory**: 29.6KB per event
- **Allocations**: 24 allocations per event
- **Bandwidth**: ~162 MB/s at full throughput

### Raw UDP Write Performance
```
BenchmarkUDPEncoder_Write-8    508477    2372 ns/op    0 B/op    0 allocs/op
```
- **Throughput**: 421,000 operations/second
- **Latency**: 2.4 microseconds per operation
- **Memory**: 0 bytes per operation
- **Allocations**: 0 allocations per operation

## Capacity Planning

### Network Bandwidth Requirements

| Scenario | Events/sec | Event Size | Bandwidth Required |
|----------|------------|------------|-------------------|
| High-frequency monitoring | 100,000 | 200B | 20 MB/s |
| Detailed process tracking | 50,000 | 1.5KB | 75 MB/s |
| Comprehensive logging | 10,000 | 9KB | 90 MB/s |
| Maximum throughput (small) | 150,000 | 200B | 30 MB/s |
| Maximum throughput (large) | 64,000 | 1.5KB | 96 MB/s |
| Maximum throughput (very large) | 18,000 | 9KB | 162 MB/s |

### Memory Usage Estimation

| Event Size | Events/sec | Memory Usage |
|------------|------------|--------------|
| Small (200B) | 100,000 | ~90 MB |
| Large (1.5KB) | 50,000 | ~285 MB |
| Very Large (9KB) | 10,000 | ~296 MB |

### CPU Usage Estimation

| Event Size | Events/sec | CPU Time (ms/sec) |
|------------|------------|-------------------|
| Small (200B) | 100,000 | 663.5 |
| Large (1.5KB) | 50,000 | 780.1 |
| Very Large (9KB) | 10,000 | 558.2 |

## Performance Recommendations

### High-Throughput Scenarios (100K+ events/sec)
- Use small events (~200 bytes)
- Configure rate limiting: `export-rate-limit: 100000`
- Monitor network bandwidth usage
- Consider UDP packet fragmentation for large events

### Medium-Throughput Scenarios (10K-100K events/sec)
- Large events acceptable (~1.5KB)
- Rate limiting: `export-rate-limit: 50000`
- Balance between event detail and performance

### Low-Throughput Scenarios (<10K events/sec)
- Very large events acceptable (~9KB)
- Rate limiting: `export-rate-limit: 10000`
- Focus on event detail over performance

## Configuration Examples

### High-Performance Configuration
```yaml
udp-output-enabled: true
udp-output-address: "192.168.1.100"
udp-output-port: 514
export-rate-limit: 100000
event-queue-size: 50000
process-cache-size: 131072
```

### High-Detail Configuration
```yaml
udp-output-enabled: true
udp-output-address: "192.168.1.100"
udp-output-port: 514
export-rate-limit: 10000
# Include all event fields and details
```

### Balanced Configuration
```yaml
udp-output-enabled: true
udp-output-address: "192.168.1.100"
udp-output-port: 514
export-rate-limit: 50000
# Moderate filtering and detail level
```

## Running Benchmarks

### Quick Benchmark
```bash
go test ./pkg/encoder -bench=BenchmarkUDP -benchmem
```

### Detailed Benchmark
```bash
go test ./pkg/encoder -bench=BenchmarkUDP -benchmem -v -count=5
```

### Specific Event Size
```bash
# Small events
go test ./pkg/encoder -bench=BenchmarkUDPEncoder_Encode -benchmem

# Large events
go test ./pkg/encoder -bench=BenchmarkUDPEncoder_EncodeLargeEvent -benchmem

# Very large events
go test ./pkg/encoder -bench=BenchmarkUDPEncoder_EncodeVeryLargeEvent -benchmem
```

## Network Considerations

### UDP Packet Sizes
- **Small events**: Fit in single UDP packets
- **Large events**: May require fragmentation depending on MTU
- **Very large events**: Likely to require fragmentation

### MTU Considerations
- Standard Ethernet MTU: 1500 bytes
- Jumbo frames: 9000 bytes
- Consider network path MTU for large events

### Packet Loss Impact
- UDP is fire-and-forget
- No retransmission mechanism
- Monitor packet loss in high-throughput scenarios
- Consider network quality and congestion

## Monitoring Metrics

### Key Performance Indicators
- `tetragon_events_exported_total`: Total events exported
- `tetragon_events_exported_bytes_total`: Total bytes exported
- `tetragon_export_ratelimit_events_dropped_total`: Dropped events

### Network Monitoring
- UDP packet loss rate
- Network bandwidth utilization
- UDP socket buffer usage
- Network interface statistics

## Troubleshooting Performance

### Common Performance Issues
1. **Network congestion**: Monitor bandwidth usage
2. **UDP buffer overflow**: Check socket buffer sizes
3. **CPU saturation**: Monitor CPU usage during high throughput
4. **Memory pressure**: Monitor memory allocation patterns

### Performance Tuning
1. **Increase event queue size**: `event-queue-size: 100000`
2. **Optimize process cache**: `process-cache-size: 262144`
3. **Adjust rate limiting**: Balance throughput vs. resource usage
4. **Network optimization**: Use dedicated network interfaces

## Version History

### Benchmark Results History
- **v1.0**: Initial benchmark results (current)
- Future versions will track performance improvements

### Performance Targets
- **Target throughput**: 200K events/sec for small events
- **Target latency**: <5μs for small events
- **Target memory**: <500B per small event 