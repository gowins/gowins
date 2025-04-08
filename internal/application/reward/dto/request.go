package dto

type CreateRewardRequest struct {
	Type        string   `json:"type"`
	Device      string   `json:"device"`
	Project     string   `json:"project"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	Deadline    int64    `json:"deadline"`
	ReviewTime  int64    `json:"review_time"`
	UnitPrice   float64  `json:"unit_price"`
	Quantity    int      `json:"quantity"`
	SingleUse   bool     `json:"single_use"`
}
