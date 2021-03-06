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
go run main.go --use-dgraph
```
| Port                          | Service   |
| ----------------------------- | --------- |
| [4000](http://localhost:4000) | Gateway   |
| [4001](http://localhost:4001) | Accounts  |
| [4002](http://localhost:4002) | Reviews   |
| [4003](http://localhost:4003) | Products  |
| [8080](http://localhost:8080) | Inventory |
| [8080](http://localhost:8080) | DGraph    |

To stop Dgraph, run:

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

### Federated Combined Schema and Queries
If you open a browser to http://localhost:4000 you should be able to run two different top level queries.

Run any of these example queries as [http://localhost:4000/playground](http://localhost:4000/playground)

```
query MyReviews{
  me {
    username
    reviews {
      body
      product {
        name
        upc
      }
    }
  }
}
```

```
query TopProducts{
  topProducts{
    name
    reviews{
      author{
        name
      }
    }
  }
}
```
The combined federated schema is this:

```graphql
type Product {
upc: String!
name: String
price: Int
weight: Int
reviews: [Review]
inStock: Boolean
shippingEstimate: Int
}

type Query {
me: User
topProducts(first: Int = 5): [Product]
}

type Review {
id: ID!
body: String
author: User
product: Product
}

type User {
id: ID!
name: String
username: String
reviews: [Review]
}
```
### 

### A word about [DGraph](https://dgraph.io/docs/get-started/)

Apollo Federation is now merged and available in the master branch of Dgraph. This is available via Docker Hub using the `dgraph/dgraph:master` Docker image. If you want to use a stable image tag (the master image always updates to the latest master), you can use `dgraph/dgraph:3642fed5`.

You can read more about how to use Apollo Federation in the docs currently in the description of [PR #7275](https://github.com/dgraph-io/dgraph/pull/7275). This adds support for the @key, @extends, and @external directives.

This is slated to be released in the official Dgraph v21.03 version in March. Please do let us know if something doesn't work for you or if there's anything we can improve.

PR extends support for the `Apollo Federation`.

## Support for Apollo federation directives
DGraph current implementation allows support for 3 directives, namely @key, @extends and @external.
[Work is in progress](https://github.com/dgraph-io/dgraph/pull/7503) to support @provides and @requires directives.

### @key directive.
This directive is used on any type and it takes one field argument inside it which is called @key field. There are some limitations on how to use @key directives.

* User can define @key directive only once for a type, Support for multiple key types is not provided yet.
* Since the @key field act as a foreign key to resolve entities from the service where it is extended, the field provided as an argument inside @key directive should be of `ID` type or having `@id` directive on it. For example:-

```
type User @key(fields: "id") {
   id: ID!
  name: String
}
```

### @extends directive.
@extends directive is provided to give support for extended definitions. Suppose the above defined `User` type is defined in some service. Users can extend it to our GraphQL service by using this keyword.

```
type User @key(fields: "id") @extends{
   id: ID! @external
  products: [Product]
}
```

### @external directive.
@external directive means that the given field is not stored on this service. It is stored in some other service. This keyword can only be used on extended type definitions. Like it is used above on the `id` field.
