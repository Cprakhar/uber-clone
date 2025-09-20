# Uber Clone (Microservices + Realtime)

A production‑style, modular clone of core Uber workflows (trip request, driver assignment, payments, realtime updates) built with Go microservices, gRPC, HTTP (Gin), Kafka (event streaming & command routing), WebSockets (live rider/driver updates), Next.js web frontend, and Kubernetes for orchestration. Tilt accelerates local iterative development.

---
## Table of Contents
1. Overview
2. Architecture
3. Services
4. Tech Stack
5. Local Development (Recommended: Tilt)
6. Manual Local Run (Without Kubernetes)
7. Kubernetes Deployment (Minikube)
8. Environment Variables
9. Kafka & Messaging Model
10. Data Flows / Sequence
11. Running Tests (placeholder)
12. Observability (logging, metrics, tracing)
13. Troubleshooting & FAQ
14. Production Hardening Checklist
15. Roadmap / Future Work

---
## 1. Overview
This repository demonstrates a domain‑driven, event‑oriented backend for on‑demand trip booking:
- Rider creates trip → Trip Service persists & emits trip.event.created
- Driver Service consumes trip events, selects & notifies candidate driver(s)
- Driver responds (accept / decline) via driver command topics
- Trip Service updates state, emits driver assignment / not interested events
- API Gateway relays domain events to connected WebSocket clients (riders & drivers)
- Payment Service creates Stripe checkout sessions on command

## 2. Architecture
[![](https://mermaid.ink/img/pako:eNqNVt9v2jAQ_lcsP21qGvGjZSEPlSpaTX1YxWDVpAmpMvZBIkicOQ6UVf3fd4mdEgcKzQOK47v7vjt_d-aVcimAhnSW5vC3gJTDXcyWiiWzlOCTMaVjHmcs1eQpB3X49Xb88J1p2LLd4d4vFWdTUJuYw-HmnYo3oM5sH34fs10Cqf7Qb6oRFL-bnZLz5c3NnmRIRgrwteJGJmXOuTa2eyP0aFAPyXIyHtWO5YaxN7-PEoNJpNrM1nOSCw3Y_QuPWLq0nBvWl4h30fIos_Bhg5n6vMIVDTgVLyNN5II4LO9L65A8wrbyJimAyAkjwlbSBHBwENisQ2vl80T4pfezapbGRW1RHckkYaloILMNi9dsvgYXtHUQv2E-lXwF-hCccQ5ZE3sNiwZ0A_O2sjSw6oPTbB_nWbT3TB3dGEiykGrLlABBtCQTNp_H-sdP8sUwK3EmkGcS2-lnAQV8rUtwVCchGSvJIc8tJ2KoMGxDjxSZKIVapZZrpovc9_2j4nF7whGPifvM8jxepp8XkW2SzAQmOVKMZXoyF6_Nwq5bunetkLxp2HfIUQR8JQts5Bqz9DJGx3K1ZtZdHAVpCc9mZStkc3s-0WZtTFukOsGnB5QeE7u6PM4gKScQsozktmF_TmuN1qidUHZJa6q9V04m2RowlrV1SnbQc5GUq33YacFL_R2faHtPzxXt0ZN1O-6iXbRW1Zu4HxeiqjTJivk6ziPTcp-U1aXD-NPgZ46a21KL061Qj6kzc38_facarzCyv1tOdGgVk7OUpCipOSxjZEI9moBKWCzwKn8tQ8yojiCBGQ3xVTC1muEV_4Z2rNByuks5DbUqwKNKFsuIhgu2znFlZo79C1Cb4L36R8rmkoav9IWGvW_-1XVn0O_1-kE3GAyHgUd3-Lnb8fu9frc_xKfbvQ6CN4_-qyJ0_KDX7Q86QTDoDAfD66ve23_1IPGQ?type=png)](https://mermaid.live/edit#pako:eNqNVt9v2jAQ_lcsP21qGvGjZSEPlSpaTX1YxWDVpAmpMvZBIkicOQ6UVf3fd4mdEgcKzQOK47v7vjt_d-aVcimAhnSW5vC3gJTDXcyWiiWzlOCTMaVjHmcs1eQpB3X49Xb88J1p2LLd4d4vFWdTUJuYw-HmnYo3oM5sH34fs10Cqf7Qb6oRFL-bnZLz5c3NnmRIRgrwteJGJmXOuTa2eyP0aFAPyXIyHtWO5YaxN7-PEoNJpNrM1nOSCw3Y_QuPWLq0nBvWl4h30fIos_Bhg5n6vMIVDTgVLyNN5II4LO9L65A8wrbyJimAyAkjwlbSBHBwENisQ2vl80T4pfezapbGRW1RHckkYaloILMNi9dsvgYXtHUQv2E-lXwF-hCccQ5ZE3sNiwZ0A_O2sjSw6oPTbB_nWbT3TB3dGEiykGrLlABBtCQTNp_H-sdP8sUwK3EmkGcS2-lnAQV8rUtwVCchGSvJIc8tJ2KoMGxDjxSZKIVapZZrpovc9_2j4nF7whGPifvM8jxepp8XkW2SzAQmOVKMZXoyF6_Nwq5bunetkLxp2HfIUQR8JQts5Bqz9DJGx3K1ZtZdHAVpCc9mZStkc3s-0WZtTFukOsGnB5QeE7u6PM4gKScQsozktmF_TmuN1qidUHZJa6q9V04m2RowlrV1SnbQc5GUq33YacFL_R2faHtPzxXt0ZN1O-6iXbRW1Zu4HxeiqjTJivk6ziPTcp-U1aXD-NPgZ46a21KL061Qj6kzc38_facarzCyv1tOdGgVk7OUpCipOSxjZEI9moBKWCzwKn8tQ8yojiCBGQ3xVTC1muEV_4Z2rNByuks5DbUqwKNKFsuIhgu2znFlZo79C1Cb4L36R8rmkoav9IWGvW_-1XVn0O_1-kE3GAyHgUd3-Lnb8fu9frc_xKfbvQ6CN4_-qyJ0_KDX7Q86QTDoDAfD66ve23_1IPGQ)

Core patterns:
- CQRS-ish separation via explicit command (driver.cmd.* / payment.cmd.*) vs event (trip.event.* / trip.event.driver_*) topics.
- At-least-once delivery with idempotent handlers & explicit offset commits.
- Back-pressure safe Kafka producers (idempotent, acks=all, compression & linger tuned).
- Single consumer poll loop per service; WebSocket fan‑out decoupled from consumption.

## 3. Services
| Service | Responsibility | Interfaces |
|---------|----------------|-----------|
| api-gateway | Public HTTP (REST), WebSockets, request routing, authentication placeholder, event fan‑out | HTTP + WS, Kafka consumer |
| trip-service | Trip lifecycle, geospatial logic placeholder, event sourcing & assignment decisions | gRPC server, Kafka producer & consumer |
| driver-service | Driver registration & selection logic, reacts to trip events & issues driver commands | gRPC server, Kafka consumer & producer |
| payment-service | Stripe session creation & payment flow orchestration | Kafka consumer (commands), Kafka producer (events future) |
| web | Next.js frontend (pages/app router) | Browser -> Gateway |
| shared | Proto (gRPC), messaging abstractions, logging, metrics, env utilities, contracts | Imported libs |

## 4. Tech Stack
- Language: Go 1.24.x
- API / Transport: gRPC, HTTP (Gin), WebSockets
- Messaging: Kafka (confluent-kafka-go / librdkafka)
- Frontend: Next.js + Tailwind
- Storage: (Placeholder for Mongo / in-memory repo currently)
- Container: Docker, multi-stage builds
- Orchestration: Kubernetes (manifests in `deployments/k8s`), Minikube local
- Dev Loop: Tilt (`Tiltfile`)
- Observability: Structured logging scaffold, metrics/tracing placeholders under `shared/observe`

## 5. Local Development (Tilt)
Prerequisites:
- Docker
- Minikube (or other K8s cluster)
- Tilt (https://tilt.dev)
- kubectl
- (Linux) Install build deps for librdkafka if building locally; otherwise rely on Docker image.

Steps:
```bash
# Start (from repo root)
tilt up

# View UI (if not auto-open)
tilt ui

# After first spin-up, port-forward / ingress or use NodePort for web if external.
```
Tilt builds images, applies K8s manifests, and streams logs. Edit code → live rebuild.

## 6. Manual Local Run (Without Kubernetes)
You can run a simplified stack locally (helpful for quick backend iteration):
1. Start Kafka + Zookeeper (e.g., docker-compose or local binary). Example (single broker):
   ```bash
   docker run -d --name zookeeper -p 2181:2181 confluentinc/cp-zookeeper:7.2.15 \
     -e ZOOKEEPER_CLIENT_PORT=2181 -e ZOOKEEPER_TICK_TIME=2000
   docker run -d --name kafka -p 9092:9092 --link zookeeper \
     -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
     -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
     -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 \
     confluentinc/cp-kafka:7.2.15
   ```
2. Run services:
   ```bash
   go run ./services/trip-service
   go run ./services/driver-service
   go run ./services/payment-service
   go run ./services/api-gateway
   (cd web && npm install && npm run dev)
   ```
3. Set required env vars (see section 8) before launching.

## 7. Kubernetes Deployment (Minikube)
```bash
minikube start --memory=6144 --cpus=4
kubectl create namespace uber-clone
# Apply core infra (zookeeper, kafka, services, web) via Tilt or:
kubectl apply -k deployments/k8s/dev
```
Check readiness:
```bash
kubectl get pods -n uber-clone
kubectl logs -n uber-clone deploy/trip-service
```
Port-forward (example):
```bash
kubectl port-forward -n uber-clone svc/api-gateway 8080:80
```

## 8. Environment Variables
(Subset – add more as they are introduced)
| Variable | Service(s) | Purpose | Default |
|----------|------------|---------|---------|
| HTTP_ADDR | api-gateway | HTTP listen address | :8080 |
| KAFKA_BROKERS | all | Comma-separated broker list | kafka:9092 |
| STRIPE_SECRET_KEY | payment-service | Stripe API secret | (none) |
| STRIPE_SUCCESS_URL | payment-service | Success redirect | appURL?payment=success |
| STRIPE_CANCEL_URL | payment-service | Cancel redirect | appURL?payment=cancel |
| APP_URL | payment-service | Base web app URL | http://localhost:3000 |

## 9. Kafka & Messaging Model
Topic naming convention:
- Events: `trip.event.created`, `trip.event.driver_assigned`, `trip.event.driver_not_interested` (example)
- Commands: `driver.cmd.trip_request`, `driver.cmd.trip_accept`, `driver.cmd.trip_decline`, `payment.cmd.create_session`
Envelope (`contracts.KafkaMessage`):
```json
{
  "entityID": "<rider|driver|trip id>",
  "data": { ... domain payload ... }
}
```
WebSocket routing uses `entityID` to map to a connection.

Consumers:
- Single poll loop per service (no multiple concurrent Poll on same consumer).
- Manual commit only after successful handler → at-least-once.

Producers:
- Idempotent, acks=all, zstd compression, linger for batching.
- Fire-and-forget + optional synchronous send (wait for delivery).

## 10. Data Flows / Sequence (Happy Path)
1. Rider requests trip → Trip Service stores & emits `trip.event.created`.
2. Driver Service consumes, selects driver, sends `driver.cmd.trip_request` to rider (via Gateway).
3. Driver accepts (`driver.cmd.trip_accept`) → Trip Service emits `trip.event.driver_assigned`.
4. API Gateway pushes assignment to rider WS.
5. Rider initiates payment command → Payment Service creates session (Stripe) → (future: emits payment events).

## 11. Running Tests
(Currently minimal / placeholder) – Add unit tests per service:
```bash
go test ./...
```
Recommend adding:
- Producer/consumer integration test (using ephemeral Kafka container)
- Trip assignment logic unit tests

## 12. Observability
Present:
- Structured logs with clear prefixes (producer/consumer/ws-consumer)
Scaffold (future):
- Metrics: Prometheus counters for produced / consumed / failures
- Traces: OpenTelemetry instrumentation around gRPC & HTTP handlers

## 13. Troubleshooting & FAQ
| Symptom | Likely Cause | Fix |
|---------|--------------|-----|
| `connection refused kafka:9092` | Kafka not ready / wrong advertised listeners | Ensure `KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092` and pod Running |
| Consumer stops receiving events | Multiple Poll loops / consumer closed by WS | Use a single TopicConsumer started in main, don’t close per connection |
| WebSocket clients not seeing driver assignment | Assignment produced before WS subscribed | Ensure WS connects earlier or cache last assignment per trip |
| SIGSEGV in confluent Poll | Poll after Close race | Never Close while another goroutine polls; single poll loop design |
| CORS / WS blocked in browser | Missing CORS headers in gateway | Add proper `Access-Control-Allow-*` and upgrade handling |
| Stripe 401 errors | Missing STRIPE_SECRET_KEY | Set key & restart payment service |

Kafka debugging:
```bash
kubectl exec -it -n uber-clone kafka-0 -- kafka-topics --bootstrap-server localhost:9092 --list
kubectl exec -it -n uber-clone kafka-0 -- kafka-console-consumer --bootstrap-server localhost:9092 --topic trip.event.created --from-beginning
```

## 14. Production Hardening Checklist
- [ ] Replace in-memory repos with persistent storage (Mongo, Postgres)
- [ ] Authentication / authorization (JWT / OAuth) at gateway
- [ ] Rate limiting & request validation
- [ ] Schema registry & versioned event payloads
- [ ] Dead-letter / retry topics for poison messages
- [ ] Metrics & tracing instrumentation (Prometheus + OTLP exporter)
- [ ] Proper health/readiness endpoints per service
- [ ] Secure Kafka (SASL/SSL) and secrets management (K8s Secrets / Vault)
- [ ] CI pipeline (lint, vet, tests, security scans, image signing)
- [ ] Canary / blue‑green deploy strategy & HPA autoscaling

## 15. Roadmap / Future Work
- Driver location streaming & proximity matching (geohash indexing)
- Surge pricing module
- Trip state machine persistence (event sourcing + snapshots)
- Payment event acknowledgments (webhook ingestion) & refund flow
- Frontend improvements (driver dashboard, trip history)
- Multi-tenant / multi-region cluster partitioning

---
## Quick Start (TL;DR)
```bash
minikube start
tilt up
# In another terminal watch logs; in browser open the web app
```