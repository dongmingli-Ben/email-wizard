# Database Schema

## Emails

| attribute     | type    | constraints       |
| ------------- | ------- | ----------------- |
| user ID       | integer | not null & unique |
| email ID      | string  | not null & unique |
| email address | string  | not null          |
| mailbox type  | enum    | not null          |
| email content | string  | -                 |

## Users

| attribute     | type    | constraints       |
| ------------- | ------- | ----------------- |
| user ID       | integer | not null & unique |
| user secret   | string  | not null & unique |
| user name     | string  | not null & unique |
| user password | string  | not null & unique |
| mailboxes     | JSON    | list of mailboxes |

Mailbox: each mailbox should have at least 1. the email address; 2. mailbox type. And can optionally have the following field depending on the mailbox type:

- mailbox password
- mailbox server address, such as the server of IMAP/POP3 mailboxes

## Events

| attribute | type    | constraints       |
| --------- | ------- | ----------------- |
| user ID   | integer | not null & unique |
| email ID  | string  | not null& unique  |
| events    | JSON    | list of events    |

Event: each event should be one of notification, registration, and activity. The format should follow those returned by [smart parser](../smart_parse/README.md).
