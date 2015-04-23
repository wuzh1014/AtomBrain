package main

import (
	"github.com/ssdb/gossdb/ssdb"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func fetchMind(tem []string, dept int, whole string) {
	time1 := time.Now().UnixNano()
	var outPut []string
	var weight float32
	for i := 0; i < len(tem); i++ {
		str, w := getLink(tem[i], dept)
		weight += w
		outPut = append(outPut, str)
	}
	var sentence string
	for _, j := range outPut {
		sentence += j
	}
	time2 := time.Now().UnixNano()
	pons("fet->", (time2 - time1))
	//	Sentence[sentence] = weight
	go setSentence(sentence, whole, weight)

}

func setSentence(sentence, whole string, weight float32) {

	db_s, err := ssdb.Connect(ip, port)
	if err != nil {
		os.Exit(1)
	}
	defer db_s.Close()
	db_s.Do("hset", whole, sentence, strconv.Itoa(int(weight)))

	goWait.Done()
}

func getLink(temJ string, dept int) (string, float32) {
	db_r, err := ssdb.Connect(ip, port)
	if err != nil {
		os.Exit(1)
	}
	defer db_r.Close()

	zsize, err := db_r.Do("zsize", temJ)
	if nil != err {
		pons(err)
	}

	if "0" != zsize[1] {
		//		time1 := time.Now().UnixNano()
		if temJ == "我" {
			pons("等等等等等等等等等等等等")
		}
		rangeVal, err := db_r.Do("zrscan", temJ, "", "", "", 1000)
		if temJ == "我" {
			pons("完完完完完完完完完完完完")
		}
		//		time2 := time.Now().UnixNano()
		//		pons(temJ)
		//		pons("get+", zsize[1], err)
		//		pons("get->", (time2 - time1))
		if nil != err {
			pons(err)
		}
		//		pons("range+", rangeVal, err)
		var count float32
		count = 0

		//数组第三个，即下标为2的元素为第一个权值
		for i := 2; i < len(rangeVal); {

			temInt, _ := strconv.Atoi(rangeVal[i])
			count += float32(temInt)
			i += 2
		}
		rand.Seed(time.Now().UnixNano())
		count = float32(rand.Int31n(int32(count)))
		//同取第一个权值
		for i := 2; i < len(rangeVal); {
			temInt, _ := strconv.Atoi(rangeVal[i])
			if count <= float32(temInt) {
				return rangeVal[i-1], float32(temInt)
			}
			count -= float32(temInt)
			i += 2
		}
	} else {
		//		<-globalLock2
		return "", 0.0
	}
	return "", 0.0
}
