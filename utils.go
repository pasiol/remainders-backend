package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func readSSLCertificates() []tls.Certificate {

	cert, err := tls.LoadX509KeyPair(os.Getenv("SSL_PUBLIC"), os.Getenv("SSL_PRIVATE"))
	if err != nil {
		log.Fatalf("reading ssl certificates failed: %s", err)
	}
	return []tls.Certificate{cert}
}

func getTSLConfig(cert string) *tls.Config {

	pool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(cert)
	if err != nil {
		log.Fatalf("reading certificate file failed: %s", err)
	}
	pool.AppendCertsFromPEM(ca)

	tslConfig := &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: true,
	}

	return tslConfig
}

func hashAndSalt(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
