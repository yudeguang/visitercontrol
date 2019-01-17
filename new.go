package visitercontrol

import (
	"github.com/yudeguang/hashset"
	"sync"
	"time"
)

//某单位时间内允许多少次访问
type Visitercontrol struct {
	ruleName                   string              //规则名称
	defaultExpiration          time.Duration       //每条访问记录需要保存的时长，也就是过期时间
	cleanupInterval            time.Duration       //多长时间需要执行一次清除操作
	maxVisitsNum               int                 //每个用户在相应时间段内最多允许访问的次数
	indexes                    sync.Map            //索引：key代表用户名或IP；value代表visitorRecords中的索引位置
	maximumNumberOfOnlineUsers int                 //单位时间最大用户数量，建议选用一个稍大于实际值的值，以减少内存分配次数
	visitorRecords             []*circleQueueInt64 //存储用户访问记录
	notUsedVisitorRecordsIndex *hashset.SetInt     //对应visitorRecords中未使用的数据的索引位置
	lock                       *sync.Mutex         //锁
}

/*
初始化
例：
vc := visitercontrol.New("every 30 minute",time.Minute*30, time.Second*5, 50, 1000)
它表示:
在30分钟内每个用户最多允许访问50次，系统每5秒针删除一次过期数据。
并且我们预计同时在线用户数量大致在1000个左右。
*/
func New(ruleName string, defaultExpiration, cleanupInterval time.Duration, maxVisitsNum, maximumNumberOfOnlineUsers int) *Visitercontrol {
	vc := createVisitercontrol(ruleName, defaultExpiration, cleanupInterval, maxVisitsNum, maximumNumberOfOnlineUsers)
	go vc.deleteExpired()
	return vc
}

func createVisitercontrol(ruleName string, defaultExpiration, cleanupInterval time.Duration, maxVisitsNum, maximumNumberOfOnlineUsers int) *Visitercontrol {
	if cleanupInterval > defaultExpiration {
		panic("每次清除访问记录的时间间隔(cleanupInterval)必须小于待统计数据时间段(defaultExpiration)")
	}
	var vc Visitercontrol
	var lock sync.Mutex
	vc.ruleName = ruleName
	vc.defaultExpiration = defaultExpiration
	vc.cleanupInterval = cleanupInterval
	vc.maxVisitsNum = maxVisitsNum
	vc.maximumNumberOfOnlineUsers = maximumNumberOfOnlineUsers
	vc.notUsedVisitorRecordsIndex = hashset.NewInt()
	vc.lock = &lock
	//根据在线用户数量初始化用户访问记录数据
	vc.visitorRecords = make([]*circleQueueInt64, vc.maximumNumberOfOnlineUsers)
	for i := range vc.visitorRecords {
		vc.visitorRecords[i] = newCircleQueueInt64(vc.maxVisitsNum)
		vc.notUsedVisitorRecordsIndex.Add(i)
	}
	return &vc

}
