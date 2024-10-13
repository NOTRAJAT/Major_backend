package types

import (
	"time"
)

type KeyPair struct{
	Key string `json:"value"`
}

type StorageInterfaceUsers interface{
	GetUserByEmail(email string) (*User,error)
	GetUserByregNo(regNO string) (*User,error)
	CreateUser(*User) error
}

type StorageInterfaceAttendence interface{
	GetAttendanceByDate(date time.Time)(*[]AttendanceDisplay,error)
	MakeAttendance(regNo string, subject string)  error
}

type StorageInterfaceProduct interface{
	GetProducts()(*ProductArray,error)
	// CreateProducts()(*Product,error)
}

type Product struct{
	Id        int    `json:"id"`
	Name 		string `json:"name"`
	Description string `json:"description"`
	Image string `json:"image"`
	Price float64 `json:"price"`
	Quantity int `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

type ProductArray struct{
	Array []*Product `json:"allProducts"`
	
}

type User struct {
	RegisterNo string `json:"registerNo" validate:"required,min=9,max=64"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Passwd    string `json:"password"`
	Address   string `json:"address" validate:"required,min=8"`
	Branch    string  `json:"branch" validate:"required"`
	Year       int  `json:"year" validate:"required,numeric,min=1,max=4"`
	CreatedAt time.Time `json:"createdAt"`

}

type RegisterUserPayload struct {
	RegisterNo string `json:"registerNo" validate:"required,min=9,max=64"`
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"` 
	Email     string `json:"email" validate:"required,email"`
	Passwd    string `json:"password" validate:"required,min=8,max=64"`
	Branch    string  `json:"branch" validate:"required"`
	Year       int64  `json:"year" validate:"required,numeric,min=1,max=4"`
	Address   string `json:"address" validate:"required,min=8"`
	
}
type LoginUserPayload struct{
	Email     string `json:"email" validate:"required,email"`
	Passwd    string `json:"password" validate:"required"`
}

type AttendanceDisplay struct{
	RegisterNo string `json:"registerNo" validate:"required,min=9,max=64"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Subject string `json:"subject" validate:"required"` 
	Branch    string  `json:"branch" validate:"required"`
	Year       int  `json:"year" validate:"required,numeric,min=1,max=4"`
	AttendanceTime time.Time

}

type EspPayload struct{
	RegisterNo string `json:"registerNo" validate:"required,min=9,max=64"`
	Subject string `json:"subject" validate:"required"` 
	ApiKey string `json:"apikey" validate:"required,min=64"` 
	AttendanceTime time.Time


}

type DateQuery struct{
	Date time.Time `json:"data" validate:"required"`
}