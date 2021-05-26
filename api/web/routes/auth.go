package routes

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ithaka/labs-pep/api/web/models"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func Login(c echo.Context) error {
	var l models.Login
	if err := c.Bind(&l); err != nil {
		c.Logger().Error("Unable to bind request: ", err)
		return err
	}

	if l.Password == viper.GetString("auth.password") {
		claims := &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Subject:   "admin",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString([]byte(viper.GetString("auth.signing_key")))
		if err != nil {
			c.Logger().Error("Unable to generate token: ", err)
			return c.JSON(http.StatusInternalServerError, models.Response{Code: http.StatusForbidden, Message: "unable to generate token"})
		}
		return c.JSON(http.StatusOK, models.Response{Code: http.StatusOK, Data: models.LoginResponse{Token: ss}})
	}
	return c.JSON(http.StatusForbidden, models.Response{Code: http.StatusForbidden, Message: "bad password"})
}
