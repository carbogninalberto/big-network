<h1 align="center">🦠😷🧼🏥SPREADING DISEASE SIMULATION🦠😷🧼🏥</h1>
<div align="center">
<img widht="256" height="256" src=".github/logo.png">
</div>

# Table of Contents
- [Introduction](#introduction)
- [Installation](#installation)
- [Simulation](#simulation)
  - [List of available flags](#list-of-available-flags)
- [Examples](#examples)
  - [Header Data](#header-data)
- [Project Timeline](#project-timeline)

# Introduction

# Installation
If you don't have already generated a Network.json file in your project folder run the following:

```
go run . -savenet=true
```
Now you will find the file in timestamp/Network.json, now you should move this file in the main project folder, or use the flag **-filenet=timestamp/Network.json**.
Assuming that you choose the first option, in order to run the simulation over the generated Network:

```
go run . -loadnet=true
```

In order to avoid running the simulation run the generation of the network with:

```
go run . -savenet=true -mctrials=0
```

# Simulation

If you want to run a Montecarlo Simulation you chan choose how many trials to run, for example for 100 trials:

```
go run . -loadnet=true -mctrials=100
```

and if you want also to compute the CI over all the three metrics [Total Cases, Healed, Dead]:

```
go run . -loadnet=true -mctrials=100 -computeCI=true
```

## List of available flags
- **-loadnet**: default value is false, if true it load a network from a file called Network.json, to change the loading file name check flag namenet
- **-savenet**: default value is false, if true saves network on timestamp/Network.json
- **-namenet**: default value is Network.json, it's the name of the network file
- **-mctrials**: default value is 1, you can choose how many trials run on the Montecarlo Simulation
- **-computeCI**: default value is false, set to true when use flag -mctrials > 1 to get Confidence Intervals of metrics
- **-runpyscript**: default valuse is false, set to true if you want to print graphs of simulation with matplotlib

# Examples
In the example folder you can find the results in .csv format for couple simulations with the same virus, but applying containment measures vs no measures.

The virus information are:
- R0: median 5, std 0.8
- Infective Epochs: 14
- Dead Rate: 5.4%

The containment measures are:
- Musk: from day 30 with prevention probability of 95% cases of contact
- Social Distance: from day 40 with contact allowed only between parents and friends

## Header Data
The exported .csv files header is the following:
- ACTIVE CASES
- NEW CASES RESPECT YESTERDAY
- TOTAL CASES
- SURVIVED CASES
- DEAD CASES

# Project Timeline
- [X] Defining Main Project Objectives
- [X] Create a Big Network Graph of the Population
- [X] Gather Information about Italian Hospitals
- [X] Gather Information about different viruses
- [ ] Coding the core-project  (**work in progress**)
- [ ] Montecarlo Simulation and CI
- [ ] Final Report (**work in progress**)