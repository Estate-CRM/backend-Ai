FROM python:3.10-slim

# Set working directory to /app inside the container
WORKDIR /app

# Copy everything from flask/ into /app
COPY . .

# Install your Kafka client
RUN pip install confluent_kafka

# Run the consumer from its path inside the container
CMD ["python", "app/consumer.py"]
