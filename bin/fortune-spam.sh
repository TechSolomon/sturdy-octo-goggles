#!/bin/bash
MQTT_HOST_DEFAULT="localhost"
MQTT_HOST="${MQTT_HOST:-$MQTT_HOST_DEFAULT}"
MQTT_TOPIC_DEFAULT="fortune"
MQTT_TOPIC="${MQTT_TOPIC:-$MQTT_TOPIC_DEFAULT}"

while true; do
    mosquitto_pub -h "$MQTT_HOST" -t "$MQTT_TOPIC" -m "$(fortune)"
    sleep 10
done
