import pandas as pd
import argparse
import math
import random
import scipy.stats as st


#---------------------------------------------------
# Settings
#---------------------------------------------------

METRIC = {'median': 'computeMedian', 'mean':'computeMean', 'gap': 'computeGap', 'jain': 'computeJFI', 'stddev': 'computeStdDev', 'variance':'computeVar', 'log_mean':'computeLogMean', 'bernoulli':'bernoulliRVBS'}

#---------------------------------------------------

def loadCSV(csv_file): #TODO: Remember to modify it according to data of exercise 2
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
    # Define the flag --path with default value file.csv
    parser.add_argument('--path', action='store', dest='path', default='file.csv')
    # parse values
    args = parser.parse_args()

    print("\nAnalysis 1")
    #--path=1591000816817596400/simulation_trials_results.csv
    #data = loadCSV(args.path)
    data = pd.read_csv(args.path, header=None)

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
    


