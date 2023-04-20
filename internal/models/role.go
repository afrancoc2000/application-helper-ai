package models

type Role int

const (
	System Role = iota
	User
	Assistant
)

// String - Creating common behavior - give the type a String function
func (w Role) String() string {
	return [...]string{"system", "user", "assistant"}[w]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (w Role) EnumIndex() int {
	return int(w)
}
