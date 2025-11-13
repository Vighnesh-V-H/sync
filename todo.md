# Real-Time Analytics Engine - 4-DAY BLITZ

**Mission**: Ship a working real-time analytics engine in 96 hours. 100k events/sec, <1s latency.

**Rules**: 
- 18-hour days (6am-12am)
- Ship or die
- No perfect code
- Parallel everything

---

## **DAY 1: FULL STACK SKELETON (18 hours)**
*Goal: End-to-end data flow working - ingest → process → store → query*

### **Hour 0-3: Project Setup (SPEED RUN)**
- [ ] Init Go project: `/cmd/{ingestor,processor,querier}`, `/internal/{event,queue,storage}`
- [ ] `go.mod` deps: `gorilla/mux`, `go-redis/redis/v9`, `clickhouse-go/v2`, `prometheus/client_golang`
- [ ] Docker Compose: Redis + ClickHouse + 3 services
- [ ] Core event struct:
  ```go
  type Event struct {
    Timestamp time.Time
    TenantID string
    EventType string
    Payload map[string]interface{}
    IdempotencyKey string
  }
  ```
- [ ] Health checks (`/health`) + Prometheus metrics init

### **Hour 3-7: Ingestor Service**
- [ ] HTTP server: `POST /v1/events` (JSON validation with `gojsonschema`)
- [ ] gRPC server: `IngestBatch(events[])` with Protobuf
- [ ] Redis Streams writer: batch 100 events OR 100ms flush
- [ ] Idempotency: LRU cache (10k capacity, 5min TTL)
- [ ] Metrics: `events_received_total`, `events_queued_total`
- [ ] **TEST**: `wrk` load test → 10k req/sec target

### **Hour 7-11: Processor Service**
- [ ] Redis Streams consumer with `XREADGROUP` (consumer groups per tenant)
- [ ] In-memory aggregator (ring buffers for 1min/5min/1hr windows):
  - Count, sum, avg per event_type
  - `sync.Map` for thread-safety
  - Flush every 10 seconds
- [ ] HyperLogLog integration (`github.com/dgryski/go-hll`) for unique counts
- [ ] Graceful shutdown (drain queue on SIGTERM)
- [ ] **TEST**: Process 100k events/sec with <500ms lag

### **Hour 11-15: ClickHouse Storage**
- [ ] Schemas:
  ```sql
  -- Raw events
  CREATE TABLE events (
    timestamp DateTime,
    tenant_id String,
    event_type String,
    payload String
  ) ENGINE=MergeTree() 
  PARTITION BY toYYYYMM(timestamp)
  ORDER BY (tenant_id, timestamp);
  
  -- Aggregates
  CREATE TABLE aggregates (
    window_start DateTime,
    tenant_id String,
    event_type String,
    metric_name String,
    value Float64
  ) ENGINE=SummingMergeTree()
  PARTITION BY toYYYYMM(window_start);
  ```
- [ ] Batch insert (1000+ events/batch) with retry logic
- [ ] TTL: 7 days raw, 90 days aggregates
- [ ] **TEST**: 500k inserts/sec sustained

### **Hour 15-18: Query API**
- [ ] REST endpoints:
  - `GET /v1/metrics?tenant_id=X&start=T1&end=T2&group_by=Y`
  - `GET /v1/events?tenant_id=X&filters=...`
- [ ] Redis cache layer (TTL = window size)
- [ ] WebSocket `/v1/stream` with Redis PubSub for live updates
- [ ] Query timeouts (5s max), pagination (cursor-based)
- [ ] **TEST**: <100ms cached, <500ms cold queries

**EOD CHECKPOINT**: Send 100k events → see them in dashboard in <2 seconds

---

## **DAY 2: SCALE + RELIABILITY (18 hours)**
*Goal: Multi-tenant, horizontally scalable, production-grade observability*

### **Hour 0-4: Multi-Tenancy & Security**
- [ ] JWT auth middleware (extract `tenant_id` from token)
- [ ] Tenant isolation: WHERE tenant_id = ? on ALL queries
- [ ] Resource quotas (Redis):
  - Max 10k events/sec per tenant
  - Max 1GB storage per tenant
- [ ] Rate limiting per tenant (100 req/sec)
- [ ] **TEST**: Verify tenant A can't access tenant B's data

### **Hour 4-8: Horizontal Scaling**
- [ ] Consistent hashing (`hashicorp/memberlist`) for processor sharding
- [ ] Shard registry in etcd (leader election for rebalancing)
- [ ] Auto-shard by tenant_id hash
- [ ] Ingestor: Load balancer ready (stateless)
- [ ] Querier: ClickHouse read replicas (1 writer, 2 readers)
- [ ] **TEST**: Kill 1 processor node → shards reassign in <5s

### **Hour 8-12: Observability Stack**
- [ ] OpenTelemetry tracing (trace context propagation)
- [ ] Jaeger setup for visualization
- [ ] Structured logging with `zap` (request_id, tenant_id, latency)
- [ ] Grafana dashboards:
  - Ingestion rate, processing lag, query latency
  - Error rates, queue depth, CPU/memory per service
- [ ] Prometheus alerts: high error rate (>1%), queue lag (>10s)
- [ ] **TEST**: Trace request end-to-end, identify bottlenecks

### **Hour 12-16: Resilience Engineering**
- [ ] Circuit breakers (`github.com/sony/gobreaker`) on ClickHouse/Redis clients
  - Trip after 5 failures, retry after 30s
- [ ] Exponential backoff for queue retries
- [ ] Dead-letter queue for poison messages
- [ ] Graceful degradation: if Redis down, ingestor returns 503 (don't crash)
- [ ] Connection pooling tuning (MaxIdleConns=100, MaxOpenConns=500)
- [ ] **CHAOS TEST**: Randomly kill ClickHouse → system recovers without data loss

### **Hour 16-18: Alerting Engine**
- [ ] Rule evaluator (YAML config):
  ```yaml
  alerts:
    - name: high_error_rate
      condition: avg(errors) > 100
      window: 5m
      action: slack_webhook
  ```
- [ ] Slack webhook integration
- [ ] Alert cooldown (5min between alerts)
- [ ] **TEST**: Simulate spike → alert fires in <15s

**EOD CHECKPOINT**: Run 100k events/sec for 1 hour straight. Zero data loss. <1% error rate.

---

## **DAY 3: ADVANCED FEATURES (18 hours)**
*Goal: ML anomalies, exports, backfill, security hardening*

### **Hour 0-4: Anomaly Detection**
- [ ] Z-score detector: flag metrics >3 std dev from 7-day baseline
- [ ] Train baseline on startup (historical aggregates)
- [ ] Real-time scoring on each aggregate flush
- [ ] Anomaly feed via WebSocket (`/v1/anomalies`)
- [ ] **TEST**: Inject outlier → detect in <1s

### **Hour 4-7: Export & Integrations**
- [ ] Webhook endpoint: POST aggregates to external URLs
- [ ] Retry logic (exponential backoff, max 5 retries)
- [ ] Kafka producer for upstream (`github.com/segmentio/kafka-go`)
- [ ] Prometheus remote write for downstream monitoring
- [ ] **TEST**: Export 10k events/sec to webhook without drops

### **Hour 7-10: Backfill & Replay**
- [ ] Replay API: `POST /v1/replay` (reprocess events from date range)
- [ ] S3 archival for cold storage (Parquet format)
- [ ] Schema migration tool (update event structure, backfill)
- [ ] **TEST**: Change schema → backfill 1M historical events in <10min

### **Hour 10-13: Security Hardening**
- [ ] mTLS on gRPC endpoints (cert rotation every 30 days)
- [ ] RBAC: read/write permissions per tenant (store in Redis)
- [ ] Audit log: all mutations (INSERT/DELETE) logged to separate ClickHouse table
- [ ] Input sanitization: prevent SQL injection, XSS
- [ ] **TEST**: Penetration test with `sqlmap`, rate limit bypass attempts

### **Hour 13-16: Performance Optimization**
- [ ] Profile with `pprof`: identify CPU/memory hotspots
- [ ] Optimize allocations (use sync.Pool for event buffers)
- [ ] Zero-copy Protobuf marshaling
- [ ] Batch Redis operations (pipeline 100 commands)
- [ ] ClickHouse query optimization (add indexes, materialized views)
- [ ] **TEST**: Push to 200k events/sec sustained

### **Hour 16-18: UI Dashboard (HTMX + Chart.js)**
- [ ] Go web server: `GET /dashboard` (serve HTML)
- [ ] HTMX polling for metrics updates (every 1s)
- [ ] Chart.js live graphs (WebSocket for real-time)
- [ ] Query builder UI (dropdowns for filters, no SQL required)
- [ ] Mobile responsive (Tailwind CSS)

**EOD CHECKPOINT**: 200k events/sec. Anomalies detected. Webhooks firing. Dashboard live.

---

## **DAY 4: PRODUCTION READY (18 hours)**
*Goal: Deploy, test at scale, document, SHIP*

### **Hour 0-4: Kubernetes Deployment**
- [ ] Helm charts for all services
- [ ] ConfigMaps for env vars, Secrets for credentials
- [ ] HPA (Horizontal Pod Autoscaler): scale on queue depth
- [ ] PodDisruptionBudget: ensure 2+ replicas always running
- [ ] Liveness/readiness probes
- [ ] **TEST**: Deploy to staging cluster, verify autoscaling

### **Hour 4-7: Disaster Recovery**
- [ ] ClickHouse backup to S3 (hourly snapshots)
- [ ] Restore procedure documented (RTO: <15min, RPO: <1hr)
- [ ] Redis persistence: AOF + RDB snapshots
- [ ] Blue-green deployment test (zero downtime rollout)
- [ ] **TEST**: Delete production DB → restore from backup in <15min

### **Hour 7-11: Load Testing at Scale**
- [ ] Generate 1M events/sec with `vegeta` or custom Go tool
- [ ] Sustain for 1 hour → measure:
  - Ingestion lag (target: <1s)
  - Query latency p99 (target: <100ms)
  - Error rate (target: <0.1%)
  - Memory/CPU usage (should be <80%)
- [ ] Chaos engineering: kill 50% of nodes → verify recovery
- [ ] **TEST**: 1B events/day sustained (11,574 events/sec average)

### **Hour 11-14: Documentation Blitz**
- [ ] Architecture diagram (use Mermaid or draw.io)
- [ ] API documentation (Swagger/OpenAPI spec)
- [ ] Operational runbook:
  - How to deploy (kubectl commands)
  - How to scale (add nodes, increase quotas)
  - How to debug (common errors, log queries)
  - Incident response playbook
- [ ] README with quickstart (Docker Compose → live in 5min)
- [ ] Write blog post: "Building a Real-Time Analytics Engine in 4 Days"

### **Hour 14-16: Security Audit & Compliance**
- [ ] OWASP Top 10 checklist (SQL injection, XSS, CSRF, etc.)
- [ ] Secrets rotation test (update JWT keys, mTLS certs)
- [ ] Penetration testing (hire ethical hacker OR use `burpsuite`)
- [ ] GDPR compliance: data retention policies, right to delete
- [ ] **SIGN-OFF**: Security review passed

### **Hour 16-18: SHIP IT 🚀**
- [ ] Tag release: `v1.0.0`
- [ ] Push to production
- [ ] Monitor dashboards (War Room mode: all eyes on Grafana)
- [ ] Send launch email/Slack announcement
- [ ] Post-mortem doc template ready (for inevitable issues)
- [ ] **CELEBRATE**: You built a production system in 96 hours

**FINAL METRICS** (measure at hour 18):
- ✅ Events ingested: X million (target: >10M/hour)
- ✅ Query latency p99: <100ms
- ✅ System uptime: 99.9%+
- ✅ Error rate: <0.1%
- ✅ Cost per 1M events: $X (optimize later)

---

## **HOURLY RITUALS**
- Every 2 hours: Stretch, hydrate, 5min walk
- Every 4 hours: Ship intermediate checkpoint to Git
- Every 6 hours: Review metrics, adjust priorities
- Hour 18 each day: Hard stop. Sleep. Repeat.

## **SURVIVAL TACTICS**
- ❌ No Slack/email (async only)
- ❌ No refactoring (make it work, not pretty)
- ❌ No rabbit holes (time-box every task: 30min max)
- ✅ Automate everything (scripts > manual steps)
- ✅ Copy-paste shamelessly (StackOverflow is your friend)
- ✅ Parallel work (use `tmux` splits for multiple terminals)

## **CRITICAL PATH** (do these first each day):
1. **Day 1**: Ingestor → Processor → Storage (data must flow)
2. **Day 2**: Observability → Scaling (can't fix what you can't see)
3. **Day 3**: Performance → Security (fast + safe)
4. **Day 4**: Deploy → Test → Document (ship or it doesn't exist)

## **EMERGENCY CONTACTS**
- Redis down? Fall back to in-memory queue (data loss acceptable for 4-day MVP)
- ClickHouse slow? Add aggressive caching, defer optimization
- Out of time? Cut scope, not quality (ship basic alerting over ML)

---

**REMEMBER**: 
- Scope creep kills. Build MVP, not perfection.
- Sleep 6 hours/day minimum (burnout fails on day 3).
- Document as you go (future-you will thank you).
- **SHIP. IT. 🚀**
