package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/uap/dto"
)

func GetUser(c *gw.Context) {
	dto := &dto.UserDto{}
	if c.Bind(dto) != nil {
		return
	}
	//user := &db
	//db := c.Store.GetDbStoreByName("user-primary")
	//biz.CreateUser(db)
}
func CreateUser(c *gw.Context) {
	c.OK(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyUser(c *gw.Context) {
	c.OK(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeleteUser(c *gw.Context) {
	c.OK(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryUser(c *gw.Context) {
	c.OK(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}
