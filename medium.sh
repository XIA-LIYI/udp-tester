#!/bin/sh
#SBATCH --time=20
#SBATCH --partition=standard
#SBATCH --nodes=2
#SBATCH --ntasks=2 --cpus-per-task=15
#SBATCH --ntasks-per-node=1
#SBATCH --nodelist=xgph13,xgph14
srun -n 2 ./client