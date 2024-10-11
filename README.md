# Basic social network project

This is a relatively small but robust project, aimed at creating a basic instagram-like social network. It uses PostgreSQL for database, Go for the backend and React Native for the frontend.

## Backend

For working with the backend, you'll need [PostgreSQL](https://www.postgresql.org) and [Go](https://go.dev).

Everything in /backend/internal/database/postgres/migrations/000001_init.sql needs to be executed manually with PostgreSQL at first, to create an initial "migrations" table, used to keep a record of migrations, and multiple functions that will be used for other tables.

Remember to create a .env file based on .env.example and adjust it for your own environment.

To run the backend, first execute the command **go mod tidy** to make sure you have the dependencies of the project installed and ready to go, then execute the command **go run ./cmd/y-net/main.go**.

Documentation is available through Swagger, go to **host:port/swagger/index.html** to access it.

## Frontend

For working with the frontend, you'll need [Node.js](https://nodejs.org).

To run the frontend, first execute the command **npm install** to install the project's dependencies. Once that's done, make sure you have an android emulator installed and running, then execute the command **npx expo run:android** to run the project on the currently running emulator.
