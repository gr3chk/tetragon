#!/bin/bash
# Script to install sample UDP output configuration files

set -e

echo "Installing Tetragon UDP output configuration files..."

# Create configuration directories
sudo mkdir -p /etc/tetragon
sudo mkdir -p /etc/tetragon/tetragon.conf.d

# Copy main configuration file
echo "Installing main UDP output configuration..."
sudo cp examples/configuration/udp-output.yaml /etc/tetragon/tetragon.yaml

# Copy drop-in configurations
echo "Installing drop-in configurations..."
sudo cp examples/configuration/udp-output-high-throughput.yaml /etc/tetragon/tetragon.conf.d/
sudo cp examples/configuration/udp-output-filtered.yaml /etc/tetragon/tetragon.conf.d/
sudo cp examples/configuration/udp-output-with-grpc.yaml /etc/tetragon/tetragon.conf.d/

# Set proper permissions
sudo chmod 644 /etc/tetragon/tetragon.yaml
sudo chmod 644 /etc/tetragon/tetragon.conf.d/*.yaml

echo "Configuration files installed successfully!"
echo ""
echo "Available configurations:"
echo "  - /etc/tetragon/tetragon.yaml (default UDP output)"
echo "  - /etc/tetragon/tetragon.conf.d/udp-output-high-throughput.yaml"
echo "  - /etc/tetragon/tetragon.conf.d/udp-output-filtered.yaml"
echo "  - /etc/tetragon/tetragon.conf.d/udp-output-with-grpc.yaml"
echo ""
echo "To use a specific configuration, copy it to /etc/tetragon/tetragon.yaml:"
echo "  sudo cp /etc/tetragon/tetragon.conf.d/udp-output-high-throughput.yaml /etc/tetragon/tetragon.yaml"
echo ""
echo "Then run Tetragon with:"
echo "  sudo tetragon" 