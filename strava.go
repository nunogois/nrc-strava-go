package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func sendActivity(activity *Activity) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)

		file, err := ioutil.TempFile("nrc-strava-go", "*.gpx")
		defer os.Remove(file.Name())
		xml, _ := xml.MarshalIndent(activity, "", " ")
		file.Write(xml)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("description", "Uploaded from NRC.")
		writer.WriteField("data_type", "gpx")
		part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
		io.Copy(part, file)
		writer.Close()

		send, err := http.NewRequest("POST", "https://www.strava.com/api/v3/uploads", body)
		send.Header.Add("Content-Type", writer.FormDataContentType())
		send.Header.Add("Authorization", "Bearer STRAVA_TOKEN") // TODO: Grab STRAVA token.

		_, err = client.Do(send)

		if err != nil {
			r <- "Something went wrong: " + err.Error()
		} else {
			r <- "Successfully uploaded activity: " + activity.Gpx.Creator
		}
	}()

	return r
}

func SendToStrava(activities *[]Activity) {
	for _, activity := range *activities {
		sendActivity(&activity)
	}
}
