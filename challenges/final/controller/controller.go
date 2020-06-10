package controller

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sonyarouje/simdb/db"

	"go.nanomsg.org/mangos"
	"go.nanomsg.org/mangos/protocol/surveyor"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"
)

type Worker struct {
	Name   string
	Tags   string
	Status string
	Usage  string
	URL    string
	Port   string
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

var workersList []Worker // Create an empty list of workers

func Start(contAddr string, db *db.Driver) {
	log.Printf("controller started")
	//var sock mangos.Socket
	//GET THESURVEYOR LISTENING
	sock, err := surveyor.NewSocket()
	if err != nil {
		die("can't get new surveyor socket: %s", err)
	}

	err = sock.Listen(contAddr)
	if err != nil {
		die("can't listen on controller: %s", err.Error())
	}

	err = sock.SetOption(mangos.OptionSurveyTime, time.Second/2)
	if err != nil {
		die("couldn't set option for survey time: %s", err.Error())
	}

	for {
		err = sock.Send([]byte("Sending message, waiting for reception"))
		if err != nil {
			die("couldn't listen to workers: %s", err.Error())
		}
		// Could also use sock.RecvMsg to get header

		for {
			msg, err := sock.Recv()
			if err != nil {
				die("Message wasn't received: %s", err.Error())
				break
			}
			workerData := strings.Split(string(msg), "|") //Recieve the msg and put it in the struct worker
			worker := Worker{}
			worker.Name = workerData[0]
			worker.Tags = workerData[1]
			worker.Status = workerData[2]
			worker.Usage = workerData[3]
			worker.URL = workerData[4]
			worker.Port = workerData[5]

			workersList = append(workersList, worker) //Add the worker to the workersList
		}

		// d := date()
		// log.Printf("Controller: Publishing Date %s\n", d)
		//
		// if err = sock.Send([]byte(d)); err != nil {
		// 	die("Failed publishing: %s", err.Error())
		// }
		// time.Sleep(time.Second * 3)

	}
}

func WorkerStatus(name string) (string, string, string, string, string) {
	for _, i := range workersList {
		if i.Name == name {
			//Return all workers attributes
			return i.Tags, i.Status, i.Usage, i.URL, i.Port
		}
	}
	return "", "", "", "", ""
}
