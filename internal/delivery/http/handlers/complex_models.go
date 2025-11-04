package handlers

type CreateComplexRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	YearBuilt     *int     `json:"year_built,omitempty"`
	LocationID    string   `json:"location_id"`
	Developer     string   `json:"developer"`
	Address       string   `json:"address"`
	StartingPrice *int64   `json:"starting_price,omitempty"`
	ImageURLs     []string `json:"image_urls"`
}

type UpdateComplexRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	YearBuilt     *int     `json:"year_built,omitempty"`
	LocationID    string   `json:"location_id"`
	Developer     string   `json:"developer"`
	Address       string   `json:"address"`
	StartingPrice *int64   `json:"starting_price,omitempty"`
	ImageURLs     []string `json:"image_urls"`
}