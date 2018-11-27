package service

import (
	"github.com/gin-gonic/gin"
)

type GinContext struct {
	*gin.Context
}

func (c *GinContext) BindWithoutId(obj interface{}) error {
	err := c.Bind(&obj)
	if err != nil {
		return err
	}
	//if obj.ID != 0 {
	//return error.New("Bind has ID")
	//}
	return nil
}

func (c *GinContext) APIJSON(data string) {
	c.JSON(200, gin.H{"data": data})
}
