set -e 

resp=$(curl -s -G -d "username=toymaker&password=123456" https://toymaker-ben.online/api/verify_user)
user_id="$(echo $resp | jq -r '.user_id')"
user_secret="$(echo $resp | jq -r '.user_secret')"

# search
curl -s -G -d "user_id=$user_id&user_secret=$user_secret&query=f" https://toymaker-ben.online/api/events

echo "benchmarking elastic search ..."
ab -n 100 -c 10 "https://toymaker-ben.online/api/events?user_id=$user_id&user_secret=$user_secret&query=f"

echo "benchmarking direct retrival from postgresql ..."
ab -n 100 -c 10 "https://toymaker-ben.online/api/events?user_id=$user_id&user_secret=$user_secret"