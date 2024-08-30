package generator

import (
	"math/rand"
)

type CampaignData struct {
	CampaignID      int
	CampaignType    string
	CampaignContent string
}

var campaignMap = make(map[int]CampaignData)

var campaignTypes = []string{"Awareness", "Engagement", "Conversion"}
var campaignTitles = []string{"Spring Sale", "Holiday Special", "New Arrivals", "Limited Time Offer"}

func CreateCampaign() CampaignData {
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
