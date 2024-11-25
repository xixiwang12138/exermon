package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const AuthUserKey = "auth_user"

func NewAuth[T any](pk string, expire time.Duration) UserAuthService[T] {
	return UserAuthService[T]{
		pk:     pk,
		expire: expire,
	}
}

type UserAuthService[T any] struct {
	jwt.RegisteredClaims
	pk string

	expire time.Duration
}

func (s UserAuthService[T]) GenerateToken(user T) (string, error) {
	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims[T]{
		UserPayload: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expire)),
		},
	}).SignedString([]byte(s.pk))
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
