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
- redis-go [Why? It's pretty simple to use and it is the official go lib for redis]



# !After implementation Concerns
I'm currently having a lot of things to deal at the same time so, I have promised quite a bit things before that I'll not be able to manage right now.
There's more issues that I can count right now and that's not a surprise, even dough I would love to complete this version and make it perfect at some level, It won't be possible because I still want to make the version 3 and other parallel things too...

Anyway, What will need to change here to be able handle the volume I proposed paying for small machines?
!Again this is just a generic idea, every system will have its own constraints and needs...

Application:
- You can keep the application as it is in terms of service setup and high level flow
  -- keep the workers running in the same machine as the main application, each machine with something like 512mb, 1CPU core, should give you enough to play[off course assuming there's nothing to heavy on it...]
  -- as you payment logic grows you will need to re-design the workers idea so you have more small units on each instead of a complete flow in only one

Redis:
- In order for the distributed lock works properly first you will need to have a cluster and not only one machine
- After the cluster been set, you now can use a library to make the lock for you or implement yourself [if you are doing it for learning purposes I personally encourage you write your own, will be really cool :)]
- remember to save the data into a volume of at least the distributed lock info

RabbitMQ:
- If you want you can keep this and just make a decent setup for it, you will not need a lot from here. But that's your choice...
- In order to this work you have to re-write/implement your own consumer/publisher/exchange code, what we have here its just few hours that i had free and decide to play with rabbitmq, I have no idea if that would work in real word,
find a good library or invest some time to write your own code...
- You will only need one avg machine with disk space for the current setup and that should be fine

Postgresql:
- For the proposed transaction per seconds, one primary machine to handle the writes + one read replicas machine should be good
- I have added some notes about this in the query file but will add it here again just in case... For performance purposes do not use view table,
instead use materialized view and only on read replica.
- I have not added any index, feel free to add it.

Distribtued lock:
- I have not locked both from-to account, only from, this will need to be fixed to avoid balance issues
- This version use set method from redis, this should be changed to setNX to guarantee it will raise an error if that key already exists

That's all remember for now...
Summary:
- 3 Apps machines/pods
- 3 redis machines - 1cluster
- 1 rebbitmq machine
- 2 postgresql machine
- 1 monolithic code base
+logging, monitoring etc...
The concept here is simple, trying to serve as much as possible with low resources, that setup will guarantee you
minimum 100/requests/second easily, but if you play around with it you will notice that this could easily handle up-to 10 times this value
if you configure and allocate resources for it properly and off course design a decent payment flow based on this primitive design.