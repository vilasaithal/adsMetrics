package models

type CampaignData struct {
	CampaignID      int
	CampaignType    string
	CampaignContent string
}

type CombinedData struct { // UserEvent
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

type BaseResp struct {
	StatusCode int32  `json:"code"`
	Message    string `json:"message"`
}

type UserEventsGenerateResp struct {
	UserEventsCount int32    `json:"user_event_count"`
	BaseResp        BaseResp `json:"base_resp`
}
