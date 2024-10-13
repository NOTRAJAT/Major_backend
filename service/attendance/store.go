package attendance

import (
	"database/sql"
	"fmt"
	"myAttendance/types"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s* Store)GetAttendanceByDate(date time.Time)(*[]types.AttendanceDisplay,error){
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// SQL query to join attendance and users
	query := `
		SELECT 
			a.registerNo, 
			u.firstName, 
			u.lastName, 
			a.subject, 
			u.branch, 
			u.year, 
			a.createdAt 
		FROM attendance a 
		JOIN users u ON a.registerNo = u.registerNo
		WHERE a.createdAt BETWEEN ? AND ?
	`

	// Execute the query
	rows, err := s.db.Query(query, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("error querying attendance: %w", err)
	}
	defer rows.Close()

	// Slice to hold the results
	var attendanceList []types.AttendanceDisplay

	// Iterate over the rows
	for rows.Next() {
		var attendance types.AttendanceDisplay
		if err := rows.Scan(
			&attendance.RegisterNo,
			&attendance.FirstName,
			&attendance.LastName,
			&attendance.Subject,
			&attendance.Branch,
			&attendance.Year,
			&attendance.AttendanceTime,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		attendanceList = append(attendanceList, attendance)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Return the attendance list
	return &attendanceList, nil
}
// MakeAttendance(regNo string, subject string)
// }

func(s* Store) MakeAttendance(regNo string, subject string) error{
	_,err := s.db.Query("INSERT INTO attendance(registerNo,subject) VALUES(?,?)",regNo,subject)
	if err!=nil{
		return err;
	}
	return nil
}


