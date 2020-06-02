echo "Running with no policies"
go run . --computeCI=true --mctrials=10 --runpyscript=true --computeSSN=true --folder=spanish_flu
echo "Running with musk policy"
go run . --computeCI=true --mctrials=10 --runpyscript=true --computeSSN=true --folder=spanish_flu_musk_policy --muskEpoch=30 --muskProb=0.25
echo "Running with social distacing policy"
go run . --computeCI=true --mctrials=10 --runpyscript=true --computeSSN=true --folder=spanish_flu_social_distacing --socDisEpoch=30
echo "Running with hard restrictions"
go run . --computeCI=true --mctrials=10 --runpyscript=true --computeSSN=true --folder=spanish_flu_hard_measures --muskEpoch=30 --muskProb=0.25 --socDisEpoch=30