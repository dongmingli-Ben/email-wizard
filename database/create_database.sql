CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    user_secret TEXT NOT NULL,
    user_name TEXT UNIQUE NOT NULL,
    user_password TEXT NOT NULL,
    mailboxes JSON
);

ALTER TABLE public.users 
    ALTER COLUMN mailboxes SET DEFAULT '{}';

CREATE TABLE IF NOT EXISTS emails (
    user_id INT NOT NULL,
    email_id TEXT NOT NULL,
    email_address TEXT NOT NULL,
    mailbox_type VARCHAR(10) NOT NULL,
    email_content TEXT,
    event_ids INT[],
    PRIMARY KEY (email_id, email_address),
    CONSTRAINT fk_emails_user_id FOREIGN KEY (user_id)
        REFERENCES users(user_id)
);

ALTER TABLE public.emails 
    ALTER COLUMN event_ids SET DEFAULT '{}';

CREATE TABLE IF NOT EXISTS events (
    event_id SERIAL PRIMARY KEY,
    email_id TEXT NOT NULL,
    email_address TEXT NOT NULL,
    event_content JSON,
    CONSTRAINT fk_events_email FOREIGN KEY 
        (email_id, email_address)
        REFERENCES emails(email_id, email_address)
);