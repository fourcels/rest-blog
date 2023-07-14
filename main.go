package main

import (
	"log"
	"rest-blog/controllers"

	"github.com/fourcels/rest"
)

func main() {
	s := rest.NewService()
	s.OpenAPI.Info.WithTitle("Rest Blog")
	s.WithHttpBearerSecurity("bearerAuth")

	auth := s.Group("/auth", rest.WithTags("Auth"))
	auth.POST("/login", controllers.Login())

	user := s.Group("/users", rest.WithTags("Users"))
	user.POST("", controllers.CreateUser())

	post := s.Group("/posts", rest.WithTags("Posts"), rest.WithSecurity("bearerAuth"))
	post.Use(controllers.JwtMiddleware)
	post.POST("", controllers.CreatePost())
	post.GET("", controllers.GetPosts())
	post.GET("/:id", controllers.GetPost())
	post.PATCH("/:id", controllers.UpdatePost())
	post.DELETE("/:id", controllers.DeletePost())
	post.POST("/:id/comments", controllers.CreateComment())
	post.GET("/:id/comments", controllers.GetComments())

	comment := s.Group("/comments", rest.WithTags("Comments"), rest.WithSecurity("bearerAuth"))
	comment.Use(controllers.JwtMiddleware)
	comment.PATCH("/:id", controllers.UpdateComment())
	comment.DELETE("/:id", controllers.DeleteComment())

	// Swagger UI endpoint at /docs.
	s.Docs("/docs", map[string]any{
		"persistAuthorization": true,
	},
	)

	// Start server.
	log.Println("http://localhost:1323/docs")
	s.Start(":1323")
}
