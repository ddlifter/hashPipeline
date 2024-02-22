package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})

	wg := &sync.WaitGroup{}

	for _, jobFunc := range jobs {
		wg.Add(1)
		go func(jobFunc job, in, out chan interface{}) {
			defer wg.Done()
			jobFunc(in, out)
			close(out)
		}(jobFunc, in, out)
		in = out
		out = make(chan interface{})
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	for n := range in {
		data := n.(int)
		dataStr := strconv.Itoa(data)
		crc32hash := DataSignerCrc32(dataStr)

		md5hash := DataSignerCrc32(DataSignerMd5(dataStr))
		res := fmt.Sprintf("%s~%s", crc32hash, md5hash)
		out <- res
	}
}

func MultiHash(in, out chan interface{}) {
	for n := range in {
		result := ""
		for th := 0; th < 6; th++ {
			i := strconv.Itoa(th)
			crc32val := DataSignerCrc32(fmt.Sprintf("%s%s", i, n.(string)))
			result += crc32val
			fmt.Println(crc32val)
		}
		fmt.Println(result)
		out <- result
	}
}

func CombineResults(in, out chan interface{}) {
	res := make([]string, len(in))
	for n := range in {
		res = append(res, n.(string))
	}
	sort.Strings(res)
	resultStr := strings.Join(res, "_")
	out <- resultStr
}
