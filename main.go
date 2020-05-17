package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	nNodes = 6000000
	nEdges = 200
	cpus   = 500
)

type person struct {
	//Edges     []*person
	Edges     []int `json:"Edges"`
	Infective bool  `json:"Infective"`
	Survived  bool  `json:"Survived"`
	Dead      bool  `json:"Dead"`
}

func generateNodesWorker(id, start, end int, network *[nNodes]person, c chan bool) {

	for i := start; i < end; i++ {
		//log.Println(rand.Intn(nNodes))
		if i%28813 == 0 { //43201 28813
			runtime.GC()
			log.Println("adding node\t", i, "\tinterval [", start, ",", end, "]\t\t", id)
		}
		newNode := person{
			Infective: false,
			Survived:  false,
			Dead:      false,
		}

		for j := 0; j < nEdges; j++ {
			if j%50 == 0 && j != 0 {
				runtime.GC()
			}
			edgeID := rand.Intn(nNodes)
			if edgeID != i {
				newNode.Edges = append(newNode.Edges, edgeID)
			}
		}

		(*network)[i] = newNode
		//log.Println("node", i, "added WOW!")
	}
	runtime.GC()
	c <- true
	//time.Sleep(100 * time.Second)
	//log.Println("[NODE GENERATOR] Service", id, "finished.")

}

/*
func generateEdgesWorker(id, start, end int, network *[nNodes]*person, c chan bool) {
	tenth := 0
	for i := start; i < end; i++ {
		//log.Println(rand.Intn(nNodes))
		if i%43201 == 0 { //43201 28813
			runtime.GC()
		}

		if i%((end-start)/10) == 0 {
			//log.Println(id, "\t", tenth*10, "%")
			tenth++
		}

		for j := 0; j < nEdges; j++ {
			if j%50 == 0 {
				runtime.GC()
			}
			edgeID := rand.Intn(nNodes)
			if edgeID != i {
				randomEdge := (*network)[edgeID]
				(*network)[i].edges = append((*network)[i].edges, randomEdge)
			}
		}
	}
	runtime.GC()
	c <- true
	//time.Sleep(100 * time.Second)
	log.Println("[EDGE GENERATOR] Service", id, "finished.")
}

*/

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	runtime.GOMAXPROCS(16)
	c := make(chan bool, cpus)

	// random seed
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	// initialize network
	var network [nNodes]person
	networkPointer := &network

	//populate network
	slices := nNodes / cpus
	for j := 0; j < cpus; j++ {
		go generateNodesWorker(j, j*slices, (j+1)*slices, networkPointer, c)

	}
	//time.Sleep(5 * time.Second)
	log.Println("Creating network...")

	for i := 0; i < cpus; i++ {
		<-c
	}
	log.Println("Network nodes allocated.")

	runtime.GC()

	file, _ := json.Marshal(network)

	_ = ioutil.WriteFile("network.json", file, 0644)

	/*
		//populate edges
		for j := 0; j < cpus; j++ {
			go generateEdgesWorker(j, j*slices, (j+1)*slices, networkPointer, c)
		}

		log.Print("Creating edges...")

		for i := 0; i < cpus; i++ {
			<-c
		}

		log.Println("Network edges allocated.")

		log.Println(*network[0])
		/*
			for node := range network {
				if len((*network[node]).edges) != 0 {
					log.Println(*network[node])
				}
			}
	*/

	//

}
