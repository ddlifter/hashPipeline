package main

import (
	"fmt"
	"sync"
)

type job func(in, out chan interface{})

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{}, 2)

	wg := &sync.WaitGroup{}

	for _, jobFunc := range jobs {
		wg.Add(1)
		go func(jobFunc job, in, out chan interface{}) {
			defer wg.Done()
			jobFunc(in, out)
		}(jobFunc, in, out)
		in = out
		out = make(chan interface{}, 2)
	}
	wg.Wait()
	close(out)
}

func SingleHash(in, out chan interface{}) {
	go func() {
		for n := range in {
			out <- n.(int) * n.(int)
		}
	}()
}

func MultiHash(in, out chan interface{}) {
	go func() {
		for n := range in {
			fmt.Println(n)
			out <- n.(int) + 1
		}
	}()
}

func CombineResults(in, out chan interface{}) {
	go func() {
		for n := range in {
			fmt.Println(n)
			out <- n.(int) + 1
		}
	}()
}

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	}
	ExecutePipeline(hashSignJobs...)
}
