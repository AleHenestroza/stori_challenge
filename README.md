# Stori Backend Challenge

This app is the solution to [Stori's Backend Challenge](Tech_Challenge_-_Software_Engineer.pdf). The application was developed using Go and Docker.

## Run the application

In order to run this application, you need to have Docker installed. You can build the application with the following command:

```sh
docker build --label storichallenge .
```

And then you can run the application with:

```sh
docker run -p 4000:4000 -e SMTP_HOST=sandbox.smtp.mailtrap.io -e SMTP_PORT=2525 -e SMTP_USERNAME=<your-username> -e SMTP_PASSWORD=<your-password> -e SMTP_SENDER="Stori Test <no-reply@storitest.com>" storichallenge
```

## SMTP

The application was developed using [mailtrap](https://mailtrap.io). Mailtrap is a free application that allows you to test sending emails by creating a private inbox to which each email will be sent. This means that no email will actually be sent to real email addresses.

You can create a free account and change the environment variables `SMTP_USERNAME` and `SMTP_PASSWORd` with your own credentials.

Alternatively, if you already own a SMTP service, you can also change the `SMTP_HOST`, `SMTP_PORT` and `SMTP_SENDER` with your own.

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

-   [ ] Add Stori logo to email template
-   [ ] Implement a database for storing summaries
-   [ ] Allow saving monthly/account summaries to a database
-   [ ] Refactor Dockerfile for multi-stage building of the application
-   [ ] Send emails to real email accounts
-   [ ] Deploy on the cloud (AWS)
