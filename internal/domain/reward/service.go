package reward

func NewReward(id string, r *Reward) *Reward {
	r.ID = id
	r.CalculateTotal()
	return r
}
