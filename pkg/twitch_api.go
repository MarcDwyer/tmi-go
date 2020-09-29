package tc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type TwitchAPI struct {
	ClientID string
	OAuth    string
}

func NewTwitchAPI(clientID string, oauth string) *TwitchAPI {
	return &TwitchAPI{
		OAuth:    oauth,
		ClientID: clientID,
	}
}

func (t TwitchAPI) FetchV5Data(url string) (*V5Streamers, error) {
	token := fmt.Sprintf("OAuth %s", t.OAuth)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Client-ID", t.ClientID)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data V5Streamers

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	if data.Error != nil {
		return nil, errors.New("Error found in twitch payload")
	}
	return &data, nil
}
