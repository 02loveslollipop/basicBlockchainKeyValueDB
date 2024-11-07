#!/bin/bash
#Auto orchestrate the blockchain cluster and add the ip of the nodes to the host file

# Get mpi_node count from argument
nodes=$1

# Orchestrate the cluster
docker compose up --scale internal-node=$nodes -d # Start the cluster with the number of nodes specified in the argument

# Get the ip of the nodes
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker ps -q) > ./shared/host # Get the ip of the nodes and write it to the host file