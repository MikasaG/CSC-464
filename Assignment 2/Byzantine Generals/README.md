# Byzantine Generals
The **Byzanting Generals** Problem (a.k.a Byzantine Fault Tolerance Problem) is the dependability of a [fault-tolerant computer system](https://en.wikipedia.org/wiki/Fault-tolerant_computer_system "Fault-tolerant computer system"), particularly [distributed computing](https://en.wikipedia.org/wiki/Distributed_computing "Distributed computing") systems, where components may fail and there is imperfect information on whether a component has failed. The difficulty is indicated by the fact that, in a "Byzantine Failure", a component such as a server can inconsistently appear both failed and functioning to failure-detection systems (In this algorithm, a faulty general will only send a wrong message if the receiver has an even ID).

# Implementation
This solution to byzantine problem is a **Golang** implementation of the [Oral Message(M) ](https://www.microsoft.com/en-us/research/uploads/prod/2016/12/The-Byzantine-Generals-Problem.pdf) algorithm that proposed by [Leslie Lamport](https://en.wikipedia.org/wiki/Leslie_Lamport) in 1982. I create a node for each message sent from a general to another, then use a depth limited BFS to build a **Message Sending Tree** up to the depth of **M** (M is the number of traitors). Finally, a recursive function is used to determine the decision of each node.

# Algorithm Performance
- Since Byzantine Problem can **only** be solved when the number of traitors *M* is less than **1/3** of the total number of generals *G*.

- When *M* is less than 1/3 of *G*, This algorithm can guarantees that the generals will come to an agreement on the order that user given. but there is an exception, when the commander (first general) is one of the traitors, this algorithm can only guarantee that all generals come to consensus, it cannot guarantee that the consensus order is same as the order that user given.

# Test Method and Result
-Prompt Information:
![Prompt Information](https://raw.githubusercontent.com/MikasaG/CSC-464/master/Assignment%202/Byzantine%20Generals/images/start.png)