package main

import (
	"crypto/rsa"

	"github.com/xenolf/lego/acme"
)

// You'll need a user or account type that implements acme.User
type MyUser struct {
	Email        string
	Registration *acme.RegistrationResource
	key          *rsa.PrivateKey
}

func (u MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *acme.RegistrationResource {
	return u.Registration
}
func (u MyUser) GetPrivateKey() *rsa.PrivateKey {
	return u.key
}
