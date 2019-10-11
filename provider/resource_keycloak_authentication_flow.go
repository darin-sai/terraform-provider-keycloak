package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakAuthenticationFlow() *schema.Resource {

	return &schema.Resource{
		Create: resourceKeycloakAuthenticationFlowCreate,
		Read:   resourceKeycloakAuthenticationRead,
		Delete: resourceKeycloakAuthenticationDelete,
		Update: resourceKeycloakAuthenticationUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakAuthenticationImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_id": {
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

func resourceKeycloakAuthenticationDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteAuthenticationFlow(realmId, id)
}

func resourceKeycloakAuthenticationUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	flow := mapToAuthenticationFlowFromData(data)

	err := keycloakClient.UpdateAuthenticationFlow(flow)
	if err != nil {
		return err
	}

	err = mapToDataFromAuthenticationFlow(data, flow)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakAuthenticationImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{flowId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func mapToDataFromAuthenticationFlow(data *schema.ResourceData, flow *keycloak.AuthenticationFlow) error {
	data.SetId(flow.Id)
	data.Set("realm_id", flow.RealmId)
	data.Set("alias", flow.Alias)
	data.Set("provider_id", flow.ProviderId)

	return nil
}

func mapToAuthenticationFlowFromData(data *schema.ResourceData) *keycloak.AuthenticationFlow {
	flow := &keycloak.AuthenticationFlow{
		Id:         data.Id(),
		RealmId:    data.Get("realm_id").(string),
		Alias:      data.Get("alias").(string),
		ProviderId: data.Get("provider_id").(string),
	}

	return flow
}
