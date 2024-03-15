package models

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Age       int
	Hobbies   string
	City      int
}

type Profile struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Age       int
	Hobbies   string
	City      string
}

type UserRegistration struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type CredentialsForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Credentials struct {
	ID       int64
	Email    string
	Password string
}
