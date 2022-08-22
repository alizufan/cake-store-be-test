# 🍰 Cake Store API

## ⚙️ Specifications

Written in Go version : 1.19
## 📚 Repo Structure
```
├── handler
├── libs
│   ├── logger
│   └── util
├── logs
├── migrations
├── repository
├── schema
├── server
│   └── middleware
└── service
```

- `handler` contains go package layer to handle requests from http (request layer)
- `libs` contains shared code that can be used on each packages
- `logs` contains logging file
- `migrations` contains migrations file
- `repository` contains go package layer to serve a requests from service (source data layer)
- `schema` contains shared code that can be used on other packages in context entity structure
- `server` contains a go http server and middleware
- `service` contains go package layer to serve a requests from handler (business logic layer)

## 🔧 Running Locally
To run this project you need some preparation :
- `create database 'cake-store'` 
- `installing migrator tools` download from [golang migrate](https://github.com/golang-migrate/migrate/releases) in release page
- `migrate -path ./migrations -database "mysql://root:secret@tcp(localhost:3306)/cake-store" -verbose up` run this command to up a migration (you can look from Makefile)
- `go mod tidy` installing a module
- `go run .` run it

To test api you can use a OpenApi extension from vs code. [Open API on VSCode](https://marketplace.visualstudio.com/items?itemName=42Crunch.vscode-openapi)

If you want easily to run this project, use with docker compose and run a migration (install migrator tools first) :
```
docker compose up -d
migrate -path ./migrations -database "mysql://root:secret@tcp(localhost:3306)/cake-store" -verbose up
```
If you wanna change a environment you can change in docker-compose.yml.

Then to run the project locally, the default port is 3000 you can change a port hard code :

```
go mod tidy
go run .
```
## 📰 Info
This project using a distroless for image, you can freely switch between a production and development :

Development : `gcr.io/distroless/static-debian11:debug` use this image you can access a shell iinteractive 
Production  : `gcr.io/distroless/static-debian11:latest` use this image you only can access logs from docker

This project has `tracker_id` to provide developer finding a error in log file easily, to look a logs file 
you can looking in docker volume and inspect it other information in docker compose file, 
or if you run manually with `go run .` you can open `./logs/logging.log` and open with text editor.

## 🔧 Deploying
1. Install docker
2. Run docker compose 
```
docker compose up -d
```

## 📦 Go Library

Using [Go Chi](https://github.com/go-chi/chi) as router for building HTTP services, looking a [Docs](https://github.com/go-chi/chi).


## 📰 Go Article

[Download Golang Binnary](https://go.dev/dl/)

[How to install Go in PC / Laptop / Server](https://go.dev/doc/install)

## 📚 Go Book

[Go Tutorial - Bahasa](https://dasarpemrogramangolang.novalagung.com/)

## 💡 Go Command

[CMD List Golang](https://go.dev/cmd/go/)

## 🧷 Recommended IDE

[Visual Studio Code](https://code.visualstudio.com/)

## 🔧 Recommended Extension Visual Studio Code
[OpenAPI on VSCode](https://marketplace.visualstudio.com/items?itemName=42Crunch.vscode-openapi)

[GO Extension on VSCode](https://marketplace.visualstudio.com/items?itemName=golang.go)
