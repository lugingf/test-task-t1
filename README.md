# talon-backend-assignment

This is a simplified backend that listens to POST request. Once it receives such a request it will push it to a PostgreSQL database.


## Instructions
Clone the repository to your local environement.

To run the program you can choose between two methods:

### Containerized

#### Prerequisites
 - [Docker](https://docs.docker.com/compose/install/) (including `docker-compose`)

#### Run
Simply type `docker-compose up` in the root folder.

Once the console logs the line:
```
api_1  | {"level":"info","time":"2020-10-10T08:16:51.431Z","msg":"Listening...","address":":8080"}
``` 
The API is available and ready to accept requests on port 8080.

### Locally

#### Prerequisites
 - [Go](https://golang.org/doc/install)

To run the go program in your local environment you need to have a database running and ready to accept connections.

After the PostgreSQL (11) database is running, configure the following environment variables to hold your database credentials:

```
"DB_USER" - the username.
"DB_PASSWORD" - the password.
"DB_NAME" - the database name.
"DB_HOST" - the database hostname.
"DB_PORT" - the port to connect to.
"DB_SSL" - SSL mode (example: require / disable).
```

You can achieve the above using the provided docker-compose file in the repository by simply run the command `docker-compose up -d db` in the root folder. Shortly there will be a running Postgresql service with the credentials provided in the [`docker-compose.yml` file](docker-compose.yml#L11-L15).

#### Run

Run the application using `go run main.go`.

### Connecting to the dockerized Postgresql service
You can connect to and inspect the database and its contents using a cli tool such as [pgcli](https://www.pgcli.com/) or a gui tool such as [TablePlus](https://tableplus.com/).

i.e. using `pgcli`:
```shell
pgcli postgres://talon:talon.one.8080@localhost/talon
```
