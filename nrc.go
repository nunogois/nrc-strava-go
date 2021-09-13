package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type NRCActivity struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	MetricTypes []string `json:"metric_types"`
	Start       int      `json:"start_epoch_ms"`
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
	}

	activity.Gpx.Metadata.Time = r.Start // TODO: Map NRCActivity properties to Activity

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
	NIKE := os.Getenv("NIKE") // TODO: This is the normal token, we should use a refresh token to generate this instead.

	// TODO: This should come from a local json file that keeps track, probably.
	lastTime := 1622505600000 // 2021-06-01

	for lastTime != 0 {
		url := fmt.Sprintf("https://api.nike.com/sport/v3/me/activities/after_time/%d?metrics=ALL", lastTime)

		get, _ := http.NewRequest("GET", url, nil)
		get.Header.Add("Authorization", "Bearer "+NIKE)

		resp, _ := client.Do(get)

		all := struct {
			Activities []NRCActivity `json:"activities"`
			Paging     struct {
				AfterTime int `json:"after_time"`
			} `json:"paging"`
		}{}

		fmt.Println(resp.Status)

		runs := []NRCActivity{}

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
