package config

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/SawitProRecruitment/UserService/model"
	"github.com/dgrijalva/jwt-go"
)

type Config struct {
	JWT JWT
}

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j JWT) Create(ttl time.Duration, user model.User) (string, error) {
	// Create a new JWT token with RS256 signing method.
	token := jwt.New(jwt.SigningMethodRS256)

	// Set claims (payload) for the token.
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = fmt.Sprint(user.UserID)
	claims["exp"] = ttl

	// Load your RSA private key for signing the token.
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		log.Println("Error loading private key:", err)
		return "", err
	}

	// Sign the token with the private key.
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return tokenString, nil
}

func (j JWT) Validate(token string) (model.User, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return model.User{}, fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return model.User{}, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return model.User{}, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return model.User{}, fmt.Errorf("validate: invalid")
	}

	fmt.Println("claims ==== ", claims["user_id"].(string))

	userID, err := strconv.Atoi(claims["user_id"].(string))
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		UserID: int32(userID),
	}, nil
}
