package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

func insertMany(client *mongo.Client, ctx context.Context,
	dataBase, col string, docs []interface{}) (*mongo.InsertManyResult, error) {

	// select database and collection ith Client.Database
	// method and Database.Collection method
	collection := client.Database(dataBase).Collection(col)

	// InsertMany accept two argument of type Context
	// and of empty interface
	result, err := collection.InsertMany(ctx, docs)
	return result, err
}

func query(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor,
	err error) {

	// select database and collection.
	collection := client.Database(dataBase).Collection(col)

	// collection has an method Find,
	// that returns a mongo.cursor
	// based on query and field.
	result, err = collection.Find(ctx, query,
		options.Find().SetProjection(field))
	return
}

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func register(w http.ResponseWriter, r *http.Request) {

	var user interface{}

	r.ParseForm()

	name := r.Form.Get("name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user = bson.D{
		{"name", name},
		{"email", email},
		{"password", password},
	}

	client, ctx, cancel, err := connect("mongodb+srv://mrinalseth3959:infant%40123@cluster0.7jw4x.mongodb.net/goDatabase?retryWrites=true&w=majority")

	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	insertOne := func(client *mongo.Client, ctx context.Context,
		dataBase, col string, doc interface{}) (*mongo.InsertOneResult,
		error) {

		// select database and collection ith Client.Database method
		// and Database.Collection method
		collection := client.Database(dataBase).Collection(col)

		// InsertOne accept two argument of type Context
		// and of empty interface
		result, err := collection.InsertOne(ctx, doc)
		return result, err
	}

	insertOneResult, err := insertOne(client, ctx,
		"testing", "user", user)

	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, "Result of InsertOne")
	fmt.Fprint(w, insertOneResult.InsertedID)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var post interface{}

	dt := time.Now()
	caption := r.Form.Get("caption")
	url := r.Form.Get("url")
	post = bson.D{
		{"caption", caption},
		{"url", url},
		{"timestamp", dt},
	}

	client, ctx, cancel, err := connect("mongodb+srv://mrinalseth3959:infant%40123@cluster0.7jw4x.mongodb.net/goDatabase?retryWrites=true&w=majority")

	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	insertOne := func(client *mongo.Client, ctx context.Context,
		dataBase, col string, doc interface{}) (*mongo.InsertOneResult,
		error) {

		// select database and collection ith Client.Database method
		// and Database.Collection method
		collection := client.Database(dataBase).Collection(col)

		// InsertOne accept two argument of type Context
		// and of empty interface
		result, err := collection.InsertOne(ctx, doc)
		return result, err
	}

	insertOneRes, err := insertOne(client, ctx, "testing", "post", post)

	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, "Result of InsertOne")
	fmt.Fprint(w, insertOneRes.InsertedID)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	client, ctx, cancel, err := connect("mongodb+srv://mrinalseth3959:infant%40123@cluster0.7jw4x.mongodb.net/goDatabase?retryWrites=true&w=majority")

	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	var filter, option interface{}

	id := r.Form.Get("id")

	filter = bson.D{
		{"id", id},
	}

	cursor, err := query(client, ctx, "testing", "post", filter, option)

	if err != nil {
		panic(err)
	}

	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {

		// handle the error
		panic(err)
	}

	// printing the result of query.
	fmt.Println("Query Result")
	for _, doc := range results {
		fmt.Println(doc)
	}

}

func handleRequest() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/post", createPost).Methods("POST")
	router.HandleFunc("/post", getPost).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func main() {

	handleRequest()
}
