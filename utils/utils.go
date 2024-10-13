package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"myAttendance/types"
	"net/http"

	"github.com/go-playground/validator/v10"
)


var Validator = validator.New()

func ParseJson[T types.RegisterUserPayload | types.LoginUserPayload |types.EspPayload |types.DateQuery](r *http.Request,
structure *T)(error){
	if err:=json.NewDecoder(r.Body).Decode(structure);err!=nil{
		if(err.Error()=="EOF" || r.Body==nil){
			return fmt.Errorf("missing body")
		}
		return err
	}
	
	return nil
}

func JsonWriter[T types.RegisterUserPayload  | types.LoginUserPayload | types.KeyPair |types.EspPayload|[]types.AttendanceDisplay](w http.ResponseWriter, status int,structure *T)(error){
	w.WriteHeader(status)
	w.Header().Add("content-type","application/json")
	return json.NewEncoder(w).Encode(structure)
}

func HandleErrorLog(w http.ResponseWriter,status int,err error){
	log.Printf("error -> %s",err.Error())
	w.WriteHeader(status)
	w.Write([]byte(""))
	
}

func ClientRequestLog(r *http.Request){
	log.Printf("[%s] %s -> %s \n",r.Method,r.RemoteAddr,r.URL)
}
