package db

import (
	"encoding/json"
	"iin_check/internal/models"
	"io/ioutil"
	"os"
)

var dbFile = "database.json"

func LoadDatabase() ([]models.Person, error) {
	var people []models.Person
	data, err := ioutil.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return people, nil
		}
		return nil, err
	}
	err = json.Unmarshal(data, &people)
	if err != nil {
		return nil, err
	}
	return people, nil
}

func SaveDatabase(people []models.Person) error {
	data, err := json.MarshalIndent(people, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dbFile, data, 0644)
}
