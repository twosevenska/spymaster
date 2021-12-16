# Spymaster

![spymaster.jpg](/spymaster.jpg)

This is a microservice that focus on managing users and notifying other services of any changes.

## How to use

### Requirements

* MongoDB 4.4 (docker-compose provided)
* GO 1.17 (when building from source)
* Make (if using the Makefile as the entrypoint, recommended)

### Starting a MongoDB container with Fixtures

```shell
make docker-up
```

The commmand above will start a MongoDB docker container in the background with two DB's in it:

* spymaster: Loaded with a single user entry
* test_spymaster: Empty for the sake of running our tests

To find more about what to expect, check `docker/mongo/fixtures/spymaster.js`

!!NOTE!!

The DB provided was never intended for production, therefore its configuration is both naive and lazy.

### Example with binaries

```shell
make bin
./build//spymaster_darwin_amd64
```

### Example from source

```shell
make run
```

### Running the tests

#### Go test

```shell
make go-test
```

#### GoConvey

[GoConvey](https://github.com/smartystreets/goconvey) will not only run the tests, but also open a web interface with notification which can become quite handy when developing.

```shell
make go-convey
```

## How to use (as a user)

In order to operate with the service we currently only support REST calls. For convenience here's a quick reference of our endpoints:

```go
r.GET("/ping", controllers.Ping)
r.GET("/health", controllers.Health)

api := r.Group("/", controllers.Pagination())
{
    api.GET("/users", controllers.ListUsers)
    api.POST("/users", controllers.CreateUser)
    api.PATCH("/users", controllers.UpdateUser)
    api.DELETE("/users", controllers.DeleteUser)
}
```

I would recommend checking the `types/types.go` file in order to easily understand the payload.

### Patch httpie example

`http PATCH '0.0.0.0:7000/users?id=61ba6382df4bec585cf60e60' first_name=omg`

## Development notes

### Lifetime of a request

API (server) -> Controllers function for validation and initial preparation of data for consumption -> spymaster for any processing needed -> Mongo (DB) -> spymaster for further processing -> Controllers to prepare and return api response

### Configuring Mongo

As seen in `cmd/spymaster/main.go` we are loading our env vars with the `SPYMASTER_` prefix. However for the Mongo instance bundled the defaults work as expected.

### Major TODOS

* Add a producer which would allow us to send notifications to other services (use Kafka, RabbitMQ, maybe Ably?)
* Async capability
  * This would rely on changing our flow a little bit
  * Relies on having a producer/dispatcher
  * Create a consumer/listener
  * Change our controllers to not call our spymaster entries directly, instead send a dispatch
  * One of our instances would have their listener consume the message and process it
  * Would need to keep track of request status

### Known issues

* Lack treating uppercase/lowercase
* Lack field validation
* Lack Input trimming
* Coverage report not working
