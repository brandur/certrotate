package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/joeshaw/envdecode"
	"github.com/xenolf/lego/acme"
)

type Conf struct {
	// Currently only staging supports the dns-01 challenge.
	AcmeURI string `env:"ACME_URI,default=https://acme-staging.api.letsencrypt.org/directory"`
	Domain  string `env:"DOMAIN"`
	//AcmeURI  string `env:"ACME_URI,default=https://acme-v01.api.letsencrypt.org/directory"`
	Email string `env:"EMAIL,required"`
}

var (
	conf Conf
)

func readConf() {
	err := envdecode.Decode(&conf)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	readConf()

	// Create a user. New accounts need an email and private key to start.
	log.Printf("Generating private RSA key for user '%v'", conf.Email)
	const rsaKeySize = 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		log.Fatal(err)
	}
	myUser := MyUser{
		Email: conf.Email,
		key:   privateKey,
	}

	// A client facilitates communication with the CA server. Bind to PORT to
	// facilitate challenge requests.
	log.Printf("Calling ACME server: %v", conf.AcmeURI)

	// It's a little unfortunate that the client initializer requires a port
	// right now even though we won't be using one (we'll be using the DNS
	// challenge). Hopefully this interface will change in the near future.
	client, err := acme.NewClient(conf.AcmeURI, &myUser, rsaKeySize, "98234")
	if err != nil {
		log.Fatal(err)
	}

	// New users will need to register; be sure to save it
	reg, err := client.Register()
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg

	// The client has a URL to the current Let's Encrypt Subscriber
	// Agreement. The user will need to agree to it.
	err = client.AgreeToTOS()
	if err != nil {
		log.Fatal(err)
	}

	// The acme library takes care of completing the challenges to obtain the certificate(s).
	// Of course, the hostnames must resolve to this machine or it will fail.
	bundle := false
	certificates, errs := client.ObtainCertificates([]string{conf.Domain}, bundle)
	for _, err := range errs {
		if err != nil {
			log.Fatal(err)
		}
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. This is where you should save them to files!
	fmt.Printf("%#v\n", certificates)
}
