package utils

import (
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Port         string `json:"port"`
	DBConnection string `json:"db_connection"`
}

func LoadConfig(defaultFile string) (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = defaultFile
	}

	fmt.Printf("Attempting to open config file: %s\n", configFile)

	file, err := os.Open(configFile)
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Printf("Error opening config file: %v (cwd: %s)\n", err, cwd)
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(config); err != nil {
		fmt.Printf("Error decoding config file: %v\n", err)
		return nil, err
	}

	return config, nil
}

func ConnectDB(connString string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlserver", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
