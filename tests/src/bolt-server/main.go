package main

import (
	"fmt"
	"log"
	"math/rand"
	"ntoolkit/transport"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func main() {

	logger := log.New(os.Stdout, "BoltServer: ", log.Ldate|log.Ltime|log.Lshortfile)

	config := transport.Config{
		MaxThreads:    2,
		AcceptTimeout: 100,
		ReadTimeout:   1000,
		Logger:        logger,
	}

	trans := transport.New(handler, &config)
	trans.Listen("127.0.0.1:0")
	logger.Printf("Listening on %d...\n", trans.Port())

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

	var world = []byte("world")
	rand.Seed(time.Now().UTC().UnixNano())

	db, err := bolt.Open("bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	value := []byte("Hello World!")

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err2 := tx.CreateBucketIfNotExists(world)
		if err2 != nil {
			return err
		}

		key := fmt.Sprintf("Key:%d", rand.Int())
		fmt.Printf("Inserting key: %s\n", key)
		err = bucket.Put([]byte(key), value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// retrieve the data
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(world)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", world)
		}

		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("A %s is %s.\n", k, v)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
