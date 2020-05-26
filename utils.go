package main

import (
	"log"
	"math"
	"math/rand"
	"time"
)

func reduceInfectiveEpochs(personPointer *person) bool {
	//log.Println("reduceInfective", personPointer.Infective, personPointer.Survived, personPointer.InfectiveEpochs)
	if personPointer.InfectiveEpochs > 1 {
		personPointer.InfectiveEpochs--
		//log.Println("personPointer.InfectiveEpochs > 1", personPointer.Infective, personPointer.Survived, personPointer.InfectiveEpochs)
	} else if personPointer.InfectiveEpochs == 1 {
		personPointer.InfectiveEpochs--
		personPointer.Infective = false
		if bernoulli(deadRate) {
			personPointer.Dead = true
		} else {
			personPointer.Survived = true
		}
		return true
		//log.Println("personPointer.InfectiveEpochs == 1", personPointer.Infective, personPointer.Survived, personPointer.InfectiveEpochs)
	} else {
		log.Panicln("ERROR", personPointer.InfectiveEpochs)
	}
	return false
}

func getInfected(networkPointer *bigNet) []int {
	infected := make([]int, 0, 1)
	for node := 0; node < nNodes; node++ {
		if (*networkPointer)[node].Infective == true &&
			(*networkPointer)[node].Survived == false &&
			(*networkPointer)[node].Dead == false {
			infected = append(infected, node)
		}
	}
	return infected
}

// countInfected counts the total number of infected people
func countInfected(networkPointer *bigNet, infective, survived, dead bool) int {
	counter := 0
	for node := 0; node < nNodes; node++ {
		if (*networkPointer)[node].Infective == infective &&
			(*networkPointer)[node].Survived == survived &&
			(*networkPointer)[node].Dead == dead {
			counter++
		}
	}
	//log.Println("INFECTED PEOPLE:", counter)
	return counter
}

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

func resetNetwork(networkPointer *bigNet) {
	for node := 0; node < nNodes; node++ {
		rand.Seed(time.Now().UnixNano())
		tmpR0 := int(math.Round(rand.NormFloat64()*stdR0 + medianR0))
		if tmpR0 < 0 {
			tmpR0 = 0
		}

		infectiveDays := make([]int8, tmpR0)

		for r := 0; r < tmpR0; r++ {
			infectiveDays[r] = int8(rand.Intn(infectiveEpochs))
		}

		(*networkPointer)[node].Infective = false
		(*networkPointer)[node].Survived = false
		(*networkPointer)[node].Dead = false
		(*networkPointer)[node].InfectiveEpochs = infectiveEpochs
		(*networkPointer)[node].InfectiveDays = infectiveDays
	}
}

func bernoulli(pSuccess float64) bool {
	if rand.Intn(10000) <= int(pSuccess*10000) {
		return true
	}
	return false
}
