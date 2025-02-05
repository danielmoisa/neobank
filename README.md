# Neobank

Neobank is a modern digital banking platform that allows users to manage their finances, make transactions, and view real-time account balances. It aims to provide an intuitive and secure banking experience, leveraging the latest technology to offer innovative features.

## Features

- **Account Management:** Create and manage multiple accounts.
- **Transaction History:** View recent transactions and detailed statements.
- **Balance Tracking:** Real-time balance updates for checking and savings accounts.
- **Security:** Bank-level security features for user data protection.
- **User Authentication:** Secure login with two-factor authentication.

### Prerequisites

Ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.17 or above)
- [PostgreSQL](https://www.postgresql.org/) (or any other preferred database)
- Git

## Getting Started

These instructions will help you set up the project locally for development and testing.
- run `docker-compose up -d` to setup the posgresql and api docker services
- run `make migrateup`
- `go run main.go` to start the app or `make serve`
- swagger localhost:8080/swagger/index.html#/