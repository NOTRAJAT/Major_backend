package api

import (
	"database/sql"

	// "ecom/service/product"
	"log"
	"myAttendance/config"
	"myAttendance/service/attendance"
	"myAttendance/service/user"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer{
	return &APIServer{
		addr: addr,
		db:db,
	}	
}
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow CORS from the specified frontend origin
		w.Header().Set("Access-Control-Allow-Origin",config.ENV.CORS )

		// Add more CORS headers if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}


func (s *APIServer) Run()  error {
	router:= mux.NewRouter()
	subRouter:= router.PathPrefix("/api/v1").Subrouter() // /reister
	// subRouter1:= router.PathPrefix("/api/v2").Subrouter() // /reister
	
	UserMysqlStore:=user.NewStore(s.db)
	// ProductMysqlStore:=product.NewStore(s.db)
	AttendanceMysqlStore:=attendance.NewStore(s.db)
	

	userHandleService:= user.NewHandler(UserMysqlStore)
	userHandleService.RegisterRoutes(subRouter)

	attendanceHandlerServcie:=attendance.NewHandler(AttendanceMysqlStore,UserMysqlStore)
	attendanceHandlerServcie.RegisterRoutes(subRouter)

	// productHandleService:=product.NewHandler(ProductMysqlStore,UserMysqlStore)
	// productHandleService.RegisterRoutes(subRouter)
  // Set frontend origin (from environment variable or fallback to default)
	


	// Wrap the router with the CORS middleware
	corsWrappedRouter := enableCORS(router)

	log.Printf("Listining on address %s",s.addr)
	return http.ListenAndServe(s.addr,corsWrappedRouter)
}