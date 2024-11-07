#!/bin/bash
#Auto orchestrate the blockchain cluster and add the ip of the nodes to the host file

# Get mpi_node count from argument
nodes=$1

# Build the image
docker compose build --no-cache

# Orchestrate the cluster
docker compose up --scale internal-node=$nodes # Start the cluster with the number of nodes specified in the argument