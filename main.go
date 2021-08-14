package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
)

var (
	line_clear string = "\r                                                                 "
)

func UNUSED(x ...interface{}) {} //needed a way to throw away t var

func wrkTimer(w, b int) {
	wrk := w
	brk := b
	tmr := "wrk" //sets timer for repeat if wanting to skip the next timer.

	// ticker setup
	wTick := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	fmt.Printf(line_clear)
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
				fmt.Printf(line_clear)
				fmt.Printf("\rMinute: %d:00  ", i)
				UNUSED(t)
			}
		}
	}()

	time.Sleep(time.Duration(wrk) * time.Second)
	wTick.Stop()
	done <- true

	nextTimer(tmr, wrk, brk)

}

func brkTimer(w, b int) {
	brk := b
	wrk := w
	timer := "brk"
	bTick := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	fmt.Printf(line_clear)
	fmt.Printf("\rBreak %d mins", brk)

	go func() {
		i := brk
		for i >= 0 {
			i = i - 1
			select {
			case <-done:
				return
			case t := <-bTick.C:
				fmt.Printf(line_clear)
				fmt.Printf("\rMinute: %d:00  ", i)
				UNUSED(t)
			}
		}
	}()

	time.Sleep(time.Duration(brk) * time.Second)
	bTick.Stop()
	done <- true

	nextTimer(timer, wrk, brk)

}

func nextTimer(t string, w, b int) {
	tmr := t
	brk := b
	wrk := w

	//var n int = 0
	//reader := bufio.NewReader(os.Stdin)
	fmt.Printf(line_clear)
	fmt.Printf("\r> Cont or restart? y(es), n(o) or r(estart)")
	//a, _ := reader.ReadString('\n')

	a, _, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}

	if a == 'y' {
		if tmr == "wrk" {
			brkTimer(wrk, brk)
		} else if tmr == "brk" {
			wrkTimer(wrk, brk)
		}
	} else if a == 'n' {
		os.Exit(3)
	} else if a == 'r' {
		fmt.Printf(line_clear)
		fmt.Printf("\rRestarting this timer")
		if tmr == "wrk" {
			wrkTimer(wrk, brk)
		} else if tmr == "brk" {
			brkTimer(wrk, brk)
		}
	} else {
		fmt.Printf(line_clear)
		fmt.Printf("\rUnexpected. Please use y, n, or r.")
		nextTimer(tmr, wrk, brk)
	}
}

func main() {
	fmt.Println("Pomodoro Timer")

	tWrk := flag.Int("w", 25, "Work minutes.")
	tBrk := flag.Int("b", 5, "Break minutes.")
	flag.Parse()

	w := *tWrk
	b := *tBrk

	if *tWrk == 25 && *tBrk == 5 {
		fmt.Printf("\rUsing Defaults")
		wrkTimer(w, b)
	} else {
		fmt.Printf(line_clear)
		fmt.Printf("\rWork timer = %d and Break timer = %d", w, b)
		wrkTimer(w, b)
	}
}
