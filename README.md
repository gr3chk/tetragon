<a href="https://tetragon.io">
  <picture>
    <source media="(prefers-color-scheme: light)" srcset="docs/assets/icons/logo.svg" width="400">
    <img src="docs/assets/icons/logo-dark.svg" width="400">
  </picture>
</a>

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://opensource.org/license/bsd-2-clause/)
[![License](https://img.shields.io/badge/license-GPL-blue.svg)](https://opensource.org/license/gpl-2-0/)

---

Cilium’s new [Tetragon](https://tetragon.io) component enables powerful
real-time, eBPF-based Security Observability and Runtime Enforcement.

Tetragon detects and is able to react to security-significant events, such as

- Process execution events
- System call activity
- I/O activity including network & file access

When used in a Kubernetes environment, Tetragon is Kubernetes-aware - that is,
it understands Kubernetes identities such as namespaces, pods and so on - so
that security event detection can be configured in relation to individual
workloads.

[![Tetragon Overview Diagram](https://github.com/cilium/tetragon/blob/main/docs/static/images/smart_observability.png)](https://tetragon.io/docs/overview/)

See more about [how Tetragon is using eBPF](https://tetragon.io/docs/overview#functionality-overview).

## UDP Output Feature

Tetragon now supports UDP output for sending events and logs to configurable destinations. This feature is ideal for integration with log aggregation systems, SIEM platforms, and custom monitoring solutions.

### Quick Start with UDP Output

```bash
# Basic UDP output
tetragon --udp-output-enabled --udp-output-address=192.168.1.100 --udp-output-port=514

# With rate limiting
tetragon --udp-output-enabled --udp-output-address=192.168.1.100 --export-rate-limit=1000
```

### Documentation

- **[UDP Output Guide](docs/content/en/docs/concepts/udp-output.md)** - Complete configuration and usage guide
- **[Test Documentation](docs/content/en/docs/testing/udp-output-tests.md)** - Comprehensive test and benchmark documentation
- **[Benchmark Results](BENCHMARK_RESULTS.md)** - Performance benchmarks and capacity planning
- **[Test Summary](TEST_SUMMARY.md)** - Complete test coverage and execution guide

### Performance Highlights

| Event Size | Throughput | Latency | Use Case |
|------------|------------|---------|----------|
| Small (200B) | 150K ops/sec | 6.6μs | High-frequency monitoring |
| Large (1.5KB) | 64K ops/sec | 15.6μs | Detailed logging |
| Very Large (9KB) | 18K ops/sec | 55.8μs | Comprehensive capture |

## UDP Minimal Mode

Tetragon now supports a **UDP Minimal Mode** that automatically disables all unnecessary services when UDP output is enabled. This creates a truly minimal deployment focused solely on event generation and UDP export.

### What Gets Automatically Disabled

- **Health Server** (Port 6789) - No health check endpoints
- **gRPC Server** (Port 54321) - No API server (unless explicitly enabled)
- **Gops Server** (Port 8118) - No debugging server
- **Metrics Server** (Port 2112) - No metrics collection
- **Pprof Server** (Port 6060) - No profiling server
- **Kubernetes API Access** - No pod association
- **Policy Filtering** - No policy management
- **Other Services** - CRI, pod info, tracing policy CRD, etc.

### Usage

```bash
# Basic minimal mode - only UDP export active
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514

# Minimal mode with custom health server
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514 --health-server-address=:9999

# Minimal mode with gRPC enabled
tetragon --udp-output-enabled --udp-output-address=127.0.0.1 --udp-output-port=514 --grpc-enabled
```

### Benefits

- **🔒 Security**: Minimal attack surface, only necessary UDP port open
- **⚡ Performance**: All resources focused on event generation and export
- **🚀 Production Ready**: Perfect for firewall-friendly deployments
- **📊 Simple Monitoring**: Single UDP stream to monitor

**Agent Metadata Export:**
On startup, Tetragon automatically exports initialization metadata over UDP, including version, hostname, kernel version, and configuration details. This provides essential context for all subsequent events.

For more details, see [UDP Output Documentation](docs/content/en/docs/concepts/udp-output.md), [UDP Minimal Mode Guide](docs/agent_changelog/UDP_MINIMAL_MODE.md), and [Agent Metadata Export Guide](docs/agent_changelog/AGENT_METADATA_EXPORT.md).

## Getting Started

Refer to the [official documentation of Tetragon](https://tetragon.io/docs/).

To get started with Tetragon, take a look at the [getting started
guides](https://tetragon.io/docs/getting-started/) to:
- [Try Tetragon on Kubernetes](https://tetragon.io/docs/getting-started/install-k8s/)
- [Try Tetragon on Linux](https://tetragon.io/docs/getting-started/install-docker/)
- [Deploy Tetragon](https://tetragon.io/docs/installation/)
- [Install the Tetra CLI](https://tetragon.io/docs/installation/tetra-cli/)

Tetragon is able to observe critical hooks in the kernel through its sensors
and generates events enriched with Linux and Kubernetes metadata:
1. **Process lifecycle**: generating `process_exec` and `process_exit` events
   by default, enabling full process lifecycle observability. Learn more about
   these events on the [process lifecycle use case page](https://tetragon.io/docs/use-cases/process-lifecycle/).
1. **Generic tracing**: generating `process_kprobe`, `process_tracepoint` and
   `process_uprobe` events for more advanced and custom use cases. Learn more
   about these events on the [TracingPolicy concept page](https://tetragon.io/docs/concepts/tracing-policy/)
   and discover [multiple use cases](https://tetragon.io/docs/use-cases/) like:
   - [🌏 network observability](https://tetragon.io/docs/use-cases/network-observability/)
   - [📂 filename access](https://tetragon.io/docs/use-cases/filename-access/)
   - [🔑 credentials monitoring](https://tetragon.io/docs/use-cases/linux-process-credentials/)
   - [🔓 privileged execution](https://tetragon.io/docs/use-cases/process-lifecycle/privileged-execution/)

See further resources:
- [Conference Talks, Books, Blog Posts, and Labs](https://tetragon.io/docs/resources/)
- [Frequently Asked Question](https://tetragon.io/docs/installation/faq/)
- [References](https://tetragon.io/docs/reference/)

## Join the community

Join the Tetragon [💬 Slack channel](https://slack.cilium.io) and the
[📅 Community Call](https://isogo.to/tetragon-meeting-notes) to chat with
developers, maintainers, and other users. This is a good first stop to ask
questions and share your experiences.

## How to Contribute

For getting started with local development, you can refer to the
[Contribution Guide](https://tetragon.io/docs/contribution-guide/). If
you plan to submit a PR, please ["sign-off"](https://tetragon.io/docs/contribution-guide/developer-certificate-of-origin/)
your commits.
