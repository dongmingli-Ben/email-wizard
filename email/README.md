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