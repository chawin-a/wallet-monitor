# Wallet Monitor
Monitoring system for tracking wallet transactions using chain explorer.

## Prerequisites
- Docker
- Go 1.21.0
- [go-migrate](https://github.com/golang-migrate/migrate)

## Getting started
1. Clone the repository
```bash

git clone https://github.com/chawin-a/wallet-monitor.git

```
2. Run docker compose to init database
```bash
docker compose up -d
```
3. Create `configs.yml` see `configs/configs.example.yml`
4. Run this script to add wallets to database
```bash
go run cmd/add/main.go
```
5. Start service
```bash
go run cmd/worker/main.go
```