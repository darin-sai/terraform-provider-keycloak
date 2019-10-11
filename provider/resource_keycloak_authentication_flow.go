package provider

import (
	// "fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	// "strings"
)

func resourceKeycloakAuthenticationFlow() *schema.Resource {

	return &schema.Resource{
		Create: resourceKeycloakAuthenticationFlowCreate,
		Read:   resourceKeycloakAuthenticationRead,
		// Delete: resourceKeycloakFlowDelete,
		// Update: resourceKeycloakFlowUpdate,
		// Importer: &schema.ResourceImporter{
		// 	State: resourceKeycloakFlowImport,
		// },
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceKeycloakAuthenticationFlowCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	flow := mapToAuthenticationFlowFromData(data)

	err := keycloakClient.NewAuthenticationFlow(flow)
	if err != nil {
		return err
	}

	data.SetId(flow.Id)

	return resourceKeycloakSamlClientRead(data, meta)
}

func resourceKeycloakAuthenticationRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	flow, err := keycloakClient.GetAuthenticationFlow(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = mapToDataFromAuthenticationFlow(data, flow)
	if err != nil {
		return err
	}

	return nil
}

func mapToDataFromAuthenticationFlow(data *schema.ResourceData, flow *keycloak.AuthenticationFlow) error {
	data.SetId(flow.Id)
	data.Set("realm_id", flow.RealmId)
	data.Set("alias", flow.Alias)

	return nil
}

func mapToAuthenticationFlowFromData(data *schema.ResourceData) *keycloak.AuthenticationFlow {
	flow := &keycloak.AuthenticationFlow{
		Id:      data.Id(),
		RealmId: data.Get("realm_id").(string),
		Alias:   data.Get("alias").(string),
	}

	return flow
}
