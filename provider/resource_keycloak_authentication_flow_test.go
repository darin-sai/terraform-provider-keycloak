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
	realmName := "terraform-r-" + acctest.RandString(10)
	config := fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id  = "${keycloak_realm.realm.id}"
	provider_id = "basic-flow"
	alias = "some alias"
	built_in = true
	top_level = false
	description = "Some kind of thing"
}`, realmName)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakAuthenticationFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
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

func TestAccKeycloakAuthenticationFlow_updateRealm(t *testing.T) {
	realmOne := "terraform-r-" + acctest.RandString(10)
	realmTwo := "terraform-r-" + acctest.RandString(10)

	template := `
	resource "keycloak_realm" "realm_1" {
		realm = "%s"
	}
	resource "keycloak_realm" "realm_2" {
		realm = "%s"
	}
	
	resource "keycloak_authentication_flow" "flow" {
		realm_id  = "${keycloak_realm.%s.id}"
		provider_id = "basic-flow"
		alias = "some alias"
	}`

	initialConfig := fmt.Sprintf(template, realmOne, realmTwo, "realm_1")
	updatedConfig := fmt.Sprintf(template, realmOne, realmTwo, "realm_2")

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakAuthenticationFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
					resource.TestCheckResourceAttr("keycloak_authentication_flow.authentication_flow", "realm_id", realmOne),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
					resource.TestCheckResourceAttr("keycloak_authentication_flow.authentication_flow", "realm_id", realmTwo),
				),
			},
		},
	})
}

func TestAccKeycloakAuthenticationFlow_updateAuthenticationFlow(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)

	aliasOne := "terraform-flow-before-" + acctest.RandString(10)
	aliasTwo := "terraform-flow-after-" + acctest.RandString(10)

	template := `
		resource "keycloak_realm" "realm" {
			realm = "%s"
		}
		resource "keycloak_authentication_flow" "flow" {
			realm_id = "${keycloak_realm.realm.id}"
			alias    = "%s"
		}`

	initialConfig := fmt.Sprintf(template, realmName, aliasOne)
	updatedConfig := fmt.Sprintf(template, realmName, aliasOne)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakAuthenticationFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
					resource.TestCheckResourceAttr("keycloak_authentication_flow.flow", "alias", aliasOne),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationFlowExists("keycloak_authentication_flow.flow"),
					resource.TestCheckResourceAttr("keycloak_authentication_flow.flow", "alias", aliasTwo),
				),
			},
		},
	})
}

func testAccCheckKeycloakAuthenticationFlowExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getAuthenticationFlowFromState(s, resourceName)
		if err != nil {
			return err
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

	flow, err := keycloakClient.GetAuthenticationFlow(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting authentication flow with id %s: %s", id, err)
	}

	return flow, nil
}

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
