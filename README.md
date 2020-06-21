# Minimal Golang+Postgres Setup With Migrations, Unit Tests, and Docker
This code is intended to provide the absolute bare minimum you need to get
started with a database-backed Go project. With this code you can write and
run unit tests that hit a real Postgres database (see main_test.go), **with a development loop
that is less than 1s long.**

It also allows you to cleanly define migrations in a migrations/ folder
so that you can document your database schema changes over time.

To get started, you can simply add tests to the main_test.go file, modify the
schema by editing the files in the migrations/ folder, or take a look at the
docker-compose.yml to understand how the go binary talks to the database. You
can also check out the Dockerfile to see how the go_test service is set up (and
note how it's super fast because it doesn't copy any data into the image).
# Quick Start
```
$ docker-compose run --rm go_test
```
**Make sure your $GOPATH is set.** This will spin up a test postgres
service called "postgres_test" defined 
in docker-compose.yml and a "go_test" service, also defined there, which
will run all of your tests that require database access.

Before running
the tests, it reads environment variables from test.env, which contains
the database config parameters.

It then connects to the database and
applies the migrations specified in the migrations/ folder (it uses
golang-migrate to apply the migrations with an "up" suffix from that
folder).

After the test is complete, it undoes the migrations by applying
the "down"-suffixed files in the migrations/ folder.

```
$ docker-compose down
```
This step is optional since the --rm parameter above will automatically
remove the go_test service once it finishes running the tests.

However, there
is a bug with --rm, which
causes it to leave dependent services running even after the main service
is complete. This means that your "postgres_test" service will stay up
until you run this command. See
[this issue for details on that](https://github.com/docker/compose/issues/2791).
Were it not for this issue, you wouldn't need to run this command.
