package user

import (
	"database/sql"
	"fmt"
	"myAttendance/types"
	"myAttendance/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store{
	return &Store{	db :db,}
}

func (s *Store) GetUserByEmail(email string) (*types.User,error){
	rows,err := s.db.Query("SELECT * FROM users WHERE email = ?",email)
	if err!=nil{
		return nil, err;
	}
	u := new(types.User)

	for rows.Next(){
		u,err= scanRowsIntoUsers(rows)

	} 
	// log.Println(u)
	if err!=nil{
		return nil, err;
	}

	 if err:= utils.Validator.Struct(u);err!=nil {
		return nil, fmt.Errorf("user not found")
	}

	return u, err;
}

func scanRowsIntoUsers(rows *sql.Rows) (*types.User, error){
		user:= new(types.User)
		if err:= rows.Scan(
			&user.RegisterNo,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Passwd,
			&user.Address,
			&user.Branch,
			&user.Year,
			&user.CreatedAt,
		);err!=nil{
			return nil,err
		}
		return user,nil
}

func (s* Store) GetUserByregNo(regNo string) (*types.User,error){

	rows,err := s.db.Query("SELECT * FROM users WHERE registerNo  = ?",regNo)
	if err!=nil{
		return nil, err;
	}
	u := new(types.User)

	for rows.Next(){
		u,err= scanRowsIntoUsers(rows)

	} 
	if err!=nil{
		return nil, err;
	}

	if err:= utils.Validator.Struct(u);err!=nil {
		return nil, fmt.Errorf("user not found")
	}


	return nil, err;
}
func (s* Store)	CreateUser(val *types.User) error{

	_,err := s.db.Query("INSERT INTO users(registerNo,firstName,lastName,email,password,address,branch,year) VALUES(?,?,?,?,?,?,?,?)",val.RegisterNo,val.FirstName,val.LastName,val.Email,val.Passwd,val.Address,val.Branch,val.Year)
	if err!=nil{
		return err;
	}
	return nil
}

