package provider

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakAuthenticationFlow_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakAuthenticationFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationFlow_basic(realmName, clientId),
				Check:  testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.authentication_flow"),
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

func testAccCheckKeycloakAuthenticationFlowExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		flow, err := getAuthenticationFlowFromState(s, resourceName)
		if err != nil {
			return err
		}

		if flow.Alias != "something" {
			return fmt.Errorf("expected flow to have alias of something, but got %s", flow.Alias)
		}

		return nil
	}
}

func getAuthenticationFlowFromState(s *terraform.State, resourceName string) (*keycloak.AuthenticationFlow, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	client, err := keycloakClient.GetAuthenticationFlow(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting flow %s: %s", id, err)
	}

	return client, nil
}

func testKeycloakAuthenticationFlow_basic(realm, clientId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "authentication_flow" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}
	`, realm, clientId)
}
