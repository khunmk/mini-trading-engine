package main

import (
	"github.com/gin-gonic/gin"
	"github.com/khunmk/mini-trading-engine/engine"
	"github.com/khunmk/mini-trading-engine/model"
)

func main() {

	router := gin.Default()

	book := engine.NewOrderBook()

	router.POST("/", func(c *gin.Context) {
		order := &model.Order{}
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(500, err.Error())
			return
		}

		if err := order.Create(model.DB); err != nil {
			c.JSON(500, err.Error())
			return
		}

		trades := book.Process(order)

		for _, trade := range trades {
			if err := trade.Create(model.DB); err != nil {
				c.JSON(500, err.Error())
				return
			}
		}
	})

	router.Run(":9000")
}
