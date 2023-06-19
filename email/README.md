# Email Retriver

This service is responsible for retrieving emails from users accounts and perform basic cleaning.

## Environment

Build the docker image:

```bash
cd email
bash build_image.sh
```

Build the docker container:

```bash
bash build_container.sh
# exit the container
docker restart email # restart the container
# then you can attach to the running container with vscode
```

## Running

### As Native Script

Place your email account configuration to `config/`. For example,

#### Config for Outlook

```json
{
    "username": "email address",
    "protocol": "outlook"
}  
```

#### Config for IMAP Mailbox

```json
{
  "imap_server": "imap.126.com",  // imap server
  "username": "email address",
  "password": "password for IMAP",
  "protocol": "IMAP"
}
```

Then run with,

```bash
python main.py --config config/email.json --n-mails 10
```

You may be prompted to grant permission in a popup window.

### As a Service

A script to setup email retrival as a microservice through gRPC is prepared.

To set up the service:

```bash
cd email
bash run_server.py
```

The API documentation about the service is available at [docs/API.md](../docs/API.md). Please check it out. Or use the test script to see how it works.

To test the service is up and running:

```bash
cd email
bash test_server.sh  # after you have setup the email config under config
```

If the service is healthy and responding, you should be able to see a list of emails in JSON format.