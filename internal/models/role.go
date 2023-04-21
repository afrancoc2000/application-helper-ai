package models

type Role int

const (
	System Role = iota
	User
	Assistant
)

// String - Creating common behavior - give the type a String function
func (r Role) String() string {
	return [...]string{"system", "user", "assistant"}[r]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (r Role) EnumIndex() int {
	return int(r)
}
