package healthz

const (
	StateUp       float64 = 0
	StateUnknown  float64 = 1
	StateDegraded float64 = 2
	StateDown     float64 = 3
)

type Checker struct {
}

func NewChecker() *Checker {
	return &Checker{}
}
