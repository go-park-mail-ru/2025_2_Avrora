package domain

type Complex struct {
	ID            int      `json:"id"`
	Description   string   `json:"description"`
	Name          string   `json:"name"`
	Address       string   `json:"address"`
	Metro         string   `json:"metro"`
	Developer     string   `json:"developer"`
	BuiltYear     int      `json:"built_year"`
	ImageURL      []string `json:"image_url"`
	StartingPrice float64  `json:"starting_price"`
}

// комплекс в ленте
type ComplexInFeed struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	StartingPrice float64 `json:"starting_price"`
	Address       string  `json:"address"`
	Metro         string  `json:"metro"`
	ImageURL      string  `json:"image_url"`
}
type ComplexesInFeed struct {
	Complexes []ComplexInFeed `json:"complexes"`
}
