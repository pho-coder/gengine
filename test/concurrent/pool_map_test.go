package concurrent

import (
	"fmt"
	"gengine/builder"
	"gengine/context"
	"gengine/engine"
	"sync"
	"testing"
	"time"
)

const map_rules = `
rule "1" "2"
begin

//a = GetRanking()
//result.Add("3",a)
a =  GetRankingMap()
result.SidResult = a
Check(result.SidResult, request.Uid)
end
rule "2" "2"
begin

//a = GetRanking()
//result.Add("3",a)
a =  GetRankingMap()
result.SidResult = a
Check(result.SidResult, request.Uid)
end
`

func Check(b map[string]*Ranking, x int64) {
	println("check-----")
	if b == nil {
		println("b is nil")
	}
	//time.Sleep(time.Duration(x) * time.Second)
}

type Result struct {
	sync.Mutex
	SidResult map[string]*Ranking
}

func (r *Result) Add(sid string, an *Ranking) {
	r.Lock()
	println("sid=", sid)
	if an == nil {
		println("an==nil")
	} else {
		println("len->", len(an.Sl))
	}

	r.SidResult[sid] = an
	r.Unlock()
}

func GetRanking() *Ranking {
	return &Ranking{Sl: []int64{1, 2, 3, 4, 5}}
}

func GetRankingMap() map[string]*Ranking {
	x := make(map[string]*Ranking)
	x["3"] = &Ranking{Sl: []int64{1, 2, 3, 4, 5}}
	x["4"] = &Ranking{Sl: []int64{1, 2, 3, 4, 5}}
	return x
}

type Ranking struct {
	Sl []int64
}

type Request struct {
	Uid int64
}

/*func Test_map_base(t *testing.T) {

	go func() {
		defer fmt.Println("hello")
		defer fmt.Println("world")
	}()
}

*/

func Test_map_bbb(t *testing.T) {

	dataContext1 := context.NewDataContext()
	dataContext1.Add("println", fmt.Println)
	dataContext1.Add("GetRanking", GetRanking)
	dataContext1.Add("GetRankingMap", GetRankingMap)
	dataContext1.Add("Check", Check)
	ruleBuilder1 := builder.NewRuleBuilder(dataContext1)
	e1 := ruleBuilder1.BuildRuleFromString(map_rules)
	if e1 != nil {
		panic(e1)
	}
	en1 := engine.NewGengine()

	//dc := *dataContext1
	//dataContext2 := &dc

	dataContext2 := context.NewDataContext()
	dataContext2.Add("println", fmt.Println)
	dataContext2.Add("GetRanking", GetRanking)
	dataContext2.Add("GetRankingMap", GetRankingMap)
	dataContext2.Add("Check", Check)
	ruleBuilder2 := builder.NewRuleBuilder(dataContext2)
	e2 := ruleBuilder2.BuildRuleFromString(map_rules)
	if e2 != nil {
		panic(e2)
	}
	en2 := engine.NewGengine()

	go func() {
		//for {
		request := &Request{Uid: 1}
		result := &Result{SidResult: make(map[string]*Ranking)}
		ruleBuilder1.Dc.Add("result", result)
		ruleBuilder1.Dc.Add("request", request)
		en1.ExecuteSelectedRulesConcurrent(ruleBuilder1, []string{"1"})
		if result == nil {
			println("result is nil")
		}
		var cache []int64

		for k, v := range result.SidResult {
			//time.Sleep(100*time.Nanosecond)
			println("1k=", k, fmt.Sprintf("v=%+v", v))
		}

		if result.SidResult == nil {
			println("yes1_1")
		}

		if result.SidResult["3"] == nil {
			println("yes1_2")
		}

		for i := 0; i < len(result.SidResult["3"].Sl); i++ {
			cache = append(cache, result.SidResult["3"].Sl[i])
		}

		//	}
	}()

	go func() {
		//for {
		request := &Request{Uid: 1}
		result := &Result{SidResult: make(map[string]*Ranking)}
		ruleBuilder2.Dc.Add("result", result)
		ruleBuilder2.Dc.Add("request", request)
		en2.ExecuteSelectedRulesConcurrent(ruleBuilder2, []string{"1"})
		if result == nil {
			println("result is nil")
		}
		var cache []int64

		for k, v := range result.SidResult {
			//time.Sleep(100*time.Nanosecond)
			println("1k=", k, fmt.Sprintf("v=%+v", v))
		}

		if result.SidResult == nil {
			println("yes1_1")
		}

		if result.SidResult["3"] == nil {
			println("yes1_2")
		}

		for i := 0; i < len(result.SidResult["3"].Sl); i++ {
			cache = append(cache, result.SidResult["3"].Sl[i])
		}
		//}
	}()

	time.Sleep(10 * time.Second)

}

//bad case
func Test_map_conc(t *testing.T) {
	//init
	apis := make(map[string]interface{})
	apis["GetRanking"] = GetRanking
	apis["GetRankingMap"] = GetRankingMap
	apis["Check"] = Check
	pool, e1 := engine.NewGenginePool(1, 2, 1, map_rules, apis)
	if e1 != nil {
		println(fmt.Sprintf("e1: %+v", e1))
	}

	go func() {
		for {
			request := &Request{Uid: 1}
			data := make(map[string]interface{})
			data["request"] = request

			result := &Result{SidResult: make(map[string]*Ranking)}
			data["result"] = result

			e2 := pool.ExecuteSelectedRulesConcurrentWithMultiInput(data, []string{"1"})
			if e2 != nil {
				panic(e2)
			}
			if result == nil {
				println("result is nil")
			}
			var cache []int64

			for k, v := range result.SidResult {
				//time.Sleep(100*time.Nanosecond)
				println("1k=", k, fmt.Sprintf("v=%+v", v))
			}

			if result.SidResult == nil {
				println("yes1_1")
			}

			if result.SidResult["3"] == nil {
				println("yes1_2")
			}

			for i := 0; i < len(result.SidResult["3"].Sl); i++ {
				cache = append(cache, result.SidResult["3"].Sl[i])
			}
		}
	}()

	go func() {
		for {
			request := &Request{Uid: 2}
			data := make(map[string]interface{})
			data["request"] = request
			data["println"] = fmt.Println

			result := &Result{SidResult: make(map[string]*Ranking)}
			data["result"] = result

			e2 := pool.ExecuteSelectedRulesConcurrentWithMultiInput(data, []string{"2"})
			if e2 != nil {
				panic(e2)
			}

			if result == nil {
				println("result is nil")
			}

			var cache []int64
			for k, v := range result.SidResult {
				//time.Sleep(200*time.Nanosecond)
				println("1k=", k, fmt.Sprintf("v=%+v", v))
			}

			if result.SidResult == nil {
				println("yes2_1")
			}

			if result.SidResult["3"] == nil {
				println("yes2_2")
			}

			for i := 0; i < len(result.SidResult["3"].Sl); i++ {
				cache = append(cache, result.SidResult["3"].Sl[i])
			}
		}

	}()

	time.Sleep(15 * time.Second)

}
