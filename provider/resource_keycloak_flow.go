package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakFlow() *schema.Resource {

	return &schema.Resource{
		Create: resourceKeycloakFlowCreate,
		Read:   resourceKeycloakFlowRead,
		Delete: resourceKeycloakFlowDelete,
		Update: resourceKeycloakFlowUpdate,
		// Importer: &schema.ResourceImporter{
		// 	// This resource can be imported using {{realm}}/{{alias}}. The required action aliases are displayed in the server info or GET realms/{{realm}}/authentication/required-actions
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authentication_executions": {
				Type:     map[string]*schema.Schema{
					"authenticator" {
						Type:     schema.TypeString,
						Optional: true,
					},
					"authenticator_config" {
						Type:     schema.TypeString,
						Optional: true,
					},
					"requirement" {
						Type:     schema.TypeString,
						Optional: true,
					},
					"authenticator_flow" {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"priority" {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"user_setup_allowed" {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
				Required: true,
			},
			"default_action": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}
