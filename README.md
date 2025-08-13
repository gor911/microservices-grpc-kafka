# microservices-grpc-kafka

There are 2 microservices [order-service](order-service) and [inventory-service](inventory-service).

Order service functionality:

* Accept POST /orders request to create order.
* gRPC call to Inventory Service to get products info.
* Save order (status = PENDING) and order items
* Insert event into outbox in same transaction
* Another goroutine runs and publish messages to kafka from outbox table

Inventory service functionality:

* gRPC server
* Returns products info
* Consume messages from kafka


Nice to have processed_events table to prevent double consuming from Kafka.

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
