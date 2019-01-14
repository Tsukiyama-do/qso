package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "./controllers"
)

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("views/*")

    router.GET("/", func(c *gin.Context) {
        controller := task.NewTask()
        tasks := controller.GetAll()

        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "title": "Welcome to JJ1POW page!",
            "tasks": tasks,    //　追記 テンプレートにタスクを渡す
        })
    })

    router.POST("/", func(c *gin.Context) {
        text := c.PostForm("text")
        ctrl := task.NewTask()
        ctrl.Create(text)

        c.Redirect(http.StatusMovedPermanently, "/")
    })


    router.Run(":8080")
}
