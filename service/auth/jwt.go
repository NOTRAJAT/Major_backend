package auth

import (
	"fmt"
	"log"
	"myAttendance/config"
	"myAttendance/types"
	"myAttendance/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateJWT(secret []byte,email string) (string, error) {
	expiration := time.Second * time.Duration(config.ENV.JwtExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"userID":(email),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString,err:= token.SignedString(secret)
	if err!=nil{
		return "",err
	}

	return tokenString,nil
}

func ValidateUser(handlerFunc http.HandlerFunc,UserStore types.StorageInterfaceUsers)(http.HandlerFunc){
	return func (w http.ResponseWriter,r*http.Request)  {
		// token1,err:= getHeaderAuth(r)
		token, err := r.Cookie(config.ENV.HeaderAuthName)
		// log.Println(token)
		// log.Println(token1)
		if(err!=nil){
			authFailed(w,r,err,err.Error())	
			return
		}
		jwtToken,err:= validateToken((token.Value))
		if err!=nil{
			authFailed(w,r,err,"Auth failed")	
			return
		}

		if !jwtToken.Valid{
			authFailed(w,r,err,"Auth failed")	
			return
		}

		claims:= jwtToken.Claims.(jwt.MapClaims)
		nowTime:= time.Now().Unix()

		Expiration:= claims["expiredAt"].(float64)
		

		if(nowTime>int64(Expiration)){
			authFailed(w,r,fmt.Errorf("token expired"),"auth failed")
			return
		}
		


		UserId:= claims["userID"].(string)
		// id,_:= strconv.ParseInt(UserId,10,32)
		_,err=UserStore.GetUserByregNo(UserId)
		if err!=nil{
			authFailed(w,r,err,err.Error())
			return
		}
		newToken,err:=CreateJWT([]byte(config.ENV.JwtSecret),UserId)
		if err!=nil{
			authFailed(w,r,err,err.Error())
			return
		}

		SendCookie(config.ENV.HeaderAuthName,newToken,w)

		handlerFunc(w,r)
	}	
}


func getHeaderAuth(r* http.Request)(string,error){
	token:= r.Header.Get(config.ENV.HeaderAuthName)
	if token==""{
		return "",fmt.Errorf("auth header not set")
	}

	return token,nil
}

func validateToken(t string)( *jwt.Token,error){
	return jwt.Parse(t,func(sign *jwt.Token) (interface{}, error){
				if _,ok:=sign.Method.(*jwt.SigningMethodHMAC);!ok{
					return nil,fmt.Errorf("unkown sign Method %s",sign.Header)
				}
			return []byte(config.ENV.JwtSecret),nil
	})

}
func authFailed(w http.ResponseWriter,r*http.Request,err error,clientError string){
	utils.ClientRequestLog(r)
	log.Println(err.Error())
	utils.JsonWriter[types.KeyPair](w,http.StatusBadRequest,&types.KeyPair{Key: clientError})		
}

func SendCookie(TokenName string ,TokenValue string ,w http.ResponseWriter ){
	expiration := time.Second * time.Duration(config.ENV.JwtExpirationInSeconds)
	http.SetCookie(w, &http.Cookie{
		Name:    TokenName,
		Value:   TokenValue,
		Expires:  time.Now().Add(expiration), // Token expiration
		HttpOnly: true,                          // Prevent access by JavaScript
		Secure:   true,                          // Only sent over HTTPS
		// SameSite: http.SameSiteStrictMode,       // Prevent CSRF by only sending the cookie in same-site requests
		// Domain: config.ENV.CORS,//Forntend
	})
}
	
	func HashPassword(password string) (string, error) {
		// Generate the bcrypt hash of the password
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
	
		return string(hash), nil
	}
	
	// AddUserToMosquittoFile appends the hashed user and password to the Mosquitto password file
	func AddUserToMosquittoFile(username, hashedPassword string) error {
		// Open the password file in append mode
		file, err := os.OpenFile("/etc/mosquitto/passwd", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer file.Close()
	
		// Format the entry: username:$6$[bcrypt-hashed-password]
		entry := fmt.Sprintf("%s:$6$%s\n", username, strings.TrimSpace(hashedPassword))
	
		// Write the new entry to the password file
		if _, err := file.WriteString(entry); err != nil {
			return err
		}
	
		return nil
	}