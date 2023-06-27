package controllers

import (
	"errors"
	"net/http"
	"rest-blog/models"
	"strconv"

	"github.com/fourcels/paginate"
	"github.com/fourcels/rest"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PathID struct {
	ID uint `path:"id"`
}

type Message struct {
	Message string `json:"message"`
}

func CreatePost() rest.Interactor {
	type input struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	return rest.NewHandler(func(c echo.Context, in input, out *models.Post) error {
		post := models.Post{
			Title:       in.Title,
			Content:     in.Content,
			Tags:        in.Tags,
			CreatedByID: getUserID(c),
		}
		if result := models.DB.Create(&post); result.Error != nil {
			return result.Error
		}
		*out = post
		return nil
	})
}
func CreateComment() rest.Interactor {
	type input struct {
		PathID
		Content string `json:"content"`
	}

	return rest.NewHandler(func(c echo.Context, in input, out *models.Comment) error {
		if _, err := getPostByID(in.ID); err != nil {
			return err
		}

		comment := models.Comment{
			Content:     in.Content,
			PostID:      in.ID,
			CreatedByID: getUserID(c),
		}
		if result := models.DB.Create(&comment); result.Error != nil {
			return result.Error
		}
		*out = comment
		return nil
	})
}
func GetComments() rest.Interactor {
	type input struct {
		PathID
		paginate.Pagination
	}

	return rest.NewHandler(func(c echo.Context, in input, out *[]models.Comment) error {
		err := setupPaginate(c, in.Pagination, out, func(db *gorm.DB) *gorm.DB {
			return db.Where("post_id = ?", in.ID)
		})
		if err != nil {
			return err
		}
		return nil
	})
}
func GetPosts() rest.Interactor {
	return rest.NewHandler(func(c echo.Context, in paginate.Pagination, out *[]models.Post) error {
		err := setupPaginate(c, in, out)
		if err != nil {
			return err
		}
		return nil
	})
}
func GetPost() rest.Interactor {
	return rest.NewHandler(func(c echo.Context, in PathID, out *models.Post) error {
		var post models.Post
		if result := models.DB.Take(&post, in.ID); result.Error != nil {
			return result.Error
		}
		*out = post
		return nil
	})
}
func UpdatePost() rest.Interactor {
	type input struct {
		PathID
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	return rest.NewHandler(func(c echo.Context, in input, out *models.Post) error {
		post, err := getPostByID(in.ID)
		if err != nil {
			return err
		}

		if result := models.DB.Model(post).Updates(&models.Post{Title: in.Title, Content: in.Content}); result.Error != nil {
			return result.Error
		}
		*out = *post
		return nil
	})
}
func DeletePost() rest.Interactor {
	return rest.NewHandler(func(c echo.Context, in PathID, out *Message) error {
		post, err := getPostByID(in.ID)
		if err != nil {
			return err
		}
		if result := models.DB.Delete(post); result.Error != nil {
			return result.Error
		}
		out.Message = "OK"
		return nil
	})
}

func getPostByID(id uint) (*models.Post, error) {
	var post models.Post
	if result := models.DB.Take(&post, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "post not found")
		}
		return nil, result.Error
	}
	return &post, nil
}

func getUserID(c echo.Context) uint {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))
	return id
}

func setupPaginate[T any](c echo.Context, p paginate.Pagination, out *[]T, query ...func(db *gorm.DB) *gorm.DB) error {
	count, err := paginate.Paginate(models.DB, p, out, query...)
	if err != nil {
		return err
	}
	c.Response().Header().Set("X-Total", strconv.FormatInt(count, 10))
	return nil
}
