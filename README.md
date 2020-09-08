Snippetbox
----------

Golang web application from book "Let’s Go", author Alex Edwards. Original code was modified with third party libraries.

List of used libraries:

* [ozzo-validataion](github.com/go-ozzo/ozzo-validation) for validation
* [gorilla-csrf](github.com/gorilla/csrf) for validation csrf tokens
* [gorilla-mux](github.com/gorilla/mux) for routing
* [gorilla-session](github.com/gorilla/sessions) for session
* [godotenv](github.com/joho/godotenv ) for reading env files
* [logrus](github.com/sirupsen/logrus) for logging


For migrations use tool [migrate](https://github.com/golang-migrate/migrate)

For running
-----------

1. Run mysql docker image for tests (`docker-compose -f docker/docker-compose-prod.yml up`)
2. If it first run, wait for loading mysql and run migrate tool (`migrate -path migrations/ -database "mysql://user:pass@tcp(127.0.0.1:3307)/snippetbox" up`)
3. `make`
4. `./snippetbox`

See `conf.env` for environment variables or set it when start app. Example: `PORT=8082 ./snippetbox`


For testing
-------

1. Run mysql docker image for tests (`docker-compose -f docker/docker-compose-tests.yml up`)
2. If it first run, wait for loading mysql and run migrate tool (`migrate -path migrations/ -database "mysql://user:pass@tcp(127.0.0.1:3308)/snippetbox_test" up`)
3. `make test`

If you want other DSN string for mysql, change variable `dsnString` in `pkg/models/mysql/test_helper.go`
