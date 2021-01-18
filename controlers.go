package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/prometheus/client_golang/prometheus"

)

// struct for storing data
type user struct {
	Name string `json:name`
	Age  int    `json:age`
}

// DbName database name
const DbName = "users"

// CollectionName collection name
const CollectionName = "user"

var userCollection = db().Database(DbName).Collection(CollectionName)


var getUsersCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_get_users_count", // metric name
		Help: "Number of get_users request.",
	},
	[]string{"status"}, // labels
)
var postUserCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_post_user_count", // metric name
		Help: "Number of post_user request.",
	},
	[]string{"status"}, // labels
)

func init() {
    // must register counter on init
	prometheus.MustRegister(getUsersCounter)
	prometheus.MustRegister(postUserCounter)
}

func createProfile(w http.ResponseWriter, r *http.Request) {

	fmt.Println("POST /user")
	w.Header().Set("Content-Type", "application/json")
	var status string
	defer func() {
        // increment the counter on defer func
		postUserCounter.WithLabelValues(status).Inc()
	}()
	var person user
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		fmt.Print(err)
	}
	insertResult, err := userCollection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(insertResult.InsertedID)
	fmt.Println("Inserted a single document: ", insertResult)

}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /users")

	w.Header().Set("Content-Type", "application/json")
	var status string
	defer func() {
        // increment the counter on defer func
		getUsersCounter.WithLabelValues(status).Inc()
	}()
	
	var results []primitive.M
	cur, err := userCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println(err)
	}

	for cur.Next(context.TODO()) {

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}

		results = append(results, elem)
	}
	cur.Close(context.TODO())
	json.NewEncoder(w).Encode(results)

}
