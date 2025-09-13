#!/bin/bash

# API Testing Script for Swiflet Backend

BASE_URL="http://localhost:8080/v1"
TOKEN=""

echo "=== Swiflet Backend API Test ==="

# Test Health Check
echo -e "\n1. Testing Health Check..."
curl -s "$BASE_URL/../health" | jq .

# Test User Registration
echo -e "\n2. Testing User Registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "phone": "+628123456789",
    "address": "Jakarta, Indonesia"
  }')

echo $REGISTER_RESPONSE | jq .

# Extract token from registration response
TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.token.token')

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo "Registration successful. Token: ${TOKEN:0:20}..."
    
    # Test User Login
    echo -e "\n3. Testing User Login..."
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
      -H "Content-Type: application/json" \
      -d '{
        "email": "test@example.com",
        "password": "password123"
      }')
    
    echo $LOGIN_RESPONSE | jq .
    
    # Update token from login response
    TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token.token')
    
    # Test Protected Endpoints
    echo -e "\n4. Testing Protected Endpoints..."
    
    # List Users
    echo -e "\n4a. List Users:"
    curl -s -X GET "$BASE_URL/users" \
      -H "Authorization: Bearer $TOKEN" | jq .
    
    # Get User by ID
    echo -e "\n4b. Get User by ID:"
    curl -s -X GET "$BASE_URL/users/1" \
      -H "Authorization: Bearer $TOKEN" | jq .
    
    # List Articles (placeholder)
    echo -e "\n4c. List Articles:"
    curl -s -X GET "$BASE_URL/articles" \
      -H "Authorization: Bearer $TOKEN" | jq .
    
    # List IoT Devices (placeholder)
    echo -e "\n4d. List IoT Devices:"
    curl -s -X GET "$BASE_URL/iot-devices" \
      -H "Authorization: Bearer $TOKEN" | jq .
    
    # List Sensors (placeholder)
    echo -e "\n4e. List Sensors:"
    curl -s -X GET "$BASE_URL/sensors" \
      -H "Authorization: Bearer $TOKEN" | jq .
    
else
    echo "Registration failed. Cannot proceed with protected endpoint tests."
fi

# Test Unauthorized Access
echo -e "\n5. Testing Unauthorized Access..."
curl -s -X GET "$BASE_URL/users" | jq .

echo -e "\n=== API Testing Complete ==="