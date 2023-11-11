# Stori Backend Challenge

This app is the solution to [Stori's Backend Challenge](Tech_Challenge_-_Software_Engineer.pdf). The application was developed using Go and Docker.

## SMTP

The application was developed using [mailtrap](https://mailtrap.io). Mailtrap is a free application that allows you to test sending emails by creating a private inbox to which each email will be sent. This means that no email will actually be sent to real email addresses.

You can create a free account and change the environment variables `SMTP_USERNAME` and `SMTP_PASSWORd` with your own credentials.

Alternatively, if you already own a SMTP service, you can also change the `SMTP_HOST`, `SMTP_PORT` and `SMTP_SENDER` with your own.

## Run the application

This application is packaged with a docker-compose solution. In order to run the application, you need to have [Docker](https://docs.docker.com/engine/) and [Docker Compose](https://docs.docker.com/compose/). Docker Compose is integrated with the new versions of [Docker Desktop](https://docs.docker.com/desktop/). If you are not using Docker Desktop, you can install both Docker and Docker Compose using a package manager like [Homebrew](https://brew.sh/), or native package managers like `apt`, `dnf`, `pacman` (for Linux distributions).

To run the application, you need to configure an `.env` file with the necessary environment variables. An `.env_example` file is provided. After setting the correct variables, you can simply run the application with the following command:

```sh
docker-compose up
```

## Endpoints

```sh
curl  -X POST \
  'http://localhost:4000/v1/transactions/summary' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "email": "alehenestroza@gmail.com",
  "name": "Alejandro Henestroza"
}'
```

This endpoint will process the provided `txns.csv` file and will send an email with the Account Summary, as detailed in the PDF.

## To-Do

-   [X] Add Stori logo to email template
-   [ ] Implement a database for storing summaries
-   [ ] Allow saving monthly/account summaries to a database
-   [ ] Refactor Dockerfile for multi-stage building of the application
-   [ ] Send emails to real email accounts
-   [ ] Deploy on the cloud (AWS)
