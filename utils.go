package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"

	jwt "github.com/dgrijalva/jwt-go"

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

func getToken(name string) (string, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error)) {
		return signingKey, nil
	})
	if err !=nil {
		return nil, err
	}
	return token.Claims, err
}