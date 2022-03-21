package main

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func CombineResults(in, out chan interface{}) {

	var result []string

	for data := range in {

		result = append(result, data.(string))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	out <- strings.Join(result, "_")
}

func pipeWorker(wg *sync.WaitGroup, pipe job, in, out chan interface{}) {

	pipe(in, out)

	wg.Done()
	close(out)
}

func ExecutePipeline(pipeForRunning ...job) {

	var in = make(chan interface{})
	out := make(chan interface{})
	wg := &sync.WaitGroup{}

	for _, pipe := range pipeForRunning {
		wg.Add(1)

		go pipeWorker(wg, pipe, in, out)

		in = out

		// новый канал для следующей передачи данных
		out = make(chan interface{})
	}

	// ожидаем окончания всех горутин
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {

	crc32Chan := make(chan string)
	md5Crc32Chan := make(chan string)
	var dataInStringType string
	var wg = &sync.WaitGroup{}

	for data := range in {

		if reflect.TypeOf(data).String() == "string" {
			dataInStringType = data.(string)
		}

		if reflect.TypeOf(data).String() == "int" {
			tmp := data.(int)
			dataInStringType = strconv.Itoa(tmp)
		}

		// чтобы не было "перегрева" считаем md5 в последовательном режиме
		md5 := DataSignerMd5(dataInStringType)

		wg.Add(2)

		go crc32Worker(wg, md5, crc32Chan)
		go crc32Worker(wg, dataInStringType, md5Crc32Chan)

		tmp1 := <-md5Crc32Chan
		tmp2 := <-crc32Chan

		out <- tmp1 + "~" + tmp2
	}

	wg.Wait()
}

func crc32Worker(wg *sync.WaitGroup, data string, out chan<- string) {

	out <- DataSignerCrc32(data)
	wg.Done()
}

func MultiHash(in, out chan interface{}) {

	wgEx := &sync.WaitGroup{}

	for data := range in {

		wgIn := &sync.WaitGroup{}
		forHash := make([]string, 6)

		wgEx.Add(1)
		go multiHashWorker(wgEx, wgIn, data, forHash, out)
	}

	wgEx.Wait()
}

func multiHashWorker(wgEx *sync.WaitGroup, wgIn *sync.WaitGroup, data interface{}, forHash []string, out chan interface{}) {

	for i := 0; i < 6; i++ {

		wgIn.Add(1)

		go func(idx int) {

			tmp := DataSignerCrc32(strconv.Itoa(idx) + data.(string))
			forHash[idx] = tmp
			wgIn.Done()

		}(i)
	}

	wgIn.Wait()

	cat := strings.Join(forHash, "")

	out <- cat

	wgEx.Done()
}

func main() {

	inputData := []int{0, 1}

	// //nolint:typecheck
	// hashSignJobs := []job{
	// 	job(func(in, out chan interface{}) {
	// 		for _, fibNum := range inputData {
	// 			out <- fibNum
	// 		}
	// 	}),
	// 	job(SingleHash),
	// 	job(MultiHash),
	// 	job(CombineResults),
	// 	job(func(in, out chan interface{}) {
	// 		dataRaw := <-in
	// 		data, ok := dataRaw.(string)
	// 		if !ok {
	// 			fmt.Println("cant convert result data to string")
	// 		}

	// 		fmt.Println(data)
	// 	}),
	// }

	// ExecutePipeline(hashSignJobs...) //nolint:typecheck

	in := make(chan interface{})
	out := make(chan interface{})

	for i := range inputData {
		go SingleHash(in, out)

		in <- i

		fmt.Println("[main]: ", <-out)
	}

}
