package reward

type Reward struct {
	ID          string
	Type        string
	Device      string
	Project     string
	Title       string
	Description string
	Steps       []string
	Deadline    int64
	ReviewTime  int64
	UnitPrice   float64
	Quantity    int
	SingleUse   bool
	TotalAmount float64
}

func (r *Reward) CalculateTotal() {
	r.TotalAmount = r.UnitPrice * float64(r.Quantity)
}
