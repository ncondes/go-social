package handlers

type Handlers struct {
	Health *Health
}

func New() *Handlers {
	return &Handlers{
		Health: &Health{},
	}
}
