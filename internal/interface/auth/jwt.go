package auth

import (
    "github.com/dgrijalva/jwt-go"

    "crypto/rsa"
    "errors"
    "fmt"
    "time"
)

// todo: key rotation
type JwtConfig struct {
    publicKey  *rsa.PublicKey
    privateKey *rsa.PrivateKey

    expire time.Duration
}

type Claims struct {
    Id string
    jwt.StandardClaims
}

func NewJwtConfig(privateBytes, publicBytes []byte, keyExpiration time.Duration) (*JwtConfig, error) {
    privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
    if err != nil {
        return nil, err
    }
    publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
    if err != nil {
        return nil, err
    }
    return &JwtConfig{
        publicKey:  publicKey,
        privateKey: privateKey,
        expire:     keyExpiration,
    }, nil
}

func (j JwtConfig) IssueToken(userId string) (string, error) {
    claims := Claims{
        Id: userId,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(j.expire).Unix(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return token.SignedString(j.privateKey)
}

func (j JwtConfig) UserIdByToken(tokenString string) (string, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected token signing method")
        }
        return j.publicKey, nil
    })
    if err != nil {
        return "", err
    }
    claims, ok := token.Claims.(*Claims)
    if !ok {
        return "", errors.New("invalid token claims")
    }
    return claims.Id, nil
}

