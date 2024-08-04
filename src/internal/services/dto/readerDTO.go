package dto

type ReaderSignInDTO struct {
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Password    string `json:"password" db:"password"`
}

type ReaderSignUpDTO struct {
	Fio         string `json:"fio" db:"fio"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Age         uint   `json:"age" db:"age"`
	Password    string `json:"password" db:"password"`
}
