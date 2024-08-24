# Ratelimiter
This is Go implement algorithms about ratelimiter.

### Installation:
The package can be installed as a Go module.

```
go get github.com/augustus281/ratelimiter
```

### TOKEN BUCKET
#### How to work
1. Imagine a bucket that holds tokens.
2. The bucket has a maximum capacity of tokens.
3. Tokens are added to the bucket at a fixed rate (e.g., 10 tokens per second).
4. When a request arrives, it must obtain a token from the bucket to proceed.
5. If there are enough tokens, the request is allowed and tokens are removed.
6. If there aren't enough tokens, the request is dropped.

### LEAKY BUCKET
#### How to work
1. Imagine a bucket with a small hole in the bottom.
2. Requests enter the bucket from the top.
3. The bucket processes ("leaks") requests at a constant rate through the hole.
4. If the bucket is full, new requests are discarded.

### FIXED WINDOW COUNTER
#### How to work
1. Time is divided into fixed windows (e.g., 1-minute intervals).
2. Each window has a counter that starts at zero.
3. New requests increment the counter for the current window.
4. If the counter exceeds the limit, requests are denied until the next window.

### SLIDING WINDOW LOG
#### How to work
1. Keep a log of request timestamps.
2. When a new request comes in, remove all entries older than the window size.
3. Count the remaining entries.
4. If the count is less than the limit, allow the request and add its timestamp to the log.
5. If the count exceeds the limit, request is denied.

### SLIDING WINDOW COUNTER
#### How to work
1. Keep track of request count for the current and previous window.
2. Calculate the weighted sum of requests based on the overlap with the sliding window.
3. If the weighted sum is less than the limit, allow the request.

