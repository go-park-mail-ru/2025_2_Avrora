package domain

type Complex struct {
	ID            int
	Description   string
	Name          string
	Address       string
	Metro         string
	Developer     string
	BuiltYear     int
	ImageURL      []string
	StartingPrice float64
}

// комплекс в ленте
type ComplexInFeed struct {
	ID            int
	Name          string
	StartingPrice float64
	Address       string
	Metro         string
	ImageURL      string
}
type ComplexesInFeed struct {
	Complexes []ComplexInFeed
}
