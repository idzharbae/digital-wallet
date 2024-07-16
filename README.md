# Go Digital Wallet

## Running migrations
```
docker-compose up -d

```

## Creating migration
```
migrate create -ext sql -dir database/migration/ -seq migration_name
```

## Running the App
```
docker-compose up -d
make run
```

