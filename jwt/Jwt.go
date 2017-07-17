package main

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	//	"os"
	"time"
)

func main() {

	secretKey := "secret-key"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "internal", // The important bit that will give you permission to
		// access the endpoints

		"iss": "email daemon or something", // Word / phrase describing who generated the token
		// (issuer). I will match this against a list of valid
		// token issuers so let me know what you put here

		"exp": time.Now().Unix() + 36000, // Expiration time

		"iat": time.Now(), // Issued at time, optional

		"nbf": time.Now().Unix()}) // "not before", meaning the token is invalid before the
	// specified time... optional

	tokenString, err := token.SignedString([]byte(secretKey))

	fmt.Printf("JWT: %s, %v", tokenString, err)
}
