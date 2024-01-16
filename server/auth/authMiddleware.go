package auth

import (
	"github.com/gofiber/fiber/v2"
)

func checkAuth(c *fiber.Ctx) (string, error) {
	accessTokenCookie := c.Cookies("access_token")
	if accessTokenCookie == "" {
		return "", fiber.ErrUnauthorized
	}

	accessTokenClaims := &AccessTokenClaims{}
	err := ParseToken(accessTokenCookie, accessTokenClaims)

	if err != nil {
		return "", fiber.ErrUnauthorized
	}

	if !accessTokenClaims.isExpired() {
		return accessTokenClaims.UserID, nil
	}

	refreshTokenCookie := c.Cookies("refresh_token")
	if refreshTokenCookie == "" {
		return "", fiber.ErrUnauthorized 
	}

	refreshTokenClaims := &RefreshTokenClaims{}
	err = ParseToken(refreshTokenCookie, refreshTokenClaims)
	
	if err != nil || refreshTokenClaims.isExpired() || accessTokenClaims.Jti != refreshTokenClaims.Ref {
		return "", fiber.ErrUnauthorized
	}

	newAccessTokenCookie, newRefreshTokenCookie, err := CreateTokenCookies(accessTokenClaims.UserID)
	if err != nil {
		return "", fiber.ErrInternalServerError
	}

	c.Cookie(newAccessTokenCookie)
	c.Cookie(newRefreshTokenCookie)
	return accessTokenClaims.UserID, nil
}

func AuthMiddleware(c *fiber.Ctx) error {

	id, authErr := checkAuth(c)

	if authErr != nil {
		c.ClearCookie("access_token")
		c.ClearCookie("refresh_token")
		return authErr
	}

	c.Locals("userID", id)

	return c.Next()
}	
