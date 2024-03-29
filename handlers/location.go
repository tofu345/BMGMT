package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tofu345/BMGMT/constants"
	"github.com/tofu345/BMGMT/db"
	"github.com/tofu345/BMGMT/utils"
)

type LocationDisplay struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func GetLocations(c echo.Context) error {
	data, err := db.Q.GetLocations(db.Ctx)
	if err != nil {
		c.String(http.StatusBadRequest, utils.PrettyDbError(err))
	}

	locations := make([]LocationDisplay, len(data))
	for i, v := range data {
		locations[i] = LocationDisplay{
			Name:    v.Name,
			Address: v.Address,
		}
	}

	return c.JSON(http.StatusOK, locations)
}

type RoomDisplay struct {
	Name string 
	// TenancyEndDate time.Time
}

func GetLocationRooms(c echo.Context) error {
	loc_id, err := strconv.Atoi(c.Param("loc_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, constants.BadRequest)
	}

	data, err := db.Q.GetLocationRooms(db.Ctx, int64(loc_id))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	rooms := make([]RoomDisplay, len(data))
	for i, v := range data {
		rooms[i] = RoomDisplay{
			Name: v.Name,
			// TenancyEndDate: v.TenancyEndDate.Time,
		}
	}

	return c.JSON(http.StatusOK, rooms)
}
