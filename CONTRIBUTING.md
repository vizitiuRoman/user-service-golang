# Participation in development

## Beginning of work

* Cloning a repository
    `` bash
    git clone https://github.com/vizitiuRoman/user-service-golang.git
    cd user-service-golang
    ``,

* Install dependencies
    `` bash
    go get. / ...
    ``,

* Linting project
    `` bash
    make lint
    ``,

* Unit tests
    `` bash
    make test
    ``,

* Running local version in hot reload mode
    `` bash
    cd docker && docker-compose up
    nodemon --exec go run cmd / main.go --signal SIGTERM
    ``,

## Code-style

The code must be formatted with the `gofmt` utility. All linter comments (except for obviously erroneous)
must be fixed before pushing to the server (or better before each commit).

## IDE
The preferred development environment is [GoLand] or [IntelliJ IDEA] from Jetbrains.

The project repository includes GoLand / IDEA configuration files, with pre-configured code auto-formatting utility
gofmt and running the linter in fast mode after every file save and connecting to the database from
local docker environment.
 
Also added a number of configurations for launch:
* Service in docker environment
* Linter
* Unit tests
* Full testing (unit + integration)

### Branch naming

Branch names consist of a prefix, issue number, if any, and a short description of the issue in English (2-4 words).
The prefix is ​​the issue type (bug, feature), followed by a trailing slash `/`.

Examples:
    
    hotfix/auth-fix-response-code
    feature/auth-new-api-method