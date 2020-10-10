# banking-ledger-system
This is a demo project/assessment for Crypto.com Ops Team Back End Engineering Coding Challenge

# Prerequisite
1. Golang version 1.13
2. MongoDB running on `localhost:27017` (sorry you have to install it manually at this stage)

# Usage
Run the API server 

`docker-compose up`

the service will be running on `localhost:3000`

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
- [ ] Fix a withdrawal or deposit transaction
- [ ] View current balance for Customer
- [x] View transaction history for Customer
- [x] View transaction history for Operation Team

- [ ] swagger doc
- [ ] test cases

## Nice to have
- [ ] dockerize API service
- [ ] dockerize mongo service
- [ ] better error encoding format
- [ ] handle atomic operation
