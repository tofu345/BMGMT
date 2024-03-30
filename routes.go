package main

import (
	"github.com/labstack/echo/v4"
	"github.com/tofu345/BMGMT/handlers"
	"github.com/tofu345/BMGMT/utils"
)

func registerRoutes(e *echo.Echo) {
	e.GET("/users", handlers.GetUsers, utils.LoginRequired, utils.AdminRequired)
	e.GET("/user", handlers.GetUserInfo, utils.LoginRequired)
	e.POST("/users", handlers.CreateUser)
	e.POST("/token", handlers.GenerateTokenPair)
	e.POST("/token/refresh", handlers.RegenerateAccessToken)

	e.POST("/locations", handlers.CreateLocation, utils.LoginRequired, utils.AdminRequired)
	e.GET("/locations", handlers.GetLocations)
	e.GET("/locations/:loc_id", handlers.GetLocationInfo)
	e.POST("/locations/:loc_id/rooms", handlers.CreateRoom, utils.LoginRequired)
}
