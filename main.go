package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var RollbarApi = "https://api.rollbar.com/api/1/items"

type QueryString struct {
	AccessToken string `json:"access_token"`
}

type RollbarResponse struct {
	Err    int `json:"err"`
	Result struct {
		Items []RollbarResponseItems `json:"items"`
	} `json:"result"`
}

type RollbarResponseItems struct {
	LastOccurrenceTimestamp int64 `json:"last_occurrence_timestamp"`
}

type AlertMessage struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
	Color   string `json:"color"`
}

const slackAlertURL = "http://slack.imm.corp/api/v1/alert"

func main() {

	cliAccessToken := os.Args[1]

	queryString := QueryString{}
	queryString.AccessToken = cliAccessToken

	fullURL := fmt.Sprintf("%s/?access_token=%s", RollbarApi, cliAccessToken)

	fmt.Println("[DEBUG] CHECKING ROLLBAR API")

	httpResponse, err := http.Get(fullURL)

	if err != nil {
		fmt.Println("Error on request")
		SlackAlert(0, true)
		panic(err)
	}

	defer httpResponse.Body.Close()

	rollbarResponse := RollbarResponse{}

	rollbarResponseDecoder := json.NewDecoder(httpResponse.Body)

	err = rollbarResponseDecoder.Decode(&rollbarResponse)

	if err != nil {
		fmt.Println("Error decoding Rollbar Response")
		SlackAlert(0, true)
	}

	latestOccurrence := rollbarResponse.Result.Items[0].LastOccurrenceTimestamp

	timeNow := time.Now()
	timeEpochNow := timeNow.Unix()

	lastOccurrenceTimeDiff := timeEpochNow - latestOccurrence

	if lastOccurrenceTimeDiff > 3600 {
		fmt.Println("Last error logged over an hour ago")
		SlackAlert(lastOccurrenceTimeDiff, false)
	} else {
		fmt.Println("[DEBUG] API LOGS OK")
		fmt.Println("[DEBUG] Last Rollbar log ingested", latestOccurrence, "seconds ago")
	}

}

func SlackAlert(timeDiff int64, e bool) {

	timeInMinutes := timeDiff / 60

	slackHttpClient := &http.Client{}

	j := AlertMessage{}
	j.Channel = "#alerts"

	if e {
		j.Message = "There was an error checking the Rollbar API. Please check Rollbar manually and add better error handling to this script"
	} else {
		j.Message = fmt.Sprintf("Respondology Rollbar last message was %d minutes ago", timeInMinutes)
	}

	j.Color = "danger"

	jsonMessage, err := json.Marshal(j)

	if err != nil {
		panic(err)
	}

	_, err = slackHttpClient.Post(slackAlertURL, "application/json", bytes.NewBuffer(jsonMessage))

	if err != nil {
		fmt.Println("Error in slackAlert", err)
		panic(err)
	}

}
