package main

import (
	"log"
	"math/rand"
	"time"
)

// muskMeasure has 2 params:
// - Psucc: probability of success on prevention in case of contact
// - Active: containment measure applyed
type muskMeasure struct {
	Active    bool
	FromEpoch int
	Psucc     float64
}

// socialDistancingMeasure is the social distance policy
// AllowContacts is the map of allowed contact relations
type socialDistancingMeasure struct {
	Active        bool
	FromEpoch     int
	AllowContacts map[string]bool
}

// middlewareContainmentMeasure apply the requested measure on spreading it return true if there is an infection
func middlewareContainmentMeasure(personPointer *person, muskPointer *muskMeasure, socialDistancePointer *socialDistancingMeasure, epoch int) (bool, int) {
	// generati a valid ID regarding socialDistancingMeasure if Active
	if (*socialDistancePointer).Active && epoch >= (*socialDistancePointer).FromEpoch && (*socialDistancePointer).FromEpoch != -1 {
		infected, infectedID := generateRandomId(personPointer, &(*socialDistancePointer).AllowContacts)
		if infected && (*muskPointer).Active {
			prevent := bernoulli((*muskPointer).Psucc)
			if ISDEBUG {
				log.Println("SOCIAL DISTACING + MUSK POLICY", !prevent, infectedID)
			}
			return !prevent, infectedID
		}

		if ISDEBUG {
			log.Println("SOCIAL DISTACING POLICY", infected, infectedID)
		}

		return infected, infectedID
	} else if (*muskPointer).Active && epoch >= (*muskPointer).FromEpoch && (*muskPointer).FromEpoch != -1 {
		// go into this branch only if musk policy
		infectedID := rand.Intn(len((*personPointer).Edges))
		prevent := bernoulli((*muskPointer).Psucc)

		if ISDEBUG {
			log.Println("MUSK POLICY ACTIVE", !prevent, infectedID)
		}
		return !prevent, infectedID

	} else {
		randomInfect := rand.Intn(len((*personPointer).Edges))
		infectedID := int((*personPointer).Edges[randomInfect])

		if ISDEBUG {
			log.Println("NO POLICY ACTIVE", infectedID)
		}

		return true, infectedID
	}
}

// generateRandomId generate a random edge on the social distance policy
func generateRandomId(personPointer *person, allowContacts *map[string]bool) (bool, int) {
	generationBaseAllowedId := make([]int, 1)
	edgeLen := len((*personPointer).Edges)

	// loop all edjes to found allowed ones
	for edge := 0; edge < edgeLen; edge++ {
		// check that relation type is allowed on the map
		if (*allowContacts)[(*personPointer).RelationType[edge]] {
			generationBaseAllowedId = append(generationBaseAllowedId, edge)
		}
		/*
			// loop over allowed edge
			for _, a := range *allowContacts {
				if a == (*personPointer).RelationType[edge] {
					generationBaseAllowedId = append(generationBaseAllowedId, edge)
					break
				}
			}
		*/
	}
	if len(generationBaseAllowedId) > 0 {
		rand.Seed(time.Now().UnixNano())
		randomInfect := int((*personPointer).Edges[rand.Intn(len(generationBaseAllowedId))])
		return true, randomInfect
	}
	return false, 0
}
