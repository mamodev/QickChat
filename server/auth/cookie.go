package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateTokenCookies (userID string) (*fiber.Cookie, *fiber.Cookie, error) {

	accessTokenClaims, refreshTokenClaims := CreateClaims(userID)
	accessToken, err := CreateToken(accessTokenClaims)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := CreateToken(refreshTokenClaims)
	if err != nil {
		return nil, nil, err
	}

	accessTokenCookie := &fiber.Cookie{
		Name: "access_token",
		Value: accessToken,
		Expires: time.Unix(accessTokenClaims.Exp, 0),
		HTTPOnly: false,
		Secure: false,
		SameSite: "None",
	}

	refreshTokenCookie := &fiber.Cookie{
		Name: "refresh_token",
		Value: refreshToken,
		Expires: time.Unix(refreshTokenClaims.Exp, 0),
		HTTPOnly: false,
		Secure: false,
		SameSite: "None",
	}
		
	return accessTokenCookie, refreshTokenCookie, nil
}