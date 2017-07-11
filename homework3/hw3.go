// hw3
package main

import (
	"fmt"
	"sync"
	"time"
)

type Task interface {
	Execute(int) (int, error)
}

type adder struct {
	augend int
}

func (a adder) Execute(addend int) (int, error) {
	result := a.augend + addend
	if result > 127 {
		return 0, fmt.Errorf("Result %d exceeds the adder threshold", a)
	}
	return result, nil
}

type PipeType struct {
	tasks []Task
}

func (m *PipeType) Execute(op int) (int, error) {

	var nextVal int
	var err error

	if len(m.tasks) == 0 {
		return 0, fmt.Errorf("no params!")
	}

	if nextVal, err = m.tasks[0].Execute(op); err != nil {

		return 0, err

	}

	for i := 1; i < len(m.tasks); i++ {

		if nextVal, err = m.tasks[i].Execute(nextVal); err != nil {

			return 0, err

		}

	}

	return nextVal, err

}

func Pipeline(tasks ...Task) Task {

	return &PipeType{tasks}

}

//second part

type lazyAdder struct {
	adder
	delay time.Duration
}

func (la lazyAdder) Execute(addend int) (int, error) {
	time.Sleep(la.delay * time.Millisecond)
	return la.adder.Execute(addend)
}

type FastType struct {
	tasks []Task
}

func Fastest(tasks ...Task) Task {

	return &FastType{tasks}

}

type Res struct {
	r int
	e error
}

func makeRes(r int, e error) Res {
	return Res{r, e}
}

func (f *FastType) Execute(op int) (int, error) {

	if len(f.tasks) == 0 {
		return 0, fmt.Errorf("no params!")
	}

	ch := make(chan Res, 1)

	for _, t := range f.tasks {
		go func() {
			select {
			case ch <- makeRes(t.Execute(op)):

			default:
			}
		}()
	}

	var ret Res = <-ch

	return ret.r, ret.e

}

//third part

type TimeType struct {
	ts Task
	tm time.Duration
}

func Timed(task Task, timeout time.Duration) Task {

	return &TimeType{task, timeout}

}

func (t TimeType) Execute(op int) (int, error) {

	ch := make(chan Res, 1)

	go func() {

		ch <- makeRes(t.ts.Execute(op))

	}()

	select {

	case res := <-ch:
		return res.r, res.e

	case <-time.After(t.tm):
		return 0, fmt.Errorf("time out!")

	}

}

//fourth part
type MapType struct {
	tasks []Task

	f func([]int) int
}

func ConcurrentMapReduce(reduce func(results []int) int, tasks ...Task) Task {

	return &MapType{tasks, reduce}

}

func (m *MapType) Execute(op int) (int, error) {

	if len(m.tasks) == 0 {
		return 0, fmt.Errorf("no params!")
	}

	var wg sync.WaitGroup

	lock := make(chan struct{}, 1) //for res slice protection
	errDetector := make(chan struct{}, 1)
	alright := make(chan struct{}, 1)

	res := make([]int, 0)

	for _, i := range m.tasks {

		wg.Add(1)
		go func() {

			lock <- struct{}{}

			if curRes, err := i.Execute(op); err != nil {

				errDetector <- struct{}{}

			} else {

				res = append(res, curRes)
				<-lock
				wg.Done()
			}

		}()
	}

	go func() {

		wg.Wait()
		alright <- struct{}{}
	}()

	select {

	case <-alright:
		//we are sure that all routines have been passed
		return m.f(res), nil

	case <-errDetector:
		return 0, fmt.Errorf("an error has occurred")

	}

}

//fifth part

var MIN int = -100000000

type GreatestType struct {
	tasks    <-chan Task
	errLimit int
	greatest int
	sync.Mutex
}

func GreatestSearcher(errorLimit int, tasks <-chan Task) Task {

	return &GreatestType{tasks, errorLimit, MIN, sync.Mutex{}}

}

func (g *GreatestType) Execute(op int) (int, error) {

	var wg sync.WaitGroup

	for {

		curTask, ok := <-g.tasks

		if !ok { //the chan is closed

			wg.Wait()

			if g.errLimit < 0 {
				return 0, fmt.Errorf("error limit exceeded")
			}

			if g.greatest == MIN { //the channel has been closed without sending any tasks
				return 0, fmt.Errorf("no tasks has been sent")
			} else {
				return g.greatest, nil
			}
		} else {

			wg.Add(1)
			go func() {

				g.Lock()

				defer g.Unlock()

				if res, err := curTask.Execute(op); err != nil {

					if g.errLimit > 0 {

						g.errLimit--

					}

				} else {

					if res > g.greatest {
						g.greatest = res
					}

				}
				wg.Done()
			}()
		}

	}

}

func main() {

	if res, err := Pipeline(adder{20}, adder{10}, adder{-50}).Execute(100); err != nil {
		fmt.Printf("The pipeline returned an error\n")
	} else {
		fmt.Printf("The pipeline returned %d\n", res)
	}

	f := Fastest(
		lazyAdder{adder{20}, 500},
		lazyAdder{adder{50}, 300},
		adder{41},
	)

	if res, err := f.Execute(1); err != nil {
		fmt.Printf("The fastest returned an error\n")
	} else {
		fmt.Printf("The fastest returned %d\n", res)
	}

	if r1, e1 := Timed(lazyAdder{adder{20}, 50}, 2*time.Millisecond).Execute(2); e1 != nil {
		fmt.Printf("The timed returned an error\n")
	} else {
		fmt.Printf("The timed returned %d\n", r1)
	}

	if r2, e2 := Timed(lazyAdder{adder{20}, 50}, 300*time.Millisecond).Execute(2); e2 != nil {
		fmt.Printf("The timed returned an error\n")
	} else {
		fmt.Printf("The timed returned %d\n", r2)
	}

	reduce := func(results []int) int {
		smallest := 128
		for _, v := range results {
			if v < smallest {
				smallest = v
			}
		}
		return smallest
	}

	mr := ConcurrentMapReduce(reduce, adder{30}, adder{50}, adder{20})
	if res, err := mr.Execute(125); err != nil {
		fmt.Printf("We got an error!\n")
	} else {
		fmt.Printf("The ConcurrentMapReduce returned %d\n", res)
	}

	tasks := make(chan Task)
	gs := GreatestSearcher(2, tasks) // Приемаме 2 грешки

	go func() {
		//tasks <- adder{4}
		//tasks <- lazyAdder{adder{22}, 20}
		//tasks <- adder{125} // Това е първата "допустима" грешка (защото 125+10 > 127)
		//time.Sleep(50 * time.Millisecond)
		//tasks <- adder{32} // Това би трябвало да "спечели"

		// Това би трябвало да timeout-не и да е втората "допустима" грешка
		//tasks <- Timed(lazyAdder{adder{100}, 2000}, 20*time.Millisecond)

		// Ако разкоментираме това, gs.Execute() трябва да върне грешка
		//	tasks <- adder{127} // трета (и недопустима) грешка

		close(tasks)
	}()

	if result, err := gs.Execute(10); err != nil {
		fmt.Printf("We got an error!\n")
	} else {
		fmt.Printf("The Greatest returned %d\n", result)
	}

}
