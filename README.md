# Basic social network project (rede social básica)

## English

This is a relatively small but robust project, aimed at creating a basic instagram-like social network. It uses PostgreSQL for the database, Go for the backend and React Native for the frontend.

### Backend

For working with the backend, you'll need [PostgreSQL](https://www.postgresql.org) and [Go](https://go.dev).

Everything in `backend/internal/database/postgres/migrations/000001_init.sql` needs to be executed manually with PostgreSQL at first, to create an initial "migrations" table, used to keep a record of migrations, and multiple functions that will be used for other tables.

Remember to create a `.env` file based on `.env.example` and adjust it for your own environment.

To run the backend, first execute the command `go mod tidy` to make sure you have the dependencies of the project installed and ready to go, then execute the command `go run ./cmd/y-net/main.go`.

Documentation is available through Swagger, go to `host:port/swagger/index.html` to access it.

### Frontend

For working with the frontend, you'll need [Node.js](https://nodejs.org).

To run the frontend, first execute the command `npm install` to install the project's dependencies. Once that's done, make sure you have an android emulator installed and running, then execute the command `npx expo run:android` to run the project on the currently running emulator.

## Português

Este é um projeto relativamente pequeno, mas robusto, voltado para a criação de uma rede social básica semelhante ao Instagram. Ele utiliza PostgreSQL para o banco de dados, Go para o backend e React Native para o frontend.

### Backend

Para trabalhar com o backend, você precisará do [PostgreSQL](https://www.postgresql.org) e do [Go](https://go.dev).

Tudo em `backend/internal/database/postgres/migrations/000001_init.sql` precisa ser executado manualmente com o PostgreSQL inicialmente, para criar uma tabela inicial de "migrations", usada para manter um registro das migrações, além de várias funções que serão usadas para outras tabelas.

Lembre-se de criar um arquivo `.env` com base no `.env.example` e ajustá-lo para o seu próprio ambiente.

Para executar o backend, primeiro execute o comando `go mod tidy` para garantir que você tenha as dependências do projeto instaladas e prontas para uso, em seguida, execute o comando `go run ./cmd/y-net/main.g`.

A documentação está disponível através do Swagger. Acesse `host:port/swagger/index.html` para visualizá-la.

### Frontend

Para trabalhar com o frontend, você precisará do [Node.js](https://nodejs.org).

Para executar o frontend, primeiro execute o comando `npm install` para instalar as dependências do projeto. Uma vez feito isso, certifique-se de que um emulador Android esteja instalado e em execução, então execute o comando `npx expo run:android` para rodar o projeto no emulador atualmente em execução.
