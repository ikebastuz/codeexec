#!/bin/bash
for i in {1..100}; do
  (
    resp=$(curl -s -w "\nRequest $i: HTTP %{http_code} - %{time_total}s\n" -X POST http://localhost:1450/run \
      -F 'language=javascript' \
      -F 'code=let sum=0;for(let i=0;i<1e8;i++){sum+=i;}console.log(sum);')
    echo "$resp"
  ) &
done
wait