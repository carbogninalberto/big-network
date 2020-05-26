# SPREADING DISEASE SIMULATION

## Preliminary states
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

## How to run the simulation 

If you want to run a Montecarlo Simulation you chan choose how many trials to run, for example for 100 trials:

```
go run . -loadnet=true -mctrials=100
```

and if you want also to compute the CI over all the three metrics [Total Cases, Healed, Dead]:

```
go run . -loadnet=true -mctrials=100 -computeCI=true
```



### List of all flags
- **-loadnet**: default value is false, if true it load a network from a file called Network.json, to change the loading file name check flag namenet
- **-savenet**: default value is false, if true saves network on timestamp/Network.json
- **-namenet**: default value is Network.json, it's the name of the network file
- **-mctrials**: default value is 1, you can choose how many trials run on the Montecarlo Simulation
- **-computeCI**: default value is false, set to true when use flag -mctrials > 1 to get Confidence Intervals of metrics
- **-runpyscript**: default valuse is false, set to true if you want to print graphs of simulation with matplotlib

### 🚀 Project Timeline
- [X] Defining Main Project Objectives
- [X] Create a Big Network Graph of the Population
- [X] Gather Information about Italian Hospitals
- [X] Gather Information about different viruses
- [ ] Coding the core-project  (**work in progress**)
- [ ] Montecarlo Simulation and CI
- [ ] Final Report (**work in progress**)