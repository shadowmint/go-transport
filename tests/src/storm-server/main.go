package main

import (
	"log"
	"ntoolkit/transport"
	"os"
	"time"

	"github.com/asdine/storm"
)

type User struct {
	ID        int    // primary key
	Group     string `storm:"index"` // this field will be indexed
	Name      string // this field will not be indexed
	CreatedAt time.Time
}

func main() {

	logger := log.New(os.Stdout, "StormServer: ", log.Ldate|log.Ltime|log.Lshortfile)

	config := transport.Config{
		MaxThreads:    2,
		AcceptTimeout: 100,
		ReadTimeout:   1000,
		Logger:        logger,
	}

	trans := transport.New(handler, &config)
	trans.Listen("0.0.0.0:0")
	logger.Printf("Listening on 0.0.0.0:%d...\n", trans.Port())

	trans.Wait()
}

type message struct {
	Message string
}

func handler(api *transport.API) {
	var msg message
	if err := api.Read(&msg); err == nil {
		processMessage(api, &msg)
	} else {
		api.Logger.Printf("Unknown message format: %s\n", api.Raw())
	}
}

func processMessage(api *transport.API, msg *message) {
	api.Logger.Printf("Got message: %s\n", msg.Message)

	db, err := storm.Open("bolt.db", storm.AutoIncrement())
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Init(&User{})
	if err != nil {
		api.Logger.Printf("Init failed: %v\n", err)
	}

	user := User{
		Group:     "staff",
		Name:      "John",
		CreatedAt: time.Now(),
	}

	err = db.Save(&user)
	if err != nil {
		api.Logger.Printf("Save failed: %v\n", err)
	}

	var users []User
	err = db.Find("Group", "staff", &users)
	if err != nil {
		api.Logger.Printf("Find failed: %v\n", err)
	}
	api.Logger.Printf("%v\n", users)
}
