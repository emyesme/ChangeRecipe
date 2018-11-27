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

# rest-api
localhost:8083
... steps to get the info of the first node answer
```
go build
./rest-api
```

# router-app
```
npm run dev
```
