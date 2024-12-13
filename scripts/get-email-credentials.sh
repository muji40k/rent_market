#! /bin/bash

function get {
    get_kv "$VAULT_TOKEN" "test/email/$1" "$2"
}

export EMAIL_2FA_HOST="$(get server host)"
export EMAIL_2FA_PORT="$(get server smtp)"
export EMAIL_2FA_SENDER="$(get sender name)"
export EMAIL_2FA_SENDER_PASSWORD="$(get sender password)"
export EMAIL_2FA_SENDER_EMAIL="$(get sender email)"

export TEST_IMAP_PORT="$(get server imap)"
export TEST_RECEPIENT_EMAIL="$(get recepient email)"
export TEST_RECEPIENT_USERNAME="$(get recepient name)"
export TEST_RECEPIENT_PASSWORD="$(get recepient password)"

