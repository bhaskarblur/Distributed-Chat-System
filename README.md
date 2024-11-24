# Distributed Chat System

A scalable, distributed chat system built with Go, Kafka, and WebSocket. This system enables real-time communication between users, leveraging Kafka for message brokering and WebSocket for maintaining persistent connections.

---

## Features

- **Real-time Messaging**: Supports instant communication using WebSocket.
- **Scalability**: Designed to handle a large number of users and messages with Kafka's horizontal scaling.
- **Distributed Architecture**: Runs multiple chat service instances seamlessly.
- **Dynamic Message Routing**: Routes messages to the correct chat service based on user connections.
- **Fault Tolerance**: Ensures reliable message delivery, even during failures.
- **Dynamic Topic Creation**: Automatically creates Kafka topics for new server instances.

---

## Architecture

1. **WebSocket Layer**:
   - Maintains persistent connections with users.
   - Forwards received messages to Kafka.

2. **Kafka (Message Broker)**:
   - Distributes messages across chat service instances.
   - Routes messages to appropriate partitions based on `ServerId`.

3. **Chat Services**:
   - Each service consumes messages from Kafka partitions.
   - Sends messages to connected users via WebSocket.

4. **User-to-Server Mapping**:
   - Maintains a mapping of which user is connected to which server.
   - Ensures message delivery to the correct server.
