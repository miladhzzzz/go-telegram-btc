#!/bin/bash

# Path to your git repository
REPO_DIR="/path/to/your/repo"
# Docker image name
IMAGE_NAME="your_docker_image_name"
# Docker container name
CONTAINER_NAME="your_container_name"
# Path to this script
SCRIPT_PATH=$(realpath "$0")

# Function to set up the cron job
setup_cron_job() {
    # Add a cron job to execute this script every 10 minutes
    (crontab -l ; echo "*/10 * * * * $SCRIPT_PATH >/dev/null 2>&1") | crontab -
    echo "Cron job set up successfully."
}

# Function to check for changes, build Docker image, and run container
check_and_update() {
    # Fetch changes from the remote repository
    cd "$REPO_DIR" || exit
    git fetch origin

    # Check if anything has changed
    if git diff --quiet HEAD FETCH_HEAD; then
        echo "No changes detected."
    else
        echo "Changes detected. Triggering Docker build..."

        # Stop and remove the existing container if it's running
        if docker ps | grep -q "$CONTAINER_NAME"; then
            docker stop "$CONTAINER_NAME"
            docker rm "$CONTAINER_NAME"
        fi

        # Build Docker image
        docker build -t "$IMAGE_NAME" .

        # Run Docker container
        docker run -d --name "$CONTAINER_NAME" "$IMAGE_NAME"
    fi
}

# Main
check_and_update
setup_cron_job
