#!/bin/sh
#SBATCH --time=20
#SBATCH --partition=standard
#SBATCH --nodes=4
#SBATCH --ntasks=4 --cpus-per-task=10
#SBATCH --ntasks-per-node=1
#SBATCH --nodelist=xgph15,xgph16,xgph17,xgph18
srun -n 4 ./client