# CSC 464 Final Project
**Project:** Investigate whether Distributed Redis Cluster can help web application to handle huge access requests.

**Solution:** Build a web application (a shopping system) using Spring Boot, comparing its performance before/after using Distributed Redis Cluster as data cache layer.

**Evaluation:**

*Correctness*

Manually tested the whole process of placing an order. The database correctly reflects the stock change

*Performance*

Using Apache JMeter to test the performance of the implementation. set up 1000 threads, each of them represents a single user, and each thread will try to access the goods list 10 times in 1 second. So, there are 10000 access requests in 1 second, Comparing QPS.



*Comprehensibility*

Using Spring boot to construct my project makes it easier for people to understand the code and configuration files clearly, the code is divided into four layers, front end layer, controller layer, service layer and data access layer.

**Future Work:**

The industrial solution for a web app always includes a combination of hashing, sharing, replication, distribution, high fail-over, key-value stores, etc. Using Redis Cluster is just a part of the concurrency optimization for a web app. In order to build a successful large-scale web app, there are many other aspects that need to be optimized. For example, the data consensus of distributed servers, the fault-tolerant technology of servers, and the data replication strategy etc. I can analyze their performance in the future.