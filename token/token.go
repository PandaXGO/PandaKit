package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId         int64
	TenantId       int64
	OrganizationId int64 //组织Id
	UserName       string
	RoleId         int64
	RoleKey        string
	DeptId         int64
	PostId         int64
	jwt.RegisteredClaims
}

type JWT struct {
	SignedKeyID  string
	SignedKey    []byte
	SignedMethod jwt.SigningMethod
}

var (
	TokenExpired          = errors.New("token is expired")
	TokenNotValidYet      = errors.New("token not active yet")
	TokenMalformed        = errors.New("that's not even a token")
	TokenInvalid          = errors.New("couldn't handle this token")
	UnsupportedSignMethod = errors.New("unsupported sign method")
)

func NewJWT(kid string, key []byte, method jwt.SigningMethod) *JWT {
	return &JWT{
		SignedKeyID:  kid,
		SignedKey:    key,
		SignedMethod: method,
	}
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	key, err := j.getKey()
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("parse error")
		}
		key, err := j.getKey()
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		// Check if the error is due to specific JWT issues
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, TokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			// Token is expired
			return nil, TokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, TokenNotValidYet
		} else {
			return nil, TokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		key, err := j.getKey()
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(60 * 60 * 24 * 7)))
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

func (j *JWT) getKey() (interface{}, error) {
	if j.isEs() {
		v, err := jwt.ParseECPrivateKeyFromPEM(j.SignedKey)
		if err != nil {
			return nil, err
		}
		return v, nil
	} else if j.isRsOrPS() {
		v, err := jwt.ParseRSAPrivateKeyFromPEM(j.SignedKey)
		if err != nil {
			return nil, err
		}
		return v, nil
	} else if j.isHs() {
		return j.SignedKey, nil
	} else {
		return nil, UnsupportedSignMethod
	}
}

func (a *JWT) isEs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "ES")
}

func (a *JWT) isRsOrPS() bool {
	isRs := strings.HasPrefix(a.SignedMethod.Alg(), "RS")
	isPs := strings.HasPrefix(a.SignedMethod.Alg(), "PS")
	return isRs || isPs
}

func (a *JWT) isHs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "HS")
}
