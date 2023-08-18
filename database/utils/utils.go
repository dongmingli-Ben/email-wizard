package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

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

func get_primary_keys(db *sql.DB, table string) ([]string, error) {
	query := fmt.Sprintf(`
		SELECT a.attname AS column_name
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = '%s'::regclass
			AND i.indisprimary;
	`, table)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pk_cols []string

	for rows.Next() {
		var col_name string
		err := rows.Scan(&col_name)
		if err != nil {
			return nil, err
		}
		pk_cols = append(pk_cols, col_name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pk_cols, nil
}

func prepare_insert_query(db *sql.DB, row map[string]interface{}, table string) (string, []interface{}, []string, error) {
	column_info, err := get_column_name_type(db, table)
	if err != nil {
		return "", nil, nil, err
	}
	pk_cols, err := get_primary_keys(db, table)
	if err != nil {
		return "", nil, nil, err
	}
	values := make([]interface{}, 0)
	query := fmt.Sprintf("INSERT INTO %s (", table)
	idx := 0
	for col_name, value := range row {
		col_type, ok := column_info[col_name]
		if !ok {
			return "", nil, nil, fmt.Errorf("unrecognized column %s in row", col_name)
		}
		if col_type == "ARRAY" {
			values = append(values, pq.Array(value))
		} else if col_type == "json" {
			json_str, err := json.Marshal(value)
			if err != nil {
				return "", nil, nil, err
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
	for i := 1; i <= len(row); i++ {
		if i != 1 {
			query += ", "
		}
		query += fmt.Sprintf("$%d", i)
	}
	query += ") RETURNING " + strings.Join(pk_cols, ", ")
	return query, values, pk_cols, nil
}

func AddRow(row map[string]interface{}, table string) (map[string]interface{}, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	insert_query, values, pk_cols, err := prepare_insert_query(db, row, table)
	if err != nil {
		return nil, err
	}
	pk_values := make([]interface{}, len(pk_cols))
	for i := 0; i < len(pk_cols); i++ {
		pk_values[i] = new(interface{})
	}
	err = db.QueryRow(insert_query, values...).Scan(pk_values...)
	if err != nil {
		return nil, err
	}
	pk_vals := make(map[string]interface{})
	for i := 0; i < len(pk_cols); i++ {
		pk_vals[pk_cols[i]] = *pk_values[i].(*interface{})
	}
	return pk_vals, nil
}

func prepare_update_query(db *sql.DB, column string, value interface{},
	condition map[string]interface{}, table string) (string, []interface{}, error) {
	column_info, err := get_column_name_type(db, table)
	if err != nil {
		return "", nil, err
	}
	query := fmt.Sprintf("UPDATE %s SET %s = $%d WHERE ", table, column, len(condition)+1)
	idx := 1
	values := make([]interface{}, 0)
	for col, val := range condition {
		if idx > 1 {
			query += " AND "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		values = append(values, val)
		idx++
	}
	if col_type, ok := column_info[column]; !ok {
		return "", nil, fmt.Errorf("%s not in table %s", column, table)
	} else if col_type == "json" {
		json_str, err := json.Marshal(value)
		if err != nil {
			return "", nil, err
		}
		values = append(values, json_str)
	} else if col_type == "ARRAY" {
		values = append(values, pq.Array(value))
	} else {
		values = append(values, value)
	}
	return query, values, nil
}

// do not support complex type condition yet
func UpdateValue(column string, value interface{}, condition map[string]interface{}, table string) error {
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	update_query, values, err := prepare_update_query(db, column, value, condition, table)
	if err != nil {
		return err
	}
	_, err = db.Exec(update_query, values...)
	if err != nil {
		return err
	}
	return nil
}

func prepare_delete_query(db *sql.DB, condition map[string]interface{}, table string) (string, []interface{}, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE ", table)
	idx := 1
	values := make([]interface{}, 0)
	for col, val := range condition {
		if idx > 1 {
			query += " AND "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		values = append(values, val)
		idx++
	}
	return query, values, nil
}

// do not support complex type condition yet
func DeleteRows(condition map[string]interface{}, table string) error {
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	delete_query, values, err := prepare_delete_query(db, condition, table)
	if err != nil {
		return err
	}
	_, err = db.Exec(delete_query, values...)
	if err != nil {
		return err
	}
	return nil
}

func prepare_select_query(db *sql.DB, columns []string, condition map[string]interface{},
	table string) (string, []interface{}, error) {
	column_info, err := get_column_name_type(db, table)
	if err != nil {
		return "", nil, err
	}
	ready_cols := make([]string, len(columns))
	copy(ready_cols, columns)
	for i := 0; i < len(ready_cols); i++ {
		if col_type, ok := column_info[ready_cols[i]]; ok && col_type == "ARRAY" {
			ready_cols[i] = fmt.Sprintf("array_to_json(%s)", ready_cols[i])
		}
	}
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(ready_cols, ", "), table)
	if len(condition) == 0 {
		return query, nil, nil
	}
	query += " WHERE "
	idx := 0
	values := make([]interface{}, 0)
	for key, val := range condition {
		if idx > 0 {
			query += " AND "
		}
		query += fmt.Sprintf("%s = $%d", key, idx+1)
		values = append(values, val)
		idx++
	}
	return query, values, nil
}

// do not support complex fields for columns yet
func Query(columns []string, condition map[string]interface{}, table string) ([]map[string]interface{}, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	column_info, err := get_column_name_type(db, table)
	if err != nil {
		return nil, err
	}
	query, values, err := prepare_select_query(db, columns, condition, table)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		record := make(map[string]interface{})
		values = make([]interface{}, len(columns))
		for i := 0; i < len(values); i++ {
			values[i] = new(interface{})
		}
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		for i := 0; i < len(columns); i++ {
			col := columns[i]
			val := *values[i].(*interface{})
			if column_info[col] == "json" {
				val_json := make(map[string]interface{})
				// fmt.Println(string(val.([]byte)))
				if err := json.Unmarshal(val.([]byte), &val_json); err != nil {
					arr_json := make([]map[string]interface{}, 0)
					if err := json.Unmarshal(val.([]byte), &arr_json); err != nil {
						return nil, err
					}
					record[col] = arr_json
				} else {
					record[col] = val_json
				}
			} else if column_info[col] == "ARRAY" {
				arr_json := make([]interface{}, 0)
				if err := json.Unmarshal(val.([]byte), &arr_json); err != nil {
					return nil, err
				}
				record[col] = arr_json
			} else {
				record[col] = val
			}
		}
		results = append(results, record)
	}
	return results, nil
}
