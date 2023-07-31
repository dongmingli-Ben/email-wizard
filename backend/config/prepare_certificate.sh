cd "$(dirname "$0")"

openssl pkcs12 -in "47-243-42-37-8080-iis-0726224204.pfx" -nocerts -out key.pem -nodes
openssl pkcs12 -in "47-243-42-37-8080-iis-0726224204.pfx" -nokeys -out test.pem
openssl rsa -in key.pem -out test.key