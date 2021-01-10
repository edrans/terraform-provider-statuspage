package statuspage

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sp "github.com/sbecker59/statuspage-api-client-go/api/v1/statuspage"
)

func resourceUserRead(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	resp, _, err := statuspageClientV1.UsersApi.GetOrganizationsOrganizationIdUsers(authV1, d.Get("organization_id").(string)).Execute()
	if err.Error() != "" {
		return translateClientError(err, "failed to get user using Status Page API")
	}

	if &resp == nil {
		log.Printf("[INFO] Statuspage could not find user with ID: %s\n", d.Id())
		d.SetId("")
		return nil
	}

	user := find(resp, d.Id())

	if user != nil {
		d.Set("email", user.GetEmail())
		d.Set("first_name", user.GetFirstName())
		d.Set("last_name", user.GetLastName())
	}

	return nil

}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	user := sp.NewPostOrganizationsOrganizationIdUsersUser()

	user.SetEmail(d.Get("email").(string))
	user.SetFirstName(d.Get("first_name").(string))
	user.SetLastName(d.Get("last_name").(string))
	user.SetPassword(d.Get("password").(string))

	o := *sp.NewPostOrganizationsOrganizationIdUsers(*user)

	result, r, err := statuspageClientV1.UsersApi.PostOrganizationsOrganizationIdUsers(authV1, d.Get("organization_id").(string)).PostOrganizationsOrganizationIdUsers(o).Execute()

	if err.Error() != "" {
		fmt.Fprintf(os.Stderr, "Error when calling `UsersApi.PostOrganizationsOrganizationIdUsers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return translateClientError(err, "failed to create User using Status Page API")
	}

	d.SetId(result.GetId())

	return resourceUserRead(d, m)

}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {

	providerConf := m.(*ProviderConfiguration)
	statuspageClientV1 := providerConf.StatuspageClientV1
	authV1 := providerConf.AuthV1

	_, _, err := statuspageClientV1.UsersApi.DeleteOrganizationsOrganizationIdUsersUserId(authV1, d.Get("organization_id").(string), d.Id()).Execute()

	if err.Error() != "" {
		return translateClientError(err, "failed to delete User using Status Page API")
	}

	return nil

}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Delete: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization Identifier",
				ForceNew:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "Email address for the team member",
				Required:    true,
				ForceNew:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password the team member uses to access the site",
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"first_name": {
				Type:        schema.TypeString,
				Description: "First name member uses to access the site",
				Required:    true,
				ForceNew:    true,
			},
			"last_name": {
				Type:        schema.TypeString,
				Description: "Last name member uses to access the site",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func find(arr []sp.User, id string) *sp.User {
	for _, a := range arr {
		if a.GetId() == id {
			return &a
		}
	}
	return nil
}
