# Recipe
Make sure you're not using localhost

# Start CockroachDB Cluster (insecure mode)
First Node
```
cockroach start --insecure --listen-addr=localhost
```
Second Node
```
cockroach start --insecure --store=node2 --listen-addr=localhost:26258 --http-addr=localhost:8081 --join=localhost:26257
```
First Node
```
cockroach start --insecure --store=node3 --listen-addr=localhost:26259 --http-addr=localhost:8082 --join=localhost:26257
```
# Create Database
Open a command line connection to the database through the port of any of the nodes in this case it's through the first node 
( port 26257)
```
cockroach sql --insecure --host=localhost:26257
```
and now the database
```
CREATE DATABASE recipes;
```
# rest-api
localhost:8083
in main.go change the variables if needed
getDB ( user, port, database )
```
go build
./rest-api
```

# router-app
```
npm run dev
```
