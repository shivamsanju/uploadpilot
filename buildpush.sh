#!/bin/bash

# Default values
DOCKER_HUB_USERNAME=""
TAG="latest"

# Parse named parameters
while [[ $# -gt 0 ]]; do
    case $1 in
        -u)
            DOCKER_HUB_USERNAME="$2"
            shift 2
            ;;
        -t)
            TAG="$2"
            shift 2
            ;;
        *)
            echo "Unknown argument: $1"
            exit 1
            ;;
    esac
done

if [ -z "$DOCKER_HUB_USERNAME" ]; then
    echo "Docker Hub username is required. Use -username <your_username>"
    exit 1
fi

# List of image names and corresponding Dockerfiles
IMAGES=("manager" "agent")


# Iterate over each image and build/push
for IMAGE in "${IMAGES[@]}"; do
    DOCKERFILE=".Dockerfile.$IMAGE"
    echo "Building $IMAGE with tag $TAG using $DOCKERFILE..."
    docker build -t "$DOCKER_HUB_USERNAME/$IMAGE:$TAG" -f "docker/$DOCKERFILE" . || {
        echo "Failed to build $IMAGE. Aborting push."
        exit 1
    }
done

# If all builds succeed, push all images
for IMAGE in "${IMAGES[@]}"; do
    echo "Pushing $IMAGE to Docker Hub with tag $TAG..."
    docker push "$DOCKER_HUB_USERNAME/$IMAGE:$TAG"
done

echo "All images processed."

# ./buildpush.sh -u abcshvm -t test