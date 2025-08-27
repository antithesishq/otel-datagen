# otel-datagen: Visual OTLP Data for Pipeline Testing

Quick visual examples for OTLP pipeline builders using the Grafana LGTM stack.

## Setup
```bash
docker run -p 3000:3000 -p 4317:4317 -p 4318:4318 --rm -ti grafana/otel-lgtm
```
Access Grafana at http://localhost:3000 (admin/admin)

---

## Metrics: From Simple to Complex

### 1. Basic Time Series
```bash
# Single service, smooth gauge progression
./otel-datagen generate metrics \
  --metric-type=gauge \
  --metric-name=cpu_usage \
  --num-metrics=20 \
  --timestamp-start="-10m" \
  --timestamp-spacing="30s" \
  --counter-min=20 \
  --counter-max=80 \
  --otlp-endpoint=localhost:4317 \
  --resource-attr="service.name=web-server"
```
**Grafana Prometheus Query** 
```promql
avg_over_time(cpu_usage{service_name="web-server"}[10m])
```

### 2. Distribution Patterns  
```bash
# Histogram with realistic response time distribution
./otel-datagen generate metrics \
  --metric-type=histogram \
  --metric-name=http_request_duration \
  --num-metrics=50 \
  --timestamp-start="-15m" \
  --timestamp-spacing="15s" \
  --counter-min=10 \
  --counter-max=500 \
  --aggro-numeric="" \
  --otlp-endpoint=localhost:4317 \
  --resource-attr="service.name=api-gateway"
```
***Grafana Prometheus Query** 
```promql
sum by (le) (rate(http_request_duration_bucket{service_name="api-gateway"}[5m]))
```

### 3. Multi-Service Correlation
```bash
# Generate load across three services simultaneously
./otel-datagen generate metrics --metric-type=counter --metric-name=requests_total --num-metrics=30 --timestamp-start="-20m" --timestamp-spacing="30s" --counter-min=100 --counter-max=1000 --otlp-endpoint=localhost:4317 --resource-attr="service.name=frontend" &

./otel-datagen generate metrics --metric-type=gauge --metric-name=memory_usage --num-metrics=30 --timestamp-start="-20m" --timestamp-spacing="30s" --counter-min=40 --counter-max=90 --otlp-endpoint=localhost:4317 --resource-attr="service.name=backend" &

./otel-datagen generate metrics --metric-type=histogram --metric-name=db_query_duration --num-metrics=30 --timestamp-start="-20m" --timestamp-spacing="30s" --counter-min=1 --counter-max=100 --aggro-numeric="" --otlp-endpoint=localhost:4317 --resource-attr="service.name=database"
```
***Grafana Prometheus Query** 
```promql
{service_name=~"frontend|backend|database"}
```

---

## Logs: Structured to Complex

### 1. Clean Log Stream
```bash
# Structured application logs
./otel-datagen generate logs \
  --num-logs=25 \
  --timestamp-start="-8m" \
  --timestamp-spacing="20s" \
  --otlp-endpoint=localhost:4317 \
  --resource-attr="service.name=user-service"
```
**Grafana Loki Query** 
```logql
{service_name="user-service"}
```


### 2. Error Detection Patterns
```bash
# Logs with aggro case errors for alerting patterns
./otel-datagen generate logs \
  --num-logs=40 \
  --timestamp-start="-12m" \
  --timestamp-spacing="15s" \
  --aggro-string="" \
  --otlp-endpoint=localhost:4317 \
  --resource-attr="service.name=payment-api"
```
**Grafana Loki Query** 
```logql
{service_name="payment-api"} | aggro_string != ""
```

### 3. Service Interaction Logs
```bash
# Multi-service log correlation showing request flow
./otel-datagen generate logs --num-logs=20 --timestamp-start="-10m" --timestamp-spacing="25s" --otlp-endpoint=localhost:4317 --resource-attr="service.name=nginx" --resource-attr="component=ingress" &

./otel-datagen generate logs --num-logs=20 --timestamp-start="-10m" --timestamp-spacing="25s" --aggro-string="" --otlp-endpoint=localhost:4317 --resource-attr="service.name=auth-service" --resource-attr="component=middleware" &

./otel-datagen generate logs --num-logs=20 --timestamp-start="-10m" --timestamp-spacing="25s" --otlp-endpoint=localhost:4317 --resource-attr="service.name=user-db" --resource-attr="component=storage"
```
**Grafana Loki Query** 
```logql
{service_name=~"user-db|nginx|auth-service"}
```


---

## Traces: Spans to Service Maps

### 1. Basic Span Timeline
```bash
# Simple operation traces
./otel-datagen generate traces \
  --num-traces=5 --num-spans=3 \
  --timestamp-start="-6m" \
  --timestamp-spacing="20s" \
  --otlp-endpoint=localhost:4317 \
  --resource-attr="service.name=order-processor"
```
**Grafana Tempo Query** 
```TraceQL
{resource.service.name = "order-processor"}
```

### 2. Complex Operation Traces
```bash
# Multi-span operations with attributes
./otel-datagen generate traces \
  --num-traces=8 --num-spans=3 \
  --num-attributes=5 \
  --timestamp-start="-10m" \
  --timestamp-spacing="15s" \
  --aggro-string="" \
  --otlp-endpoint=localhost:4317 \
  --resource-attr="service.name=checkout-flow"
```
**Grafana Tempo Query** 
```TraceQL
{resource.service.name="checkout-flow" && "aggro.string" != ""}
```

### 3. Service Dependency Mapping
```bash
# Distributed traces across service boundaries
./otel-datagen generate traces --num-traces=4 --num-spans=5 --num-attributes=3 --timestamp-start="-15m" --timestamp-spacing="30s" --otlp-endpoint=localhost:4317 --resource-attr="service.name=api-gateway" --resource-attr="tier=edge" &

./otel-datagen generate traces --num-traces=4 --num-spans=5 --num-attributes=4 --timestamp-start="-15m" --timestamp-spacing="30s" --aggro-string="" --otlp-endpoint=localhost:4317 --resource-attr="service.name=business-logic" --resource-attr="tier=application" &

./otel-datagen generate traces --num-traces=4 --num-spans=5 --num-attributes=2 --timestamp-start="-15m" --timestamp-spacing="30s" --otlp-endpoint=localhost:4317 --resource-attr="service.name=data-store" --resource-attr="tier=persistence"
```

**Grafana Tempo Query** 
```TraceQL
{resource.tier=~"edge|application|persistence"}
```

---

## Quick Navigation

- **Explore Metrics**: http://localhost:3000/explore?left=["now-1h","now","Prometheus",{}]
- **Explore Logs**: http://localhost:3000/explore?left=["now-1h","now","Loki",{}]
- **Explore Traces**: http://localhost:3000/explore?left=["now-1h","now","Tempo",{}]

## Pro Tips
- Use `--aggro-string=""` and `--aggro-numeric=""` to simulate realistic edge case scenarios
- Combine `--timestamp-start` and `--timestamp-spacing` for historical data patterns
- Layer multiple resource attributes for rich service topology
- Run commands in parallel (`&` + `wait`) for correlated multi-service data