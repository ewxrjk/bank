package bank

import (
	"database/sql"
	"net/http"

	"github.com/ewxrjk/bank/util"
)

var defaultConfig = map[string]string{
	"houseAccount": "house",
	"title":        "Test Bank",
}

// ErrNoSuchConfig is return when a nonexistent configuration item is accessed.
var ErrNoSuchConfig = util.HTTPError{"no such configuration item", http.StatusNotFound}

// GetConfigs returns the full configuration table.
func GetConfigs(tx *sql.Tx) (config map[string]string, err error) {
	var rows *sql.Rows
	if rows, err = tx.Query("SELECT key,value FROM config"); err != nil {
		return
	}
	defer rows.Close()
	config = map[string]string{}
	for key, value := range defaultConfig {
		config[key] = value
	}
	for rows.Next() {
		var key, value string
		if err = rows.Scan(&key, &value); err != nil {
			return
		}
		config[key] = value
	}
	return
}

// GetConfig returns the value of a configuration item.
func GetConfig(tx *sql.Tx, key string) (value string, err error) {
	err = tx.QueryRow("SELECT value FROM config WHERE key=?", key).Scan(&value)
	if err == sql.ErrNoRows {
		var ok bool
		if value, ok = defaultConfig[key]; ok {
			err = nil
		} else {
			err = ErrNoSuchConfig
		}
	}
	if err != nil {
		return
	}
	return
}

// PutConfig updates or stores a configuration item.
func PutConfig(tx *sql.Tx, key, value string) (err error) {
	if _, err = tx.Exec("INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)", key, value); err != nil {
		return
	}
	return
}
