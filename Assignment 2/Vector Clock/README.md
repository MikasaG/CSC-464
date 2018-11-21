# Vector Clock
A  **Vector Clock**  is an  algorithm  for generating a  partial ordering  of events in a distributed system  and detecting causality violations. Just as in Lamport timestamps, interprocess messages contain the state of the sending process's logical clock. A vector clock of a system of  _N_  processes is an array/vector of  _N_  logical clocks, one clock per process; a local "smallest possible values" copy of the global clock-array is kept in each process, with the following rules for clock updates:
-   Initially all clocks are zero.
-   Each time a process experiences an internal event, it increments its own logical clock  in the vector by one.
-   Each time a process sends a message, it increments its own logical clock in the vector by one (as in the bullet above) and then sends a copy of its own vector.
-   Each time a process receives a message, it increments its own logical clock in the vector by one and updates each element in its vector by taking the maximum of the value in its own vector clock and the value in the vector in the received message (for every element).
## Implementation
This is a **Golang** Implementation of vector clock **API**, the basic part only contain three function:
-  **newVectorClock(index  int,  thread_name  string,  threads  []string)**
>This  method  will  take  three  parameters  and  return  a  new  Vector_Clock

- **increaseOne()**
>Increments the value for a single  thread within a vector clock by 1 

- **update(v  map[string]int)**

>when  a  message  is  sent  from  one  thread  to  another,  the  receiver  call  this  method to  update  its  vector  clock.

I also implemented some other method for the convenient of testing and getting a better understanding of how vector clock identify which events have a **"happens before"** relationship, and which events might be **causally concurrent**, for example:
- **sendMessage(receiver  Vector_Clock,  ch  chan  map[string]int)**
>simulating the process of using a vector clock to send a message, the sender will  increments its own logical clock in the vector by one and then sends a copy of its own vector to the receiver



- **receiveMessage(sender  Vector_Clock,  ch  chan  map[string]int)**
>simulating the process of using a vector clock to receive a message, the receiver will  increments its own logical clock in the vector by one and updates each element in its vector by taking the maximum of the value in its own vector clock and the value in the vector in the received message 
# Test Method and Results
I wrote separate test units for all three basic functions, ensure they work properly. Furthermore, I found two typical vector clock scenarios from Internet. I implemented them and use them as a part of my testing, ensures that my vector clock behaves right in practical environment.

**Overall Testing Result:**
![Overall Testing Result](https://raw.githubusercontent.com/MikasaG/CSC-464/master/Assignment%202/Vector%20Clock/images/Ovrall%20Test%20Result.png)

**Three Unit Testing Results:**
![Three Unit Testing Results](https://raw.githubusercontent.com/MikasaG/CSC-464/master/Assignment%202/Vector%20Clock/images/3%20unit%20test%20result.png)

**Scenario 1:**
- Diagram:
![](https://raw.githubusercontent.com/MikasaG/CSC-464/master/Assignment%202/Vector%20Clock/images/s1-diagram.png)

- Testing Result:
![Testing Result for Scenario 1](https://raw.githubusercontent.com/MikasaG/CSC-464/master/Assignment%202/Vector%20Clock/images/s1%20result.png)


**Scenario 2:**
- Diagram:
![](https://upload.wikimedia.org/wikipedia/commons/thumb/5/55/Vector_Clock.svg/500px-Vector_Clock.svg.png)
- Testing Result:
![Testing Result for Scenario 2](https://raw.githubusercontent.com/MikasaG/CSC-464/master/Assignment%202/Vector%20Clock/images/s2%20result.png)
# "Happens Before" Vs. Causally Concurrent
We know that vector clock can identify both "happens before" relationship and causally concurrent relationship. That is:


 $$ a\quad happens\quad before\quad   b
\Leftrightarrow
 VC_a < VC_b $$
where VCa < VCb means each VCa[i] is less than or equal to VCb[i] (for all i, VCa[i] <= VCb[i]).

So, we can determine the relationship of  any two events by using vector clock. If we can find VCa < VCb or VCb < VCa, then we can conclude that there is a happens before relationship. If not, these two events are causally concurrent.

For example, in the scenario, I indicated  two typical pairs of these relationships.
![](https://github.com/MikasaG/CSC-464/blob/master/Assignment%202/Vector%20Clock/images/s2-diagram.svg.png?raw=true)