# Deliveries for second implementation

The deal here is to make it more reliable and also starting giving this code a more non trash look.
Since Im planning to make at least more two versions of it so, this delivery I'll try design a solution that will scale a bit better and improve our shit code.

Sleeep...Sleep...

After a few minutes trying to define a good spot for this second delivery now I have a plan :)

Objective:
- Support up to 100Req/s healthy, which is already a way better number than the previous version
- Consistency Guarantee
- Code Improvements
- If I have tim I'll provide a few kubernetes manifests so if you want to try this in your machine/cloud whatever, you can.

Which translates to:
- Still Monolith architecture - Yeah you won't need nothing more than that to hold the proposed amount of request/sec
- Ability to scale horizontally at least a little bit :)
- A Few workers to help us
- A better written code [do not be too happy, my effort will be limited to comprehension and not perfection :)]


## How would be the design

We will change the sync step-by-step for something a little bit more async

Flow:
-> Endpoint receive request -> build a message and publish to RabbitMQ Quorum Queue -> Create Transaction with Pending Status
-> Worker A -> Read msg and try to acquire the distributed lock for that accounts -> if ok publish a new message to another queue
-> Worker B -> Read msg and process the payment -> if ok update transaction status to done -> release the lock


Q&A:
- Q: Why RebbitMQ?
 - - A: Two reasons, 1 - at this moment I only know kafka and rabbitmq that could handle this for us, 
     but I have never worked with kafka before neither rabbitmq, but rabbitmq seems way easier to use. 2 - Because I want to play with rabbitmq :)
     !both will work here, even dough the effort with kafka can be re-usable in the next level... anyway your choice :)

- Q: Why use a distributed lock?
 - - A: The distributed lock will help to manage the consistency for our balance, eliminating the possibility of concurrent processing for the same accounts.

- Q: Why use redis for a distributed lock manager?
 - - A: At this point we could use the PostgreSQL to handle this job for us but this will translate to more job for our little baby that will need to handle a bunch of things already,
     So, the other alternatives that I know either would require the redis Or something like zookeeper, but since I never touched zookeeper I will stick with a cheap redis.


Deps:
- Gin [Why? I spend 2+ years using fasthttp which will be a best option for performance purposes but i want to play with Gin :)]
- Cobra [Why? At least for now is best CMD library that I know, easy to use and pretty flexible]
- env [Why? even dough we already had viper in this project I hate viper env embed struct definition, too dirty for me, feel free to use whatever you want]
- amqp091-go [Why? seems to be the most trustable lib available for rabbitmq, still has a few limitations and constraints but im planning to implement a few extensions for learning purposes in the near future]