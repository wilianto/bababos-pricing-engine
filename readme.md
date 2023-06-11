## Setup
```
docker-compose up -d
mkdir bin
go build -o bin ./...
./bin/bababos-pricing-engine db:migrate
./bin/bababos-pricing-engine db:seed
./bin/bababos-pricing-engine start
```

## Example Request
```
curl --location 'localhost:8080/price?customer_id=12&sku_id=12345&qty=100'
```

## Example Response
```
{
    "suggested_price": 11500,
    "recommended_supplier_id": 19
}
```
        