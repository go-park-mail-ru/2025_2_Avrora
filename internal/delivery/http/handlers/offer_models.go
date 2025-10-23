package handlers

type CreateOfferRequest struct {
	InHousingComplex bool     `json:"in_housing_complex"`
	HousingComplex   string   `json:"housing_complex"`
	OfferType        string   `json:"offer_type"`    // sale | rent
	PropertyType     string   `json:"property_type"` // house | apartment
	Title            string   `json:"title"`
	UserID           int      `json:"user_id"`
	Category         string   `json:"category"`
	Address          string   `json:"address"`
	Floor            int      `json:"floor"`
	TotalFloors      int      `json:"total_floors"`
	Rooms            int      `json:"rooms"`
	Area             float64  `json:"area"`
	LivingArea       float64  `json:"living_area"`
	KitchenArea      float64  `json:"kitchen_area"`
	Price            float64  `json:"price"`
	Description      string   `json:"description"`
	Deposit          float64  `json:"deposit"`
	Commission       float64  `json:"commission"`
	RentalPeriod     string   `json:"rental_period"`
	ImageURLs        []string `json:"image_urls"`
}

type FullOfferResponse struct {
	InHousingComplex bool     `json:"in_housing_complex"`
	HousingComplex   string   `json:"housing_complex"`
	OfferType        string   `json:"offer_type"`    // sale | rent
	PropertyType     string   `json:"property_type"` // house | apartment
	Title            string   `json:"title"`
	Category         string   `json:"category"`
	Address          string   `json:"address"`
	Floor            int      `json:"floor"`
	TotalFloors      int      `json:"total_floors"`
	Rooms            int      `json:"rooms"`
	Area             float64  `json:"area"`
	LivingArea       float64  `json:"living_area"`
	KitchenArea      float64  `json:"kitchen_area"`
	Price            float64  `json:"price"`
	Description      string   `json:"description"`
	Deposit          float64  `json:"deposit"`
	Commission       float64  `json:"commission"`
	RentalPeriod     string   `json:"rental_period"`
	ImageURLs        []string `json:"image_urls"`
}

type OfferInFeedResponse struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	OfferURL     string  `json:"offer_url"`
	OfferType    string  `json:"offer_type"`    // sale | rent
	PropertyType string  `json:"property_type"` // house | apartment
	Price        float64 `json:"price"`
	Area         float64 `json:"area"`
	Rooms        int     `json:"rooms"`
	Floor        int     `json:"floor"`
	TotalFloors  int     `json:"total_floors"`
	Address      string  `json:"address"`
	Metro        string  `json:"metro"`
	ImageURL     string  `json:"image_url"`
}

type UpdateOfferRequest struct {
	ID               int      `json:"id"`
	InHousingComplex bool     `json:"in_housing_complex"`
	HousingComplex   string   `json:"housing_complex"`
	OfferType        string   `json:"offer_type"`    // sale | rent
	PropertyType     string   `json:"property_type"` // house | apartment
	Title            string   `json:"title"`
	UserID           int      `json:"user_id"`
	Category         string   `json:"category"`
	Address          string   `json:"address"`
	Status           string   `json:"status"`
	Floor            int      `json:"floor"`
	TotalFloors      int      `json:"total_floors"`
	Rooms            int      `json:"rooms"`
	Area             float64  `json:"area"`
	LivingArea       float64  `json:"living_area"`
	KitchenArea      float64  `json:"kitchen_area"`
	Price            float64  `json:"price"`
	Description      string   `json:"description"`
	Deposit          int64    `json:"deposit"`
	Commission       int64    `json:"commission"`
	RentalPeriod     string   `json:"rental_period"`
	ImageURLs        []string `json:"image_urls"`
}
