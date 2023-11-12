set -e

# RUN this script outside of the container!!!

script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}
cd ..

log_file="log/email-async.log"
cat /dev/null > ${log_file}

# creds_str=$(cat config/gmail_user_credentials.json)
# req_str="{\"config\": {\"credentials\": ${creds_str}, \"protocol\": \"gmail\", \"username\": \"xxx@gmail.com\"}, \"n_mails\": 5}"
config_str=$(cat config/outlook.json)
req_str="{\"config\": ${config_str}, \"n_mails\": 5}"
str=$(echo request:${req_str})

docker exec -d kafka kafka-topics.sh \
                                --delete \
                                --bootstrap-server kafka:29092 \
                                --topic requests
docker exec -d kafka kafka-topics.sh \
    --create --topic requests \
    --partitions 1 --replication-factor 1 \
    --bootstrap-server kafka:29092

for i in {1..200}; do
    echo "Sending request ${i}"
    docker exec -it kafka bash -c "echo '$str' | kafka-console-producer.sh \
                                                    --bootstrap-server kafka:29092 \
                                                    --topic requests \
                                                    --property parse.key=true \
                                                    --property key.separator=:"
done

# then run the consumer inside email container
# python email_server_async.py
# and then look at the log to analyze