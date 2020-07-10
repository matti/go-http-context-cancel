package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type params struct {
	Amount int `uri:"amount" binding:"required"`
}

func main() {
	r := gin.Default()
	r.GET("/sleep/:amount", func(c *gin.Context) {
		var params params

		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		ctx, cancel := context.WithCancel(c)

		go func(ctx context.Context, cancel context.CancelFunc, amount int) {
			defer cancel()
			t := time.NewTicker(1 * time.Second)
			defer t.Stop()

			log.Println("sleep", amount)

			for i := 1; i < amount; i++ {
				select {
				case <-t.C:
					log.Println(i)
				case <-ctx.Done():
					log.Println("func client dropped")
					return
				}
			}
		}(ctx, cancel, params.Amount)

		select {
		case <-ctx.Done():
			log.Println("work complete")
		case <-c.Request.Context().Done():
			log.Println("main client dropped")
			cancel()
			return
		}

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("127.0.0.1:8080")
}
