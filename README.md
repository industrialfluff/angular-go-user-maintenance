# angular-go-user-maintenance

# Project: Full-Stack Application with Angular, Go, PostgreSQL, Kafka, and MongoDB

## Overview

This project is a demonstration of building a full-stack web application with the following components:

- **Frontend:** Developed using Angular 16.
- **Backend:** Implemented with Go and the Gin framework.
- **Database:** PostgreSQL for storing user and application data.
- **Message Queue:** Kafka for event-driven communication.
- **Data Storage for Kafka messages:** A separate Go module captures Kafka messages and stores them in MongoDB.

This project showcases the integration of these technologies to create a scalable, robust, and efficient web application, emphasizing proficiency in frontend, backend, and distributed systems.

## Technologies Used

1. **Frontend:**
   - Angular 16
   - Angular Material (for UI components like tables and forms)
   
2. **Backend:**
   - Go
   - Gin Web Framework
   - PostgreSQL (as the main relational database)
   - Kafka (for message brokering)

3. **Data Storage:**
   - MongoDB (for storing Kafka messages)

## Project Structure

- **Frontend (Angular):** Located in the `/angular_go` directory.
- **Backend (Go and PostgreSQL):** Located in the `/go_userlist` directory.
- **Kafka Consumer (Go and MongoDB):** Located in the `/go_mongo_kafka` directory.

## How to Run the Project

### Prerequisites

- Node.js (v16 or higher)
- Go (v1.20 or higher)
- PostgreSQL (v13 or higher)
- MongoDB (v5 or higher)
- Kafka (with Zookeeper)

### Step 1: Setting up the Frontend (Angular)

1. Navigate to the `angular_go` directory.
   ```bash
   cd angular_go
   ```

2. Install Angular dependencies.
   ```bash
   npm install
   ```

3. Run the Angular development server.
   ```bash
   ng serve
   ```

4. The application will be available at `http://localhost:4200`.

### Step 2: Setting up the Backend (Go and PostgreSQL)

1. Ensure that PostgreSQL is installed and running. Create a database for the project.

2. Update the PostgreSQL connection configuration in the `backend` directory (likely located in `config.go`).
   ```go
   const (
       Host     = "localhost"
       Port     = 5432
       User     = "your-username"
       Password = "your-password"
       Dbname   = "your-dbname"
   )
   ```

3. Navigate to the `go_userlist` directory.
   ```bash
   cd go_userlist
   ```

4. Install the required Go modules.
   ```bash
   go mod tidy
   ```

5. Run database migrations (if applicable).

6. Run the Go backend server.
   ```bash
   go run .
   ```

7. The API server will be available at `http://localhost:8080`.

### Step 3: Running Kafka and Zookeeper

1. Ensure that Kafka and Zookeeper are installed and running.

2. Start Zookeeper.
   ```bash
   zookeeper-server-start.sh config/zookeeper.properties
   ```

3. Start Kafka.
   ```bash
   kafka-server-start.sh config/server.properties
   ```

4. Create Kafka topics.
   ```bash
   kafka-topics.sh --create --topic prism-user-create --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
   kafka-topics.sh --create --topic prism-user-delete --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
   kafka-topics.sh --create --topic prism-user-update --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
   ```

### Step 4: Setting up Kafka Consumer (Go and MongoDB)

1. Ensure that MongoDB is installed and running.  Create three collections: user-delete, user-new, and user-update.,

2. Navigate to the `go_mongo_kafka` directory.
   ```bash
   cd go_mongo_kafka
   ```

3. Update MongoDB connection settings in `consumer.go`.
   ```go
   const (
       mongoURI  = "mongodb://localhost:27017"
       dbName    = "mydb"
       collection = "messages"
   )
   ```

4. Install the required Go modules.
   ```bash
   go mod tidy
   ```

5. Run the Kafka consumer to listen for messages and store them in MongoDB.
   ```bash
   go run .
   ```

### Step 5: Testing the Full System

1. Interact with the Angular frontend by performing actions that will trigger API calls to the Go backend.
2. The backend will publish events to Kafka.
3. The Kafka consumer module will capture those messages and store them in MongoDB.

## Kafka Consumer and MongoDB

The **Kafka consumer** is a separate Go module that listens to Kafka messages from a specific topic (e.g., `prism-user-update`). Upon receiving a message, it parses the message and stores it in a MongoDB collection. This demonstrates event-driven architecture and decouples message processing from the main application.

### Kafka Consumer Process:

- Connects to the Kafka broker and subscribes to a topic.
- Reads messages as they arrive.
- Processes each message and inserts it into MongoDB, preserving the message structure.
  
This separation allows for horizontal scalability where multiple consumers can be added to handle high message volumes, ensuring efficient processing and persistence of event data.

## Conclusion

This project demonstrates the creation of a full-stack application leveraging Angular, Go, PostgreSQL, Kafka, and MongoDB, integrating both synchronous and asynchronous workflows. The use of Kafka allows the application to scale, while MongoDB serves as an ideal storage solution for unstructured event data. This system can be expanded further to include more microservices and advanced features such as fault-tolerant processing, security, and more.
