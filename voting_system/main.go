package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Candidate struct {
	Name  string `json:"name" bson:"_id"`
	Votes int64  `json:"votes" bson:"count"`
}

type Vote struct {
	Candidate string             `json:"candidate" bson:"candidate"`
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func GetResults(c *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Ensure the request is a GET method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Define the aggregation pipeline
		pipeline := mongo.Pipeline{
			{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$candidate"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			}}},
			{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		}

		// Execute the aggregation
		cursor, err := c.Aggregate(context.Background(), pipeline)
		if err != nil {
			http.Error(w, "Failed to aggregate results", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())

		// Parse the aggregation results
		results := []Candidate{}
		if err := cursor.All(context.Background(), &results); err != nil {
			http.Error(w, "Failed to parse results", http.StatusInternalServerError)
			return
		}

		// Respond with the aggregated results
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func Add(c *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Ensure the request is a POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := context.Background()

		userId := r.Header.Get("user")
		// Parse the JSON body into parameter type
		var doc Vote
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		// set id on insert
		// will be returned with response.
		doc.Id = primitive.NewObjectID()
		doc.CreatedAt = time.Now()
		doc.CreatedBy = userId

		// Insert the document into the collection
		_, err := c.InsertOne(ctx, doc)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to save vote : %s", err),
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// write response back without allocating bytes
		if err := json.NewEncoder(w).Encode(&doc); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusInternalServerError)
			return
		}
	}
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://localhost:27017"),
	)

	if err != nil {
		panic(err)
	}

	defer client.Disconnect(context.Background())

	dbname := "my_database"
	collectionname := "votes"

	client.Database(dbname).Drop(context.Background())
	collection := client.Database(dbname).Collection(collectionname)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "created_by", Value: 1}}, // Ascending index on created_by
		Options: options.Index().SetUnique(true),       // Make the index unique
	}

	// Create the index
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

	http.HandleFunc("/vote", Add(collection))
	http.HandleFunc("/results", GetResults(collection))

	log.Println("Listenning at port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}
