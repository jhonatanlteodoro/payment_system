# Deliveries for first implementation

- -> Pgsql Setup with docker
- -> Table/s Definition
- -> Endpoint to make payment functional [accept a payment from an account and update account balance]



### Api Contract

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

```

### Server Initial Tech Specs [This will most likely be changed in future versions]

- Lang: Go 1.24
- DB: Pgsql
- DB lib: Pgx
- Web Framework: None - Using Net std lib



# !After implementation Concerns

If you have read the code before read this document I'm sure you are think this code is pure trash, and you are not wrong... But you will understand why it is like this.

Before continue lets anwser a few Q&A in advance:
- Q: Why write a trash code?
 - - A: The idea is to visualize the payment flow and for that we will dont need good code, only see steps of execution, and since the future versions will be improving step by step, I don't care about write good code here.

- Q: Why pgx and why I wrote the queries like that?
 - - A: I did not play with PGX in the past so I take the opportunity to play with it now

- Q: Why I did not use transaction to manage possible errors and be able to do a rollback?
- - A: I don't care about this at this moment :)

Anyway, Lets continue to what matters...
Even if you Have designed and implemented the best code possible using this kind of flow for you system you will fail.
Why? This system worth run just for a POC, if you take a closer look, we have a simple flow that 

*Receive Payment api request -> Check balances -> Create a Transaction -> Update Balances*

You could write the best code possible to archive this but when the real life production knock at your door you will fail
Cases that this kind of design will fail:
- **Service Scale**: Since this is a synchronous not locking flow, you will be forced to accept at most one request by account at time to make this transaction,
  otherwise the following will happen -> multiple transaction at the same time will not guarantee consistency, imagine a balance of 100.00 then 2 transactions to be processed of 50.00 and another of 100.00
  both can enter the flow at the same time but one of them will fail because the amount available for that account will not be enough [off course if your table has some constraints to not allow negative values]...
  
  -- So, forced to be one machine and pray for you users to not try to make payments at the same time, because one machine will handle I would guess at most 10-15Req/s, which is already a good number but this amount of requests, 
  can easily force you to scale your database, which also has a limit, and of course with this kind of implementation scale means concurrent access which means the *deadlock* is waiting to give you a great hug :)
  

Well, this is the high level problems you could face trying to use this kind of design for payment, lets continue working on it to hit somewhere better.

***Before continue, try to image this issues in an environment where 100'sReq/s are the normal load, and you have to be prepared for a BlackFriday with has potential for a few 1000`s Req/s***