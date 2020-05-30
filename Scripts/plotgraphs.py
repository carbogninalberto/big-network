import pandas as pd

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


if __name__ == "__main__":
    print("\nAnalysis 1")
    
    # TODO: load right path with flag
    data = loadCSV("data_hw1/data_ex1.csv")

    mean = computeMean(data)
    print("\tThe Mean is", mean)
    start_ci_mean_95, end_ci_mean_95 = getCIMean(data, 0.95)
    start_ci_mean_99, end_ci_mean_99 = getCIMean(data, 0.99)
    print("\t\tThe 95% CI for the Mean is [", start_ci_mean_95, ",", end_ci_mean_95, "]")
    print("\t\tThe 99% CI for the Mean is [", start_ci_mean_99, ",", end_ci_mean_99, "]")
