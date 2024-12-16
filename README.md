# Best Asset Tracker (BAM) Backend
- start by running 
    `go run ./cmd`

# about inventory tracking
- We create an item, everything about it is known except the price and age. Then we create "instances" of that item and those each has it's own price and has it's purchase date. So there're two IDs to keep track of. 1. The general item, 2. The item instance. So for example we create an Item type called hat and we give it unique id, images; weight; description; quantity, and a list of related item types. Then if we bought 10 hats and debit them to invenotry then we create and id for each of the 10 hats and each of the 10 hats can have its own price, it's own warehouse location. The bigger item type can be called item type, and the indivisual items can be called item units. Each item unit has its own count (if you bought 10 for the same price at the same date for example) and the item type will sum the quantity of all the units and report it.
- When you record journal entries you can specifiy which inventory item units are involved, by id.
  

  # TODO before deployment
  - properly administer secrets to app
  - enable SSL connection to db


check List:
1. create UI to create accounts
2. create UI to do a transaction (journal entry)
3. create UI that allows account view with total final balance
4. create UI that allows invenoty account creation
5. create UI that allows item type and item unit creation
6. create UI that allows inventory transfer
7. create UI that shows inventory quantities for each inventory account
8. do a transaction that involves inventory transfer


# Development notes
- handling the error should be done at the site the error was created, in other words calling cmn.HandleError is only done at the location where the error was created. We don't pass the error around, unless we want to use it for control flow. For example if there's an error in sql insertion, we call handleError on the sql function itself, but then pass the error to the business-component (example accounting component) where it will decide what error code to send to the client. But the actual logging and what not was done at the sql function. API -> Component -> DataAcessLayer.
- The buseinss component has all the logic, and it decides what is valid or not. The API layer only handles the request validation and communication.
- The common package is common to all other packages, all of them may import it.
- Package and component are used exchangably in documentation.
- the accounting pacakge needs data access and data access needs accouting. Since cyclic imports aren't allowed the models and common packages will hold on to the 

## Handling of money
- There should be no rounding, until the last step
- We use bigInt and store numbers as mils (like cents but with 1000 instead of 100) for example if 1.25$ is 125 in cents, it will be 1250 in mils. We don't store fractions or decimals, we store money as integer mils. 1.25$ will be stored as 1250 in the database. the UI can change the representation however they like, but from the db standpoint it will always be in mils. Meaning that the biggist amount we can store is 9223372036854776. We can't store a bigger number. If user tried to store a bigger number we issue an error. 

## handling of primary keys
- integer values can handle a max number of about 2 billion (2147483647), if you try to store a bigger number like 2147483648 you will get "integer out of range error". BIGINT can store at the biggist 9223372036854775807, if you store a bigger number you will get the same error. This is a problem for when your primary key is an integer, but you don't have to worry about it if you don't expect that many rows. It will be a problem in the future however. When you hit that range you can alter the table to accept BIGINT instead of INT for the primary key


## DB migration
we use this [tool for migration](https://github.com/golang-migrate/migrate/tree/master) version 4.16.2
 example of creating a migration
 `migrate create -ext sql -dir db/migration -seq migration_name` where migration_name is the migration name, it can be anything you want. This will create two sql files up and down.

 ### Editing the schema.
 - Everytime you edit something in the database sql files make sure to kill the db_tools job then call the setup script again.
   - you kill the job by running the command `jobs` to see the job id then run `kill %job_id` you need the `%` when passing a job id.
   -  This is because the `db_tools` will not have the new schema. It loads the schema in memory once everytime you call the script. If you make changes to an .sql file and migrate up, you will not get your changes on the database after the migration until kill the old instance and you run the script.

 ## TODO
 - Create indexes. by journal id for the transacitons table and by group id for the accounts table
   - explore adding indexes to the junction table


## Setup development environment
1. install podman and go version 1.21.
2. run `podman machine start` then pull the database image using `podman pull docker.io/postgres:15.4-bullseye`
3. run `source scripts/development_environment/dev_env_init.sh`
4. this should start the container, you can use `podman ps` to see it. If it's not there (that happens sometimes) just use `podman start pgsql` to start it.
5. You don't need to do this, but for any reason you need to access the container's database you can use `podman exec -it  pgsql bash`  to access the container, then from the container bash you can run `psql -h localhost -U $BACKEND_USER --dbname $BACKEND_DB`  or you can use the `$BACKEND_TESTING_DB` to access the testing database, if you need to access as super user do `psql -h localhost -U postgres`.
6. You can start the backend with the command `go run cmd/main.go`

## Running tests
- you can use the default testing go functinoality to run unit test
- for the tests that validates the database state you must run the app in test mode
- to run all tests you must first run the startup script. If you don't then the tests will fail.
- to run all db_state_tests use `go run ./integration_tests/database_state/cmd/main.go` the tests run one after the other since they all alter the db state.

### Integration test coverage
<!-- // integration test -->
 go build  -cover  -o out /integration_tests/database_state/cmd/main.go
 GOCOVERDIR=$GOCOVERDIR ./out 
<!-- // if you ran more the one time you will have duplicates, you can merged them like this  -->
Mkdir merged
go tool covdata merge -i=./integration_cov -o=merged
<!-- // to show your results. -->
 go tool covdata textfmt -i=integration_cov -o=cov.txt
go tool cover -func=cov.txt
go tool cover -html=cov.txt


## CI/CD
- the main CI/CD is done useing gitlab. The docker images are used for testing on a clean container without gitlab or any provider.
- There are two testing docker images. 
  - 1. `compiled_tests_env.Dockerfile`. Which requires that the program is already compiled and put in the correct directory. You can do that by running the `testing_environment/compile.sh` script before building the image
  - 2. `end_to_end_tests.Dockerfile` like the first image, but compiles the code on its own. This image is suitable for CI/CD as the compilation is done in the image itself. It uses a building container to build the go program, and then another image to run the tests.
- to compile testing use `go test -c -o ./build/test ./cmd/... -v -failfast` this will put tests in /build/tests file. the failfast will make our tests stop as soon as one test fails. and the -v will make it verbose so you know exactly which one.
- The developer should run tests locally before pushing code. the ci/cd pipeline is not for error reporting. You can run tests locally using `go test ./cmd/... -v -failfast`

### GitLab
- gitlab is used for CI/CD and the `.gitlab-ci.yml` file defines the ci jobs. 
- to test locally what happens with a particular job you can build the gitlab image with `podman build -f infra/test/gitlab_test.Dockerfile . -t my_gitlab`
- And then test the specific job with `podman run --rm -t -i my_gitlab exec shell job_build`



## Features to add:
- pagination on GET requests
- Sales component (allows you to sell items, will impact both finacial journal and inventory journal)
- Inventory forcast
- Placing orders and generating invoice (orders are either paid now or paid later)
- parent item to track different variants of the same item. this is for when we have a type of tshirt and a couple of tshirt variants. For example the same type of tshirt but 4 different sizes. Then we would have one parent item and 4 items. These 4 tiems would inherit eveything from this parent item, but then add their own "variance" which in this case is size, or maybe color. We track the item not the parent item. As we might have different numbers of each variant. This feature is todo.
- Ability to accept payments, maybe stripe. clover, square, shopify ...etc
- Reading PDFs and invoices
- reading QIF and CSV as well as producing them
- user signup and login for multiple users