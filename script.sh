#!/bin/sh
#SBATCH --time=20
#SBATCH --partition=standard
#SBATCH --nodes=8
#SBATCH --ntasks=8 --cpus-per-task=24
#SBATCH --ntasks-per-node=1
#SBATCH --nodelist=xgph0,xgph1,xgph2,xgph3,xgph4,xgph5,xgpe0,xgpe1
srun -n 8 ./client