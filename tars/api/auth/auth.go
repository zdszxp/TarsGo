package auth

import (
)

type AuthData interface{

}

type Authenticator interface{
	CheckAuth(data AuthData) bool
}

type localAuthenticator struct {
	Provider SecretProvider
}

func (la *localAuthenticator) CheckAuth(data AuthData) bool{
	return la.Provider(data) 
}

// SecretProvider is the SecretProvider function
type SecretProvider func(data AuthData) bool

func NewAuthenticator(secrets SecretProvider) Authenticator {
	return &localAuthenticator{
		Provider:secrets,
	}
}
