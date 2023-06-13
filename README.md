# Email Wizard

A wizard who structures all your emails for you.

## Description

Codes in this repo should be able to parse raw email contents into structured contents, including email title, sender, content summary, and important dates. Optionally, it may be able to pull emails from modern mail boxes (such as outlook, gmail, 126/163.com, etc.), store parsed results in a calendar format and present them to users.

## Design

### Workflow

![image](assets/workflow.drawio.svg)

## Tasks

### Stage 1

#### A1. retrieve emails

- [x] use IMAP/POP3 protocol to retrieve emails from various mailboxes
  - [x] support 126/163.com mailboxes (with IMAP)
  - [x] support outlook mailboxes

#### A2. clean email into plain texts

- [x] retrieve raw email contents from retrieved emails
- [x] transform raw email contents into plain text
  - [x] discard image (for now)
  - [x] extract and keep plain text content
  - [x] extract email subject, sender, date and time

#### B1. parse plain text for summary and dates

- [x] construct prompts for ChatGPT/GPT
  - [x] build a small dataset of emails and the structured summary
  - [ ] explore different models and different prompts
    - [x] set up evaluation metrics
    - [x] concurrent ChatGPT requests
- [x] use OpenAI API or other models to get structured summary
  - [x] obtain model output from prompts
  - [x] parse model output into structured summary

#### C1. store summary to database

- [ ] store summary with user id and email id as primary key to a database
- [ ] provide retrieval api to retrieve summary from database for a single user

#### C2. calendar view UI

- [ ] frontend UI design and development
  - [ ] UI wireframe
  - [ ] build UI with React
- [ ] connect to backend api to retrieve structured summary
  - [ ] develop backend api and deploy it to a server

### Stage 2

#### A1. retrieve emails

- [ ] support gmail mailboxes

#### A2. clean email into plain texts

- [ ] handle html content

#### C3. update calendar

- [ ] add/remove events in user calendar based on user actions
  - [ ] find libraries to access user calendar (possibly outlook calendar)
  - [ ] add links to user calendar in UI
