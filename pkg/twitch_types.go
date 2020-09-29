package tc

import "time"

type Stream struct {
	ID          int64     `json:"_id"`
	Game        string    `json:"game"`
	Viewers     int       `json:"viewers"`
	VideoHeight int       `json:"video_height"`
	AverageFps  int       `json:"average_fps"`
	Delay       int       `json:"delay"`
	CreatedAt   time.Time `json:"created_at"`
	IsPlaylist  bool      `json:"is_playlist"`
	Preview     struct {
		Small    string `json:"small"`
		Medium   string `json:"medium"`
		Large    string `json:"large"`
		Template string `json:"template"`
	} `json:"preview"`
	Channel struct {
		Mature                       bool        `json:"mature"`
		Status                       string      `json:"status"`
		BroadcasterLanguage          string      `json:"broadcaster_language"`
		DisplayName                  string      `json:"display_name"`
		Game                         string      `json:"game"`
		Language                     string      `json:"language"`
		ID                           int         `json:"_id"`
		Name                         string      `json:"name"`
		CreatedAt                    time.Time   `json:"created_at"`
		UpdatedAt                    time.Time   `json:"updated_at"`
		Partner                      bool        `json:"partner"`
		Logo                         string      `json:"logo"`
		VideoBanner                  string      `json:"video_banner"`
		ProfileBanner                string      `json:"profile_banner"`
		ProfileBannerBackgroundColor interface{} `json:"profile_banner_background_color"`
		URL                          string      `json:"url"`
		Views                        int         `json:"views"`
		Followers                    int         `json:"followers"`
	} `json:"channel"`
}

type V5Streamers struct {
	Streams []Stream `json:"streams"`
	Total   int32    `json:"_total"`
	Error   *string  `json:"error"`
}
