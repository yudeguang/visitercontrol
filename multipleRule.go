package visitercontrol

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

//用于NewMultipleVisitercontrol函数的初始化
type Rule struct {
	DefaultExpiration            time.Duration
	NumberOfAllowedAccesses      int
	EstimatedNumberOfOnlineUsers int
}
type MultipleVisitercontrol struct {
	multipleVisitercontrol []*SingleVisitercontrol
}

/*
初始化一个多重规则的频率控制策略
例：
multipleVc := NewMultipleVisitercontrol(Rule{time.Minute * 5, 20, 100}, Rule{time.Minute * 30, 50, 1000}, Rule{time.Hour * 24, 200, 10000})
它表示:
在5分钟内每个用户最多允许访问20次,并且我们预计在这5分钟内大致有100个用户会访问我们的网站
在30分钟内每个用户最多允许访问50次,并且我们预计在这30分钟内大致有1000个用户会访问我们的网站
在24小时内每个用户最多允许访问200次,并且我们预计在这24小时内大致有10000个用户会访问我们的网站
以上任何一条规则的访问次数超出，都不允许访问,并且我们在匹配规则时，是按时间段从小到大匹配
*/
func NewMultipleVisitercontrol(multipleRule ...Rule) *MultipleVisitercontrol {
	//把时间控制调整为从小到大排列，防止用户在实例化的时候，未按照预期的时间顺序添加，导致某些规则失效
	sort.Slice(multipleRule, func(i int, j int) bool { return multipleRule[i].DefaultExpiration < multipleRule[j].DefaultExpiration })
	var m []*SingleVisitercontrol
	for i := range multipleRule {
		m = append(m, NewSingleVisitercontrol(multipleRule[i].DefaultExpiration, multipleRule[i].NumberOfAllowedAccesses, multipleRule[i].EstimatedNumberOfOnlineUsers))
	}
	return &MultipleVisitercontrol{m}
}

//是否允许访问,允许访问则往访问记录中加入一条访问记录
//例: AllowVisit("usernameexample")
func (this *MultipleVisitercontrol) AllowVisit(key interface{}) bool {
	for i := range this.multipleVisitercontrol {
		if err := this.multipleVisitercontrol[i].add(key); err != nil {
			return false
		}
	}
	return true
}

//是否允许某IP的用户访问
//例: AllowVisitIP("127.0.0.1")
func (this *MultipleVisitercontrol) AllowVisitIP(ip string) bool {
	ipInt64 := Ip4StringToInt64(ip)
	if ipInt64 == 0 {
		return false
	}
	return this.AllowVisit(ipInt64)
}

//当前在线用户总数,每个时间段的在线用户总数
func (this *MultipleVisitercontrol) CurOnlineUserNum() []int {
	arr := make([]int, 0, len(this.multipleVisitercontrol))
	for i := range this.multipleVisitercontrol {
		arr = append(arr, this.multipleVisitercontrol[i].CurOnlineUserNum())
	}
	return arr
}

//剩余访问次数
//例: RemainingVisits("usernameexample")
func (this *MultipleVisitercontrol) RemainingVisits(key interface{}) []int {
	arr := make([]int, 0, len(this.multipleVisitercontrol))
	for i := range this.multipleVisitercontrol {
		arr = append(arr, this.multipleVisitercontrol[i].RemainingVisits(key))
	}
	return arr
}

//某IP剩余访问次数
//例: RemainingVisitsIP("127.0.0.1")
func (this *MultipleVisitercontrol) RemainingVisitsIP(ip string) []int {
	ipInt64 := Ip4StringToInt64(ip)
	if ipInt64 == 0 {
		return []int{}
	}
	return this.RemainingVisits(ipInt64)
}

//把在线用户数据转化成JSON输出,数据较大时，最好不要使用
func (this *MultipleVisitercontrol) OnlineUserInfoToJson() string {
	var pp []printHelper
	for i := range this.multipleVisitercontrol {
		var p printHelper
		p.DefaultExpiration = this.multipleVisitercontrol[i].defaultExpiration.String()
		p.NumberOfAllowedAccesses = this.multipleVisitercontrol[i].numberOfAllowedAccesses
		p.EstimatedNumberOfOnlineUsers = this.multipleVisitercontrol[i].estimatedNumberOfOnlineUsers
		var CurOnlineUserInfo []userInfo
		this.multipleVisitercontrol[i].indexes.Range(func(k, v interface{}) bool {
			this.multipleVisitercontrol[i].lock.Lock() //range里面不能用defer
			index := v.(int)
			//防止越界出错，理论上不存在这种情况
			if index < len(this.multipleVisitercontrol[i].visitorRecords) && index >= 0 {
				this.multipleVisitercontrol[i].visitorRecords[index].DeleteExpired()
				//删除完过期数据之后，如果该用户的所有访问记录均过期了，那么就删除该用户
				//并把该空间返还给notUsedVisitorRecordsIndex以便下次重复使用
				if this.multipleVisitercontrol[i].visitorRecords[index].UsedSize() == 0 {
					this.multipleVisitercontrol[i].indexes.Delete(k)
					this.multipleVisitercontrol[i].notUsedVisitorRecordsIndex.Add(index)
				} else {
					//加入统计数据表
					var u userInfo
					u.UserName = fmt.Sprint(k)
					u.RemainingVisits = this.multipleVisitercontrol[i].visitorRecords[index].UnUsedSize()
					u.Used = this.multipleVisitercontrol[i].visitorRecords[index].UsedSize()
					CurOnlineUserInfo = append(CurOnlineUserInfo, u)
				}
			} else {
				this.multipleVisitercontrol[i].indexes.Delete(k)
			}
			this.multipleVisitercontrol[i].lock.Unlock()
			return true
		})
		sort.Slice(CurOnlineUserInfo, func(i int, j int) bool { return CurOnlineUserInfo[i].UserName < CurOnlineUserInfo[j].UserName })
		p.CurOnlineUserInfo = CurOnlineUserInfo
		p.CurOnlineUserNum = len(p.CurOnlineUserInfo)
		pp = append(pp, p)
	}
	b, _ := json.Marshal(pp)
	return string(b)
}
