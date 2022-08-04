#!/bin/bash

dockerImage=""
azureRepoURL=""
azureRepoUsername=""
azureRepoPassword=""


docker build --tag $dockerImage .
docker login $azureRepoURL -u=$azureRepoUsername -p=$azureRepoPassword
docker push $dockerImage