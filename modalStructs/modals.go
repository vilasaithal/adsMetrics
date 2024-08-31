package modalstructs

type CampaignData struct {
	CampaignID      int
	CampaignType    string
	CampaignContent string
}

type CombinedData struct {
	CampaignID      int    `json:"campaign_id"`
	CampaignType    string `json:"campaign_type"`
	CampaignContent string `json:"campaign_content"`
	UserID          int    `json:"user_id"`
	Device          string `json:"device"`
	City            string `json:"city"`
	Age             int    `json:"age"`
	Gender          string `json:"gender"`
	EventID         string `json:"event_id"`
	Timestamp       int64  `json:"timestamp"` // Unix timestamp in seconds
	EventType       string `json:"event_type"`
}

type UserData struct {
	UserID int
	City   string
	Age    int
	Gender string
}
