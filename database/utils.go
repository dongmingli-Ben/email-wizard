package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
	host := "postgres"
	user := "postgres"
	password := "123456"
	dbname := "email-wizard-data"

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db, err
}

func get_column_name_type(db *sql.DB, table string) (map[string]string, error) {
	query := `
        SELECT column_name, data_type 
        FROM information_schema.columns 
        WHERE table_name = $1
        ORDER BY ordinal_position;
    `

	rows, err := db.Query(query, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	column_info := make(map[string]string)
	for rows.Next() {
		var col_name, col_type string
		if err := rows.Scan(&col_name, &col_type); err != nil {
			return nil, err
		}
		column_info[col_name] = col_type
	}
	return column_info, nil
}

func prepare_insert_query(db *sql.DB, row map[string]interface{}, table string) (string, []interface{}, error) {
	column_info, err := get_column_name_type(db, table)
	if err != nil {
		return "", nil, err
	}
	values := make([]interface{}, 0)
	query := fmt.Sprintf("INSERT INTO %s (", table)
	idx := 0
	for col_name, col_type := range column_info {
		value, ok := row[col_name]
		if !ok {
			return "", nil, fmt.Errorf("missing column %s in row", col_name)
		}
		if col_type == "ARRAY" {
			values = append(values, pq.Array(value))
		} else if col_type == "json" {
			json_str, err := json.Marshal(value)
			if err != nil {
				return "", nil, err
			}
			values = append(values, json_str)
		} else {
			values = append(values, value)
		}
		if idx > 0 {
			query += ", "
		}
		query += col_name
		idx++
	}
	query += ") VALUES ("
	for i := 1; i <= len(column_info); i++ {
		if i != 1 {
			query += ", "
		}
		query += fmt.Sprintf("$%d", i)
	}
	query += ")"
	return query, values, nil
}

func addRow(row map[string]interface{}, table string) error {
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	insert_query, values, err := prepare_insert_query(db, row, table)
	if err != nil {
		return err
	}
	_, err = db.Exec(insert_query, values...)
	if err != nil {
		return err
	}
	return nil
}
