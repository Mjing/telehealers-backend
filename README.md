# Telehealers Backend

## Current Features
1. Image API
2. LiveKit Access token API
3. [Current] Auth service

# Hosting Server

# Setting up Development environment

## Local DB Container Setup
```sh
$sudo docker-compose --env-file dev.env -f scripts/db_setup_scripts/db-compose.yml  up
```
*Might need sudo


# Environment Setup

```sh
export LIVEKIT_API_KEY='****'
export LIVEKIT_API_SECRET='****'
export DB_PORT=3306
export DB_DATABASENAME=telehealers
export DB_USER=""
export DB_PASS=""
export DB_ADDRESS="localhost:3306"
# Env variables only needed for development
export DB_DATA_DIR=./data_dir #Directory where sql container will save its data

```



