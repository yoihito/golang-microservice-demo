package entities

type User struct {
	Email          string `json:"email"`
	PasswordDigest string `json:"password"`
}
