## Apollo Federation Demo

This repository is a demo of using Apollo Federation to build a single schema on top of multiple services. The microservices are located under the [`./services`](./services/) folder. The gateway that composes the overall schema is in the [`gateway.js`](./gateway.js) file.

### Installation

To run this demo locally, pull down the repository then run the following commands:

```sh
npm install
```
This will install all the dependencies for the gateway (and the optional node js versions of the federated services).

### Running Locally ( Gateway + Go)

First, run Dgraph:

```sh
go run main.go --start-dgraph
```

Then, initialize the Dgraph schema (see next section for what this is doing):

```sh
go run main.go --init-dgraph
```

Then, run all the Go services and the Node.JS Apollo Gateway:
```sh
go run main.go
```
| Port                          | Service   |
| ----------------------------- | --------- |
| [4000](http://localhost:4000) | Gateway   |
| [4001](http://localhost:4001) | Accounts  |
| [4002](http://localhost:4002) | Reviews   |
| [4003](http://localhost:4003) | Products  |
| [4004](http://localhost:4004) | Inventory |

To stop Dgraph, run

```sh
go run main.go --stop-dgraph
```

### What is entailed in initializing DGraph for the inventory service?
The main.go runner does this for you, but you can also manually:
```
curl http://localhost:8080/admin/schema --upload-file ./services/inventory/schema.graphql
```

Then, insert the data for the inventory service:

```sh
curl --request POST \
  --url http://localhost:8080/graphql \
  --header 'Content-Type: application/json' \
  --data '{"query":"mutation { addProduct(upsert: true, input: [{upc: \"1\", inStock: true}, {upc: \"2\", inStock: false}, {upc: \"3\", inStock: true}]) { product { upc inStock } }}"}'
```

### Running Locally (Node JS)

```sh
npm run start-services
```

This command will run all the Node.JS microservices at once. They can be found at http://localhost:4001, http://localhost:4002, http://localhost:4003, and http://localhost:4004.

In another terminal window, run the gateway by running this command:

```sh
npm run start-gateway
```

This will start up the gateway and serve it at http://localhost:4000

### What is this?

This demo showcases four partial schemas running as federated microservices. Each of these schemas can be accessed on their own and form a partial shape of an overall schema. The gateway fetches the service capabilities from the running services to create an overall composed schema which can be queried. 

To see the query plan when running queries against the gateway, click on the `Query Plan` tab in the bottom right hand corner of [GraphQL Playground](http://localhost:4000)

To learn more about Apollo Federation, check out the [docs](https://www.apollographql.com/docs/apollo-server/federation/introduction)
