package usecase

type CreateOfferRequest struct {
	UserID      int     `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Image       string  `json:"image,omitempty"`
	Price       int     `json:"price"`
	Area        float64 `json:"area,omitempty"`
	Rooms       int     `json:"rooms,omitempty"`
	Address     string  `json:"address"`
	OfferType   string  `json:"offer_type"`
}

type OfferResponse struct {
	UserID      int     `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Image       string  `json:"image,omitempty"`
	Price       int     `json:"price"`
	Area        float64 `json:"area,omitempty"`
	Rooms       int     `json:"rooms,omitempty"`
	Address     string  `json:"address"`
	OfferType   string  `json:"offer_type"`
}

type UpdateOfferRequest struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Image       string  `json:"image,omitempty"`
	Price       int     `json:"price"`
	Area        float64 `json:"area,omitempty"`
	Rooms       int     `json:"rooms,omitempty"`
	Address     string  `json:"address"`
	OfferType   string  `json:"offer_type"`
}
