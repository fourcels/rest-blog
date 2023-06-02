package controllers

import (
	"errors"
	"net/http"
	"os"
	"rest-blog/models"
	"time"

	"github.com/fourcels/rest"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var SigningKey = []byte(os.Getenv("SECRET"))

var JwtMiddleware = echojwt.WithConfig(echojwt.Config{
	SigningKey: SigningKey,
})

type jwtCustomClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

func Login() rest.Interactor {
	// Declare input port type.
	type input struct {
		Username string `json:"username" minLength:"3" default:"admin"`
		Password string `json:"password" minLength:"6" default:"a12345"`
	}

	// Declare output port type.
	type output struct {
		Token string `json:"token"`
	}
	// jwtCustomClaims are custom claims extending default ones.
	// See https://github.com/golang-jwt/jwt for more examples
	return rest.NewHandler(func(c echo.Context, in input, out *output) error {
		errNotFound := echo.NewHTTPError(http.StatusBadRequest, "Incorrect username or password")
		user := models.User{}
		if result := models.DB.Take(&user, "username = ?", in.Username); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errNotFound
			}
			return result.Error
		}
		if !CheckPasswordHash(in.Password, user.Password) {
			return errNotFound
		}

		// Set custom claims
		claims := &jwtCustomClaims{
			user.ID,
			user.Username,
			user.Admin,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString(SigningKey)
		if err != nil {
			return err
		}
		out.Token = t
		return nil
	})
}
