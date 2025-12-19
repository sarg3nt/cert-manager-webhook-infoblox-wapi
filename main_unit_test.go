package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
)

// TestName verifies the solver name is correct
func TestName(t *testing.T) {
	solver := &customDNSProviderSolver{}
	assert.Equal(t, "infoblox-wapi", solver.Name())
}

// TestDeDot verifies trailing dot removal from FQDNs
func TestDeDot(t *testing.T) {
	solver := &customDNSProviderSolver{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "FQDN with trailing dot",
			input:    "example.com.",
			expected: "example.com",
		},
		{
			name:     "FQDN without trailing dot",
			input:    "example.com",
			expected: "example.com",
		},
		{
			name:     "subdomain with trailing dot",
			input:    "_acme-challenge.example.com.",
			expected: "_acme-challenge.example.com",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "just a dot",
			input:    ".",
			expected: "",
		},
		{
			name:     "multiple dots in middle",
			input:    "sub.domain.example.com.",
			expected: "sub.domain.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solver.DeDot(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestLoadConfig_Valid tests successful configuration parsing
func TestLoadConfig_Valid(t *testing.T) {
	configJSON := `{
		"host": "infoblox.example.com",
		"port": "8443",
		"version": "2.11",
		"view": "Internal",
		"sslVerify": true,
		"httpRequestTimeout": 90,
		"httpPoolConnections": 20,
		"ttl": 600,
		"useTtl": true
	}`

	raw := apiextensionsv1.JSON{Raw: []byte(configJSON)}
	cfg, err := loadConfig(&raw)

	require.NoError(t, err)
	assert.Equal(t, "infoblox.example.com", cfg.Host)
	assert.Equal(t, "8443", cfg.Port)
	assert.Equal(t, "2.11", cfg.Version)
	assert.Equal(t, "Internal", cfg.View)
	assert.True(t, cfg.SslVerify)
	assert.Equal(t, 90, cfg.HTTPRequestTimeout)
	assert.Equal(t, 20, cfg.HTTPPoolConnections)
	assert.Equal(t, uint32(600), cfg.TTL)
	assert.True(t, cfg.UseTTL)
}

// TestLoadConfig_Nil tests the base case with no configuration
func TestLoadConfig_Nil(t *testing.T) {
	cfg, err := loadConfig(nil)

	require.NoError(t, err)
	// Should get zero values
	assert.Equal(t, "", cfg.Host)
}

// TestLoadConfig_Empty tests with empty JSON object
func TestLoadConfig_Empty(t *testing.T) {
	raw := apiextensionsv1.JSON{Raw: []byte("{}")}
	cfg, err := loadConfig(&raw)

	require.NoError(t, err)
	assert.Equal(t, "", cfg.Host)
}

// TestLoadConfig_InvalidJSON tests error handling for malformed JSON
func TestLoadConfig_InvalidJSON(t *testing.T) {
	raw := apiextensionsv1.JSON{Raw: []byte("invalid json")}
	_, err := loadConfig(&raw)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Error decoding solver config")
}

// TestApplyDefaults tests that defaults are correctly applied
func TestApplyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    customDNSProviderConfig
		expected customDNSProviderConfig
	}{
		{
			name:  "all empty values get defaults",
			input: customDNSProviderConfig{},
			expected: customDNSProviderConfig{
				Port:                "443",
				Version:             "2.10",
				HTTPRequestTimeout:  60,
				HTTPPoolConnections: 10,
				TTL:                 300,
			},
		},
		{
			name: "custom values are preserved",
			input: customDNSProviderConfig{
				Port:                "8443",
				Version:             "2.11",
				HTTPRequestTimeout:  90,
				HTTPPoolConnections: 20,
				TTL:                 600,
			},
			expected: customDNSProviderConfig{
				Port:                "8443",
				Version:             "2.11",
				HTTPRequestTimeout:  90,
				HTTPPoolConnections: 20,
				TTL:                 600,
			},
		},
		{
			name: "partial config gets remaining defaults",
			input: customDNSProviderConfig{
				Port:    "8443",
				Version: "2.11",
			},
			expected: customDNSProviderConfig{
				Port:                "8443",
				Version:             "2.11",
				HTTPRequestTimeout:  60,
				HTTPPoolConnections: 10,
				TTL:                 300,
			},
		},
		{
			name: "negative timeout gets default",
			input: customDNSProviderConfig{
				HTTPRequestTimeout: -10,
			},
			expected: customDNSProviderConfig{
				Port:                "443",
				Version:             "2.10",
				HTTPRequestTimeout:  60,
				HTTPPoolConnections: 10,
				TTL:                 300,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.input
			applyDefaults(&cfg)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

// TestLoadConfig_WithDefaults tests that loadConfig applies defaults
func TestLoadConfig_WithDefaults(t *testing.T) {
	// Config with only host set
	configJSON := `{"host": "infoblox.example.com"}`
	raw := apiextensionsv1.JSON{Raw: []byte(configJSON)}

	cfg, err := loadConfig(&raw)

	require.NoError(t, err)
	assert.Equal(t, "infoblox.example.com", cfg.Host)
	// Check defaults were applied
	assert.Equal(t, "443", cfg.Port)
	assert.Equal(t, "2.10", cfg.Version)
	assert.Equal(t, 60, cfg.HTTPRequestTimeout)
	assert.Equal(t, 10, cfg.HTTPPoolConnections)
	assert.Equal(t, uint32(300), cfg.TTL)
}

// TestGetSecret tests secret retrieval from Kubernetes
func TestGetSecret(t *testing.T) {
	// Create fake clientset with a test secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test-namespace",
		},
		Data: map[string][]byte{
			"username": []byte("testuser"),
			"password": []byte("testpass\n"), // with newline to test trimming
		},
	}

	fakeClient := fake.NewSimpleClientset(secret)
	solver := &customDNSProviderSolver{client: fakeClient}

	tests := []struct {
		name      string
		selector  cmmeta.SecretKeySelector
		namespace string
		expected  string
		wantError bool
		errorMsg  string
	}{
		{
			name: "successful retrieval",
			selector: cmmeta.SecretKeySelector{
				LocalObjectReference: cmmeta.LocalObjectReference{
					Name: "test-secret",
				},
				Key: "username",
			},
			namespace: "test-namespace",
			expected:  "testuser",
			wantError: false,
		},
		{
			name: "newline trimming",
			selector: cmmeta.SecretKeySelector{
				LocalObjectReference: cmmeta.LocalObjectReference{
					Name: "test-secret",
				},
				Key: "password",
			},
			namespace: "test-namespace",
			expected:  "testpass",
			wantError: false,
		},
		{
			name: "secret not found",
			selector: cmmeta.SecretKeySelector{
				LocalObjectReference: cmmeta.LocalObjectReference{
					Name: "nonexistent-secret",
				},
				Key: "username",
			},
			namespace: "test-namespace",
			wantError: true,
			errorMsg:  "failed to get secret",
		},
		{
			name: "key not found in secret",
			selector: cmmeta.SecretKeySelector{
				LocalObjectReference: cmmeta.LocalObjectReference{
					Name: "test-secret",
				},
				Key: "nonexistent-key",
			},
			namespace: "test-namespace",
			wantError: true,
			errorMsg:  "key nonexistent-key not found in secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := solver.getSecret(tt.selector, tt.namespace)

			if tt.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestGetIbClient_WithSecretRefs tests client initialization with secret references
func TestGetIbClient_WithSecretRefs(t *testing.T) {
	// Create fake clientset with credentials
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "infoblox-creds",
			Namespace: "test-namespace",
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
			"password": []byte("secret123"),
		},
	}

	fakeClient := fake.NewSimpleClientset(secret)
	solver := &customDNSProviderSolver{client: fakeClient}

	cfg := customDNSProviderConfig{
		Host:    "infoblox.example.com",
		Port:    "443",
		Version: "2.10",
		UsernameSecretRef: cmmeta.SecretKeySelector{
			LocalObjectReference: cmmeta.LocalObjectReference{
				Name: "infoblox-creds",
			},
			Key: "username",
		},
		PasswordSecretRef: cmmeta.SecretKeySelector{
			LocalObjectReference: cmmeta.LocalObjectReference{
				Name: "infoblox-creds",
			},
			Key: "password",
		},
		HTTPRequestTimeout:  60,
		HTTPPoolConnections: 10,
	}

	ib, err := solver.getIbClient(&cfg, "test-namespace")

	require.NoError(t, err)
	assert.NotNil(t, ib)
}

// TestGetIbClient_NoCredentials tests error when no credentials configured
func TestGetIbClient_NoCredentials(t *testing.T) {
	solver := &customDNSProviderSolver{}

	cfg := customDNSProviderConfig{
		Host:    "infoblox.example.com",
		Port:    "443",
		Version: "2.10",
	}

	_, err := solver.getIbClient(&cfg, "test-namespace")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "No secretRefs or secretPath provided")
}

// TestGetIbClient_FromVolume tests loading credentials from volume file
func TestGetIbClient_FromVolume(t *testing.T) {
	// Create temporary credentials file
	tmpDir := t.TempDir()
	credsFile := filepath.Join(tmpDir, "creds.json")

	creds := usernamePassword{
		Username: "volumeuser",
		Password: "volumepass",
	}
	credsJSON, err := json.Marshal(creds)
	require.NoError(t, err)

	err = os.WriteFile(credsFile, credsJSON, 0o600)
	require.NoError(t, err)

	// Temporarily change SecretPath constant behavior by testing the logic
	// Note: In real scenario, we'd need to mount this at /etc/secrets/creds.json
	// For unit test, we'll skip this test or mock the file read
	t.Skip("Skipping volume test - requires mocking file system or integration test")
}

// TestGetIbClient_VolumeFileNotFound tests error when volume file doesn't exist
func TestGetIbClient_VolumeFileNotFound(t *testing.T) {
	solver := &customDNSProviderSolver{}

	cfg := customDNSProviderConfig{
		Host:              "infoblox.example.com",
		GetUserFromVolume: true,
	}

	_, err := solver.getIbClient(&cfg, "test-namespace")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

// TestGetIbClient_AppliesDefaults tests that defaults are applied defensively
func TestGetIbClient_AppliesDefaults(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "infoblox-creds",
			Namespace: "test-namespace",
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
			"password": []byte("secret123"),
		},
	}

	fakeClient := fake.NewSimpleClientset(secret)
	solver := &customDNSProviderSolver{client: fakeClient}

	// Config without defaults set
	cfg := customDNSProviderConfig{
		Host: "infoblox.example.com",
		// Port, Version, etc. not set
		UsernameSecretRef: cmmeta.SecretKeySelector{
			LocalObjectReference: cmmeta.LocalObjectReference{
				Name: "infoblox-creds",
			},
			Key: "username",
		},
		PasswordSecretRef: cmmeta.SecretKeySelector{
			LocalObjectReference: cmmeta.LocalObjectReference{
				Name: "infoblox-creds",
			},
			Key: "password",
		},
	}

	ib, err := solver.getIbClient(&cfg, "test-namespace")

	require.NoError(t, err)
	assert.NotNil(t, ib)
	// Verify defaults were applied
	assert.Equal(t, "443", cfg.Port)
	assert.Equal(t, "2.10", cfg.Version)
	assert.Equal(t, 60, cfg.HTTPRequestTimeout)
	assert.Equal(t, 10, cfg.HTTPPoolConnections)
	assert.Equal(t, uint32(300), cfg.TTL)
}

// TestInitialize tests the Initialize method
func TestInitialize(t *testing.T) {
	// Create a fake REST config
	// Note: This would normally connect to a real API server
	// For unit testing, we skip or would need more sophisticated mocking
	t.Skip("Initialize requires valid kubeClientConfig - better suited for integration test")
}

// Benchmark tests for performance-critical functions
func BenchmarkDeDot(b *testing.B) {
	solver := &customDNSProviderSolver{}
	input := "_acme-challenge.example.com."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solver.DeDot(input)
	}
}

func BenchmarkLoadConfig(b *testing.B) {
	configJSON := `{
		"host": "infoblox.example.com",
		"port": "443",
		"version": "2.10",
		"view": "Internal"
	}`
	raw := apiextensionsv1.JSON{Raw: []byte(configJSON)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = loadConfig(&raw)
	}
}

func BenchmarkApplyDefaults(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg := customDNSProviderConfig{}
		applyDefaults(&cfg)
	}
}

// TestGetSecret_ContextTimeout tests that getSecret respects context timeout
func TestGetSecret_ContextTimeout(t *testing.T) {
	// Create a fake client that simulates a slow response
	// This would require more sophisticated mocking
	t.Skip("Context timeout testing requires sophisticated mocking")
}

// Table-driven test for various configuration combinations
func TestLoadConfig_Combinations(t *testing.T) {
	tests := []struct {
		name       string
		configJSON string
		validate   func(t *testing.T, cfg customDNSProviderConfig, err error)
	}{
		{
			name:       "minimal valid config",
			configJSON: `{"host": "infoblox.local"}`,
			validate: func(t *testing.T, cfg customDNSProviderConfig, err error) {
				require.NoError(t, err)
				assert.Equal(t, "infoblox.local", cfg.Host)
				assert.Equal(t, "443", cfg.Port) // default
			},
		},
		{
			name: "config with all fields",
			configJSON: `{
				"host": "infoblox.local",
				"port": "8443",
				"version": "2.12",
				"view": "External",
				"sslVerify": true,
				"httpRequestTimeout": 120,
				"httpPoolConnections": 25,
				"getUserFromVolume": true,
				"ttl": 900,
				"useTtl": false
			}`,
			validate: func(t *testing.T, cfg customDNSProviderConfig, err error) {
				require.NoError(t, err)
				assert.Equal(t, "infoblox.local", cfg.Host)
				assert.Equal(t, "8443", cfg.Port)
				assert.Equal(t, "2.12", cfg.Version)
				assert.Equal(t, "External", cfg.View)
				assert.True(t, cfg.SslVerify)
				assert.Equal(t, 120, cfg.HTTPRequestTimeout)
				assert.Equal(t, 25, cfg.HTTPPoolConnections)
				assert.True(t, cfg.GetUserFromVolume)
				assert.Equal(t, uint32(900), cfg.TTL)
				assert.False(t, cfg.UseTTL)
			},
		},
		{
			name:       "config with zero values explicitly set",
			configJSON: `{"host": "infoblox.local", "ttl": 0}`,
			validate: func(t *testing.T, cfg customDNSProviderConfig, err error) {
				require.NoError(t, err)
				// Even though user set 0, default should apply
				assert.Equal(t, uint32(300), cfg.TTL)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := apiextensionsv1.JSON{Raw: []byte(tt.configJSON)}
			cfg, err := loadConfig(&raw)
			tt.validate(t, cfg, err)
		})
	}
}
