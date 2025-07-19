# Estate-CRM Backend AI

A scalable backend for real estate CRMs using Go and Python microservices, with Apache Kafka for data streaming and task scheduling. It supports AI-powered property matching and comparison, contract generation, and ID verification using computer vision.

---
## Problems
- Lack of market demand data
- Manual Client Matching: Agents waste time finding the right buyers or renters for properties.
- Poor Follow-up System: No structured way to monitor past recommendations or interactions.
- No Smart Comparison Tools: Agents can't easily compare multiple clients or properties.
- Paper-based Contracts: Contract generation is manual, slow, and error-prone.
- 

## Solution

## Features
- Microservices Architecture: Developed using Go and Python for modularity, scalability, and ease of maintenance.

- Real-Time Data Streaming: Employs Kafka for asynchronous data flow between services using a producer-consumer model.

- Task Scheduling: Automates background jobs such as batch processing and data synchronization through cron jobs.

- AI-Powered Contact Matching: Matches properties with potential clients using intelligent algorithms.

- Smart Comparison Engine: Generates human-like explanations comparing different clients and highlighting their similarities. This is triggered after recommendations are made.

- Automated Contract Generation: Creates real estate contracts based on property and client data.

- ID Verification: Utilizes computer vision to validate client identity documents.

## ⚙️ Components Description

### 🔹 Go Backend (Service 1)
- Built with Go (`golang`)
- Responsible for:
  - Serving secured API routes (e.g., JWT auth, validation)
  - Fetching **paginated** data from the database
  - Producing data to Kafka in **JSON chunks**
- Pagination is used to avoid memory overload and improve network efficiency

### 🔹 Apache Kafka
- Acts as a **bridge** between Go and Python
- Provides a **reliable, asynchronous messaging system**
- Enables real-time, scalable data streaming

### 🔹 Python Consumer (Service 2)
- Consumes JSON messages from Kafka
- Converts incoming contact data to **CSV format**
- CSV data is used as input to a **recommendation model** or other ML/analytics workflows

---

## 📁 Folder Structure

```
.
├── docker-compose.yml         # Orchestrates Go, Python, and Kafka services
├── go/                       # Go backend service (API, Kafka producer)
├── flask/                    # Python consumer service (Kafka consumer, CSV, ML)
│   ├── .gitignore
│   ├── Dockerfile
│   ├── requirements.txt
│   ├── run.py
│   ├── test.py
│   └── app/
└── README.md
```

- **go/**  
  _Go backend microservice (details not shown here)_

- **flask/**  
  _Python microservice for consuming Kafka, data conversion, and ML:_
  - `run.py` – Entrypoint for the Flask app
  - `test.py` – Scripts or tests for Python components
  - `requirements.txt` – Python dependencies
  - `Dockerfile` – Container setup for the Python consumer
  - `app/` – Application source code

- **docker-compose.yml**  
  _Development and orchestration setup for all services and Kafka broker_

---

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go (for standalone development)
- Python 3.8+ (for standalone development)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/Estate-CRM/backend-Ai.git
cd backend-Ai

# Start all services (Go, Python, Kafka, Zookeeper)
docker-compose up --build

python flask/run.py
```

---

## 🛠️ Development

- **Go Service:**  
  See `/go` for API and Kafka producer code.

- **Python Service:**  
  See `/flask` for Kafka consumer and CSV/ML logic.

---

## 📝 License

Distributed under the MIT License.

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!  
Feel free to check the [issues page](https://github.com/Estate-CRM/backend-Ai/issues).
