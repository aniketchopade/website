#!/bin/bash
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push $DOCKER_USERNAME/site
pip3 install --upgrade --user awscli
export PATH=~/.local/bin:$PATH
aws --version