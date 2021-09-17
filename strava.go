package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type StravaTokenBody struct {
	ClientId     string `json:"client_id"`
	GrantType    string `json:"grant_type"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

func getStravaToken() string {
	body := StravaTokenBody{
		ClientId:     os.Getenv("STRAVA_CLIENT_ID"),
		GrantType:    "refresh_token",
		ClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		RefreshToken: os.Getenv("STRAVA_REFRESH_TOKEN"),
	}

	bodyJson, _ := json.Marshal(&body)
	payload := strings.NewReader(string(bodyJson))

	post, _ := http.NewRequest("POST", "https://www.strava.com/api/v3/oauth/token", payload)
	post.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(post)

	if err != nil || resp.Status != "200 OK" {
		fmt.Println("Something went wrong requesting Strava token. Please review your .env information.")
		return ""
	}

	response := TokenResp{}

	json.NewDecoder(resp.Body).Decode(&response)
	defer resp.Body.Close()

	return response.AccessToken
}

func sendActivity(activity *Activity, token string, dir string) {

	file, err := ioutil.TempFile(dir, "*.gpx")
	defer os.Remove(file.Name())
	defer os.Remove(dir)
	// if err != nil {
	// 	fmt.Println("Error: " + err.Error())
	// }
	activityXml, err := xml.MarshalIndent(activity.Gpx, "", " ")
	activityXml = []byte(xml.Header + string(activityXml))
	file.Write(activityXml)

	fmt.Println("Dir: " + file.Name())
	fmt.Println("File: " + file.Name())
	fmt.Println("Filepath base: " + filepath.Base(file.Name()))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("description", "Uploaded from NRC.")
	writer.WriteField("data_type", "gpx")
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	send, err := http.NewRequest("POST", "https://www.strava.com/api/v3/uploads", body)
	if err != nil {
		fmt.Println("Something went wrong request: " + err.Error())
	}
	send.Header.Add("Content-Type", writer.FormDataContentType())
	send.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(send)

	fmt.Println(resp.Status)

	if err != nil {
		fmt.Println("Something went wrong: " + err.Error())
	} else {
		fmt.Println("Successfully uploaded activity: " + activity.Gpx.Trk.Name)
	}
}

func SendToStrava(activities *[]Activity) {
	dir, _ := ioutil.TempDir(".", "temp-")
	defer os.Remove(dir)

	for _, activity := range *activities {
		sendActivity(&activity, getStravaToken(), dir)
	}
}
