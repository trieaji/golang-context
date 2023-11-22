package golangcontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	//menambahkan data ke dalam context
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)

	//context get value
	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextF.Value("b"))

	//bukti kalau context itu bertanya ke child, bukan ke paranet
	fmt.Println(contextA.Value("b"))
}

// func CreateCounter() chan int {
// 	destination := make(chan int)

// 	go func() {
// 		defer close(destination)
// 		counter := 1
// 		for {
// 				destination <- counter
// 				counter++
// 		}
// 	}()

// 	return destination
// }

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <- ctx.Done():
				return
			default :
				destination <- counter
				counter++
				time.Sleep(1 *time.Second) //simulasi slow (context with timeout)
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	cancel() //mengirim sinyal cancel ke context

	time.Sleep(2 *time.Second) //untuk memastikan bahwasanya goroutinenya sudah mati

	fmt.Println("Total goroutine", runtime.NumGoroutine())
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5 * time.Second)
	defer cancel()

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 *time.Second) //untuk memastikan bahwasanya goroutinenya sudah mati

	fmt.Println("Total goroutine", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5 * time.Second))
	defer cancel()

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 *time.Second) //untuk memastikan bahwasanya goroutinenya sudah mati

	fmt.Println("Total goroutine", runtime.NumGoroutine())
}
