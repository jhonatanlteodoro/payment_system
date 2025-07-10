# Deliveries for first implementation

- -> Pgsql Setup with docker
- -> Table/s Definition
- -> Endpoint to make payment functional [accept a payment from an account and update account balance]
- -> Minimal tools for tests






# Api Contract

```bash

POST ->
Description: Allow user to make a payment from account A to accont B
route: <host>/api/<version>/payment
body:
{
    "from_account": "string-id", # uuid-7
    "to_account": "string-id", # uuid-7
    "amount": "99.01" # decimal string
    "payment_description": "this is a payment desc." 
}

Responses:

# When all good
Status code: 200
Payload: empty

# When origin account does not exist
Status code: 400
payload: {"message": "invalid origin account"}

# When target account does not exist
Status code: 400
payload: {"message": "invalid target account"}

# When balance is insuficient
Status code: 400
payload: {"message": "need to work more :) cheers!"}

# When service fails for any non expected reason
Status code: 500
payload: {"message": "the coffee machine died for some reason, try again in the future"}

```

# Server Initial Tech Specs [This most likelly to be changed in future versions]

- Lang: Go 1.24
- DB: Pgsql
- DB lib: Pgx
- Web Framework: None