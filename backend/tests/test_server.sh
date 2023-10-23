#!/bin/bash

set -e

script_dir="$(cd "$(dirname "$0")" && pwd)"
cd $script_dir

echo "adding user..."
curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -d '{"username": "toymaker", "password": "12345678"}' \
     https://toymaker-ben.online/api/users

resp=$(curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -d '{"username": "toymaker", "password": "12345678"}' \
     https://toymaker-ben.online/api/authenticate)
user_id="$(echo $resp | jq -r '.user_id')"
user_secret="$(echo $resp | jq -r '.user_secret')"

password=$(cat test_126_password.txt)

echo "adding mailbox 126.com for user..."
curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -H "X-User-Secret: $user_secret" \
    -d '{
            "type": "IMAP",
            "address": "dongmingli_Ben@126.com",
            "credentials": {
                "imap_server": "imap.126.com",
                "password": "'"$password"'"
            }
        }' \
    https://toymaker-ben.online/api/users/${user_id}/mailboxes

echo "updating event for mailbox 126.com..."
curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -H "X-User-Secret: $user_secret" \
    -d '{
            "address": "dongmingli_Ben@126.com",
            "kwargs": {}
        }' \
    https://toymaker-ben.online/api/users/$user_id/events
curl -G -H "X-User-Secret: $user_secret" \
    https://toymaker-ben.online/api/users/$user_id/events

echo "adding mailbox outlook for user..."
curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -H "X-User-Secret: $user_secret" \
    -d '{
            "type": "outlook",
            "address": "guangtouqiang@outlook.com",
            "credentials": {}
        }' \
    https://toymaker-ben.online/api/users/$user_id/mailboxes

echo "updating event for mailbox outlook..."
auth_token=$(cat test_auth_token.txt)
curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -H "X-User-Secret: $user_secret" \
    -d '{
            "address": "guangtouqiang@outlook.com",
            "kwargs": {"auth_token": "'"$auth_token"'"}
        }' \
    https://toymaker-ben.online/api/users/$user_id/events
curl -G -H "X-User-Secret: $user_secret" \
    https://toymaker-ben.online/api/users/$user_id/events
