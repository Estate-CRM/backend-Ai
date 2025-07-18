# app/app.py
from flask import Flask, jsonify
from confluent_kafka import Consumer, KafkaException
import threading
import json

app = Flask(__name__)

# Kafka consumer config
conf = {
    'bootstrap.servers': 'localhost:9092',
    'group.id': 'flask-consumer-group',
    'auto.offset.reset': 'earliest'
}

consumer = Consumer(conf)
topic = "your_topic_name_here"

@app.route('/')
def home():
    return "Kafka Flask Consumer is running"

# Optional endpoint to get last received data
latest_data = []

@app.route('/latest-contacts')
def get_contacts():
    return jsonify(latest_data)

def kafka_consumer():
    global latest_data
    consumer.subscribe([topic])
    print("📡 Kafka consumer started, waiting for messages...")
    while True:
        msg = consumer.poll(1.0)  # Poll every 1s
        if msg is None:
            continue
        if msg.error():
            print("❌ Kafka error: {}".format(msg.error()))
            continue

        try:
            contacts = json.loads(msg.value().decode('utf-8'))
            print(f"📥 Received contacts: {contacts}")
            latest_data = contacts  # Update global cache
        except Exception as e:
            print(f"❌ Failed to parse message: {e}")

# Start consumer in background thread
threading.Thread(target=kafka_consumer, daemon=True).start()

if __name__ == '__main__':
    app.run(debug=True)
