// go:build integration
//go:build integration
// +build integration

package main

import (
	"os"
	"testing"

	acmetest "github.com/cert-manager/cert-manager/test/acme"
)

var zone = os.Getenv("TEST_ZONE_NAME")

func TestRunsSuite(t *testing.T) {
	// Skip if no zone provided - allows CI to run without real DNS setup
	if zone == "" {
		t.Skip("TEST_ZONE_NAME not set, skipping integration tests")
	}

	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.
	fixture := acmetest.NewFixture(&customDNSProviderSolver{},
		acmetest.SetResolvedZone(zone),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetManifestPath("testdata/infoblox-wapi"),
		// Use SetDNSServer if testing against a real DNS server
		// acmetest.SetDNSServer("8.8.8.8:53"),
	)

	// Run the cert-manager conformance test suite
	fixture.RunConformance(t)
}
