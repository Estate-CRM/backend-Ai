# Estate-CRM Backend AI

A scalable backend system for real estate CRMs, leveraging both Go and Python microservices, Apache Kafka for real-time streaming, and ML/analytics workflows.

---

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
