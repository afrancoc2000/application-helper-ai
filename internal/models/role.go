package models

type Role int

const (
	System Role = iota
	User
	Assistant
)

func (r Role) String() string {
	return [...]string{"system", "user", "assistant"}[r]
}
