package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type GraphHelper struct {
	clientSecretCredential *azidentity.ClientSecretCredential
	graphClient            *msgraphsdk.GraphServiceClient
}

func NewGraphHelper() *GraphHelper {
	g := &GraphHelper{}
	return g
}

func (g *GraphHelper) InitializeGraph(tenantId string, clientId string, clientSecret string) error {
	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret,
		&azidentity.ClientSecretCredentialOptions{
			ClientOptions: policy.ClientOptions{
				Retry: policy.RetryOptions{
					MaxRetries:    3,
					MaxRetryDelay: time.Duration(30) * time.Second,
				},
				Logging: policy.LogOptions{
					IncludeBody: true,
				},
			},
		})
	if err != nil {
		return err
	}
	g.clientSecretCredential = credential

	// Create an auth provider using the credential
	authProvider, err := authentication.NewAzureIdentityAuthenticationProviderWithScopes(
		g.clientSecretCredential,
		[]string{"https://graph.microsoft.com/.default"},
	)
	if err != nil {
		return err
	}

	// Create a request adapter using the auth provider
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return err
	}

	// Create a Graph client using request adapter
	client := msgraphsdk.NewGraphServiceClient(adapter)
	g.graphClient = client

	return nil
}

func (g *GraphHelper) GetDevices() (models.ManagedDeviceCollectionResponseable, error) {
	ctx := context.Background()
	requestBuilder := g.graphClient.DeviceManagement().ManagedDevices()
	result, err := requestBuilder.Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting managed devices: %v", err)
	}

	return result, err
}

func (g *GraphHelper) GetDevice(deviceId string) (models.ManagedDeviceable, error) {
	ctx := context.Background()
	requestBuilder := g.graphClient.DeviceManagement().ManagedDevices().ByManagedDeviceId(deviceId)
	result, err := requestBuilder.Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting managed device: %v", err)
	}
	return result, err
}

func (g *GraphHelper) GetCompliancePolicies() (models.DeviceCompliancePolicyCollectionResponseable, error) {
	ctx := context.Background()
	requestBuilder := g.graphClient.DeviceManagement().DeviceCompliancePolicies()
	result, err := requestBuilder.Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting device compliance policies: %v", err)
	}
	return result, err
}

func (g *GraphHelper) GetDevicesWithCompliancePolicy(compliancePolicyId string) (models.DeviceComplianceDeviceStatusCollectionResponseable, error) {
	ctx := context.Background()
	requestBuilder := g.graphClient.DeviceManagement().DeviceCompliancePolicies().ByDeviceCompliancePolicyId(compliancePolicyId).DeviceStatuses()
	result, err := requestBuilder.Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting device status with compliance policies: %v", err)
	}
	return result, err
}
