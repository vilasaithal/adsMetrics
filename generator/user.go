package generator

import (
	"math/rand"
)

type UserData struct {
	UserID int
	Device string
	City   string
	Age    int
	Gender string
}

var userMap = make(map[int]UserData)

var cities = []string{"New York City", "Los Angeles", "Chicago", "Houston", "Phoenix", "New Brunswick", "Fremont", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
var devices = []string{"mobile", "desktop"}
var genders = []string{"male", "female", "other"}

func CreateUser() UserData {
	userID := rand.Intn(10000) + 1
	if user, exists := userMap[userID]; exists {
		return user
	}

	user := UserData{
		UserID: userID,
		Device: devices[rand.Intn(len(devices))],
		City:   cities[rand.Intn(len(cities))],
		Age:    rand.Intn(43) + 18, // Age between 18 and 60
		Gender: genders[rand.Intn(len(genders))],
	}
	userMap[userID] = user
	return user
}
