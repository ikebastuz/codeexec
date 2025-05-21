#!/bin/sh

cat > /tmp/code.go <<EOF
package main

import (
    "fmt"
    "math/big"
)

func main() {
    fmt.Println(func(n int64) *big.Int {
        r := big.NewInt(1)
        for i := int64(2); i <= n; i++ {
            r.Mul(r, big.NewInt(i))
        }
        return r
    }(20000))
}
EOF

for i in {1..10}; do
  (
    resp=$(curl -s -w "\nRequest $i: HTTP %{http_code} - %{time_total}s\n" --location 'localhost:1415/run' \
      --form 'language=go' \
      --form "code=< /tmp/code.go")
    echo "$resp"
  ) &
done
wait
