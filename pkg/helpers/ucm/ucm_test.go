package ucm

import (
	"log"
	"math/rand"
	"sort"
	"testing"
	"time"
)

//测试map随机性

var addrs = map[string]bool{
	"192.168.0.1":  true,
	"192.168.0.2":  true,
	"192.168.0.3":  true,
	"192.168.0.4":  true,
	"192.168.0.5":  true,
	"192.168.0.6":  true,
	"192.168.0.7":  true,
	"192.168.0.8":  true,
	"192.168.0.9":  true,
	"192.168.0.10": true,
}

func getAddr(addrs map[string]bool) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	al := []string{}
	for k, v := range addrs {
		if v {
			al = append(al, k)
		}
	}
	sort.Strings(al)
	log.Println(al)
	return al[r.Intn(len(al))]
}

func Test_tcm(t *testing.T) {
	getAddr(addrs)
	for i := 0; i < 100; i++ {
		//log.Println(getAddr(addrs))
	}
	/*
		queue := make(chan int, 10)
		select {
		case <-queue:
			log.Println("000000000000000000000")
		case <-time.After(1000000000):
			log.Println("1111111111111111111111")
		}
	*/
}

/*
func Test_delete(t *testing.T) {
	kv := map[string]string{
		"a": "11111111111111111111",
		"b": "22222222222222222222",
	}
	delete(kv, "c")
	log.Println(kv)
	tmp := kv["c"]
	log.Println(tmp)
	log.Println("00000000000000000000000")
	lis := []int{1, 2, 3, 4}
	for i := range lis {
		log.Println(i)
	}

	ti := int64(30 * time.Second)
	log.Println("time:", ti)
}
*/
