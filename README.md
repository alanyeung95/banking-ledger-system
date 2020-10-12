# banking-ledger-system
This is a demo project/assessment for Crypto.com Ops Team Back End Engineering Coding Challenge

# Prerequisite

## Using docker
1. Docker (of course)

## Using local environment
1. Golang version 1.13
2. MongoDB running on `localhost:27017` (sorry you have to install it manually at this stage)

# Usage
## Run the API server 

### Using docker
`docker-compose up`

### Using local environment

export environment variable:

```
export MONGODB_ADDRESSES=localhost:27017
export MONGODB_DATABASE=banking
```

Then build the binary file and execute it
```
go build
./banking-ledger-system
```

the service will be running on `localhost:3000`

## Run test cases

### Using docker
`docker-compose -f docker-compose.yaml -f docker-compose.test.yaml up --abort-on-container-exit`

### Using local environment
First plz start the server by following the above instruction

```
export API_TEST_DOMAIN=http://127.0.0.1:3000
```

Then run the test command
```
go test -v
```

# Assumption

## User group
1. In this system, operation team user is a kind of admin user and have more permission such as `fixing transaction`. 
2. Normally, admin user should be created by superuser. In this project for simplicity, we have another API to create such kind of user.
3. `Fix a withdrawal or deposit transaction` is one of the transaction that `Operation team` process, and it will show on the transaction history.

## Others
1. let the account id format be uuid, although in reality the format will be in different pattern
2. Assume there is not authentication service on this API service

# Priority list
## Must do
- [x] Create a new bank account
- [x] Make a withdraw
- [x] Make a deposit
- [x] Make a transfer between two accounts
- [x] Fix a withdrawal or deposit transaction (here I assume `fix` mean `undo`)
- [x] View current balance for Customer
- [x] View transaction history for Customer
- [x] View transaction history for Operation Team

- [x] swagger doc
- [x] test cases

## Nice to have
- [x] dockerize API service
- [x] dockerize mongo service
- [ ] better error encoding format
- [ ] handle atomic operation
