
# 🛰️ Kafka Client Connection Technical Specification

This document describes how external clients can connect to your Kafka server using **username and password authentication (SASL/SCRAM)**.

---

## 📍 General Information

| Description       | Value                |
|------------------|----------------------|
| **Kafka Public IP**     | `43.231.114.103`     |
| **Port**                | `9092`               |
| **Authentication**      | SASL/SCRAM (SHA-512) |
| **Username Example**    | `client1`            |
| **Password Example**    | `*************`      |
| **Mechanism**           | `SCRAM-SHA-512`      |
| **Security Protocol**   | `SASL_PLAINTEXT`     |

---

## 🔐 Security Configuration

The Kafka server uses **SASL_PLAINTEXT + SCRAM-SHA-512**. Clients must:

- Enable SASL authentication
- Use the SCRAM-SHA-512 mechanism
- Provide valid username and password

---

## **Authentication**

- Enable SASL
- Use mechanism: `SCRAM-SHA-512`
- Provide username and password

---

## 🛠️ Recommended Client SDKs

| Language | SDK                 |
|----------|---------------------|
| Go       | IBM/sarama          |
| Java     | Apache Kafka Client |
| Python   | kafka-python        |
| Node.js  | kafkajs             |

---

## 🔒 Access Control (ACL)

Each client has restricted access.

> Example: `User:client1` can only access `topic:client1-events`.

---