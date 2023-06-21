# Benchmarking the Microservices API

## email.GetEmails

Retrieve 5 emails for a user (timeout after 20s):

```
Summary:
  Count:        200
  Total:        20.11 s
  Slowest:      20.00 s
  Fastest:      914.80 ms
  Average:      10.67 s
  Requests/sec: 9.94

Response time histogram:
  914.797   [1]  |∎∎
  2823.014  [21] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  4731.231  [20] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  6639.448  [17] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  8547.665  [20] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  10455.881 [20] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  12364.098 [15] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  14272.315 [21] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  16180.532 [19] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  18088.749 [20] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  19996.966 [20] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎

Latency distribution:
  10 % in 2.51 s 
  25 % in 5.42 s 
  50 % in 10.18 s 
  75 % in 15.28 s 
  90 % in 18.17 s 
  95 % in 19.16 s 
  99 % in 19.79 s 

Status code distribution:
  [DeadlineExceeded]   6 responses     
  [OK]                 194 responses   

Error distribution:
  [6]   rpc error: code = DeadlineExceeded desc = context deadline exceeded
```

## parse.ParseEmail

Retrieve a email (timeout after 100s):

```
Summary:
  Count:        200
  Total:        154.60 s
  Slowest:      47.28 s
  Fastest:      4.04 s
  Average:      33.58 s
  Requests/sec: 1.29

Response time histogram:
  4038.180  [1]   |
  8362.642  [8]   |∎∎∎
  12687.104 [4]   |∎
  17011.566 [5]   |∎∎
  21336.028 [9]   |∎∎∎
  25660.490 [5]   |∎∎
  29984.952 [7]   |∎∎∎
  34309.414 [10]  |∎∎∎∎
  38633.876 [107] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  42958.338 [35]  |∎∎∎∎∎∎∎∎∎∎∎∎∎
  47282.800 [9]   |∎∎∎

Latency distribution:
  10 % in 17.21 s 
  25 % in 34.36 s 
  50 % in 36.35 s 
  75 % in 38.14 s 
  90 % in 41.42 s 
  95 % in 42.87 s 
  99 % in 45.61 s 

Status code distribution:
  [OK]   200 responses
```