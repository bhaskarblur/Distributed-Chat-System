### This script creates partitions for each chat service in Kafka for given topic ###

#!/bin/bash
# Configuration variables
TOPIC_NAME="chat-message"
PARTITIONS=3
REPLICATION_FACTOR=1
BROKER="kafka:9092"

# Function to check if the Kafka topic already exists
topic_exists() {
  kafka-topics --bootstrap-server $BROKER --list | grep -q "^$TOPIC_NAME$"
}

# Wait for Kafka to be ready
echo "Waiting for Kafka broker at $BROKER to be ready..."
while ! echo exit | nc kafka 9092; do
  sleep 5
done
echo "Kafka broker is ready."

# Check if the topic already exists
if topic_exists; then
  echo "Topic '$TOPIC_NAME' already exists. Skipping creation."
else
  # Create the topic
  echo "Creating topic '$TOPIC_NAME' with $PARTITIONS partitions..."
  kafka-topics --bootstrap-server $BROKER \
    --create \
    --topic $TOPIC_NAME \
    --partitions $PARTITIONS \
    --replication-factor $REPLICATION_FACTOR

  if [ $? -eq 0 ]; then
    echo "Topic '$TOPIC_NAME' created successfully."
  else
    echo "Failed to create topic '$TOPIC_NAME'."
    exit 1
  fi
fi

# Describe the topic
echo "Describing topic '$TOPIC_NAME'..."
kafka-topics --bootstrap-server $BROKER --describe --topic $TOPIC_NAME
