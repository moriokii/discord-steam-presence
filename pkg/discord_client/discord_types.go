package discordclient

type Activity struct {
	State   string `json:"state,omitempty"`
	Details string `json:"details,omitempty"`
	Assets  Assets `json:"assets,omitempty"`
}

type Assets struct {
	LargeImageID string `json:"large_image,omitempty"`
	SmallImageID string `json:"small_image,omitempty"`
	LargeText    string `json:"large_text,omitempty"`
	SmallText    string `json:"small_text,omitempty"`
}
