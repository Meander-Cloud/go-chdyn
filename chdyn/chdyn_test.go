package chdyn_test

import (
	"log"
	"testing"
	"time"

	"github.com/Meander-Cloud/go-chdyn/chdyn"
)

func Test1(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	ch := chdyn.New(
		&chdyn.Options[uint16]{
			InSize:    chdyn.InSize,
			OutSize:   chdyn.OutSize,
			LogPrefix: "Test1",
			LogDebug:  true,
		},
	)

	var t0 time.Time
	exitch := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-exitch:
				return
			case v := <-ch.Out():
				t1 := time.Now().UTC()
				log.Printf("v=%d, elapsed=%dus", v, t1.Sub(t0).Microseconds())
			}
		}
	}()

	<-time.After(time.Second)
	t0 = time.Now().UTC()
	ch.In() <- 1

	<-time.After(time.Second)
	t0 = time.Now().UTC()
	ch.In() <- 2

	<-time.After(time.Second)
	t0 = time.Now().UTC()
	ch.In() <- 3

	<-time.After(time.Second)
	t0 = time.Now().UTC()
	ch.In() <- 4

	<-time.After(time.Second)
	t0 = time.Now().UTC()
	ch.In() <- 5

	<-time.After(time.Second)
	ch.Stop()
	exitch <- struct{}{}
}

func Test2(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	ch := chdyn.New(
		&chdyn.Options[uint16]{
			InSize:    chdyn.InSize,
			OutSize:   chdyn.OutSize,
			LogPrefix: "Test2",
			LogDebug:  true,
		},
	)

	var bench uint16
	var t0 time.Time
	exitch := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-exitch:
				return
			case v := <-ch.Out():
				if v == bench {
					t1 := time.Now().UTC()
					log.Printf("v=%d, elapsed=%dus", v, t1.Sub(t0).Microseconds())
				}
			}
		}
	}()

	func() {
		<-time.After(time.Second)
		bench = 50
		t0 = time.Now().UTC()
		for i := range bench {
			ch.In() <- i + 1
		}
	}()

	func() {
		<-time.After(time.Second)
		bench = 100
		t0 = time.Now().UTC()
		for i := range bench {
			ch.In() <- i + 1
		}
	}()

	func() {
		<-time.After(time.Second)
		bench = 10000
		t0 = time.Now().UTC()
		for i := range bench {
			ch.In() <- i + 1
		}
	}()

	<-time.After(time.Second)
	ch.Stop()
	exitch <- struct{}{}
}
