package keycloak

import (
	"fmt"
)

type AuthenticationFlow struct {
	Id      string `json:"-"`
	RealmId string `json:"-"`
	Alias   string `json:"alias"`
}

func (keycloakClient *KeycloakClient) GetAuthenticationFlow(realmId, id string) (*AuthenticationFlow, error) {
	var flow AuthenticationFlow

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/authentication/flows/%s", realmId, id), &flow, nil)
	if err != nil {
		return nil, err
	}

	flow.RealmId = realmId

	return &flow, nil
}

func (keycloakClient *KeycloakClient) NewAuthenticationFlow(flow *AuthenticationFlow) error {
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/authentication/flows", flow.RealmId), flow)
	if err != nil {
		return err
	}

	flow.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) DeleteAuthenticationFlow(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/authentication/flows/%s", realmId, id), nil)
}

func (keycloakClient *KeycloakClient) UpdateAuthenticationFlow(flow *AuthenticationFlow) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/authentication/flows/%s", flow.RealmId, flow.Id), flow)
}
