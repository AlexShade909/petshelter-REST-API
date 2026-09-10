package app

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestServerStartupFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	a := New()
	a.server.Addr = listener.Addr().String()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := a.RunPetshelter(ctx); err == nil {
		t.Fatal("expected address-in-use error")
	}
}

func TestServerCancellation(t *testing.T) {
	a := New()
	a.server.Addr = "127.0.0.1:0"
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := a.RunPetshelter(ctx); err != nil {
		t.Fatal(err)
	}
}
