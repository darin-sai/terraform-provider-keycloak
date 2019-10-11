package provider

import (
	"fmt"
	// "strconv"
	// "strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakAuthenticationFlow_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	config := fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "authentication_flow" {
	realm_id  = "${keycloak_realm.realm.id}"
	provider_id = "basic-flow"
	alias = "some alias"
}`, realmName)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakAuthenticationFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr("keycloak_authentication_flow.authentication_flow", "alias", "some alias"),
			},
			{
				ResourceName:        "keycloak_authentication_flow.authentication_flow",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

// func testAccCheckKeycloakAuthenticationFlowExists(resourceName string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		_, err := getAuthenticationFlowFromState(s, resourceName)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// }

// func getAuthenticationFlowFromState(s *terraform.State, resourceName string) (*keycloak.AuthenticationFlow, error) {
// 	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

// 	rs, ok := s.RootModule().Resources[resourceName]
// 	if !ok {
// 		return nil, fmt.Errorf("resource not found: %s", resourceName)
// 	}

// 	id := rs.Primary.ID
// 	realm := rs.Primary.Attributes["realm_id"]

// 	client, err := keycloakClient.GetAuthenticationFlow(realm, id)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting flow %s: %s", id, err)
// 	}

// 	return client, nil
// }

// func testKeycloakAuthenticationFlow_basic(realm, alias string) string {
// 	return fmt.Sprintf(`
// resource "keycloak_realm" "realm" {
// 	realm = "%s"
// }

// resource "keycloak_authentication_flow" "authentication_flow" {
// 	realm_id  = "${keycloak_realm.realm.id}"
// 	provider_id = "basic-flow"
// 	alias = "%s"
// }
// 	`, realm, alias)
// }

func testAccCheckKeycloakAuthenticationFlowDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_authentication_flow" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			client, _ := keycloakClient.GetAuthenticationFlow(realm, id)
			if client != nil {
				return fmt.Errorf("authentication flow %s still exists", id)
			}
		}

		return nil
	}
}
