package main

import (
	"log"
	"math/rand"
	"runtime"
)

// spreadingDesease runs a simulation over n epochs on a bigNet ([]person)
func spreadingDesease(networkPointer *bigNet, epochs int, epochsResultsPointer *[simulationEpochs][5]int, muskPointer *muskMeasure, socialDistancePointer *socialDistancingMeasure) error {
	for epoch := 0; epoch < epochs; epoch++ {

		if epoch == 0 {
			// pick a random infect over the graph
			case0 := rand.Intn(nNodes)
			(*networkPointer)[case0].Infective = true
			log.Println("CASE 0:", case0)
			infectiveDaysLen := len((*networkPointer)[case0].InfectiveDays)

			for day := 0; day < infectiveDaysLen; day++ {
				if (*networkPointer)[case0].InfectiveDays[day] == 0 {

					isInfected, infected := middlewareContainmentMeasure(&(*networkPointer)[case0], muskPointer, socialDistancePointer, epoch)

					if isInfected {
						if (*networkPointer)[infected].InfectiveEpochs > 0 {
							(*networkPointer)[infected].Infective = true
						}
					}

					// I set to -1 in order to not consider it anymore
					(*networkPointer)[case0].InfectiveDays[day] = -1

				} else if (*networkPointer)[case0].InfectiveDays[day] > 0 {
					(*networkPointer)[case0].InfectiveDays[day]--
				}
			}
			// make time pass and reduce the remaining infective days
			_ = reduceInfectiveEpochs(&(*networkPointer)[case0])

		} else {
			infected := getInfected(networkPointer)

			for _, infectedID := range infected {

				infectiveDaysLen := len((*networkPointer)[infectedID].InfectiveDays)

				for day := 0; day < infectiveDaysLen; day++ {
					if (*networkPointer)[infectedID].InfectiveDays[day] == 0 {

						isInfected, infected := middlewareContainmentMeasure(&(*networkPointer)[infectedID], muskPointer, socialDistancePointer, epoch)

						if isInfected {
							if (*networkPointer)[infected].Infective == false &&
								(*networkPointer)[infected].Dead == false &&
								(*networkPointer)[infected].Survived == false &&
								(*networkPointer)[infected].InfectiveEpochs > 0 {
								(*networkPointer)[infected].Infective = true
							}
						}

						// I set to -1 in order to not consider it anymore
						(*networkPointer)[infectedID].InfectiveDays[day] = -1
					} else if (*networkPointer)[infectedID].InfectiveDays[day] > 0 {
						(*networkPointer)[infectedID].InfectiveDays[day]--
					}
				}

				// make time pass and reduce the remaining infective days
				_ = reduceInfectiveEpochs(&(*networkPointer)[infectedID])
			}
		}

		infectNumber := countInfected(networkPointer, true, false, false)
		log.Println("EPOCH\t", epoch, "\tINFECTED:\t", infectNumber)

		// number of infected today
		(*epochsResultsPointer)[epoch][0] = infectNumber
		// new number of infected today regards yesterday
		if epoch != 0 {
			lastInfected := (*epochsResultsPointer)[epoch-1][0]
			(*epochsResultsPointer)[epoch][1] = infectNumber - int(lastInfected)
		} else {
			(*epochsResultsPointer)[epoch][1] = infectNumber
		}

		// number of total infected
		(*epochsResultsPointer)[epoch][2] = countTotalInfected(networkPointer)
		// number of total survived
		(*epochsResultsPointer)[epoch][3] = countInfected(networkPointer, false, true, false)
		// number of total dead
		(*epochsResultsPointer)[epoch][4] = countInfected(networkPointer, false, false, true)

		runtime.GC()
	}
	return nil
}
