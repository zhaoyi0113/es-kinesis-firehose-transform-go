package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/zhaoyi0113/es-kinesis-firehose-transform-go/api"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	defer pprof.StopCPUProfile()
	f, err := os.Create("./profile.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	go func() {
		route := api.CreateRoute()
		route.Run()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
}
