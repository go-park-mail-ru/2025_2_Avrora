package domain

//go:generate easyjson -all $GOFILE
type Err string

func (e Err) Error() string { return string(e) }
