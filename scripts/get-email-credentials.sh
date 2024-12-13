#! /bin/bash

function get {
    vault kv get -field "$2" "test/email/$1"
}

echo EMAIL_2FA_HOST="$(get server host)"
echo EMAIL_2FA_PORT="$(get server smtp)"
echo EMAIL_2FA_SENDER="$(get sender name)"
echo EMAIL_2FA_SENDER_PASSWORD="$(get sender password)"
echo EMAIL_2FA_SENDER_EMAIL="$(get sender email)"

echo TEST_IMAP_PORT="$(get server imap)"
echo TEST_RECEPIENT_EMAIL="$(get recepient email)"
echo TEST_RECEPIENT_USERNAME="$(get recepient name)"
echo TEST_RECEPIENT_PASSWORD="$(get recepient password)"

