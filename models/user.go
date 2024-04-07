package models

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Age       int
	Hobbies   string
	City      int
	Password  string
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

func UserRegistrationToUser(regForm *UserRegistration) *User {
	return &User{
		Email:     regForm.Email,
		Password:  regForm.Password,
		FirstName: regForm.FirstName,
		LastName:  regForm.LastName,
	}

}

func ProfileToUser(user *Profile) *User {
	return &User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
	}
}

func UserToProfile(user User) Profile {
	return Profile{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
	}
}

//func UserRegistrationToProfile(regForm *UserRegistration) *Profile {
//	return &Profile{
//		Email:     regForm.Email,
//		FirstName: regForm.FirstName,
//		LastName:  regForm.LastName,
//	}
//}
