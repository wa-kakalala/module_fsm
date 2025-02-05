package main

import (
	"encoding/json"
	"fmt"
	"io"
	"module_fsm/dot"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/static/css/", "static/css/")
	router.Static("/static/js/", "static/js/")
	router.LoadHTMLGlob("templates/*")

	// 访问地址，处理请求 Request Response
	router.GET("/module_fsm", func(c *gin.Context) {
		c.HTML(http.StatusOK, "module_fsm.html", gin.H{
			"title": "Main website",
		})
	})

	router.POST("/fsm", func(c *gin.Context) {
		// 读取请求体
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}
		var dot_json dot.DotJson_s
		err = json.Unmarshal([]byte(string(body)), &dot_json)
		if err != nil {
			fmt.Println("parse json error")
			return
		}
//		fmt.Println(dot_json.Dot)
		pngfile, err := dot_json.GenDotFile()

		// 打开图片文件
		file, err := os.Open(pngfile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}

		fileInfo, err := file.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
			fmt.Println("Failed to get file info")
			return
		}
		fileContent := make([]byte, fileInfo.Size())

		// 读取文件内容到缓冲区
		_, err = file.Read(fileContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
			fmt.Println("Failed to read file content")
			return
		}

		c.Header("Content-Type", "image/png")

		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

		c.Data(http.StatusOK, "image/png", fileContent)
		fmt.Println("send ok")
		file.Close()
		os.Remove(pngfile)

	})
	// 服务器端口
	router.Run("0.0.0.0:9090")
}
