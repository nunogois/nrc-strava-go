package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
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

	fmt.Println(body)

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

func sendActivity(activity *Activity, token string) {
	activityXml, _ := xml.MarshalIndent(activity.Gpx, "", " ")
	activityXml = []byte(xml.Header + string(activityXml))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("description", "Uploaded from NRC.")
	writer.WriteField("data_type", "gpx")
	part, _ := writer.CreateFormFile("file", activity.Gpx.Trk.Name+".gpx")

	part.Write(activityXml)
	writer.Close()

	send, _ := http.NewRequest("POST", "https://www.strava.com/api/v3/uploads", body)

	send.Header.Add("Content-Type", writer.FormDataContentType())
	send.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(send)

	if err != nil || resp.Status != "201 Created" {
		fmt.Println("Something went wrong uploading Gpx data to Strava.")
		fmt.Println("Status: " + resp.Status)
		fmt.Println("Error: " + err.Error())
	} else {
		fmt.Println("Successfully uploaded activity: " + activity.Gpx.Trk.Name)
	}
}

func SendToStrava(activities *[]Activity) {
	token := getStravaToken()
	for _, activity := range *activities {
		sendActivity(&activity, token)
	}
}
