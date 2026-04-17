package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func main() {
	initFirebase()

	r := gin.Default()
	r.Use(cors.Default())

	// HOME
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Event backend running 🚀"})
	})

	// ROUTES
	r.POST("/events", createEvent)
	r.GET("/events", getEvents)
	r.POST("/book/:id", bookEvent)
	r.POST("/cancel/:id", cancelBooking)

	port := os.Getenv("PORT")
if port == "" {
	port = "8080"
}

r.Run(":" + port)
}

// =====================
// CREATE EVENT
// =====================
func createEvent(c *gin.Context) {
	ctx := context.Background()

	var event struct {
		Name  string `json:"name"`
		Date  string `json:"date"`
		Slots int    `json:"slots"`
	}

	if err := c.BindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	client := firestoreClient

	_, _, err := client.Collection("events").Add(ctx, map[string]interface{}{
		"name":  event.Name,
		"date":  event.Date,
		"slots": event.Slots,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Event created 🚀"})
}

// =====================
// GET EVENTS
// =====================
func getEvents(c *gin.Context) {
	ctx := context.Background()

	client := firestoreClient

	iter := client.Collection("events").Documents(ctx)

	var events []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		events = append(events, map[string]interface{}{
			"id":    doc.Ref.ID,
			"name":  doc.Data()["name"],
			"date":  doc.Data()["date"],
			"slots": doc.Data()["slots"],
		})
	}

	c.JSON(200, events)
}

// =====================
// BOOK EVENT
// =====================
func bookEvent(c *gin.Context) {
	ctx := context.Background()
	id := c.Param("id")

	client := firestoreClient

	docRef := client.Collection("events").Doc(id)

	doc, err := docRef.Get(ctx)
	if err != nil {
		c.JSON(404, gin.H{"error": "Event not found"})
		return
	}

	slots := doc.Data()["slots"].(int64)

	if slots <= 0 {
		c.JSON(400, gin.H{"error": "No slots available"})
		return
	}

	// 🔽 decrease slot
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "slots", Value: slots - 1},
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 📝 save booking
	_, _, err = client.Collection("bookings").Add(ctx, map[string]interface{}{
		"eventId": id,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Booked successfully 🎟️"})
}

// =====================
// CANCEL BOOKING
// =====================
func cancelBooking(c *gin.Context) {
	ctx := context.Background()
	id := c.Param("id")

	client := firestoreClient

	docRef := client.Collection("events").Doc(id)

	// 🔼 increase slot
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "slots", Value: firestore.Increment(1)},
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Booking cancelled 🔄"})
}