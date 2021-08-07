package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func UNUSED(x ...interface{}) {}

func workTimer(w, b int) {
	min := w
	brk := b
	wTick := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	fmt.Printf("\rWork 25min")

	go func() {
		i := min
		for i >= 0 {
			i = i - 1
			select {
			case <-done:
				return
			case t := <-wTick.C:
				fmt.Printf("\rMinute: %d:00  ", i)
				UNUSED(t)
			}
		}
	}()

	time.Sleep(time.Duration(min) * time.Second)
	wTick.Stop()
	done <- true

	n := nextTimer()

	if n == 1 {
		brkTimer(min, brk)
	} else if n == 2 {
		workTimer(min, brk)
	}
}

func brkTimer(w, b int) {
	min := b
	wrk := w
	bTick := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	fmt.Printf("\rBreak 5 mins")

	go func() {
		i := min
		for i >= 0 {
			i = i - 1
			select {
			case <-done:
				return
			case t := <-bTick.C:
				fmt.Printf("\rMinute: %d:00  ", i)
				UNUSED(t)
			}
		}
	}()

	time.Sleep(time.Duration(min) * time.Second)
	bTick.Stop()
	done <- true

	n := nextTimer()

	if n == 1 {
		workTimer(wrk, min)
	} else if n == 2 {
		brkTimer(wrk, min)
	}
}

func nextTimer() int {

	var n int = 0
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\r> Continue with next timer or restart this timer?\ny(es), n(o) or r(estart)")
	a, _ := reader.ReadString('\n')
	fmt.Println(a)

	if a == "y\n" {
		n = 1
	} else if a == "n\n" {
		os.Exit(3)
	} else if a == "r\n" {
		n = 2
		fmt.Println("Restarting this timer")
	} else {
		fmt.Println("Unexpected response. Pleas type 'y', 'n', or 'r'.")
		nextTimer()
	}

	return n
}

func main() {
	fmt.Println("Pomodoro Timer")

	tWrk := flag.Int("work", 25, "Work minutes.")
	tBrk := flag.Int("break", 5, "Break minutes.")
	tDur := flag.Int("durration", 90, "Timer Durration.")

	flag.Parse()

	w := *tWrk
	b := *tBrk
	d := *tDur

	fmt.Println("work:", w)
	fmt.Println("break:", b)
	fmt.Println("durration:", d)

	if *tDur == 90 {
		fmt.Println("Using Defaults")
		oaTimer := time.NewTimer(time.Duration(d) * time.Second)

		workTimer(w, b)
		<-oaTimer.C
	} else if *tDur <= 89 {
		fmt.Println("Durration should be greater than 2 hours")
	} else {
		oaTimer := time.NewTimer(time.Duration(d) * time.Second)
		<-oaTimer.C
		workTimer(w, b)
	}

	fmt.Println("work:", *tWrk)
	fmt.Println("break:", *tBrk)
	fmt.Println("durration:", *tDur)
}
