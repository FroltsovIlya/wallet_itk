To start the project download, then use
docker compose up --build
then you can use 

curl http://localhost:8082/api/v1/wallets/a

to test access to data base 
(also it can be 
http://localhost:8082/api/v1/wallets/b
http://localhost:8082/api/v1/wallets/c
)

you can use the tests. Use it in local or on Docker 
go test ./... -v
-v if you want to see messages not only result

in app alse exists POST method. to use it:

curl -X POST http://localhost:8082/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "walletId": "a",
    "operationType": "DEPOSIT",
    "amount": 50
  }'

you can configure your okruzhenie with .env files
