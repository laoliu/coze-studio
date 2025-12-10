# 通用领域内容生成·元工作流平台 - 部署指南 🚀

**文档版本**: V2.0  
**日期**: 2025-12-09  
**适用版本**: MetaWorkflow Platform v2.0+

---

## 📋 目录

### 快速开始
1. [系统要求](#1-系统要求)
2. [架构概览](#2-架构概览)
3. [快速部署（Docker Compose）](#3-快速部署docker-compose)

### 生产部署
4. [Kubernetes部署](#4-kubernetes部署)
5. [数据库配置](#5-数据库配置)
6. [向量数据库配置（Milvus）](#6-向量数据库配置milvus)
7. [缓存配置（Redis）](#7-缓存配置redis)
8. [对象存储配置（MinIO）](#8-对象存储配置minio)

### 监控与运维
9. [监控系统配置](#9-监控系统配置)
10. [日志管理](#10-日志管理)
11. [备份与恢复](#11-备份与恢复)
12. [性能优化](#12-性能优化)
13. [安全加固](#13-安全加固)
14. [故障排查](#14-故障排查)

### 附录
15. [环境变量参考](#15-环境变量参考)
16. [端口映射表](#16-端口映射表)
17. [常见问题FAQ](#17-常见问题faq)

---

## 1. 系统要求

### 1.1 硬件要求

#### 最小配置（开发/测试环境）
- **CPU**: 4核
- **内存**: 16GB
- **存储**: 100GB SSD
- **网络**: 100Mbps

#### 推荐配置（生产环境）
- **CPU**: 16核+
- **内存**: 64GB+
- **存储**: 500GB+ NVMe SSD
- **网络**: 1Gbps+

#### 大规模部署（企业级）
- **CPU**: 32核+ (分布式)
- **内存**: 128GB+ (分布式)
- **存储**: 2TB+ NVMe SSD (分布式)
- **网络**: 10Gbps+ (分布式)

---

### 1.2 软件要求

| 组件 | 版本要求 | 说明 |
|------|---------|------|
| **操作系统** | Ubuntu 22.04 LTS / CentOS 8+ | 推荐Ubuntu |
| **Docker** | 24.0+ | 容器运行时 |
| **Docker Compose** | 2.20+ | 编排工具 |
| **Kubernetes** | 1.28+ | 生产环境 |
| **Python** | 3.9+ | API服务 |
| **PostgreSQL** | 15+ | 关系数据库 |
| **Milvus** | 2.3+ | 向量数据库 |
| **Redis** | 7.0+ | 缓存 |
| **MinIO** | RELEASE.2024-01+ | 对象存储 |

---

### 1.3 依赖服务

| 服务 | 用途 | 是否必需 |
|------|------|---------|
| **LiteLLM** | LLM统一代理 | ✅ 必需 |
| **Prometheus** | 监控指标收集 | ⚠️ 强烈推荐 |
| **Grafana** | 监控可视化 | ⚠️ 强烈推荐 |
| **Sentry** | 错误追踪 | ⚠️ 推荐 |
| **Nginx** | 反向代理 | ⚠️ 生产必需 |
| **Certbot** | SSL证书 | ⚠️ 生产必需 |

---

## 2. 架构概览

### 2.1 系统架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        负载均衡层                                 │
│                    Nginx / ALB / Ingress                        │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────┴────────────────────────────────────┐
│                        应用层                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  API Service │  │  API Service │  │  API Service │          │
│  │  (Python)    │  │  (Python)    │  │  (Python)    │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
└─────────┼──────────────────┼──────────────────┼─────────────────┘
          │                  │                  │
┌─────────┴──────────────────┴──────────────────┴─────────────────┐
│                        数据层                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ PostgreSQL   │  │   Milvus     │  │    Redis     │          │
│  │ (关系数据)    │  │ (向量数据)    │  │   (缓存)     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐                            │
│  │    MinIO     │  │   LiteLLM    │                            │
│  │ (对象存储)    │  │  (LLM代理)   │                            │
│  └──────────────┘  └──────────────┘                            │
└──────────────────────────────────────────────────────────────────┘
          │                  │
┌─────────┴──────────────────┴─────────────────────────────────────┐
│                        监控层                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  Prometheus  │  │   Grafana    │  │    Sentry    │          │
│  │ (指标收集)    │  │  (可视化)    │  │ (错误追踪)    │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└──────────────────────────────────────────────────────────────────┘
```

---

### 2.2 数据流图

```
用户请求
    ↓
Nginx (HTTPS)
    ↓
API Service
    ↓
┌───┴───────────────────────────────┐
│                                   │
│  1. 领域路由 (Adapter)            │
│     ↓                             │
│  2. AI目标生成                    │
│     ↓                             │
│  3. 内容发现 (Milvus + 多源)      │
│     ↓                             │
│  4. 叙述生成 (LiteLLM)            │
│     ↓                             │
│  5. 结果存储 (PostgreSQL + MinIO) │
│                                   │
└───────────────────────────────────┘
    ↓
Redis缓存 (加速响应)
    ↓
返回给用户
```

---

## 3. 快速部署（Docker Compose）

### 3.1 准备工作

```bash
# 克隆仓库
git clone https://github.com/metaworkflow/platform.git
cd platform

# 创建数据目录
mkdir -p data/{postgres,milvus,redis,minio,logs}

# 复制环境配置
cp .env.example .env
```

---

### 3.2 配置环境变量

编辑 `.env` 文件：

```bash
# ==================== 基础配置 ====================
ENVIRONMENT=production
DEBUG=false
SECRET_KEY=your-secret-key-here-change-this
API_HOST=0.0.0.0
API_PORT=8000

# ==================== 数据库配置 ====================
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DB=metaworkflow
POSTGRES_USER=metaworkflow
POSTGRES_PASSWORD=change-this-password

# ==================== Milvus配置 ====================
MILVUS_HOST=milvus-standalone
MILVUS_PORT=19530
MILVUS_COLLECTION_NAME=content_embeddings

# ==================== Redis配置 ====================
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=change-this-password
REDIS_DB=0

# ==================== MinIO配置 ====================
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=change-this-password
MINIO_SECURE=false
MINIO_BUCKET_CONTENT=content-assets
MINIO_BUCKET_GENERATED=generated-content

# ==================== LLM配置 ====================
LITELLM_API_BASE=http://litellm:4000
OPENAI_API_KEY=sk-your-openai-key-here
# 或使用其他LLM服务
# AZURE_API_KEY=your-azure-key
# ANTHROPIC_API_KEY=your-claude-key

# ==================== 监控配置 ====================
PROMETHEUS_ENABLED=true
GRAFANA_ADMIN_PASSWORD=change-this-password
SENTRY_DSN=https://your-sentry-dsn@sentry.io/project-id

# ==================== 日志配置 ====================
LOG_LEVEL=INFO
LOG_FORMAT=json
```

---

### 3.3 Docker Compose配置

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  # ==================== PostgreSQL ====================
  postgres:
    image: postgres:15-alpine
    container_name: metaworkflow-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
      - ./data/postgres:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    networks:
      - metaworkflow
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  # ==================== Milvus Standalone ====================
  etcd:
    image: quay.io/coreos/etcd:v3.5.5
    container_name: metaworkflow-etcd
    restart: unless-stopped
    environment:
      - ETCD_AUTO_COMPACTION_MODE=revision
      - ETCD_AUTO_COMPACTION_RETENTION=1000
      - ETCD_QUOTA_BACKEND_BYTES=4294967296
      - ETCD_SNAPSHOT_COUNT=50000
    volumes:
      - ./data/milvus/etcd:/etcd
    command: etcd -advertise-client-urls=http://127.0.0.1:2379 -listen-client-urls http://0.0.0.0:2379 --data-dir /etcd
    networks:
      - metaworkflow

  minio-milvus:
    image: minio/minio:RELEASE.2024-01-01T16-36-33Z
    container_name: metaworkflow-minio-milvus
    restart: unless-stopped
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - ./data/milvus/minio:/minio_data
    command: minio server /minio_data --console-address ":9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3
    networks:
      - metaworkflow

  milvus-standalone:
    image: milvusdb/milvus:v2.3.3
    container_name: metaworkflow-milvus
    restart: unless-stopped
    depends_on:
      - etcd
      - minio-milvus
    environment:
      ETCD_ENDPOINTS: etcd:2379
      MINIO_ADDRESS: minio-milvus:9000
    volumes:
      - ./data/milvus/data:/var/lib/milvus
      - ./config/milvus.yaml:/milvus/configs/milvus.yaml
    ports:
      - "19530:19530"
      - "9091:9091"
    networks:
      - metaworkflow
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9091/healthz"]
      interval: 30s
      start_period: 90s
      timeout: 20s
      retries: 3

  # ==================== Redis ====================
  redis:
    image: redis:7-alpine
    container_name: metaworkflow-redis
    restart: unless-stopped
    command: >
      redis-server
      --requirepass ${REDIS_PASSWORD}
      --maxmemory 2gb
      --maxmemory-policy allkeys-lru
      --save 60 1000
      --appendonly yes
    volumes:
      - ./data/redis:/data
    ports:
      - "6379:6379"
    networks:
      - metaworkflow
    healthcheck:
      test: ["CMD", "redis-cli", "--raw", "incr", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5

  # ==================== MinIO ====================
  minio:
    image: minio/minio:RELEASE.2024-01-01T16-36-33Z
    container_name: metaworkflow-minio
    restart: unless-stopped
    environment:
      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY}
      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY}
    volumes:
      - ./data/minio:/data
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
    networks:
      - metaworkflow
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  # MinIO初始化（创建buckets）
  minio-init:
    image: minio/mc:latest
    container_name: metaworkflow-minio-init
    depends_on:
      - minio
    entrypoint: >
      /bin/sh -c "
      sleep 10;
      mc alias set myminio http://minio:9000 ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY};
      mc mb myminio/content-assets --ignore-existing;
      mc mb myminio/user-uploads --ignore-existing;
      mc mb myminio/generated-content --ignore-existing;
      mc mb myminio/backups --ignore-existing;
      mc anonymous set download myminio/content-assets;
      echo 'MinIO buckets created successfully';
      "
    networks:
      - metaworkflow

  # ==================== LiteLLM Proxy ====================
  litellm:
    image: ghcr.io/berriai/litellm:main-latest
    container_name: metaworkflow-litellm
    restart: unless-stopped
    environment:
      OPENAI_API_KEY: ${OPENAI_API_KEY}
      LITELLM_MASTER_KEY: ${SECRET_KEY}
      DATABASE_URL: postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}
    volumes:
      - ./config/litellm-config.yaml:/app/config.yaml
    command: --config /app/config.yaml --port 4000
    ports:
      - "4000:4000"
    networks:
      - metaworkflow
    depends_on:
      - postgres
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:4000/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # ==================== API Service ====================
  api:
    build:
      context: .
      dockerfile: Dockerfile
    image: metaworkflow/api:v2.0
    container_name: metaworkflow-api
    restart: unless-stopped
    depends_on:
      - postgres
      - milvus-standalone
      - redis
      - minio
      - litellm
    environment:
      # 从.env文件继承所有环境变量
      - ENVIRONMENT
      - DEBUG
      - SECRET_KEY
      - POSTGRES_HOST
      - POSTGRES_PORT
      - POSTGRES_DB
      - POSTGRES_USER
      - POSTGRES_PASSWORD
      - MILVUS_HOST
      - MILVUS_PORT
      - REDIS_HOST
      - REDIS_PORT
      - REDIS_PASSWORD
      - MINIO_ENDPOINT
      - MINIO_ACCESS_KEY
      - MINIO_SECRET_KEY
      - LITELLM_API_BASE
      - SENTRY_DSN
    volumes:
      - ./src:/app/src
      - ./config:/app/config
      - ./data/logs:/app/logs
    ports:
      - "8000:8000"
    networks:
      - metaworkflow
    command: >
      uvicorn src.main:app
      --host 0.0.0.0
      --port 8000
      --workers 4
      --log-level info
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # ==================== Prometheus ====================
  prometheus:
    image: prom/prometheus:v2.48.0
    container_name: metaworkflow-prometheus
    restart: unless-stopped
    volumes:
      - ./config/prometheus.yml:/etc/prometheus/prometheus.yml
      - ./data/prometheus:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=30d'
      - '--web.enable-lifecycle'
    ports:
      - "9090:9090"
    networks:
      - metaworkflow
    depends_on:
      - api

  # ==================== Grafana ====================
  grafana:
    image: grafana/grafana:10.2.2
    container_name: metaworkflow-grafana
    restart: unless-stopped
    environment:
      GF_SECURITY_ADMIN_PASSWORD: ${GRAFANA_ADMIN_PASSWORD}
      GF_INSTALL_PLUGINS: grafana-piechart-panel
      GF_SERVER_ROOT_URL: http://localhost:3000
    volumes:
      - ./config/grafana/dashboards:/etc/grafana/provisioning/dashboards
      - ./config/grafana/datasources:/etc/grafana/provisioning/datasources
      - ./data/grafana:/var/lib/grafana
    ports:
      - "3000:3000"
    networks:
      - metaworkflow
    depends_on:
      - prometheus

  # ==================== Nginx ====================
  nginx:
    image: nginx:1.25-alpine
    container_name: metaworkflow-nginx
    restart: unless-stopped
    volumes:
      - ./config/nginx.conf:/etc/nginx/nginx.conf
      - ./config/ssl:/etc/nginx/ssl
    ports:
      - "80:80"
      - "443:443"
    networks:
      - metaworkflow
    depends_on:
      - api
      - grafana

networks:
  metaworkflow:
    driver: bridge

volumes:
  postgres_data:
  milvus_data:
  redis_data:
  minio_data:
  prometheus_data:
  grafana_data:
```

---

### 3.4 启动服务

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f api

# 停止服务
docker-compose down

# 停止并清除数据（谨慎使用）
docker-compose down -v
```

---

### 3.5 验证部署

```bash
# 1. 检查PostgreSQL
docker-compose exec postgres psql -U metaworkflow -d metaworkflow -c "SELECT version();"

# 2. 检查Milvus
curl http://localhost:9091/healthz

# 3. 检查Redis
docker-compose exec redis redis-cli -a your-password ping

# 4. 检查MinIO
curl http://localhost:9000/minio/health/live

# 5. 检查API服务
curl http://localhost:8000/health

# 6. 检查Prometheus
curl http://localhost:9090/-/healthy

# 7. 访问Grafana
# 浏览器打开: http://localhost:3000
# 默认用户名: admin
# 密码: 在.env中设置的GRAFANA_ADMIN_PASSWORD
```

---

### 3.6 初始化数据

```bash
# 运行数据库迁移
docker-compose exec api python scripts/migrate.py

# 创建Milvus集合
docker-compose exec api python scripts/init_milvus.py

# 加载示例适配器
docker-compose exec api python scripts/load_adapters.py

# 验证初始化
curl -X GET http://localhost:8000/v2/adapters
```

---

## 4. Kubernetes部署

### 4.1 准备Kubernetes集群

```bash
# 验证集群连接
kubectl cluster-info

# 创建命名空间
kubectl create namespace metaworkflow

# 设置默认命名空间
kubectl config set-context --current --namespace=metaworkflow
```

---

### 4.2 创建ConfigMap

`k8s/configmap.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: metaworkflow-config
  namespace: metaworkflow
data:
  ENVIRONMENT: "production"
  DEBUG: "false"
  API_HOST: "0.0.0.0"
  API_PORT: "8000"
  
  POSTGRES_HOST: "postgres-service"
  POSTGRES_PORT: "5432"
  POSTGRES_DB: "metaworkflow"
  
  MILVUS_HOST: "milvus-service"
  MILVUS_PORT: "19530"
  MILVUS_COLLECTION_NAME: "content_embeddings"
  
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
  REDIS_DB: "0"
  
  MINIO_ENDPOINT: "minio-service:9000"
  MINIO_SECURE: "false"
  MINIO_BUCKET_CONTENT: "content-assets"
  MINIO_BUCKET_GENERATED: "generated-content"
  
  LITELLM_API_BASE: "http://litellm-service:4000"
  
  LOG_LEVEL: "INFO"
  LOG_FORMAT: "json"
  
  PROMETHEUS_ENABLED: "true"
```

---

### 4.3 创建Secret

```bash
# 创建Secret（命令行方式）
kubectl create secret generic metaworkflow-secrets \
  --from-literal=SECRET_KEY='your-secret-key-here' \
  --from-literal=POSTGRES_PASSWORD='your-postgres-password' \
  --from-literal=REDIS_PASSWORD='your-redis-password' \
  --from-literal=MINIO_ACCESS_KEY='minioadmin' \
  --from-literal=MINIO_SECRET_KEY='your-minio-password' \
  --from-literal=OPENAI_API_KEY='sk-your-openai-key' \
  --from-literal=SENTRY_DSN='https://your-sentry-dsn@sentry.io/project' \
  --from-literal=GRAFANA_ADMIN_PASSWORD='your-grafana-password' \
  -n metaworkflow
```

或使用YAML文件 `k8s/secret.yaml`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: metaworkflow-secrets
  namespace: metaworkflow
type: Opaque
stringData:
  SECRET_KEY: "your-secret-key-here-base64-encoded"
  POSTGRES_PASSWORD: "your-postgres-password"
  REDIS_PASSWORD: "your-redis-password"
  MINIO_ACCESS_KEY: "minioadmin"
  MINIO_SECRET_KEY: "your-minio-password"
  OPENAI_API_KEY: "sk-your-openai-key"
  SENTRY_DSN: "https://your-sentry-dsn@sentry.io/project"
  GRAFANA_ADMIN_PASSWORD: "your-grafana-password"
```

```bash
kubectl apply -f k8s/secret.yaml
```

---

### 4.4 持久化存储（PVC）

`k8s/pvc.yaml`:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: metaworkflow
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 50Gi
  storageClassName: standard  # 根据云提供商调整

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: milvus-pvc
  namespace: metaworkflow
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 100Gi
  storageClassName: standard

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: redis-pvc
  namespace: metaworkflow
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
  storageClassName: standard

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: minio-pvc
  namespace: metaworkflow
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 200Gi
  storageClassName: standard
```

```bash
kubectl apply -f k8s/pvc.yaml
```

---

### 4.5 PostgreSQL部署

`k8s/postgres-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: metaworkflow
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        ports:
        - containerPort: 5432
          name: postgres
        env:
        - name: POSTGRES_DB
          valueFrom:
            configMapKeyRef:
              name: metaworkflow-config
              key: POSTGRES_DB
        - name: POSTGRES_USER
          value: metaworkflow
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: POSTGRES_PASSWORD
        - name: PGDATA
          value: /var/lib/postgresql/data/pgdata
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "2Gi"
            cpu: "1000m"
          limits:
            memory: "4Gi"
            cpu: "2000m"
        livenessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - metaworkflow
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - metaworkflow
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: postgres-service
  namespace: metaworkflow
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
    protocol: TCP
  type: ClusterIP
```

```bash
kubectl apply -f k8s/postgres-deployment.yaml
```

---

### 4.6 Redis部署

`k8s/redis-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  namespace: metaworkflow
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:7-alpine
        ports:
        - containerPort: 6379
          name: redis
        command:
        - redis-server
        - --requirepass
        - $(REDIS_PASSWORD)
        - --maxmemory
        - 2gb
        - --maxmemory-policy
        - allkeys-lru
        - --save
        - "60 1000"
        - --appendonly
        - "yes"
        env:
        - name: REDIS_PASSWORD
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: REDIS_PASSWORD
        volumeMounts:
        - name: redis-storage
          mountPath: /data
        resources:
          requests:
            memory: "2Gi"
            cpu: "500m"
          limits:
            memory: "4Gi"
            cpu: "1000m"
        livenessProbe:
          exec:
            command:
            - redis-cli
            - ping
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          exec:
            command:
            - redis-cli
            - ping
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: redis-storage
        persistentVolumeClaim:
          claimName: redis-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: redis-service
  namespace: metaworkflow
spec:
  selector:
    app: redis
  ports:
  - port: 6379
    targetPort: 6379
    protocol: TCP
  type: ClusterIP
```

```bash
kubectl apply -f k8s/redis-deployment.yaml
```

---

### 4.7 MinIO部署

`k8s/minio-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: minio
  namespace: metaworkflow
spec:
  replicas: 1
  selector:
    matchLabels:
      app: minio
  template:
    metadata:
      labels:
        app: minio
    spec:
      containers:
      - name: minio
        image: minio/minio:RELEASE.2024-01-01T16-36-33Z
        ports:
        - containerPort: 9000
          name: api
        - containerPort: 9001
          name: console
        command:
        - minio
        - server
        - /data
        - --console-address
        - ":9001"
        env:
        - name: MINIO_ROOT_USER
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: MINIO_ACCESS_KEY
        - name: MINIO_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: MINIO_SECRET_KEY
        volumeMounts:
        - name: minio-storage
          mountPath: /data
        resources:
          requests:
            memory: "2Gi"
            cpu: "500m"
          limits:
            memory: "4Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /minio/health/live
            port: 9000
          initialDelaySeconds: 30
          periodSeconds: 20
        readinessProbe:
          httpGet:
            path: /minio/health/ready
            port: 9000
          initialDelaySeconds: 10
          periodSeconds: 10
      volumes:
      - name: minio-storage
        persistentVolumeClaim:
          claimName: minio-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: minio-service
  namespace: metaworkflow
spec:
  selector:
    app: minio
  ports:
  - port: 9000
    targetPort: 9000
    protocol: TCP
    name: api
  - port: 9001
    targetPort: 9001
    protocol: TCP
    name: console
  type: ClusterIP
```

```bash
kubectl apply -f k8s/minio-deployment.yaml
```

---

### 4.8 Milvus部署

`k8s/milvus-deployment.yaml`:

```yaml
# Etcd for Milvus
apiVersion: apps/v1
kind: Deployment
metadata:
  name: etcd
  namespace: metaworkflow
spec:
  replicas: 1
  selector:
    matchLabels:
      app: etcd
  template:
    metadata:
      labels:
        app: etcd
    spec:
      containers:
      - name: etcd
        image: quay.io/coreos/etcd:v3.5.5
        ports:
        - containerPort: 2379
          name: client
        - containerPort: 2380
          name: peer
        env:
        - name: ETCD_AUTO_COMPACTION_MODE
          value: "revision"
        - name: ETCD_AUTO_COMPACTION_RETENTION
          value: "1000"
        - name: ETCD_QUOTA_BACKEND_BYTES
          value: "4294967296"
        - name: ETCD_SNAPSHOT_COUNT
          value: "50000"
        command:
        - etcd
        - -advertise-client-urls=http://0.0.0.0:2379
        - -listen-client-urls=http://0.0.0.0:2379
        - --data-dir=/etcd
        volumeMounts:
        - name: etcd-storage
          mountPath: /etcd
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
      volumes:
      - name: etcd-storage
        emptyDir: {}

---
apiVersion: v1
kind: Service
metadata:
  name: etcd-service
  namespace: metaworkflow
spec:
  selector:
    app: etcd
  ports:
  - port: 2379
    targetPort: 2379
    protocol: TCP
  type: ClusterIP

---
# Milvus Standalone
apiVersion: apps/v1
kind: Deployment
metadata:
  name: milvus
  namespace: metaworkflow
spec:
  replicas: 1
  selector:
    matchLabels:
      app: milvus
  template:
    metadata:
      labels:
        app: milvus
    spec:
      containers:
      - name: milvus
        image: milvusdb/milvus:v2.3.3
        ports:
        - containerPort: 19530
          name: grpc
        - containerPort: 9091
          name: metrics
        env:
        - name: ETCD_ENDPOINTS
          value: "etcd-service:2379"
        - name: MINIO_ADDRESS
          value: "minio-service:9000"
        - name: MINIO_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: MINIO_ACCESS_KEY
        - name: MINIO_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: MINIO_SECRET_KEY
        volumeMounts:
        - name: milvus-storage
          mountPath: /var/lib/milvus
        resources:
          requests:
            memory: "4Gi"
            cpu: "2000m"
          limits:
            memory: "8Gi"
            cpu: "4000m"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 9091
          initialDelaySeconds: 90
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /healthz
            port: 9091
          initialDelaySeconds: 30
          periodSeconds: 10
      volumes:
      - name: milvus-storage
        persistentVolumeClaim:
          claimName: milvus-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: milvus-service
  namespace: metaworkflow
spec:
  selector:
    app: milvus
  ports:
  - port: 19530
    targetPort: 19530
    protocol: TCP
    name: grpc
  - port: 9091
    targetPort: 9091
    protocol: TCP
    name: metrics
  type: ClusterIP
```

```bash
kubectl apply -f k8s/milvus-deployment.yaml
```

---

### 4.9 API服务部署

`k8s/api-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  namespace: metaworkflow
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
      - name: api
        image: metaworkflow/api:v2.0
        ports:
        - containerPort: 8000
          name: http
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: DEBUG
          value: "false"
        - name: SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: SECRET_KEY
        - name: POSTGRES_HOST
          value: "postgres-service"
        - name: POSTGRES_PORT
          value: "5432"
        - name: POSTGRES_DB
          value: "metaworkflow"
        - name: POSTGRES_USER
          value: "metaworkflow"
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: POSTGRES_PASSWORD
        - name: MILVUS_HOST
          value: "milvus-service"
        - name: MILVUS_PORT
          value: "19530"
        - name: REDIS_HOST
          value: "redis-service"
        - name: REDIS_PORT
          value: "6379"
        - name: REDIS_PASSWORD
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: REDIS_PASSWORD
        - name: MINIO_ENDPOINT
          value: "minio-service:9000"
        - name: MINIO_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: MINIO_ACCESS_KEY
        - name: MINIO_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: MINIO_SECRET_KEY
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: OPENAI_API_KEY
        - name: SENTRY_DSN
          valueFrom:
            secretKeyRef:
              name: metaworkflow-secrets
              key: SENTRY_DSN
        resources:
          requests:
            memory: "2Gi"
            cpu: "1000m"
          limits:
            memory: "4Gi"
            cpu: "2000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 10
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: api-service
  namespace: metaworkflow
spec:
  selector:
    app: api
  ports:
  - port: 8000
    targetPort: 8000
    protocol: TCP
  type: ClusterIP

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api-ingress
  namespace: metaworkflow
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - api.metaworkflow.ai
    secretName: api-tls
  rules:
  - host: api.metaworkflow.ai
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 8000
```

```bash
kubectl apply -f k8s/api-deployment.yaml
```

---

### 4.10 验证K8s部署

```bash
# 查看所有Pods状态
kubectl get pods -n metaworkflow

# 查看所有Services
kubectl get svc -n metaworkflow

# 查看Ingress
kubectl get ingress -n metaworkflow

# 查看Pod日志
kubectl logs -f deployment/api -n metaworkflow

# 执行健康检查
kubectl exec -it deployment/api -n metaworkflow -- curl http://localhost:8000/health

# 查看资源使用情况
kubectl top pods -n metaworkflow
```

---

## 5. 数据库配置

### 5.1 PostgreSQL初始化脚本

创建 `scripts/init-db.sql`:

```sql
-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- 创建数据库（如果不存在）
SELECT 'CREATE DATABASE metaworkflow'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'metaworkflow')\gexec

-- 连接到数据库
\c metaworkflow

-- ==================== 用户表 ====================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    is_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- ==================== 领域适配器表 ====================
CREATE TABLE IF NOT EXISTS domain_adapters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    domain_type VARCHAR(50) NOT NULL,
    description TEXT,
    version VARCHAR(20) NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_adapters_domain_type ON domain_adapters(domain_type);
CREATE INDEX idx_adapters_active ON domain_adapters(is_active) WHERE is_active = true;

-- ==================== 工作流模板表 ====================
CREATE TABLE IF NOT EXISTS workflow_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    domain_type VARCHAR(50) NOT NULL,
    adapter_id UUID REFERENCES domain_adapters(id) ON DELETE CASCADE,
    config JSONB NOT NULL DEFAULT '{}',
    steps JSONB NOT NULL DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_templates_domain ON workflow_templates(domain_type);
CREATE INDEX idx_templates_adapter ON workflow_templates(adapter_id);

-- ==================== 内容源表 ====================
CREATE TABLE IF NOT EXISTS content_sources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    api_endpoint VARCHAR(500),
    api_key_encrypted VARCHAR(500),
    config JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0,
    rate_limit INTEGER DEFAULT 100,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sources_type ON content_sources(source_type);
CREATE INDEX idx_sources_active ON content_sources(is_active) WHERE is_active = true;

-- ==================== 生成任务表 ====================
CREATE TABLE IF NOT EXISTS generation_tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    adapter_id UUID REFERENCES domain_adapters(id),
    template_id UUID REFERENCES workflow_templates(id),
    input_data JSONB NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    progress INTEGER DEFAULT 0,
    result JSONB,
    error_message TEXT,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_user ON generation_tasks(user_id);
CREATE INDEX idx_tasks_status ON generation_tasks(status);
CREATE INDEX idx_tasks_created ON generation_tasks(created_at DESC);

-- ==================== 学习目标表 ====================
CREATE TABLE IF NOT EXISTS learning_objectives (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID REFERENCES generation_tasks(id) ON DELETE CASCADE,
    objective_text TEXT NOT NULL,
    bloom_level VARCHAR(50),
    cognitive_domain VARCHAR(50),
    difficulty_level VARCHAR(20),
    estimated_time INTEGER,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_objectives_task ON learning_objectives(task_id);
CREATE INDEX idx_objectives_bloom ON learning_objectives(bloom_level);

-- ==================== 内容片段表 ====================
CREATE TABLE IF NOT EXISTS content_fragments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID REFERENCES generation_tasks(id) ON DELETE CASCADE,
    source_id UUID REFERENCES content_sources(id),
    content_type VARCHAR(50) NOT NULL,
    title VARCHAR(500),
    content TEXT,
    metadata JSONB DEFAULT '{}',
    quality_score FLOAT,
    embedding_id VARCHAR(200),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fragments_task ON content_fragments(task_id);
CREATE INDEX idx_fragments_type ON content_fragments(content_type);
CREATE INDEX idx_fragments_score ON content_fragments(quality_score DESC NULLS LAST);

-- ==================== 叙述内容表 ====================
CREATE TABLE IF NOT EXISTS narrative_contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID REFERENCES generation_tasks(id) ON DELETE CASCADE,
    narrative_type VARCHAR(50) NOT NULL,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    story_arc VARCHAR(50),
    tone VARCHAR(50),
    metadata JSONB DEFAULT '{}',
    audio_url VARCHAR(1000),
    quality_metrics JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_narratives_task ON narrative_contents(task_id);
CREATE INDEX idx_narratives_type ON narrative_contents(narrative_type);

-- ==================== 审计日志表 ====================
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id UUID,
    changes JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_action ON audit_logs(action);
CREATE INDEX idx_audit_created ON audit_logs(created_at DESC);

-- ==================== 更新触发器 ====================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_adapters_updated_at BEFORE UPDATE ON domain_adapters
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_templates_updated_at BEFORE UPDATE ON workflow_templates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sources_updated_at BEFORE UPDATE ON content_sources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ==================== 初始数据 ====================
INSERT INTO users (username, email, password_hash, is_admin) VALUES
    ('admin', 'admin@metaworkflow.ai', '$2b$12$placeholder_hash', true)
ON CONFLICT (username) DO NOTHING;

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO metaworkflow;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO metaworkflow;
```

---

### 5.2 PostgreSQL性能优化

编辑 `config/postgresql.conf`:

```conf
# ==================== 连接设置 ====================
max_connections = 200
superuser_reserved_connections = 3

# ==================== 内存设置 ====================
shared_buffers = 8GB                    # 25% of total RAM
effective_cache_size = 24GB             # 75% of total RAM
maintenance_work_mem = 2GB
work_mem = 64MB
wal_buffers = 16MB

# ==================== 查询规划器 ====================
random_page_cost = 1.1                  # SSD优化
effective_io_concurrency = 200          # SSD优化
default_statistics_target = 100

# ==================== WAL设置 ====================
wal_level = replica
max_wal_size = 4GB
min_wal_size = 1GB
checkpoint_completion_target = 0.9
archive_mode = on
archive_command = 'test ! -f /var/lib/postgresql/archive/%f && cp %p /var/lib/postgresql/archive/%f'

# ==================== 复制设置 ====================
max_wal_senders = 10
max_replication_slots = 10
hot_standby = on

# ==================== 日志设置 ====================
logging_collector = on
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
log_checkpoints = on
log_connections = on
log_disconnections = on
log_duration = off
log_lock_waits = on
log_min_duration_statement = 1000       # 记录>1秒的查询

# ==================== 自动清理 ====================
autovacuum = on
autovacuum_max_workers = 4
autovacuum_naptime = 10s
```

---

### 5.3 数据库备份脚本

创建 `scripts/backup-postgres.sh`:

```bash
#!/bin/bash

# PostgreSQL备份脚本
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/postgres"
DB_NAME="metaworkflow"
RETENTION_DAYS=7

# 创建备份目录
mkdir -p ${BACKUP_DIR}

# 执行备份
echo "[$(date)] Starting PostgreSQL backup..."
docker-compose exec -T postgres pg_dump -U metaworkflow ${DB_NAME} | gzip > ${BACKUP_DIR}/backup_${DATE}.sql.gz

if [ $? -eq 0 ]; then
    echo "[$(date)] Backup completed: backup_${DATE}.sql.gz"
    
    # 删除旧备份
    find ${BACKUP_DIR} -name "backup_*.sql.gz" -mtime +${RETENTION_DAYS} -delete
    echo "[$(date)] Old backups cleaned (retention: ${RETENTION_DAYS} days)"
else
    echo "[$(date)] Backup failed!"
    exit 1
fi
```

---

## 6. 向量数据库配置（Milvus）

### 6.1 Milvus配置文件

创建 `config/milvus.yaml`:

```yaml
# Milvus配置文件 V2.3

etcd:
  endpoints:
    - etcd:2379
  rootPath: by-dev
  metaSubPath: meta
  kvSubPath: kv

minio:
  address: minio-milvus:9000
  accessKeyID: minioadmin
  secretAccessKey: minioadmin
  useSSL: false
  bucketName: milvus-bucket
  rootPath: files

storage:
  # 内存限制
  insertBufferSize: 4294967296  # 4GB
  
common:
  # 时区设置
  timezone: UTC
  
  # 安全设置
  security:
    authorizationEnabled: false

# 日志配置
log:
  level: info
  file:
    rootPath: /var/lib/milvus/logs
    maxSize: 300
    maxAge: 10
    maxBackups: 20

# 向量索引配置
indexCoord:
  minSegmentNumRowsToEnableIndex: 1024

# 查询节点配置
queryNode:
  cacheSize: 8589934592  # 8GB
  
  # 搜索优化
  gracefulStopTimeout: 30

# 数据节点配置
dataNode:
  # Flush配置
  flushInsertBufferSize: 16777216
  flushDeleteBufferBytes: 16777216

# Root协调器配置
rootCoord:
  dmlChannelNum: 16
  maxPartitionNum: 4096
  minSegmentSizeToEnableIndex: 1024

# 代理配置
proxy:
  # 时间同步
  timeTickInterval: 200
  
  # 搜索结果限制
  maxOutputSize: 104857600  # 100MB
  maxQueryResultSize: 104857600
```

---

### 6.2 Milvus集合初始化

创建 `scripts/init_milvus.py`:

```python
"""Milvus集合初始化脚本"""
from pymilvus import (
    connections,
    Collection,
    CollectionSchema,
    FieldSchema,
    DataType,
    utility
)
import os

# 连接配置
MILVUS_HOST = os.getenv("MILVUS_HOST", "localhost")
MILVUS_PORT = os.getenv("MILVUS_PORT", "19530")

def create_content_embeddings_collection():
    """创建内容嵌入向量集合"""
    
    # 定义字段
    fields = [
        FieldSchema(name="id", dtype=DataType.VARCHAR, is_primary=True, max_length=200),
        FieldSchema(name="content_id", dtype=DataType.VARCHAR, max_length=200),
        FieldSchema(name="source_type", dtype=DataType.VARCHAR, max_length=50),
        FieldSchema(name="domain_type", dtype=DataType.VARCHAR, max_length=50),
        FieldSchema(name="title", dtype=DataType.VARCHAR, max_length=500),
        FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=1536),  # OpenAI embedding
        FieldSchema(name="created_at", dtype=DataType.INT64),
    ]
    
    # 创建Schema
    schema = CollectionSchema(
        fields=fields,
        description="Content embeddings for semantic search",
        enable_dynamic_field=True
    )
    
    # 创建Collection
    collection = Collection(
        name="content_embeddings",
        schema=schema,
        using='default'
    )
    
    # 创建索引
    index_params = {
        "metric_type": "COSINE",
        "index_type": "IVF_FLAT",
        "params": {"nlist": 1024}
    }
    
    collection.create_index(
        field_name="embedding",
        index_params=index_params
    )
    
    print("✅ Collection 'content_embeddings' created successfully")
    return collection

def main():
    """主函数"""
    print("Connecting to Milvus...")
    connections.connect(
        alias="default",
        host=MILVUS_HOST,
        port=MILVUS_PORT
    )
    
    # 检查集合是否存在
    if utility.has_collection("content_embeddings"):
        print("⚠️  Collection 'content_embeddings' already exists")
        collection = Collection("content_embeddings")
    else:
        collection = create_content_embeddings_collection()
    
    # 加载集合到内存
    collection.load()
    print(f"📊 Collection stats: {collection.num_entities} entities")
    
    connections.disconnect("default")
    print("✅ Milvus initialization completed")

if __name__ == "__main__":
    main()
```

运行初始化：

```bash
python scripts/init_milvus.py
```

---

## 7. 缓存配置（Redis）

### 7.1 Redis配置优化

创建 `config/redis.conf`:

```conf
# ==================== 网络配置 ====================
bind 0.0.0.0
protected-mode yes
port 6379
tcp-backlog 511
timeout 0
tcp-keepalive 300

# ==================== 内存配置 ====================
maxmemory 4gb
maxmemory-policy allkeys-lru
maxmemory-samples 5

# ==================== 持久化配置 ====================
# RDB快照
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data

# AOF日志
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb

# ==================== 日志配置 ====================
loglevel notice
logfile /data/redis.log

# ==================== 客户端配置 ====================
maxclients 10000

# ==================== 慢查询日志 ====================
slowlog-log-slower-than 10000
slowlog-max-len 128

# ==================== 高级配置 ====================
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-size -2
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64

# ==================== 事件通知 ====================
notify-keyspace-events ""

# ==================== 高级配置 ====================
activerehashing yes
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit replica 256mb 64mb 60
client-output-buffer-limit pubsub 32mb 8mb 60
hz 10
```

---

### 7.2 Redis缓存策略

在应用代码中实现缓存策略（`src/cache/redis_manager.py`）:

```python
"""Redis缓存管理器"""
import redis
import json
from typing import Any, Optional
from datetime import timedelta

class RedisManager:
    """Redis缓存管理器"""
    
    def __init__(self, host: str, port: int, password: str, db: int = 0):
        self.client = redis.Redis(
            host=host,
            port=port,
            password=password,
            db=db,
            decode_responses=True,
            socket_connect_timeout=5,
            socket_keepalive=True,
            retry_on_timeout=True
        )
    
    # ==================== 内容发现缓存 ====================
    def cache_search_results(
        self,
        query: str,
        domain: str,
        results: list,
        ttl: int = 3600
    ):
        """缓存搜索结果（1小时）"""
        key = f"search:{domain}:{query}"
        self.client.setex(
            key,
            ttl,
            json.dumps(results, ensure_ascii=False)
        )
    
    def get_cached_search(self, query: str, domain: str) -> Optional[list]:
        """获取缓存的搜索结果"""
        key = f"search:{domain}:{query}"
        cached = self.client.get(key)
        return json.loads(cached) if cached else None
    
    # ==================== 向量嵌入缓存 ====================
    def cache_embedding(self, text: str, embedding: list, ttl: int = 86400):
        """缓存文本嵌入（24小时）"""
        key = f"embedding:{hash(text)}"
        self.client.setex(key, ttl, json.dumps(embedding))
    
    def get_cached_embedding(self, text: str) -> Optional[list]:
        """获取缓存的嵌入"""
        key = f"embedding:{hash(text)}"
        cached = self.client.get(key)
        return json.loads(cached) if cached else None
    
    # ==================== LLM响应缓存 ====================
    def cache_llm_response(
        self,
        prompt: str,
        response: str,
        model: str,
        ttl: int = 7200
    ):
        """缓存LLM响应（2小时）"""
        key = f"llm:{model}:{hash(prompt)}"
        self.client.setex(key, ttl, response)
    
    def get_cached_llm_response(self, prompt: str, model: str) -> Optional[str]:
        """获取缓存的LLM响应"""
        key = f"llm:{model}:{hash(prompt)}"
        return self.client.get(key)
    
    # ==================== 会话管理 ====================
    def set_session(self, session_id: str, data: dict, ttl: int = 3600):
        """设置会话数据"""
        key = f"session:{session_id}"
        self.client.setex(key, ttl, json.dumps(data))
    
    def get_session(self, session_id: str) -> Optional[dict]:
        """获取会话数据"""
        key = f"session:{session_id}"
        cached = self.client.get(key)
        return json.loads(cached) if cached else None
    
    # ==================== 速率限制 ====================
    def check_rate_limit(
        self,
        user_id: str,
        endpoint: str,
        max_requests: int = 100,
        window: int = 3600
    ) -> bool:
        """检查速率限制"""
        key = f"ratelimit:{endpoint}:{user_id}"
        current = self.client.incr(key)
        
        if current == 1:
            self.client.expire(key, window)
        
        return current <= max_requests
    
    # ==================== 健康检查 ====================
    def health_check(self) -> bool:
        """Redis健康检查"""
        try:
            return self.client.ping()
        except Exception:
            return False
```

---

## 8. 对象存储配置（MinIO）

### 8.1 MinIO存储桶策略

创建 `scripts/init-minio.sh`:

```bash
#!/bin/bash

# MinIO初始化脚本
MINIO_HOST="minio:9000"
MINIO_ALIAS="myminio"

# 等待MinIO启动
sleep 10

# 配置alias
mc alias set ${MINIO_ALIAS} http://${MINIO_HOST} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY}

# 创建存储桶
echo "Creating buckets..."
mc mb ${MINIO_ALIAS}/content-assets --ignore-existing
mc mb ${MINIO_ALIAS}/user-uploads --ignore-existing
mc mb ${MINIO_ALIAS}/generated-content --ignore-existing
mc mb ${MINIO_ALIAS}/backups --ignore-existing
mc mb ${MINIO_ALIAS}/temp --ignore-existing

# 设置public访问（仅content-assets）
echo "Setting bucket policies..."
mc anonymous set download ${MINIO_ALIAS}/content-assets

# 设置生命周期策略（temp bucket 7天自动删除）
cat > /tmp/lifecycle.json <<EOF
{
    "Rules": [
        {
            "ID": "DeleteTempFiles",
            "Status": "Enabled",
            "Filter": {
                "Prefix": ""
            },
            "Expiration": {
                "Days": 7
            }
        }
    ]
}
EOF

mc ilm import ${MINIO_ALIAS}/temp < /tmp/lifecycle.json

# 设置版本控制（backups bucket）
mc version enable ${MINIO_ALIAS}/backups

echo "✅ MinIO initialization completed"
```

---

### 8.2 MinIO客户端封装

创建 `src/storage/minio_client.py`:

```python
"""MinIO对象存储客户端"""
from minio import Minio
from minio.error import S3Error
from typing import Optional, BinaryIO
import os
from datetime import timedelta

class MinIOClient:
    """MinIO客户端封装"""
    
    def __init__(
        self,
        endpoint: str,
        access_key: str,
        secret_key: str,
        secure: bool = False
    ):
        self.client = Minio(
            endpoint,
            access_key=access_key,
            secret_key=secret_key,
            secure=secure
        )
        
        # 存储桶名称
        self.BUCKET_CONTENT = "content-assets"
        self.BUCKET_UPLOADS = "user-uploads"
        self.BUCKET_GENERATED = "generated-content"
        self.BUCKET_BACKUPS = "backups"
        self.BUCKET_TEMP = "temp"
    
    def upload_file(
        self,
        bucket: str,
        object_name: str,
        file_path: str,
        content_type: Optional[str] = None
    ) -> str:
        """上传文件"""
        try:
            self.client.fput_object(
                bucket,
                object_name,
                file_path,
                content_type=content_type
            )
            return f"/{bucket}/{object_name}"
        except S3Error as e:
            raise Exception(f"Upload failed: {e}")
    
    def upload_stream(
        self,
        bucket: str,
        object_name: str,
        data: BinaryIO,
        length: int,
        content_type: str = "application/octet-stream"
    ) -> str:
        """上传数据流"""
        try:
            self.client.put_object(
                bucket,
                object_name,
                data,
                length,
                content_type=content_type
            )
            return f"/{bucket}/{object_name}"
        except S3Error as e:
            raise Exception(f"Upload failed: {e}")
    
    def download_file(
        self,
        bucket: str,
        object_name: str,
        file_path: str
    ):
        """下载文件"""
        try:
            self.client.fget_object(bucket, object_name, file_path)
        except S3Error as e:
            raise Exception(f"Download failed: {e}")
    
    def get_presigned_url(
        self,
        bucket: str,
        object_name: str,
        expires: timedelta = timedelta(hours=1)
    ) -> str:
        """生成预签名URL"""
        try:
            return self.client.presigned_get_object(
                bucket,
                object_name,
                expires=expires
            )
        except S3Error as e:
            raise Exception(f"Failed to generate URL: {e}")
    
    def delete_file(self, bucket: str, object_name: str):
        """删除文件"""
        try:
            self.client.remove_object(bucket, object_name)
        except S3Error as e:
            raise Exception(f"Delete failed: {e}")
    
    def list_objects(self, bucket: str, prefix: str = "") -> list:
        """列出对象"""
        try:
            objects = self.client.list_objects(bucket, prefix=prefix)
            return [obj.object_name for obj in objects]
        except S3Error as e:
            raise Exception(f"List failed: {e}")
    
    def health_check(self) -> bool:
        """健康检查"""
        try:
            # 尝试列出存储桶
            list(self.client.list_buckets())
            return True
        except Exception:
            return False
```

---

## 9. 监控系统配置

### 9.1 Prometheus配置

创建 `config/prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'metaworkflow-prod'
    environment: 'production'

# 告警规则文件
rule_files:
  - 'alerts.yml'

# Alertmanager配置
alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - alertmanager:9093

# 抓取配置
scrape_configs:
  # Prometheus自身
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  # API服务
  - job_name: 'api'
    metrics_path: '/metrics'
    static_configs:
      - targets: ['api:8000']
    relabel_configs:
      - source_labels: [__address__]
        target_label: instance
        regex: '([^:]+)(?::\d+)?'
        replacement: '${1}'

  # PostgreSQL Exporter
  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  # Redis Exporter
  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']

  # Milvus Metrics
  - job_name: 'milvus'
    static_configs:
      - targets: ['milvus:9091']

  # MinIO Metrics
  - job_name: 'minio'
    metrics_path: '/minio/v2/metrics/cluster'
    static_configs:
      - targets: ['minio:9000']

  # Node Exporter（系统指标）
  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']

  # cAdvisor（容器指标）
  - job_name: 'cadvisor'
    static_configs:
      - targets: ['cadvisor:8080']
```

---

### 9.2 Prometheus告警规则

创建 `config/alerts.yml`:

```yaml
groups:
  # ==================== API服务告警 ====================
  - name: api_alerts
    interval: 30s
    rules:
      - alert: APIHighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "API错误率过高"
          description: "{{ $labels.instance }} 5分钟内错误率超过5%"

      - alert: APIHighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "API响应延迟过高"
          description: "{{ $labels.instance }} P95延迟超过2秒"

      - alert: APIDown
        expr: up{job="api"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "API服务宕机"
          description: "{{ $labels.instance }} 无法访问"

  # ==================== 数据库告警 ====================
  - name: database_alerts
    interval: 30s
    rules:
      - alert: PostgreSQLDown
        expr: pg_up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "PostgreSQL宕机"
          description: "PostgreSQL数据库无法连接"

      - alert: PostgreSQLHighConnections
        expr: pg_stat_database_numbackends / pg_settings_max_connections > 0.8
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "PostgreSQL连接数过高"
          description: "连接数使用率超过80%"

      - alert: PostgreSQLSlowQueries
        expr: rate(pg_stat_activity_max_tx_duration[5m]) > 60
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "PostgreSQL慢查询"
          description: "存在执行时间超过60秒的查询"

      - alert: RedisDown
        expr: redis_up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Redis宕机"
          description: "Redis无法连接"

      - alert: RedisHighMemory
        expr: redis_memory_used_bytes / redis_memory_max_bytes > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Redis内存使用率过高"
          description: "内存使用率超过90%"

  # ==================== Milvus告警 ====================
  - name: milvus_alerts
    interval: 30s
    rules:
      - alert: MilvusDown
        expr: up{job="milvus"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Milvus宕机"
          description: "Milvus向量数据库无法访问"

      - alert: MilvusHighSearchLatency
        expr: histogram_quantile(0.95, rate(milvus_search_latency_bucket[5m])) > 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Milvus搜索延迟过高"
          description: "P95搜索延迟超过1秒"

  # ==================== 系统资源告警 ====================
  - name: system_alerts
    interval: 30s
    rules:
      - alert: HighCPUUsage
        expr: 100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "CPU使用率过高"
          description: "{{ $labels.instance }} CPU使用率超过80%"

      - alert: HighMemoryUsage
        expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 85
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "内存使用率过高"
          description: "{{ $labels.instance }} 内存使用率超过85%"

      - alert: DiskSpaceLow
        expr: (node_filesystem_avail_bytes / node_filesystem_size_bytes) * 100 < 15
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "磁盘空间不足"
          description: "{{ $labels.instance }} 磁盘可用空间低于15%"
```

---

### 9.3 Grafana仪表板配置

创建 `config/grafana/dashboards/metaworkflow.json`:

```json
{
  "dashboard": {
    "title": "MetaWorkflow Platform Overview",
    "tags": ["metaworkflow", "v2.0"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "API请求速率",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 0, "y": 0}
      },
      {
        "id": 2,
        "title": "API响应时间 (P95)",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "{{endpoint}}"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 12, "y": 0}
      },
      {
        "id": 3,
        "title": "错误率",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])",
            "legendFormat": "{{status}}"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 0, "y": 8}
      },
      {
        "id": 4,
        "title": "数据库连接数",
        "type": "graph",
        "targets": [
          {
            "expr": "pg_stat_database_numbackends",
            "legendFormat": "PostgreSQL"
          },
          {
            "expr": "redis_connected_clients",
            "legendFormat": "Redis"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 12, "y": 8}
      },
      {
        "id": 5,
        "title": "Milvus向量搜索QPS",
        "type": "stat",
        "targets": [
          {
            "expr": "rate(milvus_search_total[1m])"
          }
        ],
        "gridPos": {"h": 4, "w": 6, "x": 0, "y": 16}
      },
      {
        "id": 6,
        "title": "Redis缓存命中率",
        "type": "gauge",
        "targets": [
          {
            "expr": "rate(redis_keyspace_hits_total[5m]) / (rate(redis_keyspace_hits_total[5m]) + rate(redis_keyspace_misses_total[5m]))"
          }
        ],
        "gridPos": {"h": 4, "w": 6, "x": 6, "y": 16}
      }
    ],
    "refresh": "10s",
    "schemaVersion": 36
  }
}
```

---

### 9.4 Sentry错误追踪配置

在应用中集成Sentry（`src/main.py`）:

```python
import sentry_sdk
from sentry_sdk.integrations.fastapi import FastApiIntegration
from sentry_sdk.integrations.sqlalchemy import SqlalchemyIntegration

# 初始化Sentry
sentry_sdk.init(
    dsn=os.getenv("SENTRY_DSN"),
    environment=os.getenv("ENVIRONMENT", "production"),
    release=f"metaworkflow@{__version__}",
    traces_sample_rate=0.1,  # 10%的性能追踪
    profiles_sample_rate=0.1,  # 10%的性能分析
    integrations=[
        FastApiIntegration(),
        SqlalchemyIntegration(),
    ],
    # 过滤敏感数据
    before_send=lambda event, hint: scrub_sensitive_data(event),
)

def scrub_sensitive_data(event):
    """清除敏感数据"""
    # 移除密码、token等敏感字段
    if 'request' in event:
        if 'data' in event['request']:
            sensitive_keys = ['password', 'token', 'api_key', 'secret']
            for key in sensitive_keys:
                if key in event['request']['data']:
                    event['request']['data'][key] = '[REDACTED]'
    return event
```

---

## 10. 日志管理

### 10.1 日志配置

创建 `config/logging.yaml`:

```yaml
version: 1
disable_existing_loggers: false

formatters:
  json:
    class: pythonjsonlogger.jsonlogger.JsonFormatter
    format: '%(asctime)s %(name)s %(levelname)s %(message)s'
  
  standard:
    format: '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    datefmt: '%Y-%m-%d %H:%M:%S'

handlers:
  console:
    class: logging.StreamHandler
    level: INFO
    formatter: json
    stream: ext://sys.stdout
  
  file:
    class: logging.handlers.RotatingFileHandler
    level: DEBUG
    formatter: json
    filename: /app/logs/metaworkflow.log
    maxBytes: 104857600  # 100MB
    backupCount: 10
  
  error_file:
    class: logging.handlers.RotatingFileHandler
    level: ERROR
    formatter: json
    filename: /app/logs/error.log
    maxBytes: 104857600
    backupCount: 10

loggers:
  metaworkflow:
    level: DEBUG
    handlers: [console, file, error_file]
    propagate: false
  
  uvicorn:
    level: INFO
    handlers: [console, file]
    propagate: false
  
  sqlalchemy.engine:
    level: WARNING
    handlers: [file]
    propagate: false

root:
  level: INFO
  handlers: [console, file]
```

---

### 10.2 应用日志实现

创建 `src/utils/logger.py`:

```python
"""统一日志管理"""
import logging
import logging.config
import yaml
from pathlib import Path
from typing import Optional

def setup_logging(config_path: Optional[str] = None):
    """配置日志系统"""
    if config_path and Path(config_path).exists():
        with open(config_path, 'r') as f:
            config = yaml.safe_load(f)
            logging.config.dictConfig(config)
    else:
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )

def get_logger(name: str) -> logging.Logger:
    """获取logger实例"""
    return logging.getLogger(name)

# 使用示例
logger = get_logger(__name__)
logger.info("Application started", extra={
    "environment": "production",
    "version": "2.0.0"
})
```

---

### 10.3 日志聚合（ELK Stack可选）

Docker Compose添加ELK（可选）:

```yaml
  # Elasticsearch
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - "ES_JAVA_OPTS=-Xms2g -Xmx2g"
    volumes:
      - ./data/elasticsearch:/usr/share/elasticsearch/data
    ports:
      - "9200:9200"
    networks:
      - metaworkflow

  # Logstash
  logstash:
    image: docker.elastic.co/logstash/logstash:8.11.0
    volumes:
      - ./config/logstash.conf:/usr/share/logstash/pipeline/logstash.conf
      - ./data/logs:/logs
    depends_on:
      - elasticsearch
    networks:
      - metaworkflow

  # Kibana
  kibana:
    image: docker.elastic.co/kibana/kibana:8.11.0
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch
    networks:
      - metaworkflow
```

---

## 11. 备份与恢复

### 11.1 自动化备份脚本

创建 `scripts/backup-all.sh`:

```bash
#!/bin/bash

# 全量备份脚本
set -e

BACKUP_ROOT="/backups"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

echo "========================================="
echo "MetaWorkflow Platform Backup"
echo "Date: $(date)"
echo "========================================="

# ==================== PostgreSQL备份 ====================
echo "[1/4] Backing up PostgreSQL..."
POSTGRES_BACKUP_DIR="${BACKUP_ROOT}/postgres"
mkdir -p ${POSTGRES_BACKUP_DIR}

docker-compose exec -T postgres pg_dump -U metaworkflow metaworkflow \
    | gzip > ${POSTGRES_BACKUP_DIR}/pg_backup_${DATE}.sql.gz

echo "✅ PostgreSQL backup completed"

# ==================== Milvus备份 ====================
echo "[2/4] Backing up Milvus..."
MILVUS_BACKUP_DIR="${BACKUP_ROOT}/milvus"
mkdir -p ${MILVUS_BACKUP_DIR}

# 停止Milvus写入（可选）
# docker-compose exec -T milvus python -c "from pymilvus import utility; utility.flush_all()"

# 备份Milvus数据目录
tar -czf ${MILVUS_BACKUP_DIR}/milvus_backup_${DATE}.tar.gz \
    -C ./data/milvus .

echo "✅ Milvus backup completed"

# ==================== Redis备份 ====================
echo "[3/4] Backing up Redis..."
REDIS_BACKUP_DIR="${BACKUP_ROOT}/redis"
mkdir -p ${REDIS_BACKUP_DIR}

# 触发RDB快照
docker-compose exec -T redis redis-cli -a ${REDIS_PASSWORD} BGSAVE

# 等待快照完成
sleep 5

# 复制RDB文件
docker cp metaworkflow-redis:/data/dump.rdb \
    ${REDIS_BACKUP_DIR}/redis_backup_${DATE}.rdb

echo "✅ Redis backup completed"

# ==================== MinIO备份 ====================
echo "[4/4] Backing up MinIO..."
MINIO_BACKUP_DIR="${BACKUP_ROOT}/minio"
mkdir -p ${MINIO_BACKUP_DIR}

# 使用mc mirror命令备份所有桶
docker-compose exec -T minio-init mc mirror \
    myminio ${MINIO_BACKUP_DIR}/minio_backup_${DATE}

echo "✅ MinIO backup completed"

# ==================== 清理旧备份 ====================
echo "Cleaning old backups (retention: ${RETENTION_DAYS} days)..."
find ${BACKUP_ROOT} -type f -mtime +${RETENTION_DAYS} -delete
find ${BACKUP_ROOT} -type d -empty -delete

echo "✅ Old backups cleaned"

# ==================== 备份元数据 ====================
cat > ${BACKUP_ROOT}/backup_${DATE}_metadata.json <<EOF
{
    "timestamp": "${DATE}",
    "components": [
        "postgresql",
        "milvus",
        "redis",
        "minio"
    ],
    "status": "completed",
    "retention_days": ${RETENTION_DAYS}
}
EOF

echo "========================================="
echo "Backup completed successfully!"
echo "Location: ${BACKUP_ROOT}"
echo "========================================="
```

设置定时任务（cron）:

```bash
# 编辑crontab
crontab -e

# 每天凌晨2点执行备份
0 2 * * * /path/to/scripts/backup-all.sh >> /var/log/backup.log 2>&1
```

---

### 11.2 恢复脚本

创建 `scripts/restore-all.sh`:

```bash
#!/bin/bash

# 数据恢复脚本
set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <backup_date>"
    echo "Example: $0 20251209_020000"
    exit 1
fi

BACKUP_DATE=$1
BACKUP_ROOT="/backups"

echo "========================================="
echo "MetaWorkflow Platform Restore"
echo "Backup Date: ${BACKUP_DATE}"
echo "========================================="

# 确认操作
read -p "⚠️  This will OVERWRITE current data. Continue? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "Restore cancelled"
    exit 0
fi

# ==================== 停止服务 ====================
echo "Stopping services..."
docker-compose stop api

# ==================== 恢复PostgreSQL ====================
echo "[1/4] Restoring PostgreSQL..."
PG_BACKUP="${BACKUP_ROOT}/postgres/pg_backup_${BACKUP_DATE}.sql.gz"

if [ -f "$PG_BACKUP" ]; then
    # 删除现有数据库
    docker-compose exec -T postgres psql -U postgres -c "DROP DATABASE IF EXISTS metaworkflow;"
    docker-compose exec -T postgres psql -U postgres -c "CREATE DATABASE metaworkflow;"
    
    # 恢复数据
    gunzip -c ${PG_BACKUP} | docker-compose exec -T postgres psql -U metaworkflow metaworkflow
    
    echo "✅ PostgreSQL restored"
else
    echo "❌ PostgreSQL backup not found: ${PG_BACKUP}"
    exit 1
fi

# ==================== 恢复Milvus ====================
echo "[2/4] Restoring Milvus..."
MILVUS_BACKUP="${BACKUP_ROOT}/milvus/milvus_backup_${BACKUP_DATE}.tar.gz"

if [ -f "$MILVUS_BACKUP" ]; then
    docker-compose stop milvus
    rm -rf ./data/milvus/*
    tar -xzf ${MILVUS_BACKUP} -C ./data/milvus
    docker-compose start milvus
    
    echo "✅ Milvus restored"
else
    echo "❌ Milvus backup not found: ${MILVUS_BACKUP}"
    exit 1
fi

# ==================== 恢复Redis ====================
echo "[3/4] Restoring Redis..."
REDIS_BACKUP="${BACKUP_ROOT}/redis/redis_backup_${BACKUP_DATE}.rdb"

if [ -f "$REDIS_BACKUP" ]; then
    docker-compose stop redis
    docker cp ${REDIS_BACKUP} metaworkflow-redis:/data/dump.rdb
    docker-compose start redis
    
    echo "✅ Redis restored"
else
    echo "❌ Redis backup not found: ${REDIS_BACKUP}"
    exit 1
fi

# ==================== 恢复MinIO ====================
echo "[4/4] Restoring MinIO..."
MINIO_BACKUP="${BACKUP_ROOT}/minio/minio_backup_${BACKUP_DATE}"

if [ -d "$MINIO_BACKUP" ]; then
    docker-compose exec -T minio-init mc mirror \
        ${MINIO_BACKUP} myminio --overwrite
    
    echo "✅ MinIO restored"
else
    echo "❌ MinIO backup not found: ${MINIO_BACKUP}"
    exit 1
fi

# ==================== 重启服务 ====================
echo "Restarting all services..."
docker-compose up -d

# 等待服务就绪
sleep 30

# 验证
echo "Verifying services..."
curl -f http://localhost:8000/health || echo "⚠️  API health check failed"

echo "========================================="
echo "Restore completed!"
echo "========================================="
```

---

## 12. 性能优化

### 12.1 API服务优化

#### 12.1.1 Uvicorn配置优化

```python
# src/main.py
import uvicorn

if __name__ == "__main__":
    uvicorn.run(
        "src.main:app",
        host="0.0.0.0",
        port=8000,
        workers=4,  # CPU核心数
        loop="uvloop",  # 使用uvloop提升性能
        http="httptools",  # 使用httptools
        log_level="info",
        access_log=True,
        use_colors=True,
        timeout_keep_alive=75,
        limit_concurrency=1000,
        limit_max_requests=10000,
        backlog=2048
    )
```

#### 12.1.2 数据库连接池优化

```python
# src/database.py
from sqlalchemy import create_engine
from sqlalchemy.pool import QueuePool

engine = create_engine(
    DATABASE_URL,
    poolclass=QueuePool,
    pool_size=20,  # 基础连接数
    max_overflow=40,  # 额外连接数
    pool_timeout=30,
    pool_recycle=3600,  # 1小时回收连接
    pool_pre_ping=True,  # 连接前验证
    echo=False
)
```

#### 12.1.3 异步优化

```python
# 使用异步数据库查询
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession

async_engine = create_async_engine(
    ASYNC_DATABASE_URL,
    pool_size=20,
    max_overflow=40
)

# 异步Milvus查询
async def search_vectors_async(query_embedding: list, top_k: int = 10):
    """异步向量搜索"""
    loop = asyncio.get_event_loop()
    return await loop.run_in_executor(
        None,
        lambda: collection.search(
            data=[query_embedding],
            anns_field="embedding",
            param={"metric_type": "COSINE", "params": {"nprobe": 10}},
            limit=top_k
        )
    )
```

---

### 12.2 缓存策略优化

#### 12.2.1 多层缓存架构

```python
# src/cache/multi_layer_cache.py
from functools import wraps
import hashlib
import json

class MultiLayerCache:
    """多层缓存：内存 -> Redis -> 数据库"""
    
    def __init__(self, redis_client, max_memory_items=1000):
        self.redis = redis_client
        self.memory_cache = {}
        self.max_memory_items = max_memory_items
    
    def get(self, key: str):
        """获取缓存（L1内存 -> L2 Redis）"""
        # L1: 内存缓存
        if key in self.memory_cache:
            return self.memory_cache[key]
        
        # L2: Redis缓存
        value = self.redis.get(key)
        if value:
            # 回填L1
            if len(self.memory_cache) < self.max_memory_items:
                self.memory_cache[key] = json.loads(value)
            return json.loads(value)
        
        return None
    
    def set(self, key: str, value: any, ttl: int = 3600):
        """设置缓存（同时写入L1和L2）"""
        # L1: 内存
        if len(self.memory_cache) < self.max_memory_items:
            self.memory_cache[key] = value
        
        # L2: Redis
        self.redis.setex(key, ttl, json.dumps(value, ensure_ascii=False))

# 装饰器使用
def cached(ttl=3600, key_prefix=""):
    """缓存装饰器"""
    def decorator(func):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            # 生成缓存key
            cache_key = f"{key_prefix}:{func.__name__}:{hash_args(args, kwargs)}"
            
            # 尝试从缓存获取
            cached_value = cache.get(cache_key)
            if cached_value is not None:
                return cached_value
            
            # 执行函数
            result = await func(*args, **kwargs)
            
            # 存入缓存
            cache.set(cache_key, result, ttl)
            return result
        return wrapper
    return decorator

def hash_args(*args, **kwargs):
    """生成参数哈希"""
    key_str = json.dumps([args, kwargs], sort_keys=True)
    return hashlib.md5(key_str.encode()).hexdigest()
```

---

### 12.3 数据库查询优化

#### 12.3.1 索引优化建议

```sql
-- 添加复合索引
CREATE INDEX idx_tasks_user_status ON generation_tasks(user_id, status, created_at DESC);
CREATE INDEX idx_fragments_task_score ON content_fragments(task_id, quality_score DESC);

-- 添加部分索引
CREATE INDEX idx_active_adapters ON domain_adapters(id) WHERE is_active = true;
CREATE INDEX idx_pending_tasks ON generation_tasks(id) WHERE status = 'pending';

-- 添加GIN索引（JSONB字段）
CREATE INDEX idx_adapters_config_gin ON domain_adapters USING GIN (config);
CREATE INDEX idx_tasks_input_gin ON generation_tasks USING GIN (input_data);
```

#### 12.3.2 查询优化示例

```python
# ❌ N+1查询问题
tasks = db.query(Task).all()
for task in tasks:
    adapter = db.query(Adapter).filter_by(id=task.adapter_id).first()  # N次查询

# ✅ 使用JOIN一次查询
tasks = db.query(Task).join(Adapter).all()

# ✅ 使用预加载
from sqlalchemy.orm import joinedload
tasks = db.query(Task).options(joinedload(Task.adapter)).all()
```

---

### 12.4 Milvus性能优化

```python
# 批量插入优化
def batch_insert_embeddings(embeddings: list, batch_size: int = 1000):
    """批量插入向量"""
    for i in range(0, len(embeddings), batch_size):
        batch = embeddings[i:i + batch_size]
        collection.insert(batch)
    
    # 手动触发flush
    collection.flush()

# 搜索参数优化
search_params = {
    "metric_type": "COSINE",
    "params": {
        "nprobe": 16,  # 增加探测簇数量提高召回率
        "ef": 64  # HNSW索引参数
    }
}

# 使用索引优化搜索
collection.create_index(
    field_name="embedding",
    index_params={
        "metric_type": "COSINE",
        "index_type": "HNSW",  # 使用HNSW获得更好的性能
        "params": {
            "M": 16,  # 节点连接数
            "efConstruction": 256  # 构建时搜索深度
        }
    }
)
```

---

## 13. 安全加固

### 13.1 API安全配置

#### 13.1.1 HTTPS配置（Nginx）

创建 `config/nginx.conf`:

```nginx
# ==================== HTTP -> HTTPS重定向 ====================
server {
    listen 80;
    server_name api.metaworkflow.ai;
    
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }
    
    location / {
        return 301 https://$server_name$request_uri;
    }
}

# ==================== HTTPS配置 ====================
server {
    listen 443 ssl http2;
    server_name api.metaworkflow.ai;
    
    # SSL证书
    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    
    # SSL配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256';
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_stapling on;
    ssl_stapling_verify on;
    
    # 安全头部
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Content-Security-Policy "default-src 'self'" always;
    
    # 隐藏版本信息
    server_tokens off;
    
    # 速率限制
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
    limit_req zone=api_limit burst=20 nodelay;
    
    # 代理配置
    location / {
        proxy_pass http://api:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # 缓冲设置
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
    }
    
    # 健康检查端点（无速率限制）
    location /health {
        limit_req off;
        proxy_pass http://api:8000/health;
    }
}
```

#### 13.1.2 API认证授权

```python
# src/security/auth.py
from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from jose import JWTError, jwt
from passlib.context import CryptContext
from datetime import datetime, timedelta

SECRET_KEY = os.getenv("SECRET_KEY")
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
security = HTTPBearer()

def create_access_token(data: dict, expires_delta: timedelta = None):
    """创建JWT token"""
    to_encode = data.copy()
    expire = datetime.utcnow() + (expires_delta or timedelta(minutes=15))
    to_encode.update({"exp": expire})
    return jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)

async def get_current_user(
    credentials: HTTPAuthorizationCredentials = Depends(security)
):
    """验证并获取当前用户"""
    try:
        token = credentials.credentials
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        user_id: str = payload.get("sub")
        if user_id is None:
            raise HTTPException(status_code=401, detail="Invalid token")
        return user_id
    except JWTError:
        raise HTTPException(status_code=401, detail="Invalid token")

# 使用示例
@app.get("/v2/protected")
async def protected_route(user_id: str = Depends(get_current_user)):
    return {"user_id": user_id, "message": "Access granted"}
```

---

### 13.2 数据库安全

#### 13.2.1 敏感数据加密

```python
# src/security/encryption.py
from cryptography.fernet import Fernet
import base64
import os

class FieldEncryption:
    """字段加密工具"""
    
    def __init__(self, key: str = None):
        if key is None:
            key = os.getenv("ENCRYPTION_KEY")
        self.cipher = Fernet(key.encode())
    
    def encrypt(self, plaintext: str) -> str:
        """加密"""
        encrypted = self.cipher.encrypt(plaintext.encode())
        return base64.b64encode(encrypted).decode()
    
    def decrypt(self, ciphertext: str) -> str:
        """解密"""
        encrypted = base64.b64decode(ciphertext.encode())
        return self.cipher.decrypt(encrypted).decode()

# 在模型中使用
class ContentSource(Base):
    __tablename__ = "content_sources"
    
    id = Column(UUID, primary_key=True)
    name = Column(String(200))
    api_key_encrypted = Column(String(500))  # 加密存储
    
    @property
    def api_key(self):
        """解密API密钥"""
        if self.api_key_encrypted:
            return encryption.decrypt(self.api_key_encrypted)
        return None
    
    @api_key.setter
    def api_key(self, value: str):
        """加密并存储API密钥"""
        if value:
            self.api_key_encrypted = encryption.encrypt(value)
```

#### 13.2.2 SQL注入防护

```python
# ✅ 使用参数化查询（ORM自动处理）
user = db.query(User).filter(User.username == username).first()

# ✅ 原生SQL使用参数绑定
result = db.execute(
    text("SELECT * FROM users WHERE username = :username"),
    {"username": username}
)

# ❌ 永远不要拼接SQL
# bad_query = f"SELECT * FROM users WHERE username = '{username}'"
```

---

### 13.3 容器安全

#### 13.3.1 最小权限Dockerfile

```dockerfile
# ==================== 多阶段构建 ====================
FROM python:3.11-slim as builder

WORKDIR /build
COPY requirements.txt .
RUN pip install --no-cache-dir --user -r requirements.txt

# ==================== 运行阶段 ====================
FROM python:3.11-slim

# 创建非root用户
RUN groupadd -r appuser && useradd -r -g appuser appuser

WORKDIR /app

# 复制依赖
COPY --from=builder /root/.local /home/appuser/.local
COPY --chown=appuser:appuser . .

# 设置PATH
ENV PATH=/home/appuser/.local/bin:$PATH

# 切换到非root用户
USER appuser

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD curl -f http://localhost:8000/health || exit 1

EXPOSE 8000

CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

#### 13.3.2 Docker Compose安全配置

```yaml
services:
  api:
    # ...其他配置
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
      - /var/log
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
```

---

## 14. 故障排查

### 14.1 常见问题诊断

#### 14.1.1 API服务无法启动

```bash
# 1. 检查日志
docker-compose logs -f api

# 2. 检查环境变量
docker-compose exec api env | grep -E "(POSTGRES|REDIS|MILVUS|MINIO)"

# 3. 检查依赖服务
docker-compose ps

# 4. 手动测试连接
docker-compose exec api python -c "
import psycopg2
conn = psycopg2.connect(
    host='postgres',
    database='metaworkflow',
    user='metaworkflow',
    password='your-password'
)
print('✅ PostgreSQL connected')
"
```

#### 14.1.2 Milvus连接失败

```bash
# 检查Milvus状态
curl http://localhost:9091/healthz

# 检查Etcd状态
docker-compose exec etcd etcdctl endpoint health

# 查看Milvus日志
docker-compose logs -f milvus

# 重建Milvus集合
python scripts/init_milvus.py
```

#### 14.1.3 Redis连接问题

```bash
# 测试Redis连接
docker-compose exec redis redis-cli -a your-password ping

# 检查Redis内存
docker-compose exec redis redis-cli -a your-password INFO memory

# 清空Redis缓存（谨慎使用）
docker-compose exec redis redis-cli -a your-password FLUSHDB
```

---

### 14.2 性能问题诊断

#### 14.2.1 慢查询分析

```sql
-- PostgreSQL慢查询
SELECT
    query,
    calls,
    total_time,
    mean_time,
    max_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;

-- 查看当前活动查询
SELECT
    pid,
    now() - query_start as duration,
    state,
    query
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY duration DESC;
```

#### 14.2.2 API性能分析

```python
# 使用cProfile分析
import cProfile
import pstats

profiler = cProfile.Profile()
profiler.enable()

# 执行代码
result = some_slow_function()

profiler.disable()
stats = pstats.Stats(profiler)
stats.sort_stats('cumulative')
stats.print_stats(20)
```

---

### 14.3 故障排查工具

创建 `scripts/diagnostic.sh`:

```bash
#!/bin/bash

echo "========================================="
echo "MetaWorkflow Platform Diagnostics"
echo "========================================="

# ==================== 服务状态 ====================
echo -e "\n[1] Service Status:"
docker-compose ps

# ==================== 资源使用 ====================
echo -e "\n[2] Resource Usage:"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}"

# ==================== 健康检查 ====================
echo -e "\n[3] Health Checks:"
echo -n "API: "
curl -sf http://localhost:8000/health && echo "✅ OK" || echo "❌ FAILED"

echo -n "Postgres: "
docker-compose exec -T postgres pg_isready -U metaworkflow && echo "✅ OK" || echo "❌ FAILED"

echo -n "Redis: "
docker-compose exec -T redis redis-cli -a $REDIS_PASSWORD ping | grep -q "PONG" && echo "✅ OK" || echo "❌ FAILED"

echo -n "Milvus: "
curl -sf http://localhost:9091/healthz && echo "✅ OK" || echo "❌ FAILED"

echo -n "MinIO: "
curl -sf http://localhost:9000/minio/health/live && echo "✅ OK" || echo "❌ FAILED"

# ==================== 磁盘空间 ====================
echo -e "\n[4] Disk Usage:"
df -h | grep -E "(Filesystem|/dev/)"

# ==================== 网络连接 ====================
echo -e "\n[5] Network Connections:"
docker-compose exec api netstat -an | grep ESTABLISHED | wc -l | xargs echo "Established connections:"

# ==================== 最近错误日志 ====================
echo -e "\n[6] Recent Errors (last 20):"
docker-compose logs --tail=100 api 2>&1 | grep -i error | tail -20

echo -e "\n========================================="
echo "Diagnostics completed"
echo "========================================="
```

---

## 15. 环境变量参考

### 15.1 完整环境变量列表

```bash
# ==================== 基础配置 ====================
ENVIRONMENT=production              # 运行环境: development/staging/production
DEBUG=false                         # 调试模式
SECRET_KEY=<随机密钥>              # 应用密钥（至少32字符）
API_HOST=0.0.0.0                   # API监听地址
API_PORT=8000                       # API监听端口

# ==================== PostgreSQL ====================
POSTGRES_HOST=postgres              # 数据库主机
POSTGRES_PORT=5432                  # 数据库端口
POSTGRES_DB=metaworkflow            # 数据库名称
POSTGRES_USER=metaworkflow          # 数据库用户
POSTGRES_PASSWORD=<密码>           # 数据库密码

# ==================== Milvus ====================
MILVUS_HOST=milvus-standalone       # Milvus主机
MILVUS_PORT=19530                   # Milvus端口
MILVUS_COLLECTION_NAME=content_embeddings  # 默认集合名

# ==================== Redis ====================
REDIS_HOST=redis                    # Redis主机
REDIS_PORT=6379                     # Redis端口
REDIS_PASSWORD=<密码>              # Redis密码
REDIS_DB=0                          # Redis数据库编号

# ==================== MinIO ====================
MINIO_ENDPOINT=minio:9000           # MinIO端点
MINIO_ACCESS_KEY=minioadmin         # MinIO访问密钥
MINIO_SECRET_KEY=<密码>            # MinIO密钥
MINIO_SECURE=false                  # 是否使用HTTPS
MINIO_BUCKET_CONTENT=content-assets # 内容桶
MINIO_BUCKET_GENERATED=generated-content  # 生成内容桶

# ==================== LLM服务 ====================
LITELLM_API_BASE=http://litellm:4000  # LiteLLM代理地址
OPENAI_API_KEY=sk-<密钥>           # OpenAI API密钥
AZURE_API_KEY=<密钥>                # Azure OpenAI密钥（可选）
ANTHROPIC_API_KEY=<密钥>            # Claude密钥（可选）

# ==================== 监控 ====================
PROMETHEUS_ENABLED=true             # 是否启用Prometheus
GRAFANA_ADMIN_PASSWORD=<密码>      # Grafana管理员密码
SENTRY_DSN=https://<DSN>@sentry.io/<project>  # Sentry错误追踪

# ==================== 日志 ====================
LOG_LEVEL=INFO                      # 日志级别: DEBUG/INFO/WARNING/ERROR
LOG_FORMAT=json                     # 日志格式: json/text

# ==================== 安全 ====================
ENCRYPTION_KEY=<加密密钥>          # 字段加密密钥（Fernet格式）
CORS_ORIGINS=https://app.metaworkflow.ai  # CORS允许的来源

# ==================== 性能 ====================
WORKERS=4                           # Uvicorn worker数量
MAX_CONNECTIONS=200                 # 数据库最大连接数
CACHE_TTL=3600                      # 默认缓存时间（秒）
```

---

## 16. 端口映射表

| 服务 | 内部端口 | 外部端口 | 协议 | 说明 |
|------|---------|---------|------|------|
| **API服务** | 8000 | 8000 | HTTP | REST API |
| **PostgreSQL** | 5432 | 5432 | TCP | 数据库 |
| **Milvus** | 19530 | 19530 | gRPC | 向量数据库 |
| **Milvus Metrics** | 9091 | 9091 | HTTP | 监控指标 |
| **Redis** | 6379 | 6379 | TCP | 缓存 |
| **MinIO API** | 9000 | 9000 | HTTP | 对象存储API |
| **MinIO Console** | 9001 | 9001 | HTTP | Web控制台 |
| **LiteLLM** | 4000 | 4000 | HTTP | LLM代理 |
| **Prometheus** | 9090 | 9090 | HTTP | 监控 |
| **Grafana** | 3000 | 3000 | HTTP | 可视化 |
| **Nginx** | 80 | 80 | HTTP | 反向代理 |
| **Nginx HTTPS** | 443 | 443 | HTTPS | 安全代理 |
| **Etcd** | 2379 | - | TCP | Milvus元数据 |

---

## 17. 常见问题FAQ

### Q1: 如何升级到新版本？

```bash
# 1. 备份数据
./scripts/backup-all.sh

# 2. 拉取最新代码
git pull origin main

# 3. 重建镜像
docker-compose build --no-cache

# 4. 停止服务
docker-compose down

# 5. 运行迁移
docker-compose run --rm api python scripts/migrate.py

# 6. 启动新版本
docker-compose up -d

# 7. 验证
curl http://localhost:8000/health
```

---

### Q2: 如何扩展API服务？

**Docker Compose方式：**
```bash
docker-compose up -d --scale api=5
```

**Kubernetes方式：**
```bash
kubectl scale deployment api --replicas=5 -n metaworkflow
```

---

### Q3: 如何监控资源使用？

```bash
# 实时监控
docker stats

# Prometheus查询
# CPU使用率: rate(process_cpu_seconds_total[5m])
# 内存使用: process_resident_memory_bytes
# 请求QPS: rate(http_requests_total[1m])

# 访问Grafana
http://localhost:3000
```

---

### Q4: 如何重置密码？

```bash
# PostgreSQL
docker-compose exec postgres psql -U postgres -c \
  "ALTER USER metaworkflow WITH PASSWORD 'new_password';"

# 更新.env文件
sed -i 's/POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD=new_password/' .env

# 重启服务
docker-compose restart api
```

---

### Q5: 如何清理磁盘空间？

```bash
# 清理Docker未使用资源
docker system prune -a --volumes

# 清理旧日志
find ./data/logs -name "*.log" -mtime +30 -delete

# 清理旧备份
find /backups -mtime +30 -delete

# 清理MinIO临时文件
mc rm --recursive --force myminio/temp
```

---

### Q6: 如何启用HTTPS？

```bash
# 1. 安装Certbot
docker run -it --rm --name certbot \
  -v "/etc/letsencrypt:/etc/letsencrypt" \
  -v "/var/www/certbot:/var/www/certbot" \
  certbot/certbot certonly --webroot \
  -w /var/www/certbot \
  -d api.metaworkflow.ai

# 2. 复制证书
cp /etc/letsencrypt/live/api.metaworkflow.ai/fullchain.pem ./config/ssl/
cp /etc/letsencrypt/live/api.metaworkflow.ai/privkey.pem ./config/ssl/

# 3. 重启Nginx
docker-compose restart nginx
```

---

### Q7: 如何迁移到新服务器？

```bash
# 旧服务器
# 1. 完整备份
./scripts/backup-all.sh

# 2. 打包备份
tar -czf metaworkflow-backup.tar.gz /backups /data

# 新服务器
# 3. 传输备份
scp metaworkflow-backup.tar.gz user@new-server:/tmp/

# 4. 解压并恢复
tar -xzf /tmp/metaworkflow-backup.tar.gz
./scripts/restore-all.sh <backup_date>
```

---

### Q8: 性能优化建议？

1. **数据库层**：
   - 定期执行 `VACUUM ANALYZE`
   - 添加合适的索引
   - 优化慢查询

2. **缓存层**：
   - 提高Redis内存限制
   - 使用多层缓存
   - 预加载热点数据

3. **应用层**：
   - 增加Worker数量
   - 启用异步处理
   - 使用连接池

4. **网络层**：
   - 启用HTTP/2
   - 配置CDN
   - 启用gzip压缩

---

## 📝 总结

本部署指南涵盖了MetaWorkflow Platform V2.0的完整部署流程：

✅ **快速开始**: Docker Compose一键部署  
✅ **生产部署**: Kubernetes完整配置  
✅ **数据层**: PostgreSQL + Milvus + Redis + MinIO  
✅ **监控系统**: Prometheus + Grafana + Sentry  
✅ **日志管理**: 结构化日志 + 可选ELK Stack  
✅ **备份恢复**: 自动化备份脚本  
✅ **性能优化**: 多层缓存 + 异步处理 + 索引优化  
✅ **安全加固**: HTTPS + 认证授权 + 数据加密  
✅ **故障排查**: 诊断工具 + 常见问题解决方案

---

**文档维护**:  
- 最后更新: 2025-12-09  
- 维护团队: MetaWorkflow Platform Team  
- 反馈邮箱: support@metaworkflow.ai  

© 2025 MetaWorkflow Platform - V2.0 部署指南

