#!/bin/bash

# Check if the token argument is provided
if [ -z "$1" ]; then
    echo "Error: Telegram bot token is missing. Please provide the token as an argument."
    echo "Example usage: ./init.sh YOUR_TELEGRAM_BOT_TOKEN"
    exit 1
fi


# Check if a container with the same name exists and stop and remove it
if docker ps -a --format '{{.Names}}' | grep -q "^bot-server$"; then
   echo "Found the old Container...Removing!"
   docker stop bot-server
   docker rm bot-server
fi


# Build Docker image 
echo "Building bot server image..."
docker build -t bot-server:latest . || {
    echo "failed to build docker image."
    exit 1
}

# run the created image in docker
echo "Running docker image..."
docker run -d --name bot-server -e TELEGRAM_TOKEN="$1" bot-server:latest || {
    echo "could not run the container check the logs!"
    exit 1
}
