package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xiwen1/mini-tiktok-comment/Comment/controller"
	"github.com/xiwen1/mini-tiktok-comment/Comment/service"
	"log"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)
}

func main() {
	err := service.InitComment()
	if err != nil {
		log.Println("init fail ")
		return
	}
	r := gin.Default()
	initRouter(r)
	r.Run("0.0.0.0:50051")

}
