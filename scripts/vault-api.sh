#! /bin/bash

function login_jwt {
    curl --request POST \
    --data "{\"jwt\": \"$1\", \"role\": \"${VAULT_AUTH_ROLE}\"}" \
    ${VAULT_SERVER_URL}/v1/auth/jwt/login 2>/dev/null\
    | jq -r ".auth.client_token"
}

function general {
    curl --header "X-Vault-Token:$1" \
        "${VAULT_SERVER_URL}$2" 2>/deb/null
}

function get_kv {
    curl --header "X-Vault-Token:$1" \
        "${VAULT_SERVER_URL}/v1/$2" 2>/dev/null \
        | jq -r ".data.$3"
}

apk add curl jq

