# Database Schema

## Users

| attribute     | type    | constraints       |
| ------------- | ------- | ----------------- |
| user ID (PK)  | integer | not null & unique |
| user secret   | string  | not null          |
| user name     | string  | not null & unique |
| user password | string  | not null          |
| mailboxes     | JSON    | list of mailboxes |

Mailbox: each mailbox should have at least 1. the email address (username); 2. mailbox type (protocol); 3. user credentials (credentials). User credentials may contain the following:

- mailbox password
- mailbox server address, such as the server of IMAP/POP3 mailboxes

## Emails

| attribute          | type                    | constraints |
| ------------------ | ----------------------- | ----------- |
| user ID (FK)       | integer                 | not null    |
| email ID (PK)      | string                  | not null    |
| email address (PK) | string                  | not null    |
| mailbox type       | string                  | not null    |
| email subject      | string                  | -           |
| email sender       | string                  | not null    |
| email recipients   | array of strings        | not null    |
| email datetime     | timestamp with timezone | not null    |
| email content      | string                  | -           |
| event IDs          | array of int            | -           |

Assume each email is associated with a SINGLE user for now.

## Events

| attribute          | type    | constraints   |
| ------------------ | ------- | ------------- |
| user ID (FK)       | integer | not null      |
| event ID (PK)      | integer | not null      |
| email ID (FK)      | string  | not null      |
| email address (FK) | string  | not null      |
| event content      | JSON    | event details |

Event: each event should be one of notification, registration, and activity. The format should follow those returned by [smart parser](../smart_parse/README.md).
