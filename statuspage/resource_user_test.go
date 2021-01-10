package statuspage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStatuspageUser_Basic(t *testing.T) {
	rid := acctest.RandIntRange(1, 99)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserConfig(rid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_user.default", "id"),
					resource.TestCheckResourceAttr("statuspage_user.default", "email", fmt.Sprintf("email-%d@testacc.tf", rid)),
					resource.TestCheckResourceAttr("statuspage_user.default", "first_name", "my_first_name"),
					resource.TestCheckResourceAttr("statuspage_user.default", "last_name", "my_last_name"),
				),
			},
		},
	})
}

func testAccCheckUserConfig(rand int) string {
	return fmt.Sprintf(`
	variable "email" {
		default = "email-%d@testacc.tf"
	}
	variable "organization_id" {
		default = "%s"
	}
	resource "statuspage_user" "default" {
		organization_id = var.organization_id
		email = var.email
		password = "my_password"
		first_name = "my_first_name"
		last_name = "my_last_name"
	}
	`, rand, organizationID)
}
