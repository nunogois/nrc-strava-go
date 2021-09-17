package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type NRCActivity struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	MetricTypes []string `json:"metric_types"`
	Metrics     []Metric `json:"metrics"`
	Start       int      `json:"start_epoch_ms"`
}

type Metric struct {
	Type   string        `json:"type"`
	Values []MetricValue `json:"values"`
}

type MetricValue struct {
	Start int     `json:"start_epoch_ms"`
	End   int     `json:"end_epoch_ms"`
	Value float64 `json:"value"`
}

type NRCResponse struct {
	Activities []NRCActivity `json:"activities"`
	Paging     Paging        `json:"paging"`
}

type Paging struct {
	AfterTime int `json:"after_time"`
}

type NRCTokenBody struct {
	ClientId     string `json:"client_id"`
	GrantType    string `json:"grant_type"`
	UxId         string `json:"ux_id"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResp struct {
	AccessToken string `json:"access_token"` // TODO: This is global between Strava and NRC
}

func getTime(timestamp int) time.Time {
	i := int64(timestamp / 1000)
	t := time.Unix(i, 0)

	return t
}

func dateToTime(date string) int {
	t, _ := time.Parse("2006-01-02", date)
	return int(t.Unix() * 1000)
}

func timeToISO(timestamp int) string {
	return getTime(timestamp).Format("2006-01-02T15:04:05-0700") // 2006-01-02T15:04:05.371Z
}

func timeToISODate(timestamp int) string {
	return getTime(timestamp).Format("2006-01-02")
}

func getNRCToken() string {
	body := NRCTokenBody{
		ClientId:     os.Getenv("NIKE_CLIENT_ID"),
		GrantType:    "refresh_token",
		UxId:         "com.nike.sport.running.ios.6.5.1",
		RefreshToken: os.Getenv("NIKE_REFRESH_TOKEN"),
	}

	bodyJson, _ := json.Marshal(&body)
	payload := strings.NewReader(string(bodyJson))

	post, _ := http.NewRequest("POST", "https://api.nike.com/idn/shim/oauth/2.0/token", payload)
	post.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(post)

	if err != nil || resp.Status != "200 OK" {
		fmt.Println("Something went wrong requesting NRC token. Please review your .env information.")
		return ""
	}

	response := TokenResp{}

	json.NewDecoder(resp.Body).Decode(&response)
	defer resp.Body.Close()

	return response.AccessToken
}

func NewActivity(r NRCActivity) Activity {
	activity := Activity{}

	points := []Trkpt{}

	for _, metric := range r.Metrics {
		if metric.Type == "latitude" {
			for i, lat := range metric.Values {

				var lon MetricValue
				for _, metric := range r.Metrics {
					if metric.Type == "longitude" {
						lon = metric.Values[i]
					}
				}

				point := Trkpt{
					Time: timeToISO(lat.Start),
					Lat:  lat.Value,
					Lon:  lon.Value,
				}

				points = append(points, point)
			}
		}
	}

	activity.Gpx = Gpx{
		Creator:           "StravaGPX",
		XmlnsXsi:          "http://www.w3.org/2001/XMLSchema-instance",
		XsiSchemaLocation: "http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd",
		Version:           "1.1",
		Xmlns:             "http://www.topografix.com/GPX/1/1",
		XmlnsGpxtpx:       "http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
		XmlnsGpxx:         "http://www.garmin.com/xmlschemas/GpxExtensions/v3",
		Metadata: Metadata{
			Time: timeToISO(r.Start),
		},
		Trk: Trk{
			Name: fmt.Sprintf("%s - NRC", timeToISODate(r.Start)),
			Type: 9,
			Trkseg: Trkseg{
				Trkpt: points,
			},
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
	lastTime := dateToTime(os.Getenv("START"))

	for lastTime != 0 {
		url := fmt.Sprintf("https://api.nike.com/sport/v3/me/activities/after_time/%d?metrics=ALL", lastTime)

		get, _ := http.NewRequest("GET", url, nil)
		get.Header.Add("Authorization", fmt.Sprintf("Bearer %s", getNRCToken())) // TODO: Needs refactoring, currently being requested on loop

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
