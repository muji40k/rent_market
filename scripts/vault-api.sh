#! /bin/bash

function login_jwt {
    echo --data "{\"jwt\": \"$1\", \"role\": \"${VAULT_AUTH_ROLE}\"}"
    curl --request POST \
    --data "{\"jwt\": \"$1\", \"role\": \"${VAULT_AUTH_ROLE}\"}" \
    ${VAULT_SERVER_URL}/v1/auth/jwt/login \
    | tee test.txt \
    | jq -r ".auth.client_token"
}

function general {
    curl --header "X-Vault-Token:$1" \
        "${VAULT_SERVER_URL}$2"
}

function get_kv {
    curl --header "X-Vault-Token:$1" \
        "${VAULT_SERVER_URL}/v1/$2" \
        | tee test.txt \
        | jq -r ".data.$3"
}

apk add curl jq

