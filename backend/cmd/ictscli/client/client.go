package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	adminv1connect "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
)

var (
	teamClient adminv1connect.TeamServiceClient
	baseURL    string
	useJSON    bool
)

func InitClients(url string) {
	baseURL = url
	teamClient = adminv1connect.NewTeamServiceClient(
		http.DefaultClient,
		baseURL,
	)
}

func GetBaseURL() string {
	return baseURL
}

func GetTeamClient() adminv1connect.TeamServiceClient {
	return teamClient
}

func SetJSONOutput(enabled bool) {
	useJSON = enabled
}

func IsJSONOutput() bool {
	return useJSON
}

func PrintJSON(data interface{}, fallback func()) {
	if useJSON {
		bytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			return
		}
		fmt.Println(string(bytes))
	} else if fallback != nil {
		fallback()
	}
}
