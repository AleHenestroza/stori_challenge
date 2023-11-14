# Setup Gmail as SMTP

If you want to use Gmail as your SMTP server, make sure you complete the following steps:

1. Enable 2FA on your Google account. This step is not required, but it is highly recommended.
2. Generate an app password, as described in [this article](https://support.google.com/accounts/answer/185833). If you don't have 2FA enabled, this step may not be required.
3. Use the following variables in the .env file:
    1. `SMTP_HOST=smtp.gmail.com`
    2. `SMTP_PORT=587`
    3. `SMTP_USERNAME=your_email@gmail.com`
    4. `SMTP_PASSWORD="abcd efgh ijkl mnop"` This is your app password, which is comprised of 16 characters, in groups of 4 separated by spaces.
    5. `SMTP_SENDER="Your Name <your_email@gmail.com>"`
