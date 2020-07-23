package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()
	g := r.Group("//v1")
	g.GET("/version", func(c *gin.Context) {
		c.OK( gin.H{
			"payload": "welcome.",
		})
	})
	r.Run(":9000")
}
