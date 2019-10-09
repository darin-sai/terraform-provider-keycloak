package keycloak

// import (
// 	"fmt"
// )

type AuthenticationExecution struct {
	Authenticator              string `json:"authenticator"`
	AuthenticatorConfig        string `json:"authenticator_config"`
	Requirement                string `json:"requirement"`
	AuthenticatorFlow          bool   `json:"authenticator_flow"`
	Priority                   int    `json:"priority"`
	UserSetupAllowed           bool   `json:"user_setup_allowed"`
}

type Flow struct {
    Id		                 string	                   `json:"-"`
	RealmId                  string                    `json:"-"`
	Alias                    string                    `json:"alias"`
	Description              string                    `json:"description"`
	AuthenticationExecutions []AuthenticationExecution `json;"authentication_executions"`
}

