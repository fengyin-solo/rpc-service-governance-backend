# RPCGate - 轻量 RPC 服务注册与调用治理

纯 Go 标准库实现的轻量 RPC 服务治理后台，支持服务、方法、节点、调用记录与超时规则的管理。

## 运行说明

```bash
cd origin
go build ./...
go test ./...
go run cmd/server/main.go
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## API 概览

### Service 服务
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/services | 创建服务 |
| GET | /api/services | 服务列表（支持 status/owner/keyword 筛选） |
| GET | /api/services/{id} | 获取服务详情 |
| PUT | /api/services/{id} | 更新服务（含状态机校验） |
| DELETE | /api/services/{id} | 删除服务（级联删除方法、节点、规则） |

### Method 方法
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/methods | 创建方法（外键校验） |
| GET | /api/methods | 方法列表（支持 service_id/keyword 筛选） |
| GET | /api/methods/{id} | 获取方法详情 |
| PUT | /api/methods/{id} | 更新方法 |
| DELETE | /api/methods/{id} | 删除方法 |

### Node 节点
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/nodes | 注册节点 |
| GET | /api/nodes | 节点列表（支持 service_id/region/healthy/keyword 筛选） |
| GET | /api/nodes/{id} | 获取节点详情 |
| PUT | /api/nodes/{id} | 更新节点 |
| DELETE | /api/nodes/{id} | 删除节点 |
| POST | /api/nodes/batch | 批量注册节点 |
| DELETE | /api/nodes/batch/{service_id} | 按服务批量下线节点 |

### Invocation 调用记录
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/invocations | 记录调用（仅 up 状态服务可承接） |
| GET | /api/invocations | 调用列表（支持 service_id/method_id/status/node_id 筛选） |
| GET | /api/invocations/{id} | 获取调用详情 |
| DELETE | /api/invocations/{id} | 删除调用记录 |
| POST | /api/invocations/batch-delete | 批量删除调用记录 |
| GET | /api/invocations/report | 按方法聚合报告（总次数/成功率/平均耗时） |

### TimeoutRule 超时重试规则
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/timeout-rules | 创建规则 |
| GET | /api/timeout-rules | 规则列表（支持 service_id/method_id/enabled 筛选） |
| GET | /api/timeout-rules/{id} | 获取规则详情 |
| PUT | /api/timeout-rules/{id} | 更新规则 |
| DELETE | /api/timeout-rules/{id} | 删除规则 |

### Stats 统计
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/stats/overview | 聚合概览：服务状态分组、调用状态分组、节点健康占比、Top 慢方法 |
