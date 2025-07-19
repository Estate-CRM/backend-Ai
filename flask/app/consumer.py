from confluent_kafka import Consumer, KafkaException
import json

conf = {
    'bootstrap.servers': 'kafka:9092',  # Or 'kafka:9092' if running inside Docker
    'group.id': 'python-consumer-group',
    'auto.offset.reset': 'earliest'
}

consumer = Consumer(conf)
topic = "contacts-topic"

consumer.subscribe([topic])

print(f"🚀 Listening to Kafka topic '{topic}'...")

try:
    while True:
        msg = consumer.poll(timeout=1.0)

        if msg is None:
            continue
        if msg.error():
            raise KafkaException(msg.error())
        
        # Decode and print message
        data = msg.value().decode('utf-8')
        contacts = json.loads(data)
        print("📥 Received contacts:")
        for contact in contacts:
            print(contact)

except KeyboardInterrupt:
    print("🛑 Stopping consumer...")

finally:
    consumer.close()
