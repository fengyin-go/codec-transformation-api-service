# Codec Transformation API Service

纯 Go 标准库实现的编码解码后端服务，零第三方依赖。

## 运行方式

在项目根目录下：

```bash
go run ./cmd/server
```

或先构建再运行：

```bash
go build -o server ./cmd/server
./server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## 验证

```bash
go build ./...
go vet ./...
go test ./...
```

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/algorithms | 创建算法 |
| GET | /api/algorithms | 查询算法（支持 category、reversible、keyword、分页） |
| GET | /api/algorithms/{id} | 获取单个算法 |
| PUT | /api/algorithms/{id} | 更新算法 |
| DELETE | /api/algorithms/{id} | 删除算法 |
| GET | /api/operations | 查询转换记录（支持 algorithm、mode、分页） |
| GET | /api/operations/{id} | 获取单条记录 |
| DELETE | /api/operations/{id} | 删除记录 |
| POST | /api/categories | 创建分类 |
| GET | /api/categories | 查询分类 |
| GET | /api/categories/{id} | 获取单个分类 |
| PUT | /api/categories/{id} | 更新分类 |
| DELETE | /api/categories/{id} | 删除分类 |
| POST | /api/codec/encode | 编码转换 |
| POST | /api/codec/decode | 解码转换 |
| POST | /api/codec/batch | 批量转换 |
| POST | /api/codec/compare | 哈希比对 |
| GET | /api/codec/algorithms | 支持的算法列表 |
| GET | /api/codec/algorithms/{name}/reversible | 算法是否可逆 |

## 支持的算法

| 算法 | 分类 | 可逆 |
|------|------|------|
| base64 | encoding | 是 |
| url | escape | 是 |
| html | escape | 是 |
| unicode | escape | 是 |
| hex | encoding | 是 |
| base62 | encoding | 是 |
| morse | encoding | 是 |
| md5 | hash | 否 |
| sha1 | hash | 否 |
| sha256 | hash | 否 |
| hmac_sha256 | hash | 否 |
| jwt | hash | 否（仅解析校验） |

## 统一响应格式

```json
{"code":0,"message":"ok","data":{...}}
```

错误码映射：
- 400 ValidationError
- 404 ErrNotFound
- 409 ErrConflict
- 500 其他内部错误
