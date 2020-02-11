# Mini Wallet

Mini Wallet API using Golang and Postgresql

## Requirements
- [Install Golang](https://golang.org/doc/install)
- Postgresql

## Get Started
- git clone git@github.com:jacky-htg/mini-wallet.git
- cd mini-wallet
- cp .env.example .env
- edit .env with your environment
- create database (the database name must be match with your environment)
- go mod init mini-wallet
- go run main.go migrate
- go run main.go seed
- go run main.go
- Open http://localhost:8080/api/v1/health

## Testing
- Open Postman and import mini-wallet.postman_collection.json and mini-wallet.postman_environment.json 