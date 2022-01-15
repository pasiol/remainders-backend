package main

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"os"
	"time"
)

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func (a *App) getHealthz(c echo.Context) error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err := a.Db.Client().Ping(ctx, readpref.Primary())
	if err != nil {
		return c.String(http.StatusInternalServerError, "non-operational")
		return err
	}
	return c.String(http.StatusOK, "operational")
}

func (a *App) postLogin(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(u); err != nil {
		return err
	}
	if a.Login(*u) {
		claims := &jwtCustomClaims{
			u.Username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 6).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	} else {
		return echo.ErrUnauthorized
	}
}

func (a *App) getLatest(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
