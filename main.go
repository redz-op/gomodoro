package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func UNUSED(x ...interface{}) {} //needed a way to throw away t var

func wrkTimer(w, b int) {
	wrk := w // timer setting
	brk := b // break timer pass through to brkTimer func
	t := "wrk"

	// ticker setup
	wTick := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	fmt.Printf("\rWork timer set for %d", wrk)

	// Work countdown timer.
	go func() {
		i := wrk
		for i >= 0 {
			i = i - 1
			select {
			case <-done:
				return
			case t := <-wTick.C:
				fmt.Printf("\r                     ")
				fmt.Printf("\rMinute: %d:00  ", i)
				UNUSED(t)
			}
		}
	}()

	time.Sleep(time.Duration(wrk) * time.Second)
	wTick.Stop()
	done <- true

	n := nextTimer(t, wrk, brk)

	if n == 1 {
		brkTimer(wrk, brk)
	} else if n == 2 {
		wrkTimer(wrk, brk)
	}
}

func brkTimer(w, b int) {
	brk := b
	wrk := w
	timer := "brk"
	bTick := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	fmt.Printf("\rBreak 5 mins")

	go func() {
		i := brk
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

	time.Sleep(time.Duration(brk) * time.Second)
	bTick.Stop()
	done <- true

	n := nextTimer(timer, wrk, brk)

	if n == 1 {
		wrkTimer(wrk, brk)
	} else if n == 2 {
		brkTimer(wrk, brk)
	}

}

func nextTimer(t string, w, b int) int {
	repeat := t
	brk := b
	wrk := w

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
		fmt.Println("Unexpected. Please use y, n, or r.")
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

	if *tDur == 90 {
		fmt.Printf("\rUsing Defaults")
		oaTimer := time.NewTimer(time.Duration(d) * time.Second)

		wrkTimer(w, b)
		<-oaTimer.C
	} else if *tDur <= 89 {
		fmt.Println("Durration should be greater than 2 hours")
	} else {
		oaTimer := time.NewTimer(time.Duration(d) * time.Second)
		<-oaTimer.C
		wrkTimer(w, b)
	}
}
