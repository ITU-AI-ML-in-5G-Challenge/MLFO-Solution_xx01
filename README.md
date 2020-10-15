# mlfo-minimal
This work is done as part of ITU AI/ML 5G Challenge 2020 (ITU-ML5G-PS-024 - LYIT track)

## Requires
[Docker](https://docs.docker.com/get-docker/) v19.03.13

## Usage
### Initialize docker containers

`docker-compose up`

`docker exec -it db bash -c "mysql -uroot -pmlfo1234 modelrepo < modelrepo.sql"`

### Run

For edge use case:

`docker exec -it mlfo sh -c "go run mlfo.go edge_intent.yaml"`

For central cloud use case:

`docker exec -it mlfo sh -c "go run mlfo.go cloud_intent.yaml"`

### Debugging
Most of the errors might be fixed by stopping and restarting the containers:

`docker stop mlfo db`

`docker kill mlfo db`