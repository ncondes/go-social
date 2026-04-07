package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret   string
	audience string
	issuer   string
	expiry   time.Duration
}

func NewJWTAuthenticator(secret, audience, issuer string, expiry time.Duration) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:   secret,
		audience: audience,
		issuer:   issuer,
		expiry:   expiry,
	}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (ja *JWTAuthenticator) GenerateToken(userID int64) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			Issuer:    ja.issuer,
			Audience:  jwt.ClaimStrings{ja.audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ja.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ja.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ja *JWTAuthenticator) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(ja.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(ja.audience),
		jwt.WithIssuer(ja.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
