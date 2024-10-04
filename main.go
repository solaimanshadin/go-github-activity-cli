package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type GithubEvent struct{
	ID string `json:"id"`
	Type string `json:"type"`
	Repo struct {
		ID int `json:"id"`
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Action string `json:"action"`
		Size int `json:"size"`
	}
}

type EventGroup struct {
	Pushes []struct {
		RepoName string `json:"repoName"`
		CommitCount int `json:"commitCount"`
	} `json:"pushes"`
	Stars []string `json:"stars"`
}

var client = &http.Client{
    Timeout: 10 * time.Second,
}

func getDataFromGithub(userName string) ([]GithubEvent, error) {
	url := "https://api.github.com/users/"+userName+"/events"

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: status code %d", resp.StatusCode)
	}
	
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	
	if err != nil {
		return nil , err
	}
	events := []GithubEvent{}
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil , err
	}

	return events, nil
}

func generateEventGroup(events []GithubEvent) EventGroup {
	eventGroup := EventGroup{}
	for _, event := range events{
		if event.Type == "PushEvent" {
			eventGroup.Pushes = append(eventGroup.Pushes, struct{
				RepoName string `json:"repoName"`;
				CommitCount int `json:"commitCount"`
				} {
				RepoName: event.Repo.Name,
				CommitCount: event.Payload.Size})
		} else if event.Type == "WatchEvent" && event.Payload.Action == "started" {
			eventGroup.Stars = append(eventGroup.Stars, event.Repo.Name)
		}
	}
	return eventGroup
}

func printResponse(eventGroup EventGroup) {
	for _, push := range eventGroup.Pushes {
		fmt.Printf("- Pushed %d commits into %s\n", push.CommitCount, push.RepoName)
	}
	for _, starredRepoName := range eventGroup.Stars {
		fmt.Printf("- Starred %s\n", starredRepoName)
	}
}
func main() {
	if(len(os.Args) < 2 ) {
		fmt.Println("Usage: go run main.go <github_username>")
    	return
	}
	
	events, err := getDataFromGithub(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	eventGroup := generateEventGroup(events)
	
	printResponse(eventGroup)

}