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

## Endpoint
- GET /api/v1/health
- POST /api/v1/login
- POST /api/v1/init
- GET /api/v1/wallet
- POST /api/v1/wallet
- PATCH /api/v1/wallet
- POST /api/v1/wallet/deposits
- POST /api/v1/wallet/withdrawals

## Testing
- Open Postman 
- Import mini-wallet.postman_collection.json 
- Import mini-wallet.postman_environment.json
- Get Token by POST /login . Username : rijalEwallet and Password : 12345678
- Copy Token from response of POST /login to mini-wallet environment token
- Test all Endpoint 
