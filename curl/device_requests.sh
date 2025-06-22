#!/bin/bash

echo "🔍 Get Device by ID"
curl -s http://localhost:8080/devices/11111111-1111-1111-1111-111111111111 | jq

echo "📊 Get Device Status"
curl -s http://localhost:8080/devices/11111111-1111-1111-1111-111111111111/status | jq

echo "🧪 Check Device Status (Manual Trigger)"
curl -s http://localhost:8080/devices/11111111-1111-1111-1111-111111111111/check | jq
