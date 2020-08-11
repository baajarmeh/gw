package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/admin"
	"github.com/oceanho/gw/contrib/apps/tester"
)

func main() {
	var apps []gw.App
	apps = append(apps, tester.New(), tester.NewAppRestOnly(), admin.App{})

	s1 := gw.NewServerWithNameAddr("s1", ":18080")
	s1.Register(apps...)

	s2 := gw.NewServerWithNameAddr("s2", ":18081")
	s2.Register(apps...)

	s3 := gw.NewServerWithNameAddr("s3", ":18082")
	s3.Register(apps...)
	gw.Run(s1, s2, s3)
}
