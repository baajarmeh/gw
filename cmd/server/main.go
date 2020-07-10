package main

import "github.com/oceanho/gw/contrib/app"

func main() {
	//router := gin.New()
	//router.GET("/", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	//
	//// By default it serves on :8080 unless a
	//// PORT environment variable was defined.
	//router.Run(":8000")
	//// router.Run(":3000") for a hard coded port
	server := app.New()
	server.Serve()
}
