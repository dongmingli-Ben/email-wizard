# API Definitions

## `email` service

### GetEmails

Request:

```
config: json string
n_mails: int
```

Response:

```
message: json string
```

## `parse` service

### ParseEmail

Request:

```
email: json string
additional_info: json string
```

Response:

```
message: json string
```