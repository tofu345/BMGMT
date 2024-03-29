package handlers

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/tofu345/BMGMT/constants"
	"github.com/tofu345/BMGMT/db"
	"github.com/tofu345/BMGMT/sqlc"
	"github.com/tofu345/BMGMT/utils"
)

func GetLocations(c echo.Context) error {
	data, err := db.Q.GetLocations(db.Ctx)
	if err != nil {
		c.String(http.StatusBadRequest, utils.PrettyDbError(err))
	}

	if len(data) == 0 {
		return c.JSONBlob(http.StatusOK, []byte("[]"))
	}
	return c.JSON(http.StatusOK, data)
}

type RoomDisplay struct {
	Name string `json:"name"`
	// TenancyEndDate time.Time
	User *UserDisplay `json:"user"`
}

func GetLocationInfo(c echo.Context) error {
	loc_id, err := strconv.Atoi(c.Param("loc_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, constants.BadRequest)
	}
	id := int64(loc_id)

	loc, err := db.Q.GetLocation(db.Ctx, id)
	if err != nil {
		return c.String(http.StatusBadRequest, constants.NotFound)
	}

	data, err := db.Q.GetLocationRooms(db.Ctx, id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	rooms := make([]RoomDisplay, len(data))
	for i, v := range data {
		rooms[i] = RoomDisplay{Name: v.Name, User: nil}
		if v.Email.Valid {
			rooms[i].User = &UserDisplay{
				Email:     v.Email.String,
				FirstName: v.FirstName.String,
				LastName:  v.LastName.String,
			}
		}
	}

	return c.JSON(http.StatusOK, struct {
		sqlc.Location
		Rooms []RoomDisplay
	}{Location: loc, Rooms: rooms})
}

type LocationData struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
}

func CreateLocation(c echo.Context) error {
	data := new(LocationData)
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, constants.InvalidData)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	loc, err := db.Q.CreateLocation(db.Ctx, sqlc.CreateLocationParams{
		Name:    data.Name,
		Address: data.Address,
	})
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, loc)
}

type RoomData struct {
	Name string `validate:"required"`
}

func CreateRoom(c echo.Context) error {
	data := new(RoomData)
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, constants.InvalidData)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	loc_id, err := strconv.Atoi(c.Param("loc_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, constants.BadRequest)
	}
	id := int64(loc_id)

	_, err = db.Q.GetLocation(db.Ctx, id)
	if err != nil {
		return c.String(http.StatusBadRequest, constants.NotFound)
	}

	room, err := db.Q.CreateRoom(db.Ctx, sqlc.CreateRoomParams{
		Name:       data.Name,
		LocationID: pgtype.Int8{Int64: id, Valid: true},
	})
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, RoomDisplay{Name: room.Name, User: nil})
}
