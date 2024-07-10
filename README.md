# blockchange-sentinel

## Usage

Before runnning make sure you entered getblock.io's API key into `.env` as `API_KEY`, see .env.example.

NOTE: maximum RPS is 60 for free version

1. `make run-cli` runs the program, after showing the result the program exists.

2. `make run-api` runs the server, you can make HTTP GET request to `localhost:8080/block/most_changed` and it replies with a result in a json format. 

    By default `8080` port is used, but you can add `.env` file with `HTTP_PORT` field.