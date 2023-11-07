# Setup database
1. Run the following command to initialize MySQL docker container - provide password and a local port to run the instance
`docker run --name stock-quote-db -e MYSQL_ROOT_PASSWORD=<password> -p <port>:3306 -d mysql:latest`

2. Create the stockquotedb database
`mysql --host 127.0.0.1 --port <port> -p<password> -e "CREATE DATABASE stockquotedb"`
3. Import the local dump into the database
`docker exec -i stock-quote-db sh -c 'exec mysql -uroot -P<PORT> -p<PASS> stockquotedb' < data/dump.sql`

# Start the local server
`go run . -server.port=<server_port_for_http> -db.user=root -db.pass=<pass> -db.port=<db_port>`


# Run tests
```go test ./...```
