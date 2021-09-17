package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

type Activity struct {
	Gpx Gpx `xml:"gpx"`
}

type Gpx struct {
	Creator           string   `xml:"creator,attr"`
	XmlnsXsi          string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Version           string   `xml:"version,attr"`
	Xmlns             string   `xml:"xmlns,attr"`
	XmlnsGpxtpx       string   `xml:"xmlns:gpxtpx,attr"`
	XmlnsGpxx         string   `xml:"xmlns:gpxx,attr"`
	Metadata          Metadata `xml:"metadata"`
	Trk               Trk      `xml:"trk"`
	XMLName           struct{} `xml:"gpx"`
}

type Metadata struct {
	Time string `xml:"time"`
}

type Trk struct {
	Name   string `xml:"name"`
	Type   int    `xml:"type"`
	Trkseg Trkseg `xml:"trkseg"`
}

type Trkseg struct {
	Trkpt []Trkpt `xml:"trkpt"`
}

type Trkpt struct {
	Lat  float64 `xml:"lat,attr"`
	Lon  float64 `xml:"lon,attr"`
	Time string  `xml:"time"`
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
	SendToStrava(&activities)

	fmt.Printf("Finished. Processed %d activities.", len(activities))
}
