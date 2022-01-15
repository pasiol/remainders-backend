package main

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type remainder struct {
	To        string    `bson:"to" json:"to"`
	Title     string    `bson:"title" json:"title"`
	Message   string    `bson:"message" json:"message"`
	Type      string    `bson:"type" json:"type"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type (
	User struct {
		Username string `bson:"username" json:"username" validate:"required"`
		Password string `bson:"password" json:"password" validate:"required"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

type Exception struct {
	Message string `json:"message"`
}

type Payload struct {
	Message string `json:"message"`
}
