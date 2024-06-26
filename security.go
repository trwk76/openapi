package openapi

import (
	"fmt"
	"go/token"

	"github.com/trwk76/openapi/spec"
)

func (a *API) APIKeySecurity(key string, in spec.SecurityIn, name string, setup func(s *spec.SecurityScheme)) {
	item := spec.SecurityScheme{Type: spec.SecurityAPIKey, In: in, Name: name}

	if setup != nil {
		setup(&item)
	}

	a.AddSecurity(key, item)
}

func (a *API) HTTPBasicSecurity(key string, setup func(s *spec.SecurityScheme)) {
	item := spec.SecurityScheme{Type: spec.SecurityHTTP, Scheme: "Basic"}

	if setup != nil {
		setup(&item)
	}

	a.AddSecurity(key, item)
}

func (a *API) OAuth2Security(key string, flows spec.OAuthFlows, setup func(s *spec.SecurityScheme)) {
	item := spec.SecurityScheme{Type: spec.SecurityOAuth2, Flows: &flows}

	if setup != nil {
		setup(&item)
	}

	a.AddSecurity(key, item)
}

func (a *API) OpenIDConnectSecurity(key string, url string, setup func(s *spec.SecurityScheme)) {
	item := spec.SecurityScheme{Type: spec.SecurityOpenIDConnect, OpenIDConnectURL: url}

	if setup != nil {
		setup(&item)
	}

	a.AddSecurity(key, item)
}

func (a *API) AddSecurity(key string, s spec.SecurityScheme) {
	if !token.IsIdentifier(key) {
		panic(fmt.Errorf("'%s' is not a valid key", key))
	}

	if _, ok := a.s.Components.SecuritySchemes[key]; ok {
		panic(fmt.Errorf("key '%s' is already used by another security scheme", key))
	}

	a.s.Components.SecuritySchemes[key] = spec.SecuritySchemeOrRef{Item: s}
}
