# Go Digital Wallet

## Running migrations
```
docker-compose up -d
make migrate
```

## Creating migration
```
migrate create -ext sql -dir database/migration/ -seq migration_name
```

## Running the App
```
cp .env.example .env
docker-compose up -d
make run
```

## API Endpoints
- Register: `/v1/create_user`
- Get Balance: `/v1/balance_read`
- Get User Top Transactions: `/v1/top_transactions_per_user`
- Get Top Transacting Users: `/v1/top_users`
- Balance Top Up: `/v1/balance_topup`
- Transfer: `/v1/transfer`

Postman collection also available on `postman` folder