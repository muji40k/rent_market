#! /bin/bash

function login_jwt {
    curl --request POST \
    --data "{\"jwt\": \"$1\", \"role\": \"${VAULT_AUTH_ROLE}\"}" \
    ${VAULT_SERVER_URL}/v1/auth/jwt/login 2>/dev/null \
    | jq -r ".auth.client_token" 2>/dev/null
}

function general {
    curl --header "X-Vault-Token:$1" \
        "${VAULT_SERVER_URL}$2" 2>/dev/null
}

function get_kv {
    curl --header "X-Vault-Token:$1" \
        "${VAULT_SERVER_URL}/v1/$2" 2>/dev/null \
        | jq -r ".data.$3" 2>/dev/null
}

apk add curl jq

