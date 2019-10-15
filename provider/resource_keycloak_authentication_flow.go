package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
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
				ForceNew: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"built_in": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"top_level": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"authentication_execution": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authenticator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authenticator_config": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"requirement": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_setup_allowed": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"autheticator_flow": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
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

	return resourceKeycloakAuthenticationRead(data, meta)
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
	data.Set("description", flow.Description)
	data.Set("built_in", flow.BuiltIn)
	data.Set("top_level", flow.TopLevel)

	executionCount := len(flow.AuthenticationExecutions)
	fmt.Println("Return trip exec count", executionCount)
	data.Set("authentication_execution.#", executionCount)
	for i := 0; i < executionCount; i++ {
		prefix := fmt.Sprintf(`authentication_execution.%d.`, i)
		data.Set(prefix+"authenticator", flow.AuthenticationExecutions[i].Authenticator)
		data.Set(prefix+"priority", flow.AuthenticationExecutions[i].Priority)
		data.Set(prefix+"requirement", flow.AuthenticationExecutions[i].Requirement)
		data.Set(prefix+"user_setup_allowed", flow.AuthenticationExecutions[i].UserSetupAllowed)
		data.Set(prefix+"autheticator_flow", flow.AuthenticationExecutions[i].AutheticatorFlow)
	}
	return nil
}

func mapToAuthenticationFlowFromData(data *schema.ResourceData) *keycloak.AuthenticationFlow {
	executionCount := data.Get("authentication_execution.#").(int)
	authenticationExecutions := make([]keycloak.AuthenticationExecution, executionCount)

	for i := 0; i < executionCount; i++ {
		prefix := fmt.Sprintf(`authentication_execution.%d.`, i)
		authenticationExecutions[i] = keycloak.AuthenticationExecution{
			Authenticator:    data.Get(prefix + "authenticator").(string),
			Priority:         data.Get(prefix + "priority").(int),
			Requirement:      data.Get(prefix + "requirement").(string),
			UserSetupAllowed: data.Get(prefix + "user_setup_allowed").(bool),
			AutheticatorFlow: data.Get(prefix + "autheticator_flow").(bool),
		}
		fmt.Println("Built authenticator ", i)
	}

	flow := &keycloak.AuthenticationFlow{
		Id:                       data.Id(),
		RealmId:                  data.Get("realm_id").(string),
		Alias:                    data.Get("alias").(string),
		ProviderId:               data.Get("provider_id").(string),
		Description:              data.Get("description").(string),
		BuiltIn:                  data.Get("built_in").(bool),
		TopLevel:                 data.Get("top_level").(bool),
		AuthenticationExecutions: authenticationExecutions,
	}

	return flow
}

// func mapToAuthenticationExecutionsFromData(data *schema.ResourceData) *[]keycloak.AuthenticationExecution {

// }
