# Bima Framework

This repository contain only the framework, you can refer to [skeleton](https://github.com/KejawenLab/skeleton) for implementation

Bima Framework utilize [gRPC Gateway](https://grpc-ecosystem.github.io/grpc-gateway) to make REST easier and added some features for productivity. Bima Framework designed to running behind api gateway or proxy

## Features

- [gRPC Gateway](https://grpc-ecosystem.github.io/grpc-gateway)

- CRUD Generator

- Authentication and Authorization (using Middleware)

- Support RDBMS (using [Gorm](https://gorm.io)) and MongoDB (using [mgm](https://github.com/Kamva/mgm))

- Soft Deletable (only for RDBMS)

- Support [Elasticsearch](https://github.com/olivere/elastic) and [AMQP (Queue)](https://github.com/ThreeDotsLabs/watermill) Out of The Box

- [Event Dispatcher](https://en.wikipedia.org/wiki/Observer_pattern)

- Auto documentation (Swagger)

- [Dependency Injection](https://github.com/sarulabs/dingo)

- Two Level Middleware (Global or Per Route)

- Better Log Management

- Support Custom Route

- Support [HTTP Compression](https://github.com/CAFxX/httpcompression) Out of The Box

- Health Check

- Easy to Custom

## Testing

`go test -coverprofile /tmp/coverage ./... -v`

## Todo

- Testing testing testing

- Documentation
