package main

import "testing"

func TestReduceInfectiveEpochs(t *testing.T) {
	p1 := person{
		Infective:       false,
		Survived:        false,
		Dead:            false,
		InfectiveEpochs: 2,
	}
	p1Result := reduceInfectiveEpochs(&p1)
	if p1Result != false {
		t.Errorf("Result inccorect, got: %t, want: %t.", p1Result, false)
	}
	p2 := person{
		Infective:       false,
		Survived:        false,
		Dead:            false,
		InfectiveEpochs: 1,
	}
	p2Result := reduceInfectiveEpochs(&p2)
	if p2Result != true {
		t.Errorf("Result inccorect, got: %t, want: %t.", p2Result, true)
	}
}
