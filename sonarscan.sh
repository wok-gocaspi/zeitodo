#!/bin/bash

go test ./... -coverprofile=coverage.out
sonar-scanner
