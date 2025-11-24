package controller

import (
	"github.com/nicitapa/cooking-service/internal/middleware"
	"github.com/nicitapa/cooking-service/internal/service"
)

func (ctrl *Controller) RegisterRoutes() {
	r := ctrl.router

	// health
	r.GET("/ping", ctrl.Ping)

	// auth
	authSvc := service.NewAuthService()
	authCtrl := NewAuthController(authSvc)
	r.POST("/auth/login", authCtrl.Login)
	r.POST("/auth/refresh", authCtrl.Refresh)

	// recipes (защищённые маршруты)
	recipes := r.Group("/recipes")
	recipes.Use(middleware.AuthRequired())
	{
		recipes.GET("", ctrl.GetAll)
		recipes.GET("/get", ctrl.GetByID)
		recipes.POST("/create", ctrl.Create)
		recipes.PUT("/update", ctrl.Update)
		recipes.DELETE("/delete", ctrl.Delete)
		recipes.POST("/search", ctrl.SearchByIngredientsAndTags)
	}
}
