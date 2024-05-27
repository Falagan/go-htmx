#!/bin/bash\

#################################################
#       DEVELOPMENT ENVIRONMENT SETUP           #
#################################################

# SET ENVIRONMENT VARIABLES

echo "1. [Setting up env variables]"

echo "1.1. APP ENV"

export APP_ENV=development

echo "1.2. SERVER ADDRESS"

export SERVER_ADDRESS=""

echo "1.3. SERVER PORT"

export SERVER_PORT=4000

echo "1.4. DATABASE URL"

export DATABASE_URL="postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"

# RUN DOCKER CONTAINERS

echo "2. [Running docker containers]"

## DATABASE

echo "2.1. Running database container"

sudo docker remove --force /ar-db-dev 
sudo docker run --name ar-db-dev -e POSTGRES_PASSWORD=mysecretpassword -p "5432:5432" -d postgres
