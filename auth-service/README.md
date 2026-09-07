# Authentication Service
Authentication microservice, the goal of this project is to continue learning and exploring Golang, it will be updated as new version of Golang or the tools using here like "Gin Web Framework" are released.

It is a microservice for user administration and authenication, using JWT tokens and encrypting the information within it, different programming approaches are used, that is why there may be different implementations performing the same functions.

Requirements
------------
* Golang      1.25.6+
* Gin         1.12.0+
* Gorm        1.31
* Golang-jwt  5.3.1

Project Structure
----------
```text
.
├── code            # shared base code (services, models, db, etc)
│
├── gin-framework   # REST API built using Gin and Gorm
│
├── native          # REST API using Go Standard Libraries
│
├── go.work
├── go.work.sum
└── README.md
```

How to use it
-------------

Running from the terminal located in the root folder "auth-service".

* Download the packages
```bash
go mod tidy
```
* Run one of the projects

```bash
# native project
go run ./native/main.go

# gin project
go run ./gin-framework/main.go
```

 