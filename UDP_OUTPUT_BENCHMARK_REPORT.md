# UDP Output Benchmark Report

## Overview

This report presents the performance benchmarks for the new UDP output feature in Tetragon. The benchmarks measure CPU and memory usage for different log sizes and throughput scenarios.

## Test Environment

- **OS**: Linux 6.14.0-28-generic
- **Architecture**: amd64
- **CPU**: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
- **Go Version**: 1.21+

## Benchmark Results

### Event Encoding Performance

| Event Size | Operations/sec | Latency (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|------------|----------------|-----------------|---------------|-------------------------|
| Small (~200 bytes) | 150,000 | 6,635 | 904 | 23 |
| Large (~1,500 bytes) | 64,000 | 15,601 | 5,772 | 24 |
| Very Large (~9,000 bytes) | 18,000 | 55,823 | 29,597 | 24 |

### Raw UDP Write Performance

| Metric | Value |
|--------|-------|
| Operations/sec | 421,000 |
| Latency (ns/op) | 2,372 |
| Memory (B/op) | 0 |
| Allocations (allocs/op) | 0 |

## Performance Analysis

### CPU Usage

1. **Small Events (200 bytes)**: 
   - High throughput at 150K ops/sec
   - Low latency of ~6.6μs per event
   - Suitable for high-frequency event streams

2. **Large Events (1,500 bytes)**:
   - Moderate throughput at 64K ops/sec
   - Latency increases to ~15.6μs per event
   - Memory usage increases significantly (5.7KB per event)

3. **Very Large Events (9,000 bytes)**:
   - Lower throughput at 18K ops/sec
   - Higher latency of ~55.8μs per event
   - Significant memory usage (29.6KB per event)

### Memory Usage

- **Small events**: Minimal memory overhead (~904 bytes per event)
- **Large events**: Moderate memory usage (~5.7KB per event)
- **Very large events**: High memory usage (~29.6KB per event)

### Throughput Recommendations

Based on the benchmark results:

1. **High-throughput scenarios** (100K+ events/sec): Use small events
2. **Medium-throughput scenarios** (10K-100K events/sec): Large events acceptable
3. **Low-throughput scenarios** (<10K events/sec): Very large events acceptable

## Network Considerations

### UDP Packet Sizes

- **Small events**: Fit comfortably in single UDP packets
- **Large events**: May require fragmentation depending on network MTU
- **Very large events**: Likely to require fragmentation

### Network Bandwidth

Estimated bandwidth requirements:
- Small events: ~30 MB/s at 150K events/sec
- Large events: ~96 MB/s at 64K events/sec  
- Very large events: ~162 MB/s at 18K events/sec

## Recommendations

### Production Deployment

1. **Monitor network bandwidth**: Ensure sufficient bandwidth for expected event volumes
2. **Use rate limiting**: Implement rate limiting for high-volume deployments
3. **Consider event size**: Balance between event detail and performance
4. **Network monitoring**: Monitor UDP packet loss and fragmentation

### Configuration Guidelines

1. **High-performance scenarios**:
   ```bash
   tetragon --udp-output-enabled --udp-output-address=destination:514 --export-rate-limit=100000
   ```

2. **High-detail scenarios**:
   ```bash
   tetragon --udp-output-enabled --udp-output-address=destination:514 --export-rate-limit=10000
   ```

3. **Balanced scenarios**:
   ```bash
   tetragon --udp-output-enabled --udp-output-address=destination:514 --export-rate-limit=50000
   ```

## Conclusion

The UDP output feature provides excellent performance for event streaming with:

- **High throughput**: Up to 150K events/sec for small events
- **Low latency**: Sub-10μs latency for small events
- **Efficient memory usage**: Minimal overhead for small events
- **Scalable**: Performance degrades gracefully with event size

The feature is suitable for production deployments across various use cases, from high-frequency monitoring to detailed event logging. 