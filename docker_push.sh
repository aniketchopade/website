#!/bin/bash
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push $DOCKER_USERNAME/site
pip install --upgrade pip
pip install --upgrade --user awscli
export PATH=~/.local/bin:$PATH
aws ecs update-service --cluster my-cluster --service aniketswebsite --force-new-deployment --region us-east-1 