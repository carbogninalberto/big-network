# SPREADING DISEASE SIMULATION

## How to run the simulation 

### Preliminary states
If you don't have already generated a Network.json file in your project folder run the following:

```
go run . -savenet=true
```
Now you will find the file in timestamp/Network.json, now you should move this file in the main project folder, or use the flag **-filenet=timestamp/Network.json**.
Assuming that you choose the first option, in order to run the simulation over the generated Network:

```
go run . -loadnet=true
```

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
- **-namenet**: Network.json", "default value is Network.json, it's the name of the network file
- **-mctrials**: default value is 1, you can choose how many trials run on the Montecarlo Simulation
- **-computeCI**:
- **-**:

### 🚀 Project Timeline
- [ ] Defining Main Project Objectives
- [ ] Reading Papers about the topic
- [ ] Gather Data
- [ ] Clean Data
- [ ] Statistical Inference on data
- [ ] Test [***windpowerlib***](https://github.com/wind-python/windpowerlib) as fitness function for power generation of wind turbines
- [ ] Implement fitness function for wind turbines placement
- [ ] Gather Data about cost and activation function the choosen wind turbines
- [ ] Euclidian distance fitness function for placing the **power plant**
- [ ] Collect **matplotlib** graphs, Pareto fronts in particular
- [ ] Final Report