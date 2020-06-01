package main

import (
	"log"
	"math/rand"
	"runtime"
)

// nationalHealthcareSystem SSN
type nationalHealthcareSystem struct {
	intensiveCare                   int
	subIntensiveCare                int
	intensiveCareHospitalization    []uint32
	subIntensiveCareHospitalization []uint32
}

// spreadingDesease runs a simulation over n epochs on a bigNet ([]person)
func spreadingDesease(networkPointer *bigNet, epochs int, epochsResultsPointer *[simulationEpochs][5]int, muskPointer *muskMeasure, socialDistancePointer *socialDistancingMeasure, ssnPointer *nationalHealthcareSystem, trialsResultsPointer *[][3]int, ssnEpochResults *[simulationEpochs][2]int, trial int) error {
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
							// Check if Healthcare is neeeded
							requireHealthcare, typeHealthcare := bernoulliHealthcare(pIntensiveCare, pSubIntensiveCare)
							if requireHealthcare {
								// if Healthcare needed, check if there are bed available
								addedToSSN := addToSSN(ssnPointer, uint32(infected), typeHealthcare)
								if addedToSSN {
									// if added to SSN can still infect others
									(*networkPointer)[infected].Infective = true
									(*networkPointer)[infected].InfectiveEpochs += hospitalDays
								} else {
									// if not possible to add to SSN the patient is dead
									//log.Println("NO BED AVAILABLE")
									(*networkPointer)[infected].InfectiveEpochs = 0
									(*networkPointer)[infected].Dead = true
								}
							} else {
								(*networkPointer)[infected].Infective = true
							}
						}
					}

					// I set to -1 in order to not consider it anymore
					(*networkPointer)[case0].InfectiveDays[day] = -1

				} else if (*networkPointer)[case0].InfectiveDays[day] > 0 {
					(*networkPointer)[case0].InfectiveDays[day]--
				}
			}
			// make time pass and reduce the remaining infective days
			_ = reduceInfectiveEpochs(&(*networkPointer)[case0], ssnPointer, uint32(case0))

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

								// Check if Healthcare is neeeded
								requireHealthcare, typeHealthcare := bernoulliHealthcare(pIntensiveCare, pSubIntensiveCare)
								if requireHealthcare {
									// if Healthcare needed, check if there are bed available
									addedToSSN := addToSSN(ssnPointer, uint32(infected), typeHealthcare)
									if addedToSSN {
										// if added to SSN can still infect others
										(*networkPointer)[infected].Infective = true
										(*networkPointer)[infected].InfectiveEpochs += hospitalDays
									} else {
										//log.Println("NO BED AVAILABLE")
										// if not possible to add to SSN the patient is dead
										(*networkPointer)[infected].InfectiveEpochs = 0
										(*networkPointer)[infected].Dead = true
									}
								} else {
									(*networkPointer)[infected].Infective = true
								}

							}
						}

						// I set to -1 in order to not consider it anymore
						(*networkPointer)[infectedID].InfectiveDays[day] = -1
					} else if (*networkPointer)[infectedID].InfectiveDays[day] > 0 {
						(*networkPointer)[infectedID].InfectiveDays[day]--
					}
				}

				// make time pass and reduce the remaining infective days
				_ = reduceInfectiveEpochs(&(*networkPointer)[infectedID], ssnPointer, uint32(infectedID))
			}
		}

		infectNumber := countInfected(networkPointer, true, false, false)
		log.Println("EPOCH\t", epoch,
			"\tACTIVE:\t", infectNumber,
			"\t\tINT.CARE:\t", len((*ssnPointer).intensiveCareHospitalization),
			"\tSUB.INT.CARE:", len((*ssnPointer).subIntensiveCareHospitalization))

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
		// number of total recovered
		(*epochsResultsPointer)[epoch][3] = countInfected(networkPointer, false, true, false)
		// number of total deaths
		(*epochsResultsPointer)[epoch][4] = countInfected(networkPointer, false, false, true)

		// number of intensive care
		(*ssnEpochResults)[epoch][0] = len((*ssnPointer).intensiveCareHospitalization)
		// number of sub intensive care
		(*ssnEpochResults)[epoch][1] = len((*ssnPointer).subIntensiveCareHospitalization)

		runtime.GC()
	}

	// assign number of total infected to col 0 of trial
	(*trialsResultsPointer)[trial][0] = countTotalInfected(networkPointer)
	// assign number of total recovered to col 1 of trial
	(*trialsResultsPointer)[trial][1] = countInfected(networkPointer, false, true, false)
	// assign number of total deaths to col 2 of trial
	(*trialsResultsPointer)[trial][2] = countInfected(networkPointer, false, false, true)

	return nil
}
