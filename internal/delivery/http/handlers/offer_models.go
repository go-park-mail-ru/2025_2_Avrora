package handlers

//go:generate easyjson -all $GOFILE
//easyjson:json
type CreateOfferRequest struct {
	InHousingComplex bool     `json:"in_housing_complex"`
	HousingComplex   string   `json:"housing_complex"`
	OfferType        string   `json:"offer_type"`    // sale | rent
	PropertyType     string   `json:"property_type"` // house | apartment
	Title            string   `json:"title"`
	UserID           string   `json:"user_id"`
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
	Deposit          int64    `json:"deposit"`
	Commission       int64    `json:"commission"`
	RentalPeriod     string   `json:"rental_period"`
	ImageURLs        []string `json:"image_urls"`
}

//easyjson:json
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
	Metro            string   `json:"metro"`
	ImageURLs        []string `json:"image_urls"`
}

//easyjson:json
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

//easyjson:json
type UpdateOfferRequest struct {
	ID               int      `json:"id"`
	InHousingComplex bool     `json:"in_housing_complex"`
	HousingComplex   string   `json:"housing_complex"`
	OfferType        string   `json:"offer_type"`    // sale | rent
	PropertyType     string   `json:"property_type"` // house | apartment
	Title            string   `json:"title"`
	UserID           string   `json:"user_id"`
	Category         string   `json:"category"`
	Address          string   `json:"address"`
	Status           string   `json:"status"` // active | sold | archived
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

// WebhookRequest represents the root structure of the incoming webhook payload.
//
//easyjson:json
type WebhookRequest struct {
	Event  string  `json:"event"`
	Type   string  `json:"type"`
	Object Payment `json:"object"`
}

// Payment represents the payment details in the webhook payload.
//
//easyjson:json
type Payment struct {
	Amount               AmountDetails        `json:"amount"`
	AuthorizationDetails AuthorizationDetails `json:"authorization_details"`
	CreatedAt            string               `json:"created_at"`
	Description          string               `json:"description"`
	ExpiresAt            string               `json:"expires_at"`
	ID                   string               `json:"id"`
	Metadata             map[string]string    `json:"metadata"`
	Paid                 bool                 `json:"paid"`
	PaymentMethod        PaymentMethod        `json:"payment_method"`
	Recipient            Recipient            `json:"recipient"`
	Refundable           bool                 `json:"refundable"`
	Status               string               `json:"status"`
	Test                 bool                 `json:"test"`
}

// AmountDetails represents the amount details in the payment.
//
//easyjson:json
type AmountDetails struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

// AuthorizationDetails represents the authorization details in the payment.
//
//easyjson:json
type AuthorizationDetails struct {
	AuthCode     string              `json:"auth_code"`
	RRN          string              `json:"rrn"`
	ThreeDSecure ThreeDSecureDetails `json:"three_d_secure"`
}

// ThreeDSecureDetails represents the 3D Secure details.
//
//easyjson:json
type ThreeDSecureDetails struct {
	Applied            bool   `json:"applied"`
	ChallengeCompleted bool   `json:"challenge_completed"`
	MethodCompleted    bool   `json:"method_completed"`
	Protocol           string `json:"protocol"`
}

// PaymentMethod represents the payment method details.
//
//easyjson:json
type PaymentMethod struct {
	Card   CardDetails `json:"card"`
	ID     string      `json:"id"`
	Saved  bool        `json:"saved"`
	Status string      `json:"status"`
	Title  string      `json:"title"`
	Type   string      `json:"type"`
}

// CardDetails represents the card details in the payment method.
//
//easyjson:json
type CardDetails struct {
	CardProduct   CardProductDetails `json:"card_product"`
	CardType      string             `json:"card_type"`
	ExpiryMonth   string             `json:"expiry_month"`
	ExpiryYear    string             `json:"expiry_year"`
	First6        string             `json:"first6"`
	IssuerCountry string             `json:"issuer_country"`
	Last4         string             `json:"last4"`
}

// CardProductDetails represents the card product details.
//
//easyjson:json
type CardProductDetails struct {
	Code string `json:"code"`
}

// Recipient represents the recipient details in the payment.
//
//easyjson:json
type Recipient struct {
	AccountID string `json:"account_id"`
	GatewayID string `json:"gateway_id"`
}
