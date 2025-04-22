package main

import (
	"time"
	"sync"
	"fmt"
)

var (
    balance int
    mu      sync.Mutex  // Declare a mutex
)

func deposit(amount int) {
    mu.Lock()          // 🔒 Lock
    balance += amount  // Safe modification
    mu.Unlock()        // 🔓 Unlock
}

func main() {
    go deposit(100)  // Goroutine 1
    go deposit(200)  // Goroutine 2
    time.Sleep(time.Second)
    fmt.Println(balance)  
}