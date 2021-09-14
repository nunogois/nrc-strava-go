package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type NRCActivity struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	MetricTypes []string `json:"metric_types"`
	Start       int      `json:"start_epoch_ms"`
}

type NRCResponse struct {
	Activities []NRCActivity `json:"activities"`
	Paging     Paging        `json:"paging"`
}

type Paging struct {
	AfterTime int `json:"after_time"`
}

type TokenBody struct {
	ClientId     string `json:"client_id"`
	GrantType    string `json:"grant_type"`
	UxId         string `json:"ux_id"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResp struct {
	AccessToken string `json:"access_token"`
}

func getTime(timestamp int) time.Time {
	i := int64(timestamp / 1000)
	t := time.Unix(i, 0)

	return t
}

func timeToISO(timestamp int) (date string) {
	return getTime(timestamp).Format("2006-01-02T15:04:05-0700")
}

func timeToISODate(timestamp int) (date string) {
	return getTime(timestamp).Format("2006-01-02")
}

func getToken() (token string) {
	body := TokenBody{
		ClientId:     os.Getenv("NIKE_CLIENT_ID"),
		GrantType:    "refresh_token",
		UxId:         "com.nike.sport.running.ios.6.5.1",
		RefreshToken: os.Getenv("NIKE_REFRESH_TOKEN"),
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(&body)

	post, _ := http.NewRequest("POST", "https://api.nike.com/idn/shim/oauth/2.0/token", payloadBuf)

	resp, err := client.Do(post)

	if err != nil {
		return
	}

	response := TokenResp{}

	json.NewDecoder(resp.Body).Decode(&response)
	defer resp.Body.Close()

	return response.AccessToken
}

func NewActivity(r NRCActivity) Activity {
	activity := Activity{}

	activity.Gpx = Gpx{
		Creator:           "StravaGPX",
		XmlnsXsi:          "http://www.w3.org/2001/XMLSchema-instance",
		XsiSchemaLocation: "http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd",
		Version:           "1.1",
		Xmlns:             "http://www.topografix.com/GPX/1/1",
		XmlnsGpxtpx:       "http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
		XmlnsGpxx:         "http://www.garmin.com/xmlschemas/GpxExtensions/v3",
		Metadata: Metadata{
			Time: timeToISODate(r.Start),
		},
		Trk: Trk{
			Name: fmt.Sprintf("%s - NRC", timeToISO(r.Start)),
			Type: 9,
		},
	}

	return activity
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetNRCActivities(activities *[]Activity) {
	// TODO: This should come from a local json file that keeps track, probably.
	lastTime := 1622505600000 // 2021-06-01

	for lastTime != 0 {
		url := fmt.Sprintf("https://api.nike.com/sport/v3/me/activities/after_time/%d?metrics=ALL", lastTime)

		get, _ := http.NewRequest("GET", url, nil)
		get.Header.Add("Authorization", fmt.Sprintf("Bearer %s", getToken()))

		resp, _ := client.Do(get)

		all := NRCResponse{}
		runs := []NRCActivity{}

		fmt.Println(resp.Status)

		json.NewDecoder(resp.Body).Decode(&all)
		defer resp.Body.Close()

		for _, a := range all.Activities {
			if a.Type == "run" && contains(a.MetricTypes, "latitude") && contains(a.MetricTypes, "longitude") {
				runs = append(runs, a)
			}
		}

		for _, r := range runs {
			*activities = append(*activities, NewActivity(r))
		}

		fmt.Printf("Total run activities found: %d\n", len(runs))
		lastTime = all.Paging.AfterTime
	}
}
