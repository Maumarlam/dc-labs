package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Maumarlam/dc-labs/challenges/final/controller"
	"github.com/Maumarlam/dc-labs/challenges/final/scheduler"

	"github.com/sonyarouje/simdb/db" //Simple database for go (JSON)
)

func main() {
	log.Println("Welcome to the Distributed and Parallel Image Processing System")

	controllerAddress := "tcp://localhost:40899"
	port := ":8080"
	DB := "workers"

	//Start the database
	db, err := db.New(DB) //Create the database
	if err != nil {
		panic(err)
	}

	// Start Controller
	go controller.Start(controllerAddress, db)

	// Start Scheduler
	jobs := make(chan scheduler.Job)
	go scheduler.Start(jobs)

	//Start api
	//Que le necesito mandar a mi api?
	go api.apiStart()

	///////////////////////////////////////////////////
	// Send sample jobs
	sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: "hello"}

	for {
		sampleJob.RPCName = fmt.Sprintf("hello-%v", rand.Intn(10000))
		jobs <- sampleJob
		time.Sleep(time.Second * 5)
	}
	// API
	// Here's where your API setup will be
}
