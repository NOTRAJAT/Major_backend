package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Hasher(pass string) (string,error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass),bcrypt.DefaultCost)

	if err!=nil{
		return "" ,err
	}

	return string(hash), err
}

func ComparePasswords(hashed string,plain string)bool{
	err:= bcrypt.CompareHashAndPassword([]byte(hashed),[]byte(plain))
	log.Println(err)
	return err==nil
}