# Interface for producers and consumers of Kafka messages

## `kafka.requests`

Producer: `backend`
Consumer: `email`

Example message:

```yaml
key: request
value: '{
    "user_id": 1,
    "config":
        {
            "credentials": credentials,
            "protocol": "gmail",
            "username": "xxx@gmail.com"
        },
    "n_mails": 5
}'
```

## `kafka.emails`

Producer: `email`
Consumer: `data`

Example message:

```yaml
key: email
value: '{
    "user_id": 1,
    "email_id": "idxxx",
    "address": "xxx@gmail.com",
    "protocol: "gmail",
    "item": {
        "subject": "subject",
        ...,
        "content": "..."
    }
}'
```

## `kafka.new_emails`

Producer: `data`
Consumer: `parse`

Example message:

```yaml
key: new_email
value: '{
    "user_id": 1,
    "email_id": "idxxx",
    "address": "xxx@gmail.com",
    "protocol: "gmail",
    "item": {
        "subject": "subject",
        ...,
        "content": "..."
    }
}'
```

## `kafka.events`

Producer: `parse`
Consumer: `data`, `backend`

Example message:

```yaml
key: event
value: '{
    "user_id": 1,
    "email_id": "idxxx",
    "address": "xxx@gmail.com",
    "event" {
        "event_type": "...",
        ...,
        "summary": "..."
    }
}
```
