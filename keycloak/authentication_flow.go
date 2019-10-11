package keycloak

import (
	"fmt"
)

type AuthenticationFlow struct {
	Id         string `json:"id,omitempty"`
	RealmId    string `json:"-"`
	Alias      string `json:"alias"`
	ProviderId string `json:"providerId"`
}

func (keycloakClient *KeycloakClient) GetAuthenticationFlow(realmId, id string) (*AuthenticationFlow, error) {
	var flow AuthenticationFlow

	url := fmt.Sprintf("/realms/%s/authentication/flows/%s", realmId, id)
	fmt.Println(url)
	err := keycloakClient.get(url, &flow, nil)
	if err != nil {
		return nil, err
	}

	flow.RealmId = realmId
	fmt.Println(flow.Id, flow.RealmId, flow.Alias)

	return &flow, nil
}

func (keycloakClient *KeycloakClient) NewAuthenticationFlow(flow *AuthenticationFlow) error {
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/authentication/flows", flow.RealmId), flow)
	if err != nil {
		return err
	}
	fmt.Println(location)

	flow.Id = getIdFromLocationHeader(location)
	fmt.Println(flow.Id)

	return nil
}

func (keycloakClient *KeycloakClient) DeleteAuthenticationFlow(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/authentication/flows/%s", realmId, id), nil)
}

func (keycloakClient *KeycloakClient) UpdateAuthenticationFlow(flow *AuthenticationFlow) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/authentication/flows/%s", flow.RealmId, flow.Id), flow)
}
