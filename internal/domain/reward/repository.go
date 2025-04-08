package reward

type Repository interface {
	Save(*Reward) error
	FindByID(id string) (*Reward, error)
}
