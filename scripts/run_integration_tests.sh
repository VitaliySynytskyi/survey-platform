#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Colors for better output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting integration tests for Survey App...${NC}"

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: docker-compose is not installed. Please install docker-compose first.${NC}"
    exit 1
fi

# Build and start all services in detached mode
echo -e "${BLUE}Starting all services...${NC}"
docker-compose up -d

# Wait for services to be fully operational
echo -e "${BLUE}Waiting for services to start...${NC}"
sleep 10

# Run integration tests
echo -e "${BLUE}Running integration tests...${NC}"
go test -v ./integration_tests/...

# Optional: Run API tests with Newman if installed
if command -v newman &> /dev/null; then
    echo -e "${BLUE}Running API tests with Newman...${NC}"
    newman run docs/survey-app-api-tests.postman_collection.json -e docs/local-environment.postman_environment.json
else
    echo -e "${BLUE}Newman not found. Skipping API tests. Install Newman with 'npm install -g newman' to run API tests.${NC}"
fi

# Cleanup: Stop and remove containers
echo -e "${BLUE}Stopping all services...${NC}"
docker-compose down

echo -e "${GREEN}Integration tests completed!${NC}" 