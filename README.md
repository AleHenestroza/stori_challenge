# Stori Backend Challenge

This app is the solution to [Stori's Backend Challenge](Tech_Challenge_-_Software_Engineer.pdf). The application was developed using Go and Docker.

## Endpoints

```sh
curl http://localhost:4000/v1/transactions/summary
```

This endpoint will process the provided `txns.csv` file and will send an email with the Account Summary, as detailed in the PDF.

## To-Do

- [ ] Add Stori logo to email template
- [ ] Implement a database for storing summaries
- [ ] Allow saving monthly/account summaries to a database
- [ ] Refactor Dockerfile for multi-stage building of the application
- [ ] Send emails to real email accounts
- [ ] Deploy on the cloud (AWS)
