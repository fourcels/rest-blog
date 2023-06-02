package controllers

import (
	"errors"
	"net/http"
	"rest-blog/models"

	"github.com/fourcels/rest"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UpdateComment() rest.Interactor {
	type input struct {
		PathID
		Content string `json:"content"`
	}

	return rest.NewHandler(func(c echo.Context, in input, out *models.Comment) error {
		comment, err := getCommentByID(in.ID)
		if err != nil {
			return err
		}

		if result := models.DB.Model(comment).Updates(&models.Comment{Content: in.Content}); result.Error != nil {
			return result.Error
		}
		*out = *comment
		return nil
	})
}
func DeleteComment() rest.Interactor {

	return rest.NewHandler(func(c echo.Context, in PathID, out *Message) error {
		comment, err := getCommentByID(in.ID)
		if err != nil {
			return err
		}
		if result := models.DB.Delete(comment); result.Error != nil {
			return result.Error
		}
		out.Message = "OK"
		return nil
	})
}

func getCommentByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	if result := models.DB.Take(&comment, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "comment not found")
		}
		return nil, result.Error
	}
	return &comment, nil
}
