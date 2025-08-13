# microservices-grpc-kafka

There are 2 microservices [order-service](order-service) and [inventory-service](inventory-service).

To run microservices use:

`docker compose up -d`

 To create order use CURL:

`curl --location 'http://localhost:8080/orders' 
--header 'Content-Type: application/json' 
--header 'Accept: application/json' 
--data '{
  "user_id": 3,
  "items": [
    {
      "product_id": 1,
      "qty": 2
    },
    {
      "product_id": 2,
      "qty": 4
    }
  ]
}'`
