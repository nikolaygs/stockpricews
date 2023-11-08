# Overview
The service provides a simplistic RESTFUL API that calculates the maximum profit that could have been realized by trading 
specific stock in a given historical time slice.

### Endpoints
The service exposes just 1 endpoint:
`GET /maxprofit` that requires three query params in order to return a response:
* `stock` - the symbol of the stock (string with length between 1-4 chars)
* `begin` - the begin date point of the time slice (in unix secs)
* `end` - the end date point of the time slice (in unix secs)

### Sample usage:
```curl "http://localhost:8080/maxprofit?symbol=UBER&begin=1696934700&end=1699443780"```

### Response
* If service computes the client query it returns 200 OK with the following body structure:
```json
{
   "buyPoint":{
      "price":40.62,
      "date":"2023-10-26T00:00:00Z"
   },
   "sellPoint":{
      "price":47.75,
      "date":"2023-11-03T00:00:00Z"
   }
}
```
* If server fails to process the query a client or server error message is returned:
```bash
HTTP/1.1 400 Bad Request
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Wed, 08 Nov 2023 09:00:17 GMT
Content-Length: 67

{"message":"begin period is after the end period: bad request"}  
```

# Start the service locally
`go run . -server.port=<server_port_for_http> -db.user=root -db.pass=<pass> -db.port=<db_port>`

### Usage
```
  -server.port int
        port to listen for incoming http requests (default 8080)
  -db.user string
        username to access the local mysql instance (default "root")
  -db.pass string
        password to access the local mysql instance (default "")
  -db.port int
        port of the local mysql instance (default 3306)
```

# Setup database
The repo comes with hardcoded predefined dump `data/dump.sql` if you want to use it for test purposes please follow the steps:
1. Run the following command to initialize MySQL docker container - provide password and a local port to run the instance

   ```docker run --name stock-quote-db -e MYSQL_ROOT_PASSWORD=<password> -p <port>:3306 -d mysql:latest```

2. Create the stockquotedb database

   `mysql --host 127.0.0.1 --port <port> -p<password> -e "CREATE DATABASE stockquotedb"`

3. Import the local dump into the database

   `docker exec -i stock-quote-db sh -c 'exec mysql -uroot -P<PORT> -p<PASS> stockquotedb' < data/dump.sql`


# Run tests
```go test ./...```
