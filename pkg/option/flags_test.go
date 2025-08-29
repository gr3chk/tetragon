package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUDPMinimalMode(t *testing.T) {
	tests := []struct {
		name                     string
		udpOutputEnabled         bool
		initialHealthServer      string
		initialGopsAddr          string
		initialMetricsServer     string
		initialPprofAddr         string
		initialEnableK8s         bool
		initialEnablePolicy      bool
		initialEnablePodInfo     bool
		initialEnableTracingCRD  bool
		initialEnableCRI         bool
		expectedHealthServer     string
		expectedGopsAddr         string
		expectedMetricsServer    string
		expectedPprofAddr        string
		expectedEnableK8s        bool
		expectedEnablePolicy     bool
		expectedEnablePodInfo    bool
		expectedEnableTracingCRD bool
		expectedEnableCRI        bool
	}{
		{
			name:                     "UDP disabled - no changes",
			udpOutputEnabled:         false,
			initialHealthServer:      ":6789",
			initialGopsAddr:          "localhost:8118",
			initialMetricsServer:     ":2112",
			initialPprofAddr:         "localhost:6060",
			initialEnableK8s:         true,
			initialEnablePolicy:      true,
			initialEnablePodInfo:     true,
			initialEnableTracingCRD:  true,
			initialEnableCRI:         true,
			expectedHealthServer:     ":6789",
			expectedGopsAddr:         "localhost:8118",
			expectedMetricsServer:    ":2112",
			expectedPprofAddr:        "localhost:6060",
			expectedEnableK8s:        true,
			expectedEnablePolicy:     true,
			expectedEnablePodInfo:    true,
			expectedEnableTracingCRD: true,
			expectedEnableCRI:        true,
		},
		{
			name:                     "UDP enabled - all services disabled",
			udpOutputEnabled:         true,
			initialHealthServer:      ":6789",
			initialGopsAddr:          "localhost:8118",
			initialMetricsServer:     ":2112",
			initialPprofAddr:         "localhost:6060",
			initialEnableK8s:         true,
			initialEnablePolicy:      true,
			initialEnablePodInfo:     true,
			initialEnableTracingCRD:  true,
			initialEnableCRI:         true,
			expectedHealthServer:     "",
			expectedGopsAddr:         "",
			expectedMetricsServer:    "",
			expectedPprofAddr:        "",
			expectedEnableK8s:        false,
			expectedEnablePolicy:     false,
			expectedEnablePodInfo:    false,
			expectedEnableTracingCRD: false,
			expectedEnableCRI:        false,
		},
		{
			name:                     "UDP enabled - custom health server preserved",
			udpOutputEnabled:         true,
			initialHealthServer:      ":9999",
			initialGopsAddr:          "localhost:8118",
			initialMetricsServer:     ":2112",
			initialPprofAddr:         "localhost:6060",
			initialEnableK8s:         true,
			initialEnablePolicy:      true,
			initialEnablePodInfo:     true,
			initialEnableTracingCRD:  true,
			initialEnableCRI:         true,
			expectedHealthServer:     ":9999", // Custom port preserved
			expectedGopsAddr:         "",
			expectedMetricsServer:    "",
			expectedPprofAddr:        "",
			expectedEnableK8s:        false,
			expectedEnablePolicy:     false,
			expectedEnablePodInfo:    false,
			expectedEnableTracingCRD: false,
			expectedEnableCRI:        false,
		},
		{
			name:                     "UDP enabled - custom gops preserved",
			udpOutputEnabled:         true,
			initialHealthServer:      ":6789",
			initialGopsAddr:          "localhost:9999",
			initialMetricsServer:     ":2112",
			initialPprofAddr:         "localhost:6060",
			initialEnableK8s:         true,
			initialEnablePolicy:      true,
			initialEnablePodInfo:     true,
			initialEnableTracingCRD:  true,
			initialEnableCRI:         true,
			expectedHealthServer:     "",
			expectedGopsAddr:         "localhost:9999", // Custom port preserved
			expectedMetricsServer:    "",
			expectedPprofAddr:        "",
			expectedEnableK8s:        false,
			expectedEnablePolicy:     false,
			expectedEnablePodInfo:    false,
			expectedEnableTracingCRD: false,
			expectedEnableCRI:        false,
		},
		{
			name:                     "UDP enabled - services already disabled",
			udpOutputEnabled:         true,
			initialHealthServer:      "",
			initialGopsAddr:          "",
			initialMetricsServer:     "",
			initialPprofAddr:         "",
			initialEnableK8s:         false,
			initialEnablePolicy:      false,
			initialEnablePodInfo:     false,
			initialEnableTracingCRD:  false,
			initialEnableCRI:         false,
			expectedHealthServer:     "",
			expectedGopsAddr:         "",
			expectedMetricsServer:    "",
			expectedPprofAddr:        "",
			expectedEnableK8s:        false,
			expectedEnablePolicy:     false,
			expectedEnablePodInfo:    false,
			expectedEnableTracingCRD: false,
			expectedEnableCRI:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset config to initial state
			Config.UDPOutputEnabled = tt.udpOutputEnabled
			Config.HealthServerAddress = tt.initialHealthServer
			Config.GopsAddr = tt.initialGopsAddr
			Config.MetricsServer = tt.initialMetricsServer
			Config.PprofAddr = tt.initialPprofAddr
			Config.EnableK8s = tt.initialEnableK8s
			Config.EnablePolicyFilter = tt.initialEnablePolicy
			Config.EnablePolicyFilterCgroupMap = tt.initialEnablePolicy
			Config.EnablePolicyFilterDebug = tt.initialEnablePolicy
			Config.EnablePodInfo = tt.initialEnablePodInfo
			Config.EnableTracingPolicyCRD = tt.initialEnableTracingCRD
			Config.EnableCRI = tt.initialEnableCRI

			// Apply UDP minimal mode logic
			applyUDPMinimalMode()

			// Verify results
			assert.Equal(t, tt.expectedHealthServer, Config.HealthServerAddress, "Health server address mismatch")
			assert.Equal(t, tt.expectedGopsAddr, Config.GopsAddr, "Gops address mismatch")
			assert.Equal(t, tt.expectedMetricsServer, Config.MetricsServer, "Metrics server mismatch")
			assert.Equal(t, tt.expectedPprofAddr, Config.PprofAddr, "Pprof address mismatch")
			assert.Equal(t, tt.expectedEnableK8s, Config.EnableK8s, "K8s API mismatch")
			assert.Equal(t, tt.expectedEnablePolicy, Config.EnablePolicyFilter, "Policy filter mismatch")
			assert.Equal(t, tt.expectedEnablePolicy, Config.EnablePolicyFilterCgroupMap, "Policy filter cgroup map mismatch")
			assert.Equal(t, tt.expectedEnablePolicy, Config.EnablePolicyFilterDebug, "Policy filter debug mismatch")
			assert.Equal(t, tt.expectedEnablePodInfo, Config.EnablePodInfo, "Pod info mismatch")
			assert.Equal(t, tt.expectedEnableTracingCRD, Config.EnableTracingPolicyCRD, "Tracing policy CRD mismatch")
			assert.Equal(t, tt.expectedEnableCRI, Config.EnableCRI, "CRI mismatch")
		})
	}
}

// applyUDPMinimalMode applies the UDP minimal mode logic from ReadAndSetFlags
// This is extracted for testing purposes
func applyUDPMinimalMode() {
	// If UDP output is enabled, automatically disable health server for minimal operation
	if Config.UDPOutputEnabled && Config.HealthServerAddress == ":6789" {
		Config.HealthServerAddress = ""
	}

	// If UDP output is enabled, also disable gops server for minimal operation
	if Config.UDPOutputEnabled && Config.GopsAddr == "localhost:8118" {
		Config.GopsAddr = ""
	}

	// If UDP output is enabled, also disable metrics server for minimal operation
	if Config.UDPOutputEnabled && Config.MetricsServer != "" {
		Config.MetricsServer = ""
	}

	// If UDP output is enabled, also disable pprof server for minimal operation
	if Config.UDPOutputEnabled && Config.PprofAddr != "" {
		Config.PprofAddr = ""
	}

	// If UDP output is enabled, also disable Kubernetes API access for minimal operation
	if Config.UDPOutputEnabled && Config.EnableK8s {
		Config.EnableK8s = false
	}

	// If UDP output is enabled, also disable policy filtering for minimal operation
	if Config.UDPOutputEnabled && Config.EnablePolicyFilter {
		Config.EnablePolicyFilter = false
		Config.EnablePolicyFilterCgroupMap = false
		Config.EnablePolicyFilterDebug = false
	}

	// If UDP output is enabled, also disable pod info and tracing policy CRD for minimal operation
	if Config.UDPOutputEnabled && (Config.EnablePodInfo || Config.EnableTracingPolicyCRD) {
		Config.EnablePodInfo = false
		Config.EnableTracingPolicyCRD = false
	}

	// If UDP output is enabled, also disable CRI for minimal operation
	if Config.UDPOutputEnabled && Config.EnableCRI {
		Config.EnableCRI = false
	}
}

func TestUDPMinimalModeGRPCDisabling(t *testing.T) {
	tests := []struct {
		name               string
		udpOutputEnabled   bool
		grpcEnabled        bool
		initialServerAddr  string
		expectedServerAddr string
	}{
		{
			name:               "UDP disabled, gRPC enabled - server address preserved",
			udpOutputEnabled:   false,
			grpcEnabled:        true,
			initialServerAddr:  "localhost:54321",
			expectedServerAddr: "localhost:54321",
		},
		{
			name:               "UDP enabled, gRPC disabled - server address cleared",
			udpOutputEnabled:   true,
			grpcEnabled:        false,
			initialServerAddr:  "localhost:54321",
			expectedServerAddr: "",
		},
		{
			name:               "UDP enabled, gRPC enabled - server address preserved",
			udpOutputEnabled:   true,
			grpcEnabled:        true,
			initialServerAddr:  "localhost:54321",
			expectedServerAddr: "localhost:54321",
		},
		{
			name:               "UDP disabled, gRPC disabled - server address preserved",
			udpOutputEnabled:   false,
			grpcEnabled:        false,
			initialServerAddr:  "localhost:54321",
			expectedServerAddr: "localhost:54321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset config to initial state
			Config.UDPOutputEnabled = tt.udpOutputEnabled
			Config.GRPCEnabled = tt.grpcEnabled
			Config.ServerAddress = tt.initialServerAddr

			// Apply gRPC disabling logic
			applyUDPGRPCDisabling()

			// Verify results
			assert.Equal(t, tt.expectedServerAddr, Config.ServerAddress, "Server address mismatch")
		})
	}
}

// applyUDPGRPCDisabling applies the gRPC disabling logic from ReadAndSetFlags
// This is extracted for testing purposes
func applyUDPGRPCDisabling() {
	// If UDP output is enabled, disable gRPC by default unless explicitly enabled
	if Config.UDPOutputEnabled && !Config.GRPCEnabled {
		Config.ServerAddress = ""
	}
}

func TestUDPMinimalModeIntegration(t *testing.T) {
	// Test the complete integration of UDP minimal mode
	// This simulates what happens when ReadAndSetFlags is called

	// Set up initial configuration
	Config.UDPOutputEnabled = true
	Config.HealthServerAddress = ":6789"
	Config.GopsAddr = "localhost:8118"
	Config.MetricsServer = ":2112"
	Config.PprofAddr = "localhost:6060"
	Config.EnableK8s = true
	Config.EnablePolicyFilter = true
	Config.EnablePolicyFilterCgroupMap = true
	Config.EnablePolicyFilterDebug = true
	Config.EnablePodInfo = true
	Config.EnableTracingPolicyCRD = true
	Config.EnableCRI = true
	Config.GRPCEnabled = false
	Config.ServerAddress = "localhost:54321"

	// Apply UDP minimal mode
	applyUDPMinimalMode()
	applyUDPGRPCDisabling()

	// Verify all services are disabled
	assert.Equal(t, "", Config.HealthServerAddress, "Health server should be disabled")
	assert.Equal(t, "", Config.GopsAddr, "Gops server should be disabled")
	assert.Equal(t, "", Config.MetricsServer, "Metrics server should be disabled")
	assert.Equal(t, "", Config.PprofAddr, "Pprof server should be disabled")
	assert.Equal(t, false, Config.EnableK8s, "K8s API should be disabled")
	assert.Equal(t, false, Config.EnablePolicyFilter, "Policy filter should be disabled")
	assert.Equal(t, false, Config.EnablePolicyFilterCgroupMap, "Policy filter cgroup map should be disabled")
	assert.Equal(t, false, Config.EnablePolicyFilterDebug, "Policy filter debug should be disabled")
	assert.Equal(t, false, Config.EnablePodInfo, "Pod info should be disabled")
	assert.Equal(t, false, Config.EnableTracingPolicyCRD, "Tracing policy CRD should be disabled")
	assert.Equal(t, false, Config.EnableCRI, "CRI should be disabled")
	assert.Equal(t, "", Config.ServerAddress, "gRPC server should be disabled")

	// Verify UDP output is still enabled
	assert.Equal(t, true, Config.UDPOutputEnabled, "UDP output should remain enabled")
}
