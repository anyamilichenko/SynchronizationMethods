package main

import (
	"fmt"
	"sync"
	"time"
)

// sync.Mutex
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func demoMutex() {
	fmt.Println("1. sync.Mutex")
	var wg sync.WaitGroup
	counter := &Counter{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Println("Final counter value:", counter.Value())
}

// sync.RWMutex
type SafeMap struct {
	mu sync.RWMutex
	m  map[string]int
}

func (s *SafeMap) Read(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[key]
}

func (s *SafeMap) Write(key string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func demoRWMutex() {
	fmt.Println("2. sync.RWMutex")
	s := SafeMap{m: make(map[string]int)}
	s.Write("a", 1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("Read a:", s.Read("a"))
	}()

	go func() {
		defer wg.Done()
		s.Write("a", 42)
	}()

	wg.Wait()
}

// Channels
func worker(data <-chan int, done chan<- bool) {
	sum := 0
	for num := range data {
		sum += num
	}
	fmt.Println("Sum:", sum)
	done <- true
}

func demoChannels() {
	fmt.Println("3. Channels")
	data := make(chan int)
	done := make(chan bool)

	go worker(data, done)

	for i := 1; i <= 5; i++ {
		data <- i
	}
	close(data)
	<-done
}

// sync.WaitGroup
func printMessage(msg string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(msg)
}

func demoWaitGroup() {
	fmt.Println("4. sync.WaitGroup")
	var wg sync.WaitGroup
	messages := []string{"Hello", "Anna", "its string"}

	for _, msg := range messages {
		wg.Add(1)
		go printMessage(msg, &wg)
	}
	wg.Wait()
}

// sync.Once
var once sync.Once

func initialize() {
	fmt.Println("Initialized")
}

func demoOnce() {
	fmt.Println("5. sync.Once")
	for i := 0; i < 3; i++ {
		go func() {
			once.Do(initialize)
		}()
	}
	time.Sleep(1 * time.Second)
}

// sync.Cond
func demoCond() {
	fmt.Println("6. sync.Cond")
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	ready := false

	go func() {
		mu.Lock()
		for !ready {
			cond.Wait()
		}
		fmt.Println("Received signal")
		mu.Unlock()
	}()

	time.Sleep(1 * time.Second)
	mu.Lock()
	ready = true
	cond.Signal()
	mu.Unlock()
	time.Sleep(1 * time.Second)
}

func main() {
	demoMutex()
	demoRWMutex()
	demoChannels()
	demoWaitGroup()
	demoOnce()
	demoCond()
}
