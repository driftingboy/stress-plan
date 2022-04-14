
## 命令目录
- run 执行压测
- plan 定制、查看压测计划


curl -X POST "http://119.3.106.151:10100/v1/app/evidences" -H "accept: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzUxMiJ9.eyJleHAiOjE2NDAwOTUzNDYsImlhdCI6MTY0MDA4ODE0NiwianRpIjoianI1NnZxanludmtxN3AiLCJzdWIiOiJ1aWQtdGVuYW50In0.y-xNb2DDCi2cU1JQlO9HAxoG_AyjYha8I3wfcv5x9dBnDVLwgDSdzIYl9BlzHyww3fOIj4VImA-w26n2LMPATQ" -H "Content-Type: application/json" -d "{ \"tenant_id\": \"tid-yuhu1\", \"title\": \"stress-test-01\", \"content\": \"c3RyZXNzLXRlc3QtMDE=\", \"evidence_type\": \"text\"}"