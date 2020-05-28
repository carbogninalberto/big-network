package main

import (
	"log"
	"math"
	"math/rand"
	"time"
)

// reduceInfectiveEpochs returns true if person is no more infective
func reduceInfectiveEpochs(personPointer *person, ssnPointer *nationalHealthcareSystem, id uint32) bool {
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

		_ = removeFromSSN(ssnPointer, id)

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
	// random seed
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(10000) < int(pSuccess*10000) {
		return true
	}
	return false
}

// addToSSN returns false if it cannot be added to SSN
func addToSSN(ssnPointer *nationalHealthcareSystem, ID uint32, intensiveCare bool) bool {

	if intensiveCare &&
		(len((*ssnPointer).intensiveCareHospitalization)+1 <= (*ssnPointer).intensiveCare) {
		// add to intensive care if is not full
		(*ssnPointer).intensiveCareHospitalization = append((*ssnPointer).intensiveCareHospitalization, ID)
		return true
	} else if !intensiveCare &&
		(len((*ssnPointer).subIntensiveCareHospitalization)+1 <= (*ssnPointer).subIntensiveCare) {
		// add to sub intensive care if is not full
		(*ssnPointer).subIntensiveCareHospitalization = append((*ssnPointer).subIntensiveCareHospitalization, ID)
		return true
	}
	return false
}

// removeFromSSN
func removeFromSSN(ssnPointer *nationalHealthcareSystem, ID uint32) bool {
	intensiveLen := len((*ssnPointer).intensiveCareHospitalization)
	subIntensiveLen := len((*ssnPointer).subIntensiveCareHospitalization)
	// check if ID is on intensive care
	for i := 0; i < intensiveLen; i++ {
		if (*ssnPointer).intensiveCareHospitalization[i] == ID {
			(*ssnPointer).intensiveCareHospitalization = removeElementId((*ssnPointer).intensiveCareHospitalization, i)
			return true
		}
	}
	// ceck if ID is on sub intensive care
	for i := 0; i < subIntensiveLen; i++ {
		if (*ssnPointer).subIntensiveCareHospitalization[i] == ID {
			(*ssnPointer).subIntensiveCareHospitalization = removeElementId((*ssnPointer).subIntensiveCareHospitalization, i)
			return true
		}
	}
	// not found on SSN
	return false
}

// bernoulliHealthcare compute if Hospitalization is required and what kind
// Type of hospitalization: true intensive, false subintensive
func bernoulliHealthcare(pIntensive, pSubIntensive float64) (bool, bool) {
	rand.Seed(time.Now().UnixNano()) //set a new casual seed
	// check if requires intensive Care
	intensiveHospitalization := bernoulli(pIntensive)
	if !intensiveHospitalization {
		// if not intensive care, check if sub Intensive care
		subIntensiveHospitalization := bernoulli(pSubIntensive)
		if !subIntensiveHospitalization {
			// if any kind of hospitalization required
			return false, false
		}
		return true, false
	}
	return true, true
}

// remove element at certain ID
func removeElementId(s []uint32, i int) []uint32 {
	return append(s[:i], s[i+1:]...)
}
