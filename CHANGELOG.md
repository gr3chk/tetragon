# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **UDP Output Support**: Added new UDP output option that sends all logs and events to a configurable UDP address and port
  - New command line options: `--udp-output-enabled`, `--udp-output-address`, `--udp-output-port`
  - New configuration file options for UDP output
  - Kubernetes Helm chart support for UDP output configuration
  - UDP output supports rate limiting, filtering, and aggregation features
  - Events are sent as JSON over UDP with no acknowledgment expected
  - Comprehensive unit tests and benchmarks for UDP functionality

### Changed
- **gRPC Output**: Made gRPC optional and only active if explicitly configured
  - gRPC is disabled by default when UDP output is enabled
  - New `--grpc-enabled` flag to explicitly enable gRPC
  - Backward compatibility maintained for existing deployments

### Technical Details
- Implemented modular UDP encoder following SOLID principles
- Added UDP exporter that implements server.Listener interface
- Integrated UDP output with existing rate limiting and filtering mechanisms
- Added comprehensive test coverage including unit tests and benchmarks
- Updated Kubernetes deployment templates and Helm charts
- Added detailed documentation with configuration examples and troubleshooting guide

## [Previous Releases] 