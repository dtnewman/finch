package main

import "time"

func IsReady(what string, seconds int64, ch chan<- bool) {
    time.Sleep(time.Duration(seconds * 1e9)) // nanoseconds
    println(what, "is ready!")
    ch <- true
}

func main() {
    ch := make(chan bool)
    println("Let's go!")
    println("Steep the coffee for 6 seconds...")
    go IsReady("Coffee", 6, ch)
    println("Steep the tea for 2 seconds...")
    go IsReady("Tea", 2, ch)
    println("Now the coffee and tea are steeping.")

    for i := 0; i < 2; i++ {
        <-ch
    }
}
