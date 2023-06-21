package jwtUser

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

var jwtKey = []byte("supersecretkeyvdjwbdhwjdbiwuhdqwihdiq")

type JWTClaim struct {
	UserPhone string `json:"userPhone"`
	UserEmail string `json:"userEmail"`
	UserId    string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(phone string, email string, id string) (tokenString string, err error) {
	expTime := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	tokenclaim := &JWTClaim{
		UserPhone: phone,
		UserEmail: email,
		UserId:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expTime,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "Token issue	",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenclaim)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func StripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

var UserJWTData *JWTClaim

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		log.Fatal("error occurred during parsing the token", err.Error())
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		log.Fatal("error occurred during parsing the token", err.Error())
		return
	}

	if token.Valid {
		UserJWTData = claims
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		fmt.Println("That's not even a token")
		return
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		err = errors.New("token is either expired or not active yet")
		return
	} else {
		fmt.Println("Couldn't handle this token:", err)
		return
	}

	return
}

func Auth() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Please provide auth token", http.StatusUnauthorized)
				return
			}
			tokenString, terror := StripBearerPrefixFromTokenString(tokenString)
			if terror != nil {
				http.Error(w, "Error while parsing the authorization token", http.StatusGone)
				return
			}
			fmt.Println("Got Token:", tokenString)

			err := ValidateToken(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusGone)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
