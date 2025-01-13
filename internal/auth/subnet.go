package auth

import (
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
)

const realIPHeaderName = "X-Real-IP"

// NewTrustedSubnetMiddleware provides middleware that checks if user header IP in trusted subnet
func NewTrustedSubnetMiddleware(trustedSubnet *net.IPNet) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if trustedSubnet != nil {
				if ip := net.ParseIP(c.Request().Header.Get(realIPHeaderName)); ip != nil && trustedSubnet.Contains(ip) {
					return next(c)
				}
			}
			return c.NoContent(http.StatusForbidden)
		}
	}
}
