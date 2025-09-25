package chdyn

import (
	"log"
	"sync"
)

const (
	// defaults for when not provided in Options
	InSize  uint16 = 64
	OutSize uint16 = 64
)

type Options[V any] struct {
	InSize  uint16
	OutSize uint16

	LogPrefix string
	LogDebug  bool
}

type Chan[V any] struct {
	*Options[V]

	exitwg sync.WaitGroup
	exitch chan struct{}

	in  chan V
	out chan V
}

func New[V any](options *Options[V]) *Chan[V] {
	var inSize, outSize uint16
	if options.InSize == 0 {
		inSize = InSize
	} else {
		inSize = options.InSize
	}
	if options.OutSize == 0 {
		outSize = OutSize
	} else {
		outSize = options.OutSize
	}

	c := &Chan[V]{
		Options: options,
		exitwg:  sync.WaitGroup{},
		exitch:  make(chan struct{}, 1),
		in:      make(chan V, inSize),
		out:     make(chan V, outSize),
	}

	c.exitwg.Add(1)
	go c.bridge()

	return c
}

func (c *Chan[V]) Stop() {
	if c.LogDebug {
		log.Printf("%s: synchronized stop starting", c.LogPrefix)
	}

	select {
	case c.exitch <- struct{}{}:
	default:
		if c.LogDebug {
			log.Printf("%s: exitch already signaled", c.LogPrefix)
		}
	}

	c.exitwg.Wait()
	if c.LogDebug {
		log.Printf("%s: synchronized stop done", c.LogPrefix)
	}
}

func (c *Chan[V]) bridge() {
	if c.LogDebug {
		log.Printf("%s: bridge goroutine starting", c.LogPrefix)
	}

	defer func() {
		if c.LogDebug {
			log.Printf("%s: bridge goroutine exiting", c.LogPrefix)
		}
		c.exitwg.Done()
	}()

	list := NewList[V]()

	for {
		if list.Size() == 0 {
			select {
			case <-c.exitch:
				if c.LogDebug {
					log.Printf("%s: exitch received", c.LogPrefix)
				}
				return // exit
			case v := <-c.in:
				list.Push(v)
			}
		} else {
			select {
			case <-c.exitch:
				if c.LogDebug {
					log.Printf("%s: exitch received", c.LogPrefix)
				}
				return // exit
			case v := <-c.in:
				list.Push(v)
			case c.out <- list.Head().Value():
				list.Remove()
			}
		}
	}
}

func (c *Chan[V]) In() chan<- V {
	return c.in
}

func (c *Chan[V]) Out() <-chan V {
	return c.out
}
