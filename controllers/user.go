package controllers

import (
	"errors"
	"net/http"
	"rest-blog/models"

	"github.com/fourcels/rest"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser() rest.Interactor {
	type input struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	return rest.NewHandler(func(c echo.Context, in input, out *models.User) error {

		if exist, err := checkUsername(in.Username); err != nil {
			return err
		} else if exist {
			return echo.NewHTTPError(http.StatusBadRequest, "username already registered")
		}

		user := models.User{
			Username: in.Username,
			Admin:    in.Username == "admin",
			Password: HashPassword(in.Password),
		}
		if result := models.DB.Create(&user); result.Error != nil {
			return result.Error
		}
		*out = user
		return nil
	})
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func checkUsername(username string) (bool, error) {
	var user models.User
	if result := models.DB.Take(&user, "username = ?", username); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
