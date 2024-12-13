#! /bin/bash

function get {
    get_kv "$VAULT_TOKEN" "test/email/$1" "$2"
}

echo export EMAIL_2FA_HOST="$(get server host)"
echo export EMAIL_2FA_PORT="$(get server smtp)"
echo export EMAIL_2FA_SENDER="$(get sender name)"
echo export EMAIL_2FA_SENDER_PASSWORD="$(get sender password)"
echo export EMAIL_2FA_SENDER_EMAIL="$(get sender email)"

echo export TEST_IMAP_PORT="$(get server imap)"
echo export TEST_RECEPIENT_EMAIL="$(get recepient email)"
echo export TEST_RECEPIENT_USERNAME="$(get recepient name)"
echo export TEST_RECEPIENT_PASSWORD="$(get recepient password)"

