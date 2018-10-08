package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const listen = ":3826" // ".tao"

const root = "web/html"
const examples = "web/examples"

func main() {
	router := gin.Default()
	router.GET("/html/", func(c *gin.Context) {
		c.File(root + "/index.html")
	})
	router.GET("/html/:path", func(c *gin.Context) {
		path := "/" + c.Param("path")
		path = filepath.Clean(path)
		path = filepath.Join(root, path)
		c.File(path)
	})
	router.POST("/v1/execute", func(c *gin.Context) {
		var data struct {
			Source string
		}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "./bin/tao")
		cmd.Stdin = strings.NewReader(data.Source)
		output, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}
		c.String(http.StatusOK, "%s", string(output))
	})
	router.GET("/v1/examples", func(c *gin.Context) {
		files, err := ioutil.ReadDir(examples)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}
		names := make([]string, 0, len(files))
		for _, file := range files {
			ext := filepath.Ext(file.Name())
			if ext == ".tao" {
				names = append(names, file.Name())
			}
		}
		c.JSON(http.StatusOK, names)
	})
	router.GET("/v1/examples/:path", func(c *gin.Context) {
		path := "/" + c.Param("path")
		path = filepath.Clean(path)
		path = filepath.Join(examples, path)
		source, err := ioutil.ReadFile(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}
		c.String(http.StatusOK, "%s", string(source))
	})
	router.Run(listen)
}
