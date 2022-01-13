package probes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

// Time to wait for the server to start before sending requests
const wait = time.Second / 2

func TestProbe_LivenessSucceed(t *testing.T) {
	const want = liveResponse

	addr := freeAddr(t)
	p := New(addr)

	go func() {
		err := p.Start()
		if err != nil && err != http.ErrServerClosed {
			panic("could not start probes server")
		}
	}()

	p.LivenessSucceed()

	// Wait for the server to start
	time.Sleep(wait)

	resp, err := http.Get("http://" + addr + "/live")
	if err != nil {
		panic("liveness probe failed: " + err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read /live response body: " + err.Error())
	}

	if string(body) != want {
		t.Errorf("got: %q, want: %q", string(body), want)
	}
}

func TestProbe_LivenessFail(t *testing.T) {
	const want = "error, check the logs"

	addr := freeAddr(t)
	p := New(addr)

	go func() {
		err := p.Start()
		if err != nil && err != http.ErrServerClosed {
			panic("could not start probes server")
		}
	}()

	p.LivenessFail()

	// Wait for the server to start
	time.Sleep(wait)

	respLive, err := http.Get("http://" + addr + "/live")
	if err != nil {
		panic("liveness probe failed: " + err.Error())
	}

	liveBody, err := io.ReadAll(respLive.Body)
	if err != nil {
		t.Fatalf("failed to read /live response body: " + err.Error())
	}

	if string(liveBody) != want {
		t.Errorf("got: %q, want: %q", string(liveBody), want)
	}
}

func TestProbe_ReadinessSucceed(t *testing.T) {
	const want = readyResponse

	addr := freeAddr(t)
	p := New(addr)

	go func() {
		err := p.Start()
		if err != nil && err != http.ErrServerClosed {
			panic("could not start probes server")
		}
	}()

	p.ReadinessSucceed()

	// Wait for the server to start
	time.Sleep(wait)

	resp, err := http.Get("http://" + addr + "/ready")
	if err != nil {
		panic("liveness probe failed: " + err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read /live response body: " + err.Error())
	}

	if string(body) != want {
		t.Errorf("got: %q, want: %q", string(body), want)
	}
}

func TestProbe_ReadinessFail(t *testing.T) {
	const want = "error, check the logs"

	addr := freeAddr(t)
	p := New(addr)

	go func() {
		err := p.Start()
		if err != nil && err != http.ErrServerClosed {
			panic("could not start probes server")
		}
	}()

	p.ReadinessFail()

	// Wait for the server to start
	time.Sleep(wait)

	respLive, err := http.Get("http://" + addr + "/ready")
	if err != nil {
		panic("liveness probe failed: " + err.Error())
	}

	liveBody, err := io.ReadAll(respLive.Body)
	if err != nil {
		t.Fatalf("failed to read /live response body: " + err.Error())
	}

	if string(liveBody) != want {
		t.Errorf("got: %q, want: %q", string(liveBody), want)
	}
}

func TestProbe_Shutdown(t *testing.T) {
	addr := freeAddr(t)
	p := New(addr)

	go func() {
		err := p.Start()
		if err != nil && err != http.ErrServerClosed {
			panic("could not start probes server")
		}
	}()

	// Wait for the server to start
	time.Sleep(wait)

	err := p.Shutdown(context.Background())

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		t.Errorf("want nil err, got: %v", err)
	}
}
func freeAddr(t *testing.T) string {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			t.Fatalf(fmt.Sprintf("failed to find a free port: %v", err))
		}
	}
	defer l.Close()

	return l.Addr().String()
}

func TestProbe_raceCondition(t *testing.T) {
	const want = "error, check the logs"

	addr := freeAddr(t)
	p := New(addr)

	go func() {
		err := p.Start()
		if err != nil && err != http.ErrServerClosed {
			panic("could not start probes server")
		}
	}()

	go func() {
		p.ReadinessFail()
	}()

	p.ReadinessFail()

	// Wait for the server to start
	time.Sleep(wait)

	respLive, err := http.Get("http://" + addr + "/ready")
	if err != nil {
		panic("liveness probe failed: " + err.Error())
	}

	liveBody, err := io.ReadAll(respLive.Body)
	if err != nil {
		t.Fatalf("failed to read /live response body: " + err.Error())
	}

	if string(liveBody) != want {
		t.Errorf("got: %q, want: %q", string(liveBody), want)
	}
}
