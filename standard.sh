#!/bin/sh
#SBATCH --time=20
#SBATCH --partition=standard
#SBATCH --nodes=6
#SBATCH --ntasks=6 --cpus-per-task=1
#SBATCH --ntasks-per-node=1
#SBATCH --nodelist=xgph0,xgph1,xgph2,xgph3,xgph4,xgph5
srun -n 6 ./client