# Email Wizard

A wizard who structures all your emails for you.

## Description

Codes in this repo should be able to parse raw email contents into structured contents, including email title, sender, content summary, and important dates. Optionally, it may be able to pull emails from modern mail boxes (such as outlook, gmail, 126/163.com, etc.), store parsed results in a calendar format and present them to users.

## Design

### Workflow

![image](../assets/workflow.drawio.svg)

### Architecture

![image](../assets/architecture.drawio.svg)

Architecture for asynchronous services (with kafka streams marked with `/`)

![image](../assets/async-kafka.drawio.svg)

## Tasks

### Stage 1 (basic functionality)

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
- [x] deploy with microservice (together with A1)
  - [x] benchmark microservice performance

#### B1. parse plain text for summary and dates

- [x] construct prompts for ChatGPT/GPT
  - [x] build a small dataset of emails and the structured summary
  - [ ] explore different models and different prompts
    - [x] set up evaluation metrics
    - [x] concurrent ChatGPT requests
- [x] use OpenAI API or other models to get structured summary
  - [x] obtain model output from prompts
  - [x] parse model output into structured summary
- [ ] migrate to langchain?
- [x] deploy with microservice
  - [x] benchmark microservice performance

#### C1. store summary to database

- [x] design API protocol
  - [x] CRUD operations interface
- [x] Implement APIs
  - [x] database design
    - [x] choose database engine
    - [x] design database schema
  - [x] CRUD for email db
    - [x] `add_row`
    - [x] `update_item`
    - [x] `delete_row`
    - [x] `query`
  - [x] CRUD for user db
    - [x] `add_row`
    - [x] `update_item`
    - [x] `delete_row`
    - [x] `query`
  - [x] CRUD for event db
    - [x] `add_row`
    - [x] `update_item`
    - [x] `delete_row`
    - [x] `query`
    <!-- - [ ] `execute` -->
  - [x] testing
- [x] deploy as a microservice
  - [x] testing
  - [ ] benchmark microservice performance

#### C2. API requests

- [x] backend API interface design
  - [x] refresh events, query events, and more
- [x] handle requests by interacting with microservices
  - [x] choose microservice framework
  - [x] implement logic with prepared microservices
    - [x] retrieve emails from user mailboxes
    - [x] extract new emails
    - [x] parse to events
    - [x] store results to database
    - [x] query database for results
    - [x] user registration
- [ ] testing
  - [x] script testing
  - [ ] API tools testing

#### C3. calendar view UI

- [x] frontend UI design and development
  - [x] UI wireframe
  - [x] build UI with React
    - [x] implement main calendar page (use fake user_id and secret for now)
    - [x] implement add email account (use fake autheticate for now)
    - [x] implement login page (get user_id and secret)
    - [x] implement register page (get user_id and secret)
    - [x] implement logout mechanism
    - [x] implement intro page
    - [x] implement account authorization grants
    - [x] integrate with backend API
      - [x] integrate with updated APIs
- [x] deploy to production server

### Stage 2 (minimum useful requirements)

#### A1. retrieve emails

- [x] support gmail mailboxes

#### C2. API requests

- [x] refine API interface
  - [x] migrate to RESTful APIs

#### C3. calendar view UI

- [x] refine API interface
- [x] add remove mailbox functionality
  - [x] revoke authorized mailbox access
- [x] add event search (pure frontend for now)
- [x] make the UI more beautiful
  - [x] iterate with @mui

### Stage 3 (performance optimization)

Following the [best practices for building scalable backend system](https://www.linkedin.com/pulse/10-best-practices-building-scalable-backend-systems-sandeep-shah/)

#### A1. retrieve emails

- [x] refactor the microservice to listen from kafka logs
- [x] async / multi-threading retrieving emails

#### B1. parse plain text for summary and dates

- [x] optimize QPS
  - [x] refactor the microservice to listen from kafka logs
  - [x] async / multi-threading calling LLM API

#### C1. store summary to database

- [x] optimize query

#### C2. API requests

- [ ] implement caching
- [ ] implement auto-scaling
- [ ] monitor systems
- [ ] use CDN
- [x] improve responsiveness via message queues
  - [x] set up kafka servers
  - [x] refactor retrieve-parse-display flow to be asynchronous
    - [x] listen to kafka logs
    - [x] return events via websocket

#### C3. calendar view UI

- [ ] faster source retrival

### Stage 4 (more features and supports)

#### A1. retrieve emails

- [ ] update outlook oauth to auth code flow
- [ ] support exchange protocol

#### A2. clean email into plain texts

- [ ] handle html content

#### B1. parse plain text for summary and dates

- [ ] better precision and recall
- [ ] timezone
- [ ] refine API interface

#### C1. store summary to database

- [ ] refine API interface

#### C2. API requests

- [ ] API testing
- [ ] microservices testing
- [ ] migrate to full docker orchestration
  - [ ] migrate to docker compose
  - [ ] deploy with k8s (depending on machine resources)

#### C3. calendar view UI

- [ ] update outlook auth to auth code flow
- [ ] add forget password functionaility
- [ ] make the UI more beautiful

#### C4. update calendar

- [ ] add/remove events in user calendar based on user actions
  - [ ] find libraries to access user calendar (possibly outlook calendar)
  - [ ] add links to user calendar in UI

## Environments

### Benchmarking

Use `ghz` with Docker for benchmarking gRPC microservice performance:

```bash
DOCKER_BUILDKIT=1 sudo docker build --output=/usr/local/bin --target=ghz-binary-built https://github.com/bojand/ghz.git
```
