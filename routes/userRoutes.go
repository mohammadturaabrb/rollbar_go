package routes

import (
    controller "rollbar_go/controllers"

    "github.com/gin-gonic/gin"
)

//UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
    incomingRoutes.POST("/users/signup", controller.SignUp())
    incomingRoutes.POST("/users/login", controller.Login())
}