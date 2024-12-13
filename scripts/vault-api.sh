#! /bin/bash

function login_jwt {
    curl \
    --request POST \
    --data "{\"jwt\": \"$1\", \"role\": \"${VAULT_AUTH_ROLE}\"}" \
    ${VAULT_SERVER_URL}/v1/auth/jwt/login | jq -r ".auth.client_token"
}

function general {
    curl -v \
        --header "X-Vault-Token:$1" \
        "${VAULT_SERVER_URL}$2"
}

function get_kv {
    curl \
        --header "X-Vault-Token:$1" \
        "${VAULT_SERVER_URL}/v1/$2" | jq -r ".data.$3"
}

apk add curl jq

