package db

import (
	"database/sql"
	"log"
)

var dbFile = "firewall.db" // Use your actual path to the database file

// GetAllRequests retrieves all rows from the requests table.
func GetAllRequests() ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Println("Failed to open database:", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM requests")
	if err != nil {
		log.Println("Failed to query database:", err)
		return nil, err
	}
	defer rows.Close()

	var requests []map[string]interface{}
	for rows.Next() {
		var id, no, length int
		var time, source, destination, protocol, info string
		if err := rows.Scan(&id, &no, &time, &source, &destination, &protocol, &length, &info); err != nil {
			log.Println("Failed to scan row:", err)
			return nil, err
		}
		request := map[string]interface{}{
			"id":          id,
			"no":          no,
			"time":        time,
			"source":      source,
			"destination": destination,
			"protocol":    protocol,
			"length":      length,
			"info":        info,
		}
		requests = append(requests, request)
	}

	return requests, nil
}
