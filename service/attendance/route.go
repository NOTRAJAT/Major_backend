package attendance

import (
	"errors"
	"myAttendance/config"
	"myAttendance/service/auth"
	"myAttendance/types"
	"myAttendance/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	StoreInstanceAttendance types.StorageInterfaceAttendence //for database
	StoreInstanceUser types.StorageInterfaceUsers //for database
}

func NewHandler(StoreInstanceAttendance types.StorageInterfaceAttendence,	StoreInstanceUser types.StorageInterfaceUsers) *Handler {
	return &Handler{
		StoreInstanceAttendance: StoreInstanceAttendance,
		StoreInstanceUser: StoreInstanceUser ,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/espPoint", h.addAttendance).Methods("POST")
	router.HandleFunc("/getAttendanceByData",auth.ValidateUser( h.getAttendanceDate,h.StoreInstanceUser)).Methods("GET")


}

func (h *Handler) addAttendance(w http.ResponseWriter,r *http.Request){
	SpecialHashEsp:= config.ENV.SpecialHashEsp
	payload:= new(types.EspPayload)
//payload parsing 
if err := utils.ParseJson[types.EspPayload](r,payload);err!=nil{
	utils.HandleErrorLog(w,http.StatusBadRequest,err)
	return
}
//payload validation
if err:= utils.Validator.Struct(payload);err!=nil{
	utils.HandleErrorLog(w,http.StatusBadRequest,err)
	return
}
//checking api key

if SpecialHashEsp!=payload.ApiKey{
	err := errors.New("bad Auth")


	utils.HandleErrorLog(w,http.StatusBadRequest,err)
	return
}


err:=h.StoreInstanceAttendance.MakeAttendance(payload.RegisterNo,payload.Subject)
if err!=nil{
	utils.HandleErrorLog(w,http.StatusBadRequest,err)
	return
}
//


}

func(h * Handler) getAttendanceDate(w http.ResponseWriter,r *http.Request){
	payload:= new(types.DateQuery)
	//payload parsing 
	if err := utils.ParseJson[types.DateQuery](r,payload);err!=nil{
		utils.HandleErrorLog(w,http.StatusBadRequest,err)
		return
	}
	//payload validation
	if err:= utils.Validator.Struct(payload);err!=nil{
		utils.HandleErrorLog(w,http.StatusBadRequest,err)
		return
	}

	AllEntries,err:= h.StoreInstanceAttendance.GetAttendanceByDate(payload.Date)
	if err!=nil{
		utils.HandleErrorLog(w,http.StatusBadRequest,err)
	return
	}
	utils.ClientRequestLog(r)
	utils.JsonWriter(w,200,AllEntries)
}

