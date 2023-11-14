# Stori Backend Challenge

## Table of Contents

1. [Introduction](#introduction)
2. [SMTP](#smtp)
3. [Run the application](#run-the-application)
4. [Continuous Integration with GitHub Actions](#continuous-integration-with-github-actions)
    1. [Key Workflow Steps](#key-workflow-steps)
5. [Endpoints](#endpoints)
    1. [Send summary from local .csv file](#send-summary-from-local-csv-file)
    2. [Register user](#register-user)
    3. [Activate user](#activate-user)
    4. [Request activation token](#request-activation-token)
    5. [Authenticate](#authenticate)
    6. [Save transaction](#save-transaction)
    7. [Save transactions in bulk](#save-transactions-in-bulk)
    8. [Send email summary - v2](#send-email-summary---v2)

## Introduction <a name="introduction"></a>

This app is the solution to [Stori's Backend Challenge](Tech_Challenge_-_Software_Engineer.pdf). The application was developed using Go 1.21, PostgreSQL 16.0, Docker and docker-compose.

## SMTP <a name="smtp"></a>

The application was developed using [mailtrap](https://mailtrap.io). Mailtrap is a free application that allows you to test sending emails by creating a private inbox to which each email will be sent. This means that no email will actually be sent to real email addresses.

You can create a free account and change the environment variables `SMTP_USERNAME` and `SMTP_PASSWORd` with your own credentials.

Alternatively, you can use a real SMTP server (like Gmail). In order to send real emails, you should also change `SMTP_HOST`, `SMTP_PORT` and `SMTP_SENDER` with your own. In the [.env_example](.env_example) file, there is an example of how you can use Gmail SMTP server and send real emails.

Step-by-step instructions can be found using the following links:

-   [Mailtrap](smtp_setup/mailtrap.md)
-   [Gmail](smtp_setup/gmail.md)

> [!NOTE]
> If using Gmail as your SMTP server, please make sure you have read [this article](https://support.google.com/accounts/answer/185833). In it, it is explained how to generate an app password. Regular passwords (the one used for logins) won't work.

## Run the application <a name="run-the-application"></a>

This application is packaged with a docker-compose solution. In order to run the application, you need to have [Docker](https://docs.docker.com/engine/) and [Docker Compose](https://docs.docker.com/compose/). Docker Compose is integrated with the new versions of [Docker Desktop](https://docs.docker.com/desktop/). If you are not using Docker Desktop, you can install both Docker and Docker Compose using a package manager like [Homebrew](https://brew.sh/), or native package managers like `apt`, `dnf`, `pacman` (for Linux distributions).

To run the application, you need to configure an `.env` file with the necessary environment variables. An `.env_example` file is provided. After setting the correct variables, you can simply run the application with the following command:

```sh
docker-compose up
```

> [!NOTE]  
> If you installed Docker by installing Docker Desktop, Compose will already be packaged with it, and the command `docker-compose` won't be valid. Instead, you can use the `docker compose` command.

> [!IMPORTANT]  
> The .env file MUST be generated with the appropriate variables. Failing to do so will prevent the application from functioning correctly.

> [!IMPORTANT]  
> In some cases, MailTrap configuration will difer from inbox to inbox in more than just the username and password. In some cases, an alternative port of 25 will need to be used. Before setting the environment variables, and if using Mailtrap as the SMTP server, please double check that the host and the port are also correct.

> [!WARNING]
> In some cases, if a PostgreSQL image has already been downloaded and run, there may be a volume of data that can interfere with database initialization. In order to fix this issue, you can run `docker-compose down -v` and then `docker-compose up`.

## Continuous Integration with GitHub Actions <a name="continuous-integration-with-github-actions"></a>

This repository utilizes GitHub Actions for Continuous Integration (CI) purposes. The workflow configurations are defined in the `.github/workflows/` directory, managing tasks such as code compilation, automated testing, and more.

### Key Workflow Steps <a name="key-workflow-steps"></a>

1. **Checkout Code (checkout.yml):** This step ensures that the latest code from the main branch is pulled into the CI environment.

2. **Environment Setup (go.yml):** Copies the `.env-ci` file to `.env` for configuring the environment variables needed during the build and test processes.

3. **Build Container (docker-compose.yml):** Uses the [isbang/compose-action](https://github.com/isbang/compose-action) action to build the Docker containers defined in the `docker-compose.yml` file. The `--volumes` flag ensures that any associated volumes are removed during the process.

4. **Set Up Go (go.yml):** Configures the Go runtime environment using the [actions/setup-go](https://github.com/actions/setup-go) action with Go version 1.21.

5. **Build Application (go.yml):** Compiles the Go application using the `go build` command, ensuring a verbose output for detailed build information.

6. **Test Application (go.yml):** Executes the `go test` command to run the suite of tests associated with the Go application, providing verbose output for comprehensive test results.

## Endpoints <a name="endpoints">

A [Postman Collection](collection.json) is provided to facilitate testing.

### Send summary from local .csv file <a name="send-summary-from-local-csv-file"></a>

```sh
curl  -X POST \
  'http://localhost:4000/v1/transactions/summary' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "email": "alehenestroza@gmail.com",
  "name": "Alejandro Henestroza"
}'
```

This endpoint will process the `txns.csv` file and will send an email with the Account Summary, as detailed in the PDF.

### Register user <a name="register-user"></a>

```sh
curl  -X POST \
  'http://localhost:4000/v2/users/register' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "name": "user user",
  "email": "username@email.com",
  "password": "password"
}'
```

This endpoint will create the user in the database and return an activation token. Users need to be activated before performing any operation (save transaction, receive summary).

### Activate user <a name="activate-user"></a>

```sh
curl  -X POST \
  'http://localhost:4000/v2/users/activate' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "token": "JMEXK44ITK7WIVGEMDIXZYQOPY"
}'
```

With this endpoint, you can activate your user. The token must be an activation token, generated by registering a new user, or by requesting a new token (in case of expiration).

### Request activation token <a name="request-activation-token"></a>

```sh
curl  -X POST \
  'http://localhost:4000/v2/auth/activation' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "email": "username@email.com"
}'
```

If a new user was registered, but failed to activate before the activation token expired, they can request a new activation token with this endpoint.

### Authenticate <a name="authenticate"></a>

```sh
curl  -X POST \
  'http://localhost:4000/v2/auth/authenticate' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "email": "username@email.com",
  "password": "password"
}'
```

In order to get an auth token, an activated user needs to send a request to this endpoint. Then they need to use the token as a Bearer token for all future requests.

### Save transaction <a name="save-transaction"></a>

```sh
curl  -X POST \
  'http://localhost:4000/v2/transactions' \
  --header 'Accept: */*' \
  --header 'Authorization: Bearer JCB5JYNLXFTUGLQ437NFO3NPLY' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "transaction_date": "2023/11/09",
  "amount": -23.57
}'
```

This will create a new transaction and save it to the database. The transaction will be assigned a userId, based on the token provided (only activated users will be valid).

### Save transactions in bulk <a name="save-transactions-in-bulk"></a>

```sh
curl  -X POST \
  'http://localhost:4000/v2/transactions/bulk' \
  --header 'Accept: */*' \
  --header 'Authorization: Bearer JCB5JYNLXFTUGLQ437NFO3NPLY' \
  --form 'file=@/path/to/your/file/txns.csv'
```

With this endpoint, you can save multiple transactions at once, by uploading a .csv file, like [the one provided](txns.csv) in this repository.

### Send email summary - v2 <a name="send-email-summary---v2"></a>

```sh
curl  -X POST \
  'http://localhost:4000/v2/transactions/summary' \
  --header 'Accept: */*' \
  --header 'Authorization: Bearer JCB5JYNLXFTUGLQ437NFO3NPLY'
```

This endpoint works exactly as the first one, it will send an email with the account summary, but instead of reading a .csv file, it will look for all the user' transactions in the DataBase.
