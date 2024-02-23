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
	wg := &sync.WaitGroup{}
	for n := range in {
		data := n.(int)
		dataStr := strconv.Itoa(data)
		md5dataStr := DataSignerMd5(dataStr)

		wg.Add(1)
		go func(dataStr, md5dataStr string) {
			defer wg.Done()
			dataStr = DataSignerCrc32(dataStr)
			md5hash := DataSignerCrc32(md5dataStr)
			res := fmt.Sprintf("%s~%s", dataStr, md5hash)
			out <- res
		}(dataStr, md5dataStr)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for n := range in {
		data := n.(string)
		result := make([]string, 6)

		for th := 0; th < 6; th++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				crc32val := DataSignerCrc32(fmt.Sprintf("%d%s", index, data))
				result[index] = crc32val
			}(th)
		}

		wg.Wait()

		finalResult := ""
		for _, res := range result {
			finalResult += res
		}
		out <- finalResult
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
