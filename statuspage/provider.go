package statuspage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
	providerVersion "github.com/sbecker59/terraform-provider-statuspage/version"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	retryablehttp "github.com/sbecker59/terraform-provider-statuspage/statuspage/internal/go-retryablehttp"
)

var (
	statuspageProvider *schema.Provider
)

func Provider() *schema.Provider {
	statuspageProvider = &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"STATUSPAGE_API_KEY", "SP_API_KEY"}, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"statuspage_component":       resourceComponent(),
			"statuspage_component_group": resourceComponentGroup(),
			"statuspage_incident":        resourceIncident(),
			"statuspage_metric":          resourceMetric(),
			"statuspage_metric_provider": resourceMetricProvider(),
			"statuspage_subscriber":      resourceSubscriber(),
			"statuspage_page_access_group": resourcePageAccessGroup(),
			"statuspage_page_access_user": resourcePageAccessUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"statuspage_component_groups": dataSourceComponentGroups(),
			"statuspage_components":       dataSourceComponents(),
		},
		ConfigureFunc: providerConfigure,
	}

	return statuspageProvider
}

type ProviderConfiguration struct {
	StatuspageClientV1 *sp.APIClient
	AuthV1             context.Context

	now func() time.Time
}

func (p *ProviderConfiguration) Now() time.Time {
	return p.now()
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing Status Page client")
	apiKey := d.Get("api_key").(string)

	if apiKey == "" {
		return nil, errors.New("api_key must be set unless")
	}

	authV1 := context.WithValue(
		context.Background(),
		sp.ContextAPIKeys,
		map[string]sp.APIKey{
			"api_key": {
				Key:    apiKey,
				Prefix: "oauth",
			},
		},
	)

	config := sp.NewConfiguration()
	config.HTTPClient = retryablehttp.NewClient().StandardClient()
	config.UserAgent = getUserAgent(config.UserAgent)
	statuspageClientV1 := sp.NewAPIClient(config)

	return &ProviderConfiguration{
		StatuspageClientV1: statuspageClientV1,
		AuthV1:             authV1,
		now:                time.Now,
	}, nil

}

func TranslateClientErrorDiag(err error, msg string) error {
	if msg == "" {
		msg = "an error occurred"
	}

	if apiErr, ok := err.(sp.GenericOpenAPIError); ok {
		return fmt.Errorf(msg+": %v: %s", err, apiErr.Body())
	}
	if errUrl, ok := err.(*url.Error); ok {
		return fmt.Errorf(msg+" (url.Error): %s", errUrl)
	}

	return fmt.Errorf(msg+": %s", err.Error())
}

// TranslateClientErrorDiagDiag returns client error as type diag.Diagnostics
func TranslateClientErrorDiagDiag(err error, msg string) diag.Diagnostics {
	return diag.FromErr(TranslateClientErrorDiag(err, msg))
}

func getUserAgent(clientUserAgent string) string {
	return fmt.Sprintf("terraform-provider-statuspage/%s (terraform %s; terraform-cli %s) %s",
		providerVersion.ProviderVersion,
		meta.SDKVersionString(),
		statuspageProvider.TerraformVersion,
		clientUserAgent)
}
