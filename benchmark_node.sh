#!/bin/sh
for i in {1..10}; do
  (
    resp=$(curl -s -w "\nRequest $i: HTTP %{http_code} - %{time_total}s\n" --location 'localhost:1450/run' \
--form 'language="javascript"' \
--form 'code="const f=n=>{let r=1n;for(let i=2n;i<=n;i++)r*=i;return r};f(100000n);"')
    echo "$resp"
  ) &
done
wait
