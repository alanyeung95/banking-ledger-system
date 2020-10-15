# Banking Ledger System
This is a demo project/assessment for Crypto.com Ops Team Back End Engineering Coding Challenge

Demo link: https://youtu.be/1J2_9MB0beg

-------------------------
- [Requirements](#requirements)
- [Setup / Running Instructions](#setup--running-instructions)
- [Project Assumptions](#project-assumptions)
- [Project Requirements](#project-requirements)
- [Development Notes](#development-notes)

-------------------------
## Requirements

### Using docker
1. Docker (of course)

### Using local environment
1. Golang version 1.13 or above
2. MongoDB running on `localhost:27017` (sorry you have to install it manually at this stage)

## Setup / Running Instructions
### Run the API server 

The API server will be running at `localhost:3000` and the mongo service will use port `27017`, so please reserve these two ports for this project

#### Using docker
```
docker network create network
docker-compose up
```

or just
```
make run
```

#### Using local environment

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

### Running test cases

#### Using docker
`docker-compose -f docker-compose.yaml -f docker-compose.test.yaml up --abort-on-container-exit`

or just
```
make test
```

#### Using local environment
First please start the server by following the above instruction

Then export this environment variable
```
export API_TEST_DOMAIN=http://127.0.0.1:3000
```

Finally run the test command
```
go test -v
```

-------------------------
## Project Assumptions

### User group
1. In this system, operation team user is a kind of admin user and have more permission such as `fixing transaction`. 
2. Normally, admin user should be created by superuser. In this project for simplicity, we have another API (`localhost:3000/accounts/create-admin`) to create such kind of user.
3. `Fix a withdrawal or deposit transaction` is one of the transaction that `Operation team` process, and it will show on the transaction history.

### Others
1. For simplicity, this project doesn't have login service. So instead of using `jwt` for authentication service, you just need to include the `account_id` on the request header.

-------------------------
## Project Requirements
### Basic Requirements
Please check the swagger.yaml for API description
- [x] Create a new bank account `POST /accounts/create` and `POST /accounts/create-admin`
- [x] Make a withdraw `POST /transactions`
- [x] Make a deposit `POST /transactions`
- [x] Make a transfer between two accounts `POST /transactions`
- [x] Fix a withdrawal or deposit transaction (here I assume `fix` mean `undo`) `POST /transactions/{id}/undo`
- [x] View current balance for Customer `GET /accounts/{id}`
- [x] View transaction history for Customer `GET /transactions`
- [x] View transaction history for Operation Team `GET /transactions`
- [x] Swagger doc
- [x] Test cases
  - [x] Happy path
  - [x] Negative test cases (see next next section)

### Nice to have
- [x] dockerize API service
- [x] dockerize mongo service
- [x] return error status code
- [ ] handle atomic operation
- [x] screenshoot or video recording of project demo
- [x] Makefile to simplify startup/test case command (docker only)

### Negative test cases
  - [x] Create account when account name or password is empty
  - [x] Withdraw when the account balance is 0
  
-------------------------
## Development Notes
Given the limited time, I chose CRUD over Event-Sourcing. But Event-Sourcing is a much better design option for this project.

Compare with CRUD design,  the single source of truth in Event-Sourcing will be the trasanction events. If we use CRUD, we need to maintain both account and transcation record and we need to make sure these two records are in sync.
