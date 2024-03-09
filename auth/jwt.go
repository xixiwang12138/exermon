package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

func NewAuth[T any](pk string) UserAuthService[T] {
	return UserAuthService[T]{
		pk: pk,
	}
}

type UserAuthService[T any] struct {
	jwt.RegisteredClaims
	pk string
}

func (s UserAuthService[T]) GenerateToken(claim UserClaims[T]) (string, error) {
	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(s.pk))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (s UserAuthService[T]) ParseToken(token string) (*T, error) {
	userClaims := &UserClaims[T]{}
	t, err := jwt.ParseWithClaims(token, userClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.pk), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return &userClaims.UserPayload, nil

}

type UserClaims[T any] struct {
	UserPayload T
	jwt.RegisteredClaims
}
