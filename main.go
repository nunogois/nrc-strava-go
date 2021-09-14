package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

type Activity struct {
	Gpx Gpx `json:"gpx"`
}

type Gpx struct {
	Creator           string   `json:"@creator"`
	XmlnsXsi          string   `json:"@xmlns:xsi"`
	XsiSchemaLocation string   `json:"@xsi:schemaLocation"`
	Version           string   `json:"@version"`
	Xmlns             string   `json:"@xmlns"`
	XmlnsGpxtpx       string   `json:"@xmlns:gpxtpx"`
	XmlnsGpxx         string   `json:"@xmlns:gpxx"`
	Metadata          Metadata `json:"metadata"`
	Trk               Trk      `json:"trk"`
}

type Metadata struct {
	Time string `json:"time"`
}

type Trk struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file. Read the README.md for more information.")
	}
}

func main() {
	fmt.Printf("Started...\n")
	loadEnv()

	activities := make([]Activity, 0)

	GetNRCActivities(&activities)
	fmt.Println(activities)
	//SendToStrava(&activities)

	fmt.Printf("Finished. Processed %d activities", len(activities))
}
