# Email Wizard

A wizard who structures all your emails for you and deliver them to you in calendar.

## Disclaimer

Preview! The service is under active development. Things may change quickly and drastically. Security and Privacy may be vulnerable. We currently store all parsed emails in our server. Email parser is implemented with OpenAI's ChatGPT API, which means that your emails will be sent to OpenAI. OpenAI claims that they will not store and use data sent by API but there is not way to guarentee.

## Try It Out Yourself!

Our service is available at [https://www.toymaker-ben.online/](https://www.toymaker-ben.online/).

Note that the service is under development, so it may be inaccessible sometimes.

See the [demo](https://acsweb.ucsd.edu/~dol031/projects/email-wizard)!

## Feature

Extract events from various emails and present them in calendar

- support search
- support Outlook, Gmail, IMAP/POP3 mailboxes
- support different types of events

See our upcoming feature [plan](./docs/plan.md).

## Hosting on your machine (NOT recommended)

For data privacy issues, you may want to host the service on your machines. Here is a brief guide on how to do that:

### Setting up OpenAI API key

```bash
touch smart_parse/setup_openai.sh
echo "export OPENAI_API_KEY=<your-openai-api-key>" > smart_parse/setup_openai.sh
```

### Setting up request forwarding and domain name

Backend endpoints are hosted on port 8080. We use Nginx to forward each request starting with `/api` to port 8080. This part is already configured in `deploy/nginx.conf`.
To use your own machine with your own domain name, modify corresponding Nginx configuration and SSL certificates in `backend/server.go`.

### (Re)Starting the service

```bash
docker compose up
```

For the first time after the service starts, you need to create the Kafka topics. Afterwards, you do not need to do so. Even if the service is down, it will automatically recognize previous topics and data after restart.

Use the script to initialize kafka topics:

```bash
bash backend/kafka_init.sh
```

### Using the service with browsers

The webpage is not built automatically on container starts for development purpose. To build the webpage, use:

```bash
docker exec -d frontend bash -c "cd /mnt/frontend && bash run.sh"
```

Then you can view the webpage using your server url specified in `deploy/nginx.conf`.

### Shutting down the service

```bash
docker compose down
```
