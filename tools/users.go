package tools

import (
	"encoding/json"
	"io"
	"os"
	"webtools/encryption"
)

var Users = make([]User, 0)

type User struct {
	Login    string
	Password encryption.PasswordObject
}

func SaveUsersToJSON() {
	//Create file
	Logger.Log(2, "Saving users.json...")
	os.Remove("users.json")
	file, err := os.Create("users.json")
	if err != nil {
		Logger.Log(3, "Error saving users.json: "+err.Error())
		return
	}

	//Convert to JSON
	data, err := json.Marshal(Users)
	if err != nil {
		Logger.Log(3, "Error marshalling users.json: "+err.Error())
		return
	}

	//Write data
	_, err = file.Write(data)
	if err != nil {
		Logger.Log(3, "Error writing to users.json: "+err.Error())
		return
	}
	Logger.Log(2, "Saving users.json done!")
}

func LoadUsersFromJSON() {
	//Open file
	Logger.Log(2, "Loading users.json...")
	file, err := os.Open("users.json")
	if err != nil {
		Logger.Log(3, "Error loading users.json: "+err.Error())
		return
	}

	//Read file
	data, err := io.ReadAll(file)
	if err != nil {
		Logger.Log(3, "Error reading users.json: "+err.Error())
		return
	}

	//Convert to data
	err = json.Unmarshal(data, &Users)
	if err != nil {
		Logger.Log(3, "Error unmarshalling users.json: "+err.Error())
		return
	}
	Logger.Log(2, "Loading users.json done!")
}
