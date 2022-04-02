# Git Mirror
Provides simple mechanism to mirror the Git Repo received from webhooks configured.

## Requirements
Packages required for development, specific for Ubuntu 20.04

    apt-get install -y build-essential
    snap install go

Packages required for test validation, specific for Ubuntu 20.04

    snap install ngrok

## Building
Simple build procedure

    make

## Test and Validation
If you are using Ngrok, sign up for an account.  Then generate an auth token.  Once you have the token, configure Ngrok to use your token.

    ngrok authtoken <token>

Start the ngrok on port 4000, which is the default port Git-Mirror runs on.  If you have changed the port, adjust the port below as needed.

    ngrok http 4000

Now run the Git-Mirror service

    ./build/git-mirror-linux-amd64 --storage-path=/tmp/repo

## NOTES
This currently is built with emphasis on supporting Bitbucket for now only.  We review other Git Hosting Platform later.

## References
Bitbucket Webhook Payloads - https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#EventPayloads-Push
