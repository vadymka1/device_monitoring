#!/bin/bash

# Replace {id} with actual device UUID

echo "📡 Get Device by ID"
curl -X GET http://localhost:8080/devices/{id}
echo -e "\n"

echo "📡 Get Device Status"
curl -X GET http://localhost:8080/devices/{id}/status
echo -e "\n"
