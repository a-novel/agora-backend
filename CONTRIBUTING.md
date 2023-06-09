# Installation

Download and install:
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go](https://go.dev/dl/)
- [Mockery](https://github.com/vektra/mockery)

`mockery --name=Repository --with-expecter --case underscore --inpackage`

You can use the below command to install (required) psql toolset on your mac:

```shell
brew install postgresql
```

## Clone the project

```shell
# Clone and move inside repository.
git clone git@github.com:a-novel/agora-backend.git
cd ./agora-backend

# Setup local environment (you can rerun this step anytime if needed).
make setup
```

Ask your admin for the Sendgrid API key, and add the following line at the end of your `.envrc` file:

```shell
export SENDGRID_API_KEY="api secret key"
```

> The `.envrc` file is reset each time you run `make setup`. Normally, the script is able to automatically
> copy the API key once set, but you might want to look at `.envrc.old` file in case you need to restore
> anything.

## Try the application

```shell
make run

# go run ./cmd/main/main.go
# [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
#  - using env:   export GIN_MODE=release
#  - using code:  gin.SetMode(gin.ReleaseMode)
# 
# [GIN-debug] GET /ping --> github.com/a-novel/agora-backend/api/base.PingPong (5 handlers)
# ...
```
