package generator

import (
	"math/rand"
	"sync"
)

type CampaignData struct {
	CampaignID      int
	CampaignType    string
	CampaignContent string
}

type UserData struct {
	UserID int
	City   string
	Age    int
	Gender string
}

var (
	campaignMap = make(map[int]CampaignData)
	userMap     = make(map[int]UserData)
	mapMutex    sync.Mutex
)

var campaignTypes = []string{"Awareness", "Engagement", "Conversion"}
var campaignTitles = []string{"Spring Sale", "Holiday Special", "New Arrivals", "Limited Time Offer"}
var cities = []string{"New York City", "Los Angeles", "Chicago", "Houston", "Phoenix", "New Brunswick", "Fremont", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
var genders = []string{"male", "female", "other"}

func CreateCampaign() CampaignData {
	mapMutex.Lock() // Lock before accessing the map
	defer mapMutex.Unlock()
	campaignID := rand.Intn(50) + 1
	if campaign, exists := campaignMap[campaignID]; exists {
		return campaign
	}

	campaign := CampaignData{
		CampaignID:      campaignID,
		CampaignType:    campaignTypes[rand.Intn(len(campaignTypes))],
		CampaignContent: campaignTitles[rand.Intn(len(campaignTitles))],
	}
	campaignMap[campaignID] = campaign
	return campaign
}

func CreateUser() UserData {
	mapMutex.Lock() // Lock before accessing the map
	defer mapMutex.Unlock()
	userID := rand.Intn(10000) + 1
	if user, exists := userMap[userID]; exists {
		return user
	}

	user := UserData{
		UserID: userID,
		City:   cities[rand.Intn(len(cities))],
		Age:    rand.Intn(43) + 18, // Age between 18 and 60
		Gender: genders[rand.Intn(len(genders))],
	}
	userMap[userID] = user
	return user
}
