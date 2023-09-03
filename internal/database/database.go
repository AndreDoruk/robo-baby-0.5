package database

import (
	"encoding/json"
	"log"
	"os"
)

func LoadJson(filepath string, database any) error {
	data, err := os.ReadFile(filepath)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &database)

	return err
}

func SaveJson(filepath string, database any) {
	bytes, err := json.Marshal(database)

	if err != nil {
		log.Fatalln(err)
	}
	os.WriteFile(filepath, bytes, 0666)
}
