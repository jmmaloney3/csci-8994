package sim

import "fmt"

type Tribe struct {
  agents []Agent
  numAgents int
  totalPayouts int32
}

func MakeTribe(numAgents int) Tribe {
  agents := make([]Agent, numAgents)
  // create agents
  for i := 0; i < numAgents; i++ {
    fmt.Println("    make agent ", i)
    agents[i] = MakeAgent()
  }
  return Tribe { agents: agents, numAgents: numAgents, totalPayouts: 0 }
}
