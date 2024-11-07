# Basic Blockchain Key-Value Database (BBD)

> A basic blockchain-based key-value database with a CLI for interacting, a Python class wrapper for the API, and Docker support for orchestrating a cluster of nodes.

## Intro

BBD is a simple blockchain-based key-value database that allows users to store key-value pairs in a blockchain. The project includes a Go application for the blockchain node, a Python CLI for interacting with the blockchain, and Docker support for orchestrating a cluster of nodes. One of the blockchain nodes is exposed to the host machine to interact with the user, but all of the nodes are capable of adding, storing and retrieving key-value pairs if needed.

## Project Structure

- `src/cmd/main/main.go`: Main Go application for the blockchain node.
- `cli/cli.py`: Python CLI for interacting with the blockchain.
- `.github/workflows/docker-image.yml`: GitHub Actions workflow for building the Docker image.
- `run.sh`: Script to orchestrate the blockchain cluster detached and add the IP of the nodes to the host file.
- `rebuild.sh`: Script to rebuild the Docker image and orchestrate the blockchain cluster **NON-DETACHED**.
- `docker-compose.yml`: Docker Compose configuration for the blockchain cluster.
- `Dockerfile`: Dockerfile for building the blockchain node image.

## Getting Started

### Prerequisites

- Docker (tested with version 27.3.1)
- Docker Compose (tested with version 1.29.7)
- Go 1.23.2+ (tested on version 1.23.2)(for local development)
- Python 3.10.12+ (tested from python 3.10.12 to 3.12.7)(for CLI and API wrapper)

### 1.a Building and Running the Blockchain Node
The recommended and tested way to run the blockchain node is compiling the container from the Dockerfile and running the cluster with Docker Compose.

1. **Build the Docker image:**

    ```sh
    docker compose build --no-cache
    ```

2. **Run the blockchain cluster:**

    ```sh
    ./run.sh <number_of_nodes>
    ```

3. **Rebuild the Docker image and run the cluster:**

    ```sh
    ./rebuild.sh <number_of_nodes>
    ```

### 1.b Retrieving container from Docker Hub
It is also possible to retrieve the container from Docker Hub and run the cluster with Docker Compose. In this case, the `run.sh` and `rebuild.sh` scripts may not work as expected.

1. **Pull the Docker image:**

    ```sh
    docker pull 02loveslollipop/bbd:latest
    ```

2. **Run the blockchain cluster:**

    ```sh
    docker compose up --scale node=<number_of_nodes>
    ```

### 2. Interacting with the Blockchain

Use the Python CLI to interact with the blockchain.

1. **Install dependencies:**

    ```sh
    pip install -r requirements.txt
    ```

2. **Add a key-value pair to the blockchain:**

    ```sh
    python cli/cli.py add --data "key":"value" --url <server_url> --port <server_port>
    ```

### API Endpoints

The blockchain node exposes the following API endpoints:

- `GET /blockchain`: Get the current blockchain in the node.
- `POST /append`: Add a key-value pair to the blockchain.
- `GET /ping`: Check if the node is running.

### GitHub Actions

The project includes a GitHub Actions workflow to build the Docker image on push or pull request to the `main` branch.

## Known Issues

- If a node joins the network after the blockchain has been initialized, it might not be able to fetch the current blockchain from the other nodes.