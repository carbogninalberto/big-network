import pandas as pd
import argparse
import math
import random
import scipy.stats as st
import matplotlib.pyplot as plt
import numpy as np
from matplotlib.ticker import FuncFormatter

#---------------------------------------------------
# Settings
#---------------------------------------------------

METRIC = {'median': 'computeMedian', 'mean':'computeMean', 'gap': 'computeGap', 'jain': 'computeJFI', 'stddev': 'computeStdDev', 'variance':'computeVar', 'log_mean':'computeLogMean', 'bernoulli':'bernoulliRVBS'}

#---------------------------------------------------

def loadCSV(csv_file): 
    dataset = pd.read_csv(csv_file, header=None) #here we are working with Pandas DataFrame
    #print(dataset.shape)

    data_points = []
    column_dim = dataset.shape[1]  # save the number of columns
    row_dim = dataset.shape[0]  # save the number of columns
    if (column_dim == 1):  # check if we have data organized in a single column
        data_points = dataset[0].values.tolist()  # convert the values of the column into a list of values
    elif (column_dim > 1):  # check if data is expressed as a matrix of values (multiple columns)
        data_points = dataset.values.tolist()
        if (row_dim == 1):  # check if data is expressed as a single row
            data_points = data_points[0]
    return data_points

def getCIMean(data, ci_value):
    eta = computeEta(ci_value, data, 't') # version for t distribution
    #eta = computeEta(ci_value, data, 'normal') # version for normal distribution

    #np_mean = np.mean(data) #numpy mean
    mean = computeMean(data)
    #np_sem = st.sem(data) #scipy stats standard error mean
    std_error_mean = computeStdDev(data) / math.sqrt(len(data))
    conf_interval = eta * std_error_mean

    start_int = mean - conf_interval
    end_int = mean + conf_interval

    return start_int, end_int

def computeStdDev(x):
    mean = computeMean(x)
    sum_squares = 0
    n = len(x)
    for i in range(0,n):
        sum_squares += pow((x[i] - mean), 2)
    std_dev = math.sqrt(sum_squares / (n - 1))
    return std_dev

def computeMean(x):
    sum_x = 0
    n = len(x)
    for i in range(0,n):
        sum_x += x[i]
    mean = sum_x / n
    return mean

def computeStdDev(x):
    mean = computeMean(x)
    sum_squares = 0
    n = len(x)
    for i in range(0,n):
        sum_squares += pow((x[i] - mean), 2)
    std_dev = math.sqrt(sum_squares / (n - 1))
    return std_dev

def computeEta(ci_value, data, distribution):
    if distribution == 't':
        eta = st.t.ppf((1 + ci_value) / 2, len(data) - 1)
    elif distribution == 'normal':
        eta = st.norm.ppf((1 + ci_value) / 2)
    return eta

def bootstrapAlgorithm(dataset, accuracy=25, ci_level=0.95, metric='mean'):
    ds_length = len(dataset)
    samples_metric = []

    samples_metric.append(globals()[METRIC[metric]](dataset))
    R = math.ceil(2 * (accuracy / (1-ci_level))) - 1

    for r in range(R):
        tmp_dataset = []
        for i in range(ds_length):
            tmp_dataset.append(dataset[random.randrange(0, ds_length, 1)])
        samples_metric.append(globals()[METRIC[metric]](tmp_dataset)) # load the desired metric function

    samples_metric.sort()
    #print('sample_metric_len:', len(samples_metric), 'range len:', len(samples_metric[accuracy:(R+1-accuracy)]))
    return samples_metric[accuracy:(R+1-accuracy)]

if __name__ == "__main__":
    # Define the parser
    parser = argparse.ArgumentParser(description="Analyze the result data of the spreading disease")
    # Define the flag --trials with default value file.csv
    parser.add_argument('--trialsFile', action='store', dest='trialsFile', default='file.csv')
     # Define the flag --ssn with default value empty
    parser.add_argument('--ssnFile', action='store', dest='ssnFile', default='')
     # Define the flag --simulation with default value simulation.csv
    parser.add_argument('--simulationFile', action='store', dest='simulationFile', default='simulation.csv')
    # Define the flag --folder with default value of graphs
    parser.add_argument('--folder', action='store', dest='folder', default='graphs/')
    # parse values
    args = parser.parse_args()

    print("\nAnalysis 1")
    data = pd.read_csv(args.folder+args.trialsFile, header=None)

    # Total Infected Population
    print("\nTOTAL INFECTED")
    total_infected = data.iloc[:, 0]
    total_infected_mean = computeMean(total_infected)
    print("\n\tThe Mean of Total Infected is", total_infected_mean)
    total_infected_start_ci_mean_95, total_infected_end_ci_mean_95 = getCIMean(total_infected, 0.95)
    total_infected_start_ci_mean_99, total_infected_end_ci_mean_99 = getCIMean(total_infected, 0.99)
    print("\t1. Confidence Intervals using Asymptotic Formulas")
    print("\t\tThe 95% CI for the Mean of Total Infected is [", total_infected_start_ci_mean_95, ",", total_infected_end_ci_mean_95, "]")
    print("\t\tThe 99% CI for the Mean of Total Infected is [", total_infected_start_ci_mean_99, ",", total_infected_end_ci_mean_99, "]")
    total_infected_bs_95 = bootstrapAlgorithm(dataset=total_infected)
    print("\t2. Confidence Intervals using Bootstrap Algorithm")
    print('\t\tThe 95% CI for Mean of Total Infected is [{}, {}]'.format(total_infected_bs_95[0], total_infected_bs_95[len(total_infected_bs_95)-1]))

    # Total Recovered
    print("\nTOTAL RECOVERED")
    total_recovered = data.iloc[:, 1]
    total_recovered_mean = computeMean(total_recovered)
    print("\n\tThe Mean of Total Recovered is", total_recovered_mean)
    total_recovered_start_ci_mean_95, total_recovered_end_ci_mean_95 = getCIMean(total_recovered, 0.95)
    total_recovered_start_ci_mean_99, total_recovered_end_ci_mean_99 = getCIMean(total_recovered, 0.99)
    print("\t1. Confidence Intervals using Asymptotic Formulas")
    print("\t\tThe 95% CI for the Mean of Total Recovered is [", total_recovered_start_ci_mean_95, ",", total_recovered_end_ci_mean_95, "]")
    print("\t\tThe 99% CI for the Mean of Total Recovered is [", total_recovered_start_ci_mean_99, ",", total_recovered_end_ci_mean_99, "]")
    total_recovered_bs_95 = bootstrapAlgorithm(dataset=total_recovered)
    print("\t2. Confidence Intervals using Bootstrap Algorithm")
    print('\t\tThe 95% CI for Mean of Total Recovered is [{}, {}]'.format(total_recovered_bs_95[0], total_recovered_bs_95[len(total_recovered_bs_95)-1]))
    
    # Total Deaths
    print("\nTOTAL DEATHS")
    total_deaths = data.iloc[:, 2]
    total_deaths_mean = computeMean(total_deaths)
    print("\n\tThe Mean of Total Deaths is", total_deaths_mean)
    total_deaths_start_ci_mean_95, total_deaths_end_ci_mean_95 = getCIMean(total_deaths, 0.95)
    total_deaths_start_ci_mean_99, total_deaths_end_ci_mean_99 = getCIMean(total_deaths, 0.99)
    print("\t1. Confidence Intervals using Asymptotic Formulas")
    print("\t\tThe 95% CI for the Mean of Total Deaths is [", total_deaths_start_ci_mean_95, ",", total_deaths_end_ci_mean_95, "]")
    print("\t\tThe 99% CI for the Mean of Total Deaths is [", total_deaths_start_ci_mean_99, ",", total_deaths_end_ci_mean_99, "]")
    total_deaths_bs_95 = bootstrapAlgorithm(dataset=total_deaths)
    print("\t2. Confidence Intervals using Bootstrap Algorithm")
    print('\t\tThe 95% CI for Mean of Total Deaths is [{}, {}]'.format(total_deaths_bs_95[0], total_deaths_bs_95[len(total_deaths_bs_95)-1]))

    metrics = ["Total \nInfected \nAsymptotic", "Total \nRecovered \nAsymptotic", "Total \nDeaths \nAsymptotic"]
    mean_values = [total_infected_mean, total_recovered_mean, total_deaths_mean]
    error_values = [total_infected_mean-total_infected_start_ci_mean_95, total_recovered_mean-total_recovered_start_ci_mean_95, total_deaths_mean-total_deaths_start_ci_mean_95]
    
    metrics_bs = ["Total \nInfected \nBootstrap", "Total \nRecovered \nBootstrap", "Total \nDeaths \nBootstrap"]
    error_values_bs = [total_infected_mean-total_infected_bs_95[0], total_recovered_mean-total_recovered_bs_95[0], total_deaths_mean-total_deaths_bs_95[0]]

    mean_values[:] = [ x / 1000 for x in mean_values] #in thousands
    error_values[:] = [ x / 1000 for x in error_values] #in thousands
    error_values_bs[:] = [ x / 1000 for x in error_values_bs] #in thousands

    # image settings
    my_dpi = 92 # setting my dpi
    plt.figure(figsize=(650/my_dpi, 650/my_dpi), dpi=my_dpi)

    #plt.subplot(1, 1, 1)
    plt.title('Confidence Intervals of Total Infected with 2 techniques')
    plt.errorbar([metrics[0], metrics_bs[0]], [mean_values[0], mean_values[0]], yerr=[error_values[0], error_values_bs[0]], linestyle='None', marker='.', lw=1, fmt='.k', capsize=3, fillstyle='full')
    plt.ylabel("Value/1000")
    plt.xlabel("Total Infected")
    plt.savefig(args.folder+"results_total_infected.png", dpi=my_dpi*3)
    plt.clf()

    #plt.subplot(3, 1, 2)
    plt.title('Confidence Intervals of Total Recovered with 2 techniques')
    plt.errorbar([metrics[1], metrics_bs[1]], [mean_values[1], mean_values[1]], yerr=[error_values[1], error_values_bs[1]], linestyle='None', marker='.', lw=1, fmt='.k', capsize=3)
    plt.ylabel("Value/1000")
    plt.xlabel("Total Recovered")
    plt.savefig(args.folder+"results_total_recovered.png", dpi=my_dpi*3)
    plt.clf()

    #plt.subplot(3, 1, 3)
    plt.title('Confidence Intervals of Total Deaths with 2 techniques')
    plt.errorbar([metrics[2], metrics_bs[2]], [mean_values[2], mean_values[2]], yerr=[error_values[2], error_values_bs[2]], linestyle='None', marker='.', lw=1, fmt='.k', capsize=3)
    plt.ylabel("Value/1000")
    plt.xlabel("Total Deaths")
    plt.savefig(args.folder+"results_total_deaths.png", dpi=my_dpi*3)
    plt.clf()
    
    #plt.show()

    x = np.arange(3)
    only_metrics = ("Total \nInfected", "Total \nRecovered", "Total \nDeaths")

    def millions(x, pos):
        'The two args are the value and tick position'
        return x 


    formatter = FuncFormatter(millions)

    fig, ax = plt.subplots()
    #ax.yaxis.set_major_formatter(formatter)
    plt.title("Metrics Results Averages")
    plt.bar(only_metrics, mean_values)
    plt.xticks(x, only_metrics)
    plt.savefig(args.folder+"results_.png", dpi=my_dpi*3)
    plt.clf()
    #plt.show()


    # ANALYSIS 2

    simulationData = pd.read_csv(args.folder+args.simulationFile, header=None)

    # generate x for number of epochs
    t = np.arange(0, len(simulationData), 1)
    active_infected = simulationData.iloc[:, 0]

    fig, ax = plt.subplots()
    ax.plot(t, active_infected, label="Active Infected")

    ax.set(xlabel='epochs (day)', ylabel='number of people',
        title='Spreading of the disease')
    ax.grid()
    
    plt.legend()

    fig.savefig(args.folder+"results_epidemic.png", dpi=my_dpi*3)
    plt.clf()

    if args.ssnFile != '':
        # Data for plotting
        ssnData = pd.read_csv(args.folder+args.ssnFile, header=None)

        # generate x for number of epochs
        t = np.arange(0, len(ssnData), 1)
        intensiveCare = ssnData.iloc[:, 0]
        subIntensiveCare = ssnData.iloc[:, 1]

        fig, ax = plt.subplots()
        ax.plot(t, intensiveCare, label="Intensive Beds")
        ax.plot(t, subIntensiveCare, label="Other Beds")

        ax.set(xlabel='epochs (day)', ylabel='number of people',
            title='National Healthcare Systems usage because of epidemic')
        ax.grid()

        plt.legend()

        fig.savefig(args.folder+"results_ssn_epidemic.png", dpi=my_dpi*3)
        plt.clf()



    
