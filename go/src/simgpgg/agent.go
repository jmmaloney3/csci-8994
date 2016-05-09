package simgpgg

type Agent struct {
  payouts float64
  cooperate bool // true for cooperators and false for defectors
}

func NewAgent(cooperate bool) *Agent {
  return &Agent { cooperate: cooperate, payouts: 0 }
}
