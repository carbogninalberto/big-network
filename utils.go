package main

import (
	"log"
	"math"
	"math/rand"
	"time"
)

// reduceInfectiveEpochs returns true if person is no more infective
func reduceInfectiveEpochs(personPointer *person) bool {
	if personPointer.InfectiveEpochs > 1 {
		personPointer.InfectiveEpochs--
	} else if personPointer.InfectiveEpochs == 1 {
		personPointer.InfectiveEpochs--
		personPointer.Infective = false
		// at the end of the illness make it die with a certain probability
		if bernoulli(deadRate) {
			personPointer.Dead = true
		} else {
			personPointer.Survived = true
		}
		return true
	} else {
		log.Println("ERROR", personPointer.InfectiveEpochs)
	}
	return false
}

// getInfected return a list of infected persons
func getInfected(networkPointer *bigNet) []int {
	// infected is the returned slice
	infected := make([]int, 0, 1)
	for node := 0; node < nNodes; node++ {
		// this condition allows to return only the actual infected people
		if (*networkPointer)[node].Infective == true &&
			(*networkPointer)[node].Survived == false &&
			(*networkPointer)[node].Dead == false {
			infected = append(infected, node)
		}
	}
	return infected
}

// countInfected counts the number of infected people
// that satisfy the logic formula infective && survived && dead
func countInfected(networkPointer *bigNet, infective, survived, dead bool) int {
	counter := 0
	for node := 0; node < nNodes; node++ {
		if (*networkPointer)[node].Infective == infective &&
			(*networkPointer)[node].Survived == survived &&
			(*networkPointer)[node].Dead == dead {
			counter++
		}
	}
	return counter
}

// countTotalInfected counts the total number of infected people
func countTotalInfected(networkPointer *bigNet) int {
	counter := 0
	for node := 0; node < nNodes; node++ {
		if (*networkPointer)[node].Infective == true ||
			(*networkPointer)[node].Survived == true ||
			(*networkPointer)[node].Dead == true {
			counter++
		}
	}
	return counter
}

// resetNetwork reset the graph to random initial values without touching the shape
func resetNetwork(networkPointer *bigNet) {
	for node := 0; node < nNodes; node++ {
		// random seed
		rand.Seed(time.Now().UnixNano())
		// generate a random RO following a Normal Distribution
		tmpR0 := int(math.Round(rand.NormFloat64()*stdR0 + medianR0))
		if tmpR0 < 0 {
			tmpR0 = 0
		} // check if has been generated a negative number

		// generating infective days array
		infectiveDays := make([]int8, tmpR0)
		for r := 0; r < tmpR0; r++ {
			infectiveDays[r] = int8(rand.Intn(infectiveEpochs))
		} // infect tmpR0 people during the infectiveEpochs period

		// init parameters
		(*networkPointer)[node].Infective = false
		(*networkPointer)[node].Survived = false
		(*networkPointer)[node].Dead = false
		(*networkPointer)[node].InfectiveEpochs = infectiveEpochs
		(*networkPointer)[node].InfectiveDays = infectiveDays
	}
}

// bernoulli is a simple implementation that returns true with a certain pSuccess probability
func bernoulli(pSuccess float64) bool {
	if rand.Intn(10000) < int(pSuccess*10000) {
		return true
	}
	return false
}
