# Technical Documentation
In this section are written the implementation specification for the version 0.1 of the toolchain.

# Table of Contents
- [Introduction](#introduction)
- [Customizable Parameters](#customizable-parameters)
  - [Graph](#graph)
  - [Virus](#virus)
  - [National healthcare system](#national-healthcare-system)
  - [Containment measures](#containment-measures)
  - [Simulation](#simulation)
- [Requirements](#requirements)

# Introduction
In the project has been consider a graph with 4905854 nodes, equal to the number of people in the Italy's region of Veneto; moreover 150 edges per node are randomly added following [Dunbar number](#https://en.wikipedia.org/wiki/Dunbar%27s_number) as reference and assign to each edge a relation type:
- 10 edges: family
- 10 edges: close friends
- 30 edges: colleagues, people with low interactions
- 100 edges: strangers 

# Customizable Parameters
The tool aims to guarantee a high customizable set of parameters to conduct different kind of simulation, with different viruses and conditions, taking in account of:
- **Graph**: you can generate a custom graph where to run the simulation and save it
- **Virus**: multiple parameters to set
- **National Healthcare System**: intensive and sub-intensive healthcare
- **Containment measures**: musk related and social distacing policies
- **Simulation**: trials, data export, python scripts... ecc

## Graph
For the graph you can choose the number of edges per node with **nEdges** and the number of nodes with **nNodes**.

## Virus
For the virus is possible to modify the following parameters:
- **medianRO**: this is the median value of r0 for the virus
- **stdR0**: standard deviation of R0 assumin a norma distribution
- **infectiveEpochs**: the number of infective days of the virus (including the incubation period)
- **deadRate**: the percentage of people of infected people that die because of the virus 
- **incubationEpochs**: the incubation period in days of the virus
- **pIntensiveCare**: percentage of people that require intensive care beds
- **pSubIntensiveCare**: percentage of people that require sub-intensive care beds

## National healthcare system
For the healthcare system is possible to modify the following parameters:
- **bedIntensiveCare**: number of people that can use the intensive care beds
- **bedSubIntensiveCare**: number of people that can use the sub-intensive care beds

## Containment measures
Regards containment measures the following parameters are customizable:
- **muskEpoch**: from which epoch start using the musk policy
- **muskProb**: no infection spreading probability in case of contact
- **socDisEpoch**: from which epoch start using the social distacing policy
Note that is also possible to edit the map of the social distacing struct in order to exclude the kind of relation you want on the simulation.

### Musk Policy
The musk policy is used by every node, in case of contact with infected node the probability of being infected is 1 - **muskProb**. You can combine this policy with social distacing.

### Social Distacing
The social distacing policy allows the simulation to generate contact only with the allowed edges. You can combine this policy with the musk one.

## Simulation
Regards the simulation you can decide the following parameters:
- **trials**: number of trials of the simulation, useful to gather data for calculating Confidence Intervals with a python script
- **simulationEpochs**: number of epochs to run the simulation, each epoch corresponds to a day

# Requirements
In order to run the simulation on big graphs a certain hardware is needed, here some examples:
- 400k nodes, 150 edges/node => 12 GB RAM
- 4mln nodes, 150 edges/node => 54 GB RAM (note that you can use a server running linux)