ghz --insecure --proto ./protos/email.proto \
    --async \
    --timeout=20s \
    --call email.EmailHelper.GetEmails \
    -d '{   "config": "{ \"username\": \"guangtouqiang1964@outlook.com\",    \"protocol\": \"outlook\"}", "n_mails": 5  }' \
    127.0.0.1:50051