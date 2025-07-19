# Estate-CRM Backend AI

A scalable backend for real estate CRMs using Go and Python microservices, with Apache Kafka for data streaming and task scheduling. It supports AI-powered property matching and comparison, contract generation, and ID verification using computer vision.

---
## Problems
- Lack of market demand data
- Manual Client Matching: Agents waste time finding the right buyers or renters for properties.
- Poor Follow-up System: No structured way to monitor past recommendations or interactions.
- No Smart Comparison Tools: Agents can't easily compare multiple clients or properties.
- Paper-based Contracts: Contract generation is manual, slow, and error-prone.


## Solution
- Smart Recommendation System
- Client Comparison Engine
- Kafka-Driven Real-Time Architectur

  
## Features
- Microservices Architecture: Developed using Go and Python for modularity, scalability, and ease of maintenance.

- Real-Time Data Streaming: Employs Kafka for asynchronous data flow between services using a producer-consumer model.

- Task Scheduling: Automates background jobs such as batch processing and data synchronization through cron jobs.

- AI-Powered Contact Matching: Matches properties with potential clients using intelligent algorithms.

- Smart Comparison Engine: Generates human-like explanations comparing different clients and highlighting their similarities. This is triggered after recommendations are made.

- Contract Generation: Creates real estate contracts based on property and client data.

- ID Verification: Utilizes computer vision to validate client identity documents.

---

## 🧠 Architecture Overview

The backend is divided into **two main services**:

### 1. Go API Server (`go-backend`)
- Implements routes for:
  - Authentication (`/api/auth`)
  - Contact management (`/api/contact`)
  - Property management (`/api/property`)
  - Matchmaking (`/api/match`)
- Connects to PostgreSQL for data storage.
- Retrieves **contact data in paginated chunks** to optimize performance and avoid memory issues.



### 2. Go Producer (`go-producer`)
- Sends paginated JSON data to the Kafka topic (e.g., `contact-data`) for processing by Python.


### 3. Go Consumer (`go-consumer`)
- consume the recomondation data sent by python producer
- Handles and stores the response for further use or integration with the main backend.


### 4. Python Consumer (`python-consumer`)
- Consumes Kafka messages (JSON contact chunks).
- Converts them to CSV files, which are then used for a recommendation model (e.g., ML-based property recommendations).
- consume propritety batch and use teh recomandation system for each one 
- Simple Flask setup is used for potential status reporting or lightweight endpoints.


### 5. Python Producer (`python-producer`)
- Produces recommendation results and natural-language explanations for each property test case.
- Sends results back to the Kafka topic consumed by the Go backend (go-consumer).


---
## 🔁 Data Flow

1. **Client** → Sends a request to the Go API to create or update a contact.
2. **Go Backend** → Fetches contact data from the database in paginated chunks and publishes each chunk to Kafka.
3. **Kafka Broker** → Buffers and streams data reliably between services.
4. **Python Consumer** → Receives JSON chunks, converts them to CSV, and prepares the data for the ML recommendation model.
5. **Go Producer** → Sends a batch of properties to the Python Consumer to trigger the recommendation system.
6. **Python Producer** → After generating recommendations and explanations, sends the results back to the Go Consumer via Kafka.
7. **Go Consumer** → Receives the recommendation results and integrates them into the backend system.


---

## 🛠️ Run with Docker

### 1. Build & Run All Services except python service
bash

```
cd backend-Ai
docker-compose up --build
python flask/run.py
