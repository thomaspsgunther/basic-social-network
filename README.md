# Basic social network project

This is a relatively small but robust project, aimed at creating a basic instagram-like social network. It uses PostgreSQL for database, Go for the backend and React Native for the frontend.

## Backend

For working with the backend, you'll need [PostgreSQL](https://www.postgresql.org) and [Go](https://go.dev).

Everything in /backend/internal/database/postgres/migrations/000001_init.sql needs to be executed manually with PostgreSQL at first, to create an initial "migrations" table, used to keep a record of migrations, and multiple functions that will be used for other tables.

Remember to create a .env file based on .env.example and adjust it for your own environment.

To run the backend, first execute the command **go mod tidy** to make sure you have the dependencies of the project installed and ready to go, then just execute the command **go run server.go**, or, more likely, **go run ./cmd/app/server.go**, it depends on what directory you're running it from.
