#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
GREEN="\033[0;32m"
RED="\033[0;31m"
RESET="\033[0m"

pass=0
fail=0

check() {
    local label=$1
    local expected=$2
    local actual=$3

    if [ "$actual" -eq "$expected" ]; then
        echo -e "${GREEN}[PASS]${RESET} $label (HTTP $actual)"
        ((pass++))
    else
        echo -e "${RED}[FAIL]${RESET} $label (expected $expected, got $actual)"
        ((fail++))
    fi
}

# --- Auth ---
echo "=== Auth ==="

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","password":"testpass1","nom":"Test","prenom":"User","telephone":"0600000000"}')
check "POST /auth/register" 201 "$STATUS"

LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","password":"testpass1"}')
LOGIN_STATUS=$(echo "$LOGIN_RESPONSE" | grep -o '"[0-9]\{3\}"' | head -1)
LOGIN_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","password":"testpass1"}')
check "POST /auth/login" 200 "$LOGIN_STATUS"

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d: -f2 | tr -d '"')

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Token non récupéré, abandon des tests peers${RESET}"
    exit 1
fi

AUTH_HEADER="Authorization: Bearer $TOKEN"

# --- Autonomous Systems ---
echo ""
echo "=== Autonomous Systems ==="

AS1_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/as/" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"asn":65001,"name":"AS-65001","router_id":"10.0.0.11","description":"AS local 1"}')
CREATE_AS1_STATUS=$(echo "$AS1_RESPONSE" | tail -1)
AS1_RESPONSE=$(echo "$AS1_RESPONSE" | head -1)
check "POST /as/ (AS65001)" 201 "$CREATE_AS1_STATUS"
AS1_ID=$(echo "$AS1_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
AS1_ID=${AS1_ID:-1}

AS2_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/as/" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{"asn":65002,"name":"AS-65002","router_id":"10.0.0.12","description":"AS local 2"}')
CREATE_AS2_STATUS=$(echo "$AS2_RESPONSE" | tail -1)
AS2_RESPONSE=$(echo "$AS2_RESPONSE" | head -1)
check "POST /as/ (AS65002)" 201 "$CREATE_AS2_STATUS"
AS2_ID=$(echo "$AS2_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
AS2_ID=${AS2_ID:-2}

# --- Peers ---
echo ""
echo "=== Peers ==="

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/peers/all" \
    -H "$AUTH_HEADER")
check "GET /peers/all" 200 "$STATUS"

# Peer 1 : AS65001 → voisin AS65002
PEER1_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/peers/create" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d "{\"local_as_id\":$AS1_ID,\"remote_asn\":65002,\"peer_ip\":\"10.0.0.12\",\"description\":\"AS65001 -> AS65002\"}")
CREATE_PEER1_STATUS=$(echo "$PEER1_RESPONSE" | tail -1)
PEER1_RESPONSE=$(echo "$PEER1_RESPONSE" | head -1)
check "POST /peers/create (AS65001->AS65002)" 201 "$CREATE_PEER1_STATUS"
PEER1_ID=$(echo "$PEER1_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
PEER1_ID=${PEER1_ID:-1}

# Peer 2 : AS65002 → voisin AS65001
PEER2_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/peers/create" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d "{\"local_as_id\":$AS2_ID,\"remote_asn\":65001,\"peer_ip\":\"10.0.0.11\",\"description\":\"AS65002 -> AS65001\"}")
CREATE_PEER2_STATUS=$(echo "$PEER2_RESPONSE" | tail -1)
PEER2_RESPONSE=$(echo "$PEER2_RESPONSE" | head -1)
check "POST /peers/create (AS65002->AS65001)" 201 "$CREATE_PEER2_STATUS"
PEER2_ID=$(echo "$PEER2_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
PEER2_ID=${PEER2_ID:-2}

PEER_ID=$PEER1_ID

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/peers/$PEER_ID" \
    -H "$AUTH_HEADER")
check "GET /peers/:peerID" 200 "$STATUS"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/peers/$PEER_ID/sessions" \
    -H "$AUTH_HEADER")
check "GET /peers/:peerID/sessions" 200 "$STATUS"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/peers/sync" \
    -H "$AUTH_HEADER")
check "POST /peers/sync" 200 "$STATUS"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/peers/$PEER1_ID" \
    -H "$AUTH_HEADER")
check "DELETE /peers/:peerID (peer1)" 200 "$STATUS"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/peers/$PEER2_ID" \
    -H "$AUTH_HEADER")
check "DELETE /peers/:peerID (peer2)" 200 "$STATUS"

# --- Prefixes ---
echo ""
echo "=== Prefixes ==="

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/bgp/create/prefix" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d "{\"prefix\":\"192.168.1.0/24\",\"asn\":65001,\"next_hop\":\"10.0.0.11\",\"local_pref\":100}")
check "POST /bgp/create/prefix (AS65001)" 201 "$STATUS"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/bgp/create/prefix" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d "{\"prefix\":\"192.168.2.0/24\",\"asn\":65002,\"next_hop\":\"10.0.0.12\",\"local_pref\":100}")
check "POST /bgp/create/prefix (AS65002)" 201 "$STATUS"

# --- Résumé ---
echo ""
echo "=== Résumé ==="
echo -e "${GREEN}PASS: $pass${RESET} | ${RED}FAIL: $fail${RESET}"
