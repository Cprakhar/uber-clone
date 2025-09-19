package handler

import (
	"encoding/json"
	"log"

	grpcclient "github.com/cprakhar/uber-clone/services/api-gateway/grpc-client"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/messaging"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
	"github.com/cprakhar/uber-clone/shared/proto/driver"
	"github.com/gin-gonic/gin"
)

var (
	connManager = messaging.NewConnectionManager()
)

// RidersWSHandler handles WebSocket connections for riders
func RidersWSHandler(ctx *gin.Context, kfc *kafka.KafkaClient) {
	conn, err := connManager.Upgrade(ctx.Writer, ctx.Request)
	if err != nil {
		log.Printf("Websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	riderID := ctx.Query("riderID")
	if riderID == "" {
		log.Println("No riderID provided")
		return
	}

	// Add the connection to the manager
	connManager.Add(riderID, conn)
	defer connManager.Remove(riderID)

	topics := []string{
		contracts.TripEventNoDriversFound,
	}

	topicConsumer := messaging.NewTopicConsumer(kfc, connManager, topics)

	go func() {
		if err := topicConsumer.Consume(ctx); err != nil {
			log.Printf("Error consuming topics: %v", err)
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("Received message from rider %s: %s", riderID, message)
	}
}

// DriversWSHandler handles WebSocket connections for drivers
func DriversWSHandler(ctx *gin.Context, kfc *kafka.KafkaClient) {
	conn, err := connManager.Upgrade(ctx.Writer, ctx.Request)
	if err != nil {
		log.Printf("Websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	driverID := ctx.Query("driverID")
	if driverID == "" {
		log.Println("No driverID provided")
		return
	}

	packageSlug := ctx.Query("packageSlug")
	if packageSlug == "" {
		log.Println("No packageSlug provided")
		return
	}

	// Add the connection to the manager
	connManager.Add(driverID, conn)

	driverService, err := grpcclient.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		connManager.Remove(driverID)

		driverService.Client.UnregisterDriver(ctx, &driver.RegisterDriverRequest{
			DriverID:    driverID,
			PackageSlug: packageSlug,
		})

		driverService.Close()
		log.Printf("Driver unregistered %s: ", driverID)
	}()

	driver, err := driverService.Client.RegisterDriver(ctx, &driver.RegisterDriverRequest{
		DriverID:    driverID,
		PackageSlug: packageSlug,
	})
	if err != nil {
		log.Printf("Failed to register driver: %v", err)
		return
	}

	msg := contracts.WSMessage{
		Type: contracts.DriverCmdRegister,
		Data: driver.Driver,
	}

	if err := connManager.SendMessage(driverID, msg); err != nil {
		log.Printf("Failed to send register message to driver %s: %v", driverID, err)
		return
	}

	topics := []string{
		contracts.DriverCmdTripRequest,
	}

	topicConsumer := messaging.NewTopicConsumer(kfc, connManager, topics)

	go func() {
		if err := topicConsumer.Consume(ctx); err != nil {
			log.Printf("Error consuming topics: %v", err)
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		type driverMessage struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}

		var dm driverMessage
		if err := json.Unmarshal(message, &dm); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		switch dm.Type {
		case contracts.DriverCmdLocation:
			// Update driver location in the system
			continue
		case contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline:
			// Notify trip service about trip acceptance
			if err := kfc.Producer.SendMessage(dm.Type, &contracts.KafkaMessage{
				EntityID: driverID,
				Data:     dm.Data,
			}); err != nil {
				log.Printf("Failed to send message to trip service: %v", err)
			}
		default:
			log.Printf("Unknown message type from driver %s: %s", driverID, dm.Type)
		}
	}
}
