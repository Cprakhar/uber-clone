package handler

import (
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RidersWSHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	riderID := ctx.Query("rider_id")
	if riderID == "" {
		log.Println("no rider_id provided")
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}
		log.Printf("received message from rider %s: %s", riderID, message)
	}
}

func DriversWSHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	driverID := ctx.Query("driver_id")
	if driverID == "" {
		log.Println("no driver_id provided")
		return
	}

	packageSlug := ctx.Query("package_slug")
	if packageSlug == "" {
		log.Println("no package_slug provided")
		return
	}

	type Driver struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ProfilePic  string `json:"profile_pic"`
		CarPlate    string `json:"car_plate"`
		PackageSlug string `json:"package_slug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			ID:          driverID,
			Name:        "Driver " + driverID,
			ProfilePic:  util.GetRandomProfilePic(rand.IntN(9) + 1),
			CarPlate:    "MH12 AB 1234",
			PackageSlug: packageSlug,
		},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("error sending message: %v", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}
		log.Printf("received message from driver %s: %s", driverID, message)
	}
}
