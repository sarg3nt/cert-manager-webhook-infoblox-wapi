package main

// Webhook example plugin:
// https://github.com/cert-manager/webhook-example

// cspell:ignore cmapi cmacme klog
import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	webhook "github.com/cert-manager/cert-manager/pkg/acme/webhook"
	whapi "github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

const SecretPath = "/etc/secrets/creds.json"

// var _ webhook.Solver = (*customDNSProviderSolver)(nil)

// GroupName is the API group name for the webhook, set via GROUP_NAME environment variable.
var GroupName = os.Getenv("GROUP_NAME")

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	klog.InfoS("CMI: GroupName is", "GroupName", GroupName)
	// This will register our custom DNS provider with the webhook serving
	// library, making it available as an API under the provided groupName.
	// You can register multiple DNS provider implementations with a single
	// webhook, where the Name() method will be used to disambiguate between
	// the different implementations.

	cmd.RunWebhookServer(GroupName, &customDNSProviderSolver{})
}

// Validate the customDNSProviderSolver satisfies the interface that cert manager expects.
var _ webhook.Solver = (*customDNSProviderSolver)(nil)

// customDNSProviderSolver implements the provider-specific logic needed to
// 'present' an ACME challenge TXT record for your own DNS provider.
// To do so, it must implement the `github.com/cert-manager/cert-manager/pkg/acme/webhook.Solver`
// interface.
type customDNSProviderSolver struct {
	// If a Kubernetes 'clientset' is needed, you must:
	// 1. uncomment the additional `client` field in this structure below
	// 2. uncomment the "k8s.io/client-go/kubernetes" import at the top of the file
	// 3. uncomment the relevant code in the Initialize method below
	// 4. ensure your webhook's service account has the required RBAC role
	//    assigned to it for interacting with the Kubernetes APIs you need.
	client *kubernetes.Clientset
}

// customDNSProviderConfig is a structure that is used to decode into when
// solving a DNS01 challenge.
// This information is provided by cert-manager, and may be a reference to
// additional configuration that's needed to solve the challenge for this
// particular certificate or issuer.
// This typically includes references to Secret resources containing DNS
// provider credentials, in cases where a 'multi-tenant' DNS solver is being
// created.
// If you do *not* require per-issuer or per-certificate configuration to be
// provided to your webhook, you can skip decoding altogether in favour of
// using CLI flags or similar to provide configuration.
// You should not include sensitive information here. If credentials need to
// be used by your provider here, you should reference a Kubernetes Secret
// resource and fetch these credentials using a Kubernetes clientset.
type customDNSProviderConfig struct {
	// Change the two fields below according to the format of the configuration
	// to be decoded.
	// These fields will be set by users in the
	// `issuer.spec.acme.dns01.providers.webhook.config` field.

	Host                string                   `json:"host"`
	Port                string                   `json:"port"                default:"443"`
	Version             string                   `json:"version"             default:"2.10"`
	UsernameSecretRef   cmmeta.SecretKeySelector `json:"usernameSecretRef"`
	PasswordSecretRef   cmmeta.SecretKeySelector `json:"passwordSecretRef"`
	View                string                   `json:"view"`
	SslVerify           bool                     `json:"sslVerify"           default:"false"`
	HTTPRequestTimeout  int                      `json:"httpRequestTimeout"  default:"60"`
	HTTPPoolConnections int                      `json:"httpPoolConnections" default:"10"`
	GetUserFromVolume   bool                     `json:"getUserFromVolume"   default:"false"`
	TTL                 uint32                   `json:"ttl"                 default:"300"`
	UseTTL              bool                     `json:"useTtl"              default:"true"`
}

type usernamePassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource.
// This should be unique **within the group name**, i.e. you can have two
// solvers configured with the same Name() **so long as they do not co-exist
// within a single webhook deployment**.
// For example, `cloudflare` may be used as the name of a solver.
func (c *customDNSProviderSolver) Name() string {
	return "infoblox-wapi"
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (c *customDNSProviderSolver) Present(ch *whapi.ChallengeRequest) error {
	klog.InfoS("CMI: Presenting DNS record", "DNS", ch.DNSName)
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		klog.InfoS("CMI: Error loading config", "error", err.Error())
		return err
	}

	// Initialize ibclient
	ib, err := c.getIbClient(&cfg, ch.ResourceNamespace)
	if err != nil {
		klog.InfoS("CMI: Error getting Infoblox client", "error", err.Error())
		return err
	}

	// Find or create TXT record
	recordName := c.DeDot(ch.ResolvedFQDN)
	klog.InfoS("CMI: Record name", "name", recordName)

	klog.InfoS("CMI: Getting current txt record.", "key", ch.Key)
	recordRef, err := c.GetTXTRecord(ib, recordName, ch.Key, cfg.View)
	klog.InfoS("CMI: Record ref after getting current txt record", "recordRef", recordRef)

	if err != nil {
		klog.InfoS("CMI: Error getting TXT record", "name", recordName, "error", err.Error())
		return err
	}

	// If record exists, delete it first to ensure we create with correct value
	// This ensures idempotent behavior as required by cert-manager
	if recordRef != "" {
		klog.InfoS("CMI: TXT record already exists, deleting before recreating", "name", recordName, "ref", recordRef)
		err = c.DeleteTXTRecord(ib, recordRef)
		if err != nil {
			klog.InfoS("CMI: Error deleting existing TXT record", "name", recordName, "error", err.Error())
			return err
		}
		klog.InfoS("CMI: Deleted existing TXT record", "name", recordName, "ref", recordRef)
	}

	// Create the TXT record
	klog.InfoS("CMI: Creating TXT record", "name", recordName)
	recordRef, err = c.CreateTXTRecord(ib, recordName, ch.Key, cfg.View, cfg.TTL, cfg.UseTTL)
	klog.InfoS("CMI: Record ref after creating txt record", "recordRef", recordRef)

	if err != nil {
		klog.InfoS("CMI: Error creating TXT record", "name", recordName, "error", err.Error())
		return err
	}

	klog.InfoS("CMI: Successfully created TXT record", "name", recordName, "ref", recordRef)

	klog.InfoS("CMI: Done presenting for DNS record", "DNS", ch.DNSName)
	return nil
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (c *customDNSProviderSolver) CleanUp(ch *whapi.ChallengeRequest) error {
	klog.InfoS("CMI: Cleaning up")
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}

	// Initialize ibclient
	ib, err := c.getIbClient(&cfg, ch.ResourceNamespace)
	if err != nil {
		return err
	}

	// Find and delete TXT record
	recordName := c.DeDot(ch.ResolvedFQDN)

	recordRef, err := c.GetTXTRecord(ib, recordName, ch.Key, cfg.View)
	if err != nil {
		return err
	}

	if recordRef == "" {
		klog.InfoS("CMI: TXT record not found, skipping deletion", "name", recordName, "text", ch.Key)
		return nil
	}

	err = c.DeleteTXTRecord(ib, recordRef)
	if err != nil {
		return err
	}
	klog.InfoS("CMI: Deleted TXT record", "name", recordName, "ref", recordRef)

	return nil
}

// Initialize will be called when the webhook first starts.
// This method can be used to instantiate the webhook, i.e. initialising
// ibections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes
// client that can be used to fetch resources from the Kubernetes API, e.g.
// Secret resources containing credentials used to authenticate with DNS
// provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases
// where a SIGTERM or similar signal is sent to the webhook process.
func (c *customDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	klog.InfoS("CMI: Initializing k8s client")
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		klog.InfoS("CMI: Error initializing k8s client.", "error", err.Error())
		return err
	}
	klog.InfoS("CMI: Initialized k8s client")
	c.client = cl

	return nil
}

// loadConfig is a small helper function that decodes JSON configuration into
// the typed config struct.
func loadConfig(cfgJSON *apiextensionsv1.JSON) (customDNSProviderConfig, error) {
	klog.InfoS("CMI: Loading config")

	cfg := customDNSProviderConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("CMI: Error decoding solver config: %v", err)
	}

	return cfg, nil
}

// Initialize and return infoblox client connector
// Configuration can be set in the webhook `config` section.
// Two secretRefs are needed to securely pass infoblox credentials
func (c *customDNSProviderSolver) getIbClient(cfg *customDNSProviderConfig, namespace string) (ibclient.IBConnector, error) {
	var username, password string
	hasConfig := false

	klog.InfoS("CMI: Getting Infoblox User Data")
	if cfg.UsernameSecretRef.Key != "" && cfg.PasswordSecretRef.Key != "" {
		klog.InfoS("CMI: Getting Infoblox User and Password from secret")
		hasConfig = true
		var err error
		// Find secret credentials
		username, err = c.getSecret(cfg.UsernameSecretRef, namespace)
		if err != nil {
			return nil, err
		}

		password, err = c.getSecret(cfg.PasswordSecretRef, namespace)
		if err != nil {
			return nil, err
		}
		klog.InfoS("CMI: Infoblox User", "username", username)
	}

	if cfg.GetUserFromVolume && !hasConfig {
		klog.InfoS("CMI: Getting Infoblox User and Password from volume")
		hasConfig = true

		if _, err := os.Stat(SecretPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("CMI: File %s does not exist", SecretPath)
		}

		fileData, err := os.ReadFile(SecretPath)
		if err != nil {
			return nil, err
		}

		var creds usernamePassword
		if err := json.Unmarshal(fileData, &creds); err != nil {
			return nil, err
		}

		username = creds.Username
		password = creds.Password
		klog.InfoS("CMI: Infoblox User", "username", username)
	}

	if !hasConfig {
		return nil, fmt.Errorf("CMI: No secretRefs or secretPath provided")
	}

	// Set default values if needed
	if cfg.Port == "" {
		cfg.Port = "443"
	}
	if cfg.Version == "" {
		cfg.Version = "2.10"
	}
	if cfg.HTTPRequestTimeout <= 0 {
		cfg.HTTPRequestTimeout = 60
	}
	if cfg.HTTPPoolConnections <= 0 {
		cfg.HTTPPoolConnections = 10
	}
	if cfg.TTL == 0 {
		cfg.TTL = 300
	}

	// Initialize ibclient
	hostConfig := ibclient.HostConfig{
		Host:    cfg.Host,
		Version: cfg.Version,
		Port:    cfg.Port,
	}

	// Initialize the auth config for ibclient
	authConfig := ibclient.AuthConfig{
		Username: username,
		Password: password,
	}

	transportConfig := ibclient.NewTransportConfig(strconv.FormatBool(cfg.SslVerify), cfg.HTTPRequestTimeout, cfg.HTTPPoolConnections)
	requestBuilder := &ibclient.WapiRequestBuilder{}
	requestor := &ibclient.WapiHttpRequestor{}

	ib, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
	if err != nil {
		klog.InfoS("CMI: Error creating Infoblox client", "error", err.Error())
		return nil, err
	}

	return ib, nil
}

// Resolve the value of a secret given a SecretKeySelector with name and key parameters
func (c *customDNSProviderSolver) getSecret(sel cmmeta.SecretKeySelector, namespace string) (string, error) {
	klog.InfoS("CMI: Getting secret", "name", sel.Name, "namespace", namespace)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	secret, err := c.client.CoreV1().Secrets(namespace).Get(ctx, sel.Name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s/%s: %w", namespace, sel.Name, err)
	}

	secretData, ok := secret.Data[sel.Key]
	if !ok {
		return "", fmt.Errorf("key %s not found in secret %s/%s", sel.Key, namespace, sel.Name)
	}

	return strings.TrimSuffix(string(secretData), "\n"), nil
}

// Get the ref for TXT record in InfoBlox given its name, text and view
func (c *customDNSProviderSolver) GetTXTRecord(ib ibclient.IBConnector, name string, text string, view string) (string, error) {
	klog.InfoS("CMI: Getting TXT record", "name", name)
	var records []ibclient.RecordTXT
	recordTXT := ibclient.NewEmptyRecordTXT()
	params := map[string]string{
		"name": name,
		"text": text,
		"view": view,
	}
	err := ib.GetObject(recordTXT, "", ibclient.NewQueryParams(false, params), &records)
	klog.InfoS("CMI: Number of records is", "number", strconv.Itoa(len(records)))

	if len(records) > 0 {
		klog.InfoS("CMI: Found TXT record")
		return records[0].Ref, nil
	}

	// No records found - check if it's a NotFoundError (expected) or real error
	if _, ok := err.(*ibclient.NotFoundError); ok {
		klog.InfoS("CMI: No TXT record found. This can be normal for the first run.")
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return "", nil
}

// Create a TXT record in Infoblox
func (c *customDNSProviderSolver) CreateTXTRecord(ib ibclient.IBConnector, name string, text string, view string, ttl uint32, useTTL bool) (string, error) {
	klog.InfoS("CMI: Creating TXT record", "name", name)

	recordTXT := ibclient.NewRecordTXT(view, "", name, text, ttl, useTTL, "", nil)
	klog.InfoS("CMI: RecordTXT", "recordTXT", recordTXT)
	return ib.CreateObject(recordTXT)
}

// Delete a TXT record in Infoblox by ref
func (c *customDNSProviderSolver) DeleteTXTRecord(ib ibclient.IBConnector, ref string) error {
	klog.InfoS("CMI: Deleting TXT record", "ref", ref)
	_, err := ib.DeleteObject(ref)

	return err
}

// Remove trailing dot
func (c *customDNSProviderSolver) DeDot(string string) string {
	klog.InfoS("CMI: Removing trailing dot")
	result := strings.TrimSuffix(string, ".")

	return result
}
