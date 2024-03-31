#!/bin/bash

# Variables
REGION=us-east-1
REPOSITORY_NAME=golang
IMAGE_TAG=regroup-service
ACCOUNT_ID=799934209842

# Get the login command from ECR and execute it directly
docker login -u AWS -p $(aws ecr get-login-password --region $REGION) $ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/$REPOSITORY_NAME

# Build your Docker image locally
docker build -t $ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/$REPOSITORY_NAME:$IMAGE_TAG --platform=linux/amd64 .

# Tag your Docker image to match the repository name
# docker tag $REPOSITORY_NAME:latest $ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/$REPOSITORY_NAME:$IMAGE_TAG

# Push this image to your newly created AWS repository
docker push $ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/$REPOSITORY_NAME:$IMAGE_TAG