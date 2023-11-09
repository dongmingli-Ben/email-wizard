# Interface for producers and consumers of Kafka messages

## `kafka.requests`

Producer: `backend`
Consumer: `email`

Example message:

```yaml
key: request
value: '{"config": {"credentials": credentials, "protocol": "gmail", "username": "xxx@gmail.com"}, "n_mails": 5}'
```

## `kafka.emails`

Producer: `email`
Consumer: `data`

Example message:

```yaml
key: email
value: '{
    "email_id": "idxxx",
    "address": "xxx@gmail.com",
    "item": {
        "subject": "subject",
        ...,
        "content": "..."
    }
}'
```
