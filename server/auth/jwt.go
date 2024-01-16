package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var REFRESH_TOKEN_DURATION = time.Hour * 24 * 7
var ACCESS_TOKEN_DURATION = time.Minute * 24 * 60
var SECRET = []byte("secret")

type BaseClaims struct {
	Exp int64 `json:"exp"`
	Iat int64 `json:"iat"`
	Jti string `json:"jti"`
}

type AccessTokenClaims struct {
	BaseClaims
	UserID string `json:"user_id"`
}

type RefreshTokenClaims struct {
	BaseClaims
	Ref string `json:"ref"`
}

func (c *BaseClaims) isExpired() bool {
	return time.Unix(c.Exp, 0).Before(time.Now())
}

func (c *BaseClaims) Valid() error {
	if c.Exp == 0 {
		return jwt.NewValidationError("exp is 0", jwt.ValidationErrorClaimsInvalid)
	}

	if c.Iat == 0 {
		return jwt.NewValidationError("iat is 0", jwt.ValidationErrorClaimsInvalid)
	}

	if c.Jti == "" {
		return jwt.NewValidationError("jti is empty", jwt.ValidationErrorClaimsInvalid)
	}

	if c.Exp < c.Iat {
		return jwt.NewValidationError("malformed jwt", jwt.ValidationErrorClaimsInvalid)
	}

	return nil
}

func (c *AccessTokenClaims) Valid() error {

	if err := c.BaseClaims.Valid(); err != nil {
		return err
	}

	if c.UserID == "" {
		return jwt.NewValidationError("user_id is empty", jwt.ValidationErrorClaimsInvalid)
	}

	return nil
}

func (c *RefreshTokenClaims) Valid() error {
	if err := c.BaseClaims.Valid(); err != nil {
		return err
	}

	if c.Ref == "" {
		return jwt.NewValidationError("ref is empty", jwt.ValidationErrorClaimsInvalid)
	}

	return nil
}


func CreateClaims (userID string) (*AccessTokenClaims, *RefreshTokenClaims) {

	now := time.Now()
	accessTokenClaims := &AccessTokenClaims{
		BaseClaims: BaseClaims{
			Exp: now.Add(ACCESS_TOKEN_DURATION).Unix(),
			Iat: now.Unix(),
			Jti: uuid.New().String(),
		},
		UserID: userID,
	} 

	refreshTokenClaims := &RefreshTokenClaims{
		BaseClaims: BaseClaims{
			Exp: now.Add(REFRESH_TOKEN_DURATION).Unix(),
			Iat: now.Unix(),
			Jti: uuid.New().String(),
		},
		Ref: accessTokenClaims.Jti,
	}

	return accessTokenClaims, refreshTokenClaims
}


func CreateToken (claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET)
}

func ParseToken (tokenString string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})

	if err != nil {
		return err
	}

	if err := claims.Valid(); err != nil {
		return err
	}

	return nil
}
