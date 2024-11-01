package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	database "github.com/encall/cpeevent-backend/src/database"
	models "github.com/encall/cpeevent-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection = database.OpenCollection(database.Client, "posts")

func NewPost(post models.Post) interface{} {
	switch post.Kind {
	case "post":
		// Create and return a PPost
		return models.PPost{Post: post}
	case "vote":
		// Create and return a PVote with questions
		return models.PVote{Post: post, Questions: post.Questions}
	default:
		// Handle unknown post kinds, return nil or an error if needed
		return nil
	}
}


func GetPostFromEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel() // Ensure cancel is called to release resources

		// Get the eventID from the URL parameters
		eventID := c.Param("eventID")
		if eventID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "eventID is required"})
			return
		}

		// Parse eventID as an ObjectID
		objectID, err := primitive.ObjectIDFromHex(eventID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventID format"})
			return
		}

		// Query the event by its ID
		var event models.Event
		if err := eventCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&event); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Event not found"})
			return
		}

		// Query the posts collection with a single query using $in
		var posts []models.Post
		cursor, err := postCollection.Find(ctx, bson.M{"_id": bson.M{"$in": event.PostList}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving posts"})
			return
		}
		defer cursor.Close(ctx)

		// Decode all the posts from the cursor
		if err = cursor.All(ctx, &posts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding posts"})
			return
		}

		// Create a slice to hold specific post types
		var specificPosts []interface{}

		// Convert each post to its specific type based on the Kind
		for _, post := range posts {
			specificPost := NewPost(post) // Convert to specific type
			if specificPost == nil {
				continue // Or handle unknown kind if needed
			}
			specificPosts = append(specificPosts, specificPost)
		}

		// Respond with the specific posts data
		c.JSON(http.StatusOK, gin.H{"success": true, "data": specificPosts})
	}
}