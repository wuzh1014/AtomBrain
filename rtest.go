package main

import (
	"bufio"
	"fmt"
	. "github.com/dgryski/dgobloom"
	"github.com/ssdb/gossdb/ssdb"
	"hash/fnv"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var mark, all int
var Sentence map[string]float32

func pons(str ...interface{}) {
	fmt.Println(str)
}

func Sub(whole, str string, pre []string, deep, dept int) int {
	rs := []rune(str)
	rl := len(rs)
	end := rl
	for i := 1; i < rl; i++ {
		if deep > dept-1 {
			tem := append(pre, string(rs[0:i]), string(rs[i:end]))
			j := 0
			k := 1

			for k < len(tem) {

				if linkMap.Exists([]byte(tem[j] + "/" + tem[k])) {
					//					mark++
					//					pons("", mark, "/", all, float64(mark)/float64(all))
					//					continue
				} else {
					all++
					linkMap.Insert([]byte(tem[j] + "/" + tem[k]))
					//					goWait.Add(1)
					go GoLink(tem[j], tem[k], dept)
				}
				j++
				k++
			}
			goWait.Add(1)
			go fetchMind(tem, dept, whole)

		}

		if deep < dept {
			Sub(whole, string(rs[i:end]), append(pre, string(rs[0:i])), deep+1, dept)
		}

	}
	//	goWait.Done()
	return dept
}

func GoLink(temJ, temK string, dept int) {
	db_w, err := ssdb.Connect(ip, port)
	if err != nil {
		os.Exit(1)
	}
	defer db_w.Close()

	_, _ = db_w.Do("zincr", temJ, temK, 1)
	//	goWait.Done()
}

//var globalLock, globalLock2 chan int
//var temLinkMap map[string]map[string]float32
var linkMap BloomFilter

const CAPACITY = 100000
const ERRPCT = 0.01

var goWait sync.WaitGroup

var db *ssdb.Client
var val interface{}
var err error
var ip string
var port int

func main() {
	fh, _ := os.Open("速度与激情6.txt")
	buf := bufio.NewReader(fh)
	//初始化数据库
	ip = "192.168.228.131"
	port = 8888

	db, err = ssdb.Connect(ip, port)
	if err != nil {
		os.Exit(1)
	}
	defer db.Close()

	//	for {
	l, _, err := buf.ReadLine()
	if err != nil {
		//		break
	}

	sen := string(l)
	senRune := []rune(sen)
	if len(sen) > 30 {
		senRune = senRune[:10]
	}
	sentence := string(senRune)
	pons("from:", sentence)
	Sentence = make(map[string]float32)

	mark = 0
	all = 0

	/*
	 * Bloom Filter
	 */

	saltsNeeded := SaltsRequired(CAPACITY, ERRPCT)
	salts := make([]uint32, saltsNeeded)
	for i := uint(0); i < saltsNeeded; i++ {
		salts[i] = rand.Uint32()
	}
	linkMap = NewBloomFilter(CAPACITY, ERRPCT, fnv.New32a(), salts)

	for i := 1; i < 4; i++ {
		//			goWait.Add(1)
		go Sub(sentence, sentence, []string{}, 1, i)
		//			pons(linkMap.GetElements())
	}

	linkMap = NewBloomFilter(CAPACITY, ERRPCT, fnv.New32a(), salts)
	goWait.Add(1)
	go goSay(sentence)
	//	}
	goWait.Wait()
}
func goSay(sentence string) {
	db_g, err := ssdb.Connect(ip, port)
	if err != nil {
		os.Exit(1)
	}
	defer db_g.Close()
	hgetall, _ := db_g.Do("hgetall", sentence)
	pons(hgetall)
	var count float32
	count = 1
	//	for _, weight := range hgetall {
	for i := 2; i < len(hgetall); {
		tem, _ := strconv.Atoi(hgetall[i])
		count += float32(tem)
		i += 2
	}

	rand.Seed(time.Now().UnixNano())
	count = float32(rand.Int63n(int64(count))) - 1
	//	for outStr, weight := range hgetall {
	for i := 2; i < len(hgetall); {
		tem, _ := strconv.Atoi(hgetall[i])
		if count <= float32(tem) {
			pons("I say:", hgetall[i-1])
			break
		}
		count -= float32(tem)
		i += 2
	}
	goWait.Done()
	//	Sentence = make(map[string]float32)
}
