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

#### Example

Request:

config:

```json
{
  "username": "xxx@outlook.com",
  "protocol": "outlook",
  "credentials": {}
}
```

```json
{
  "imap_server": "imap.126.com",
  "username": "xxx@126.com",
  "credentials": {
    "password": "xxx",
    "protocol": "IMAP"
  }
}
```

n_mails:

```
2
```

Response:

```json
{
  "items": [
    {
      "email_id": "AQMkADAwATNiZmYAZC01YTBlLTVlMjMtMDACLTAwCgBGAAADWf5YPRULfUyDmeRFNdjmKgcAnOT21DDXNEqBQujyKJX_pAAAAgEMAAAAnOT21DDXNEqBQujyKJX_pAAAAHzUVt8AAAA=",
      "item": {
        "subject": "API Updates: GPT-3.5 Turbo, Function Calling, Longer Context, Lower Prices",
        "sender": "noreply@email.openai.com",
        "date": "2023-06-13T23:13:32Z",
        "recipient": ["xxx@outlook.com"],
        "content": "More capable, more steerable OpenAI models We're excited to announce a few updates to the OpenAI developer platform. GPT-3.5 Turbo This model has been updated with a new version: gpt-3.5-turbo-0613 which is more steerable with the system message and includes a new capability: function calling. By describing functions in your prompts, the model can intelligently output a JSON object containing arguments to call these functions based on user input — perfect for integrating with other tools or APIs. Learn more in our function calling documentation . Plus enjoy a 25% cost reduction for input tokens on GPT-3.5 Turbo (now $0.0015 per 1K input tokens), effective immediately. Longer Context We're also introducing gpt-3.5-turbo-16k . This model offers four times the context length of the 4k base model and is priced at $0.003 per 1K input tokens and $0.004 per 1K output tokens. Model Transitioning You can begin using the new gpt-3.5-turbo-0613 model today. On June 27th, the stable gpt-3.5-turbo will be automatically upgraded to this new version. If you need more time to transition, you can specify gpt-3.5-turbo-0301 to keep using the older version, which will remain available until September 13th as part of our upgrade and deprecation process . —The OpenAI team OpenAI 548 Market Street, PMB 97273 , San Francisco , CA 94104-5401 Unsubscribe - Unsubscribe Preferences"
      }
    },
    {
      "email_id": "AQMkADAwATNiZmYAZC01YTBlLTVlMjMtMDACLTAwCgBGAAADWf5YPRULfUyDmeRFNdjmKgcAnOT21DDXNEqBQujyKJX_pAAAAgEMAAAAnOT21DDXNEqBQujyKJX_pAAAAHwqDUwAAAA=",
      "item": {
        "subject": "Microsoft account unusual sign-in activity",
        "sender": "account-security-noreply@accountprotection.microsoft.com",
        "date": "2023-06-12T10:10:00Z",
        "recipient": ["xxx@outlook.com"],
        "content": "Microsoft account Unusual sign-in activity We detected something unusual about a recent sign-in to the Microsoft account gu**4@outlook.com . Sign-in details Country/region: China IP address: 218.81.240.163 Date: 6/12/2023 10:10 AM (GMT) Platform: Android Browser: Android Please go to your recent activity page to let us know whether or not this was you. If this wasn't you, we'll help you secure your account. If this was you, we'll trust similar activity in the future. Review recent activity To opt out or change where you receive security notifications, click here . Thanks, The Microsoft account team"
      }
    }
  ]
}
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

### Example

Request:

email:

```json
{
  "subject": "重要通知 ......",
  "sender": "someone@example.com",
  "date": "2023-06-07T13:58:41Z",
  "recipient": ["whoami@outlook.com"],
  "content": "亲爱的 同学： 你好！ 经由学院、书院推荐，......"
}
```

additional_info:

```json
{
  "timezone": "Asia/Shanghai"
}
```

Response:

```json
{
  "events": [
    {
      "event_type": "registration",
      "end_time": "2023-04-06T00:00:00 Asia/Shanghai",
      "summary": "2023大学......",
      "venue": "https://....../vm/YVgulbu.aspx"
    }
  ]
}
```

## Database Services

### Email database

The email database should store users' emails. It should support queries on users.

### User database

The user database should store users' profile.

### Event database

The event database should store the parsed events from users' emails.

| endpoint    | service           | input                          | output                      |
| ----------- | ----------------- | ------------------------------ | --------------------------- |
| add_row     | add a row         | JSON                           | error message & primary key |
| update_item | update an element | JSON                           | error message               |
| delete_row  | delete a row      | JSON                           | error message               |
| query       | query database    | JSON of fields and constraints | JSON                        |

## Backend API Endpoint

### Passing user credentials

After the APIs are redesigned to follow REST, user id is passed in the URL and user secret should be passed in the headers with `X-User-Secret`. Here is an example:

```bash
curl -G -H "X-User-Secret: $user_secret" \
    https://toymaker-ben.online/api/users/$user_id/events
```

### `users/:id/events` endpoint - GET

This API is meant for querying events already in DB.

Endpoint: `https://public-ip:8080/users/:id/events`

Params:

```yaml
user_id: int
user_secret: string
```

Response:

A json string of events (something like the below)

```json
[
  {
    "event_type": "registration",
    "end_time": "2023-04-06T12:00:00 Asia/Shanghai",
    "summary": "2023大学......",
    "venue": "https://....../vm/YVgulbu.aspx"
  }
]
```

### `users/:id/events` endpoint - POST

This API is meant for updating new events of a mailbox of a user.

Endpoint: `https://public-ip:8080/users/:id/events`

Params:

```yaml
address: string
kwargs: JSON # to be deprecated
```

Response:

A json of error message

```json
{ "errMsg": "empty string if there is no error" }
```

### `verify_email` endpoint

This API is meant to verify whether our service can retrieve emails from users' IMAP or POP3 mailboxes.

Endpoint: `https://public-ip:8080/verify_email`

Method: GET

Params:

```yaml
username: string
password: string
imap_server: string # for IMAP
pop3_server: string # for POP3
type: string # must be IMAP or POP3
```

Response:

If it is successful, the response will be similar to `GetEmails` service.
If there is an internal error, it will return:

```json
{
  "errMsg": "error message"
}
```

### `users/:id/mailboxes` endpoint

This API is meant to add users' mailbox info to the user DB.

Endpoint: `https://public-ip:8080/users/:id/mailboxes`

Method: POST

Payload:

```yaml
type: string # mailbox type
address: string # mailbox address
credentials: JSON # other necessary fields such as password, IMAP/POP3 server, auth code, etc
```

Response:

If the request is not successful:

```yaml
errMsg: string
```

### `users/:id/profile` endpoint

This API is meant to get user name and their added mailboxes information for display.

Endpoint: `https://public-ip:8080/users/:id/profile`

Method: GET

No payload.

Response:

If the request is not successful:

```yaml
errMsg: string
```

If the request is successful:

```yaml
user_name: string
user_accounts:
  - username: string # mailbox address
  - protocal: string # mailbox type
```

### `authenticate` endpoint

This API is meant to autheticate user name and password then return the user id and user secret for future API call.

Endpoint: `https://public-ip:8080/authenticate`

Method: POST

Payload:

```yaml
username: string
password: string
```

Response:

If the request is not successful:

```yaml
errMsg: string
```

If the request is successful:

```yaml
user_id: string
user_secret: string
```

### `users` endpoint

This API is meant to add user name and password to backend database.

Endpoint: `https://public-ip:8080/users`

Method: POST

Payload:

```yaml
username: string
password: string
```

Response:

If the request is not successful:

```yaml
errMsg: string
```
