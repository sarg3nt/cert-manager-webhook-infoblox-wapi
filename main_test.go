package main

// TODO:  Test suite is not currently working.
// See: https://github.com/cert-manager/webhook-example/tree/master

import (
	"os"
	"testing"

	acmetest "github.com/cert-manager/cert-manager/test/acme"
	//"github.com/cert-manager/webhook-example/example"
	//"github.com/sarg3nt/cert-manager-webhook-infoblox-wapi"
)

var zone = os.Getenv("TEST_ZONE_NAME")

func TestRunsSuite(t *testing.T) {
	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.
	//

	fixture := acmetest.NewFixture(&customDNSProviderSolver{},
		acmetest.SetResolvedZone(zone),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetManifestPath("testdata/infoblox-wapi"),
	)

	// solver := example.New("59351")
	// fixture := acmetest.NewFixture(solver,
	// 	acmetest.SetResolvedZone("example.com."),
	// 	acmetest.SetManifestPath("testdata/my-custom-solver"),
	// 	acmetest.SetDNSServer("127.0.0.1:59351"),
	// 	acmetest.SetUseAuthoritative(false),
	// )
	// need to uncomment and  RunConformance delete runBasic and runExtended once https://github.com/cert-manager/cert-manager/pull/4835 is merged
	// fixture.RunConformance(t)
	fixture.RunBasic(t)
	fixture.RunExtended(t)
}
