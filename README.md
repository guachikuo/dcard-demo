# dcard-demo

## 執行方法
1. 把環境建立起來
```
sh run.sh
```
2. 打 api 測試
```
GET localhost:8080/api/demo/data
```
3. response header 會提供
```
X-RateLimit-Current: 目前的 request 量
X-RateLimit-Remaining: reset 前，剩下的 quota
```
4. 如果超過了 quota，會回 429 Too Many Requests
```
{
    "message": "Error"
}
```

## 簡單介紹
1. 利用 redis 來實作，在分散式系統下，ratelimit 就可以是 global 的阻擋，而不是每台機器處理自己的
2. 用 redis lua script 來達到 atomicity 的效果 (https://redis.io/commands/eval#atomicity-of-scripts)

## 補充
rate limit 還有 token bucket 和 leaky bucket 兩種方法
1. leaky bucket (https://github.com/uber-go/ratelimit)
2. token bucket (https://github.com/juju/ratelimit)

不過這兩個 package 都是單機版的應用，要支援分散式系統可以利用 redis 來實作出來
