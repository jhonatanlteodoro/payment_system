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
- Still Monolith architecture - Yeah you wont need nothing more than that to hold the proposed amount of request/sec
- Ability to scale horizontally at least a little bit :)
- A Few workers to help us
- A better written code [do not be too happy, my effort will be limited to comprehension and not perfection :)]