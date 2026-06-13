package main

import (
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	})
}

func InitServer(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	r.Any("/params", func(c *gin.Context) {
		params := map[string]string{}
		err := c.ShouldBindQuery(&params)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, params)
	})
	r.Any("/json", func(c *gin.Context) {
		body := map[string]string{}
		err := c.ShouldBindJSON(&body)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, body)
	})

	r.Any("/form", func(c *gin.Context) {
		form := map[string]string{}
		err := c.ShouldBind(&form)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, form)
	})

	r.POST("/files", func(c *gin.Context) {
		files, err := c.FormFile("files")
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		dst := filepath.Join("./files-delete/", filepath.Base(files.Filename))
		c.SaveUploadedFile(files, dst)
		c.JSON(200, gin.H{
			"message": "File uploaded successfully",
			"data": gin.H{
				"filename": files.Filename,
				"filepath": dst,
			},
		})
	})
}

func main() {
	server := gin.Default()
	server.Use(CorsMiddleware())
	InitServer(server)
	server.Run(":8080")
}
