package user

import (

	// "ecom/service/auth"

	"myAttendance/config"
	"myAttendance/service/auth"
	"myAttendance/types"
	"myAttendance/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	StoreInstance types.StorageInterfaceUsers //for database 
}

func NewHandler(StoreInstance types.StorageInterfaceUsers) *Handler {
	return &Handler{
		StoreInstance: StoreInstance,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router){
	(router.HandleFunc("/login",(h.handleLogin)).Methods("POST"))
	router.HandleFunc("/register",(h.handleRegister)).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter,r *http.Request){
	var payload types.LoginUserPayload
	//payload parsing 
	if err := utils.ParseJson[types.LoginUserPayload](r,&payload);err!=nil{
		utils.HandleErrorLog(w,http.StatusBadRequest,err)
		return
	}
	//payload validation
	if err:= utils.Validator.Struct(payload);err!=nil{
		utils.HandleErrorLog(w,http.StatusBadRequest,err)
		return
	}
	// log.Println(payload)

	//checks if the user exist by email
	usersStorage,err := h.StoreInstance.GetUserByEmail(payload.Email)

	if err!=nil{
		// log.Println("email")
		utils.JsonWriter[types.KeyPair](w,http.StatusBadGateway,&types.KeyPair{Key: "invalid email or password"})
		return
	}

	//password check

	if(!auth.ComparePasswords(usersStorage.Passwd,payload.Passwd)){
		// log.Println("pass")
		utils.JsonWriter[types.KeyPair](w,http.StatusBadGateway,&types.KeyPair{Key: "invalid email or password"})
		return
	}

	// jwt token
	token,err:= auth.CreateJWT([]byte(config.ENV.JwtSecret),usersStorage.Email)
	if (err!=nil){
	utils.HandleErrorLog(w,http.StatusInternalServerError,err);
		return
	}

	auth.SendCookie(config.ENV.HeaderAuthName,token,w)

	utils.ClientRequestLog(r)
	utils.JsonWriter[types.KeyPair](w,http.StatusOK,&types.KeyPair{Key: token})

}
func (h *Handler) handleRegister(w http.ResponseWriter,r *http.Request){
	var payload types.RegisterUserPayload
	//payload parsing 
	if err:= utils.ParseJson[types.RegisterUserPayload](r,&payload); err!=nil{
		utils.HandleErrorLog(w,http.StatusBadRequest,err)
		return
	}
	// validate payload

	if err:=utils.Validator.Struct(payload); err!=nil{
		utils.HandleErrorLog(w,http.StatusBadRequest,err)
		return
	}

		//checking existing user
	_,err:=h.StoreInstance.GetUserByEmail(payload.Email)
	if err==nil{
		utils.ClientRequestLog(r)
		utils.JsonWriter[types.KeyPair](w,http.StatusOK,&types.KeyPair{
			Key: "User Already Exists",
		})
		return
	}
		//checking existing user
	_,err=h.StoreInstance.GetUserByregNo(payload.RegisterNo)
	if err==nil{
		utils.ClientRequestLog(r)
		utils.JsonWriter[types.KeyPair](w,http.StatusOK,&types.KeyPair{
			Key: "User Already Exists",
		})
		return
	}
	//hasing plain passsword
	hashPassword,err := auth.Hasher(payload.Passwd)
	if err!=nil{
		utils.HandleErrorLog(w,404,err)
		return
	}
	//adding user to table
		err = h.StoreInstance.CreateUser(&types.User{
		RegisterNo: payload.RegisterNo,
		Email: payload.Email,
		FirstName: payload.FirstName,
		LastName: payload.LastName,
		Address: payload.Address,
		Passwd: hashPassword,
		Branch: payload.Branch,
	Year: int(payload.Year),})

	if err!=nil{
			utils.HandleErrorLog(w,http.StatusBadRequest,err)
			return
	}

	utils.ClientRequestLog(r)

	utils.JsonWriter[types.KeyPair](w,http.StatusCreated,&types.KeyPair{
		Key: "User Created",
	})

		

	

	payload.Passwd=""
	utils.JsonWriter(w,200,&payload)
}