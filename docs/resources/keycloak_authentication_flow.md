# keycloak_authentication_flow

Allows for creating and managing authentication flows.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_authentication_flow" "browser-copy-flow" {
	alias    = "browserCopyFlow"
	realm_id = "${keycloak_realm.id}"
	description = "browser based authentication"
    provider_id = "basic-flow"
}
```

## Argument Reference
- `alias`
- `realm_id`
- `provider_id` - (Optional) boolean
- `description` - (Optional) string
