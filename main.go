package main

import (
	"embed"
	"flag"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var frontend embed.FS

func main() {
	var port = 7256
	flag.IntVar(&port, "port", 7256, "指定代理服务器端口")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 允许所有 CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 代理接口
	r.Any("/proxy", func(c *gin.Context) {
		targetURL := c.Query("url")
		if targetURL == "" {
			c.JSON(400, gin.H{"error": "缺少 url 参数"})
			return
		}

		// 构建请求
		req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 复制请求头（跳过 Host，保留其他头）
		for k, v := range c.Request.Header {
			if strings.ToLower(k) != "host" {
				req.Header[k] = v
			}
		}

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": "请求失败: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		// 复制响应头
		for k, v := range resp.Header {
			c.Header(k, strings.Join(v, ", "))
		}

		// 复制响应状态码和 Body
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	})

	r.GET("/", func(c *gin.Context) {
		f, err := frontend.ReadFile("dist/index.html")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "text/html", f)
	})
	r.GET("/assets/*filepath", func(c *gin.Context) {
		assertFile := path.Join("dist", "assets", c.Param("filepath"))
		f, err := frontend.ReadFile(assertFile)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctype := mime.TypeByExtension(path.Ext(assertFile))
		if ctype == "" {
			ctype = "application/octet-stream"
		}
		c.Data(200, ctype, f)
	})

	log.Println("🚀 代理服务器运行在 http://localhost:" + strconv.Itoa(port))
	r.Run(":" + strconv.Itoa(port))
}
