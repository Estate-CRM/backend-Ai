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
- Sends paginated JSON data to the Kafka topic (e.g., `contact-data`) for processing by Python.

### 2. Python Consumer (`python-consumer`)
- Consumes Kafka messages (JSON contact chunks).
- Converts them to CSV files, which are then used for a recommendation model (e.g., ML-based property recommendations).
- Simple Flask setup is used for potential status reporting or lightweight endpoints.

---

## 🔁 Data Flow

1. **Client** → Sends requests to Go API (e.g., `/api/contact/getAll`).
2. **Go Backend** → Fetches contacts from DB in chunks, publishes each chunk to Kafka.
3. **Kafka Broker** → Buffers and streams data reliably.
4. **Python Consumer** → Receives JSON chunks, converts to CSV, and prepares data for ML model.

---

## 🛠️ Run with Docker

### 1. Build & Run All Services except python service
bash

```
cd backend-Ai
docker-compose up --build
python flask/run.py
