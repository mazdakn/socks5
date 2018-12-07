package socks5

import (
    "testing"
)

func TestSOCKS5_Start(t *testing.T) {
    server := Socks5Server{}

    err := server.Start("127.0.0.1:8000")
    if err != nil {
        t.Fatal(err)
    }
}
