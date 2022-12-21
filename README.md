# Drive N Go'
## About program
### 2 groups of users:
1. Passengers
2. Drivers
### It is a ride-sharing platform using microservice architecture. 

## Design Considerations

To start off with the design of the Drive N Go platform, it mainly uses 3 microservices, the Console microservice, the Passenger microservice and the Driver microservice. This design is inspired from the Domain-Driven Design that we had to research about. 

These microservices act as separate programs and they are not interconnected such that they each can operate on their own. Hence, when working with Passenger or the Driver, they both can run separately as individual sub-domains. This make sures that they are loosely coupled and when one microservice fails, the other is not dependent on it. Also, they both can be independently worked on.

Each of the microservice has its own sql script, their own table to retrieve and store data from. 
## Architecture Diagram
<img width="518" alt="ArchitectureDiagram" src="https://user-images.githubusercontent.com/74234483/208900763-633102d8-8400-443f-b708-8140f5500070.png">
<br>For the architecture of the program, it starts out by the user interacting with the console. Hence, our console is the frontend.<br/>
<br>Then, the console interacts with the Passenger microservice and the Driver microservice with API reference so do some actions that the user specifies.<br/>
<br>Next, when the microservices needs to access their database, they will use queries to retrieve and send messages from the sql databases. Thus, the whole program is working together independently.<br/>

## Instructions
- Passenger
**<br>Action 1: Login<br/>**
A passenger is able to login through the main menu.<br/>
**<br>Action 2: Create account<br/>**
A passenger is able to create account through the main menu. Then he will select the 'Passenger' role and input his information.<br/>
**<br>Action 3: Update account details<br/>**
A passenger is able to update account details by first logging in. Then he will select the Update account option and input his new information.<br/>
**<br>Action 4: Book a ride<br/>**
A passenger is able to book a ride after logging in. Then choose the book a ride option. There, he needs to input the postal codes to and fro, and the system will check for an available driver.<br/>
**<br>Action 5: View bookings<br/>**
A passenger is able to view his previous bookings after logging in. Then choose the view bookings option and it will be displayed.

- Driver
**<br>Action 1: Login<br/>**
A driver is able to login through the main menu.<br/>
**<br>Action 2: Create account<br/>**
A driver is able to create account through the main menu. Then he will select the 'Driver' role and input his information.<br/>
