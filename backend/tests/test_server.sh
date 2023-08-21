#!/bin/bash

set -e

script_dir="$(cd "$(dirname "$0")" && pwd)"
cd $script_dir

echo "adding user..."
curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -d '{"username": "toymaker", "password": "12345678"}' https://toymaker-ben.online/api/add_user

resp=$(curl -G -d "username=toymaker&password=12345678" https://toymaker-ben.online/api/verify_user)
user_id="$(echo $resp | jq -r '.user_id')"
user_secret="$(echo $resp | jq -r '.user_secret')"

echo "adding mailbox 126.com for user..."
curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -d '{
            "type": "IMAP",
            "userId": '$user_id',
            "userSecret": "'"$user_secret"'",
            "address": "dongmingli_Ben@126.com",
            "imap_server": "imap.126.com",
            "password": "JPOSQUNLDYZXBPRO"
        }' \
    https://toymaker-ben.online/api/add_mailbox

echo "updating event for mailbox 126.com..."
curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -d '{
            "userId": '$user_id',
            "userSecret": "'"$user_secret"'",
            "address": "dongmingli_Ben@126.com",
            "kwargs": {}
        }' \
    https://toymaker-ben.online/api/events
curl -G -d "user_id=$user_id&secret=$user_secret" https://toymaker-ben.online/api/events

echo "adding mailbox outlook for user..."
curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -d '{
            "type": "outlook",
            "userId": '$user_id',
            "userSecret": "'"$user_secret"'",
            "address": "guangtouqiang@outlook.com"
        }' \
    https://toymaker-ben.online/api/add_mailbox

echo "updating event for mailbox outlook..."
auth_token=$(cat test_auth_token.txt)
curl -X POST -H "Content-Type: application/json;charset=UTF-8" \
    -H "Accept: application/json, text/plain, */*" \
    -d '{
            "userId": '$user_id',
            "userSecret": "'"$user_secret"'",
            "address": "guangtouqiang@outlook.com",
            "kwargs": {"auth_token": "'"$auth_token"'"}
        }' \
    https://toymaker-ben.online/api/events
curl -G -d "user_id=$user_id&secret=$user_secret" https://toymaker-ben.online/api/events
