package main

import (
	"github.com/labstack/echo/v4"
	"github.com/tofu345/BMGMT/handlers"
	"github.com/tofu345/BMGMT/utils"
)

func registerRoutes(e *echo.Echo) {
	e.GET("/users", handlers.GetUsers, utils.SuperUserRequired)
	e.GET("/user", handlers.GetUserInfo, utils.LoginRequired)
	e.POST("/users", handlers.CreateUser)
	e.POST("/token", handlers.GenerateTokenPair)
	e.POST("/token/refresh", handlers.RegenerateAccessToken)

	e.POST("/locations", handlers.CreateLocation, utils.SuperUserRequired)
	e.GET("/locations", handlers.GetLocations)
	e.POST("/locations/:loc_id/admin", handlers.CreateLocationAdmin, utils.SuperUserRequired)
	e.GET("/locations/:loc_id", handlers.GetLocationInfo)
	e.POST("/locations/:loc_id/rooms", handlers.CreateRoom, utils.LoginRequired)
	e.GET("/locations/:loc_id/issues", handlers.GetLocationIssues, utils.LocationAdminRequired)
	e.POST("/locations/:loc_id/issues", handlers.CreateLocationIssue, utils.LoginRequired)
}
