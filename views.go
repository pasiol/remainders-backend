package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func (a *App) getHealthz(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	err := a.Db.Client().Ping(ctx, readpref.Primary())
	defer cancel()
	if err != nil {
		return c.String(http.StatusInternalServerError, "non-operational")
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
	expirationTime := time.Hour * 2
	if a.Debug {
		expirationTime = time.Minute * 10
	}
	if u.Login(a.Db) {

		claims := &jwtCustomClaims{
			u.Username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(expirationTime).Unix(),
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
	remainders, err := find(a.Db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, remainders)
}

func (a *App) getSearch(c echo.Context) error {
	filter := c.Param("filter")
	remainders, err := search(filter, a.Db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, remainders)
}
