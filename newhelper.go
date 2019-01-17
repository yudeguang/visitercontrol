package visitercontrol

import (
	"time"
)

//用于更加方便的生成对象
type newHelper struct {
}

type ruleNameHelper struct {
	ruleName string
}
type defaultExpirationHelper struct {
	ruleName          string
	defaultExpiration time.Duration
}

type cleanupIntervalHelper struct {
	ruleName          string
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

type maxVisitsNumHelper struct {
	ruleName          string
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	maxVisitsNum      int
}

type maximumNumberOfOnlineUsersHelper struct {
	ruleName                   string
	defaultExpiration          time.Duration
	cleanupInterval            time.Duration
	maxVisitsNum               int
	maximumNumberOfOnlineUsers int
}

/*
通过链式调用生成流量限制规则,为New()函数的替代品
因为New()函数中，传入的参数太多，很多时候记不清楚
各参数的顺序，通过对NewHelper()以及后续相应方法的
调用可以在一定程度上减少记忆量，减少出错的概率。
但同时也增加了代码的长度

vc := visitercontrol.New("every 30 minute",time.Minute*30, time.Second*5, 50, 1000)
vc := NewHelper().SetRuleName("every 30 minute").
	SetDefaultExpiration(time.Minute * 30).
	SetcleanupInterval(time.Second * 5).
	SetMaxVisitsNum(50).
	SetMaximumNumberOfOnlineUsers(1000).
	CreateVisitercontrol()
*/
func NewHelper() *newHelper {
	return new(newHelper)
}

//设置规则名
func (this *newHelper) SetRuleName(ruleName string) *ruleNameHelper {
	var v ruleNameHelper
	v.ruleName = ruleName
	return &v
}

//设置默认过期时间
func (this *ruleNameHelper) SetDefaultExpiration(defaultExpiration time.Duration) *defaultExpirationHelper {
	var v defaultExpirationHelper
	v.ruleName = this.ruleName
	v.defaultExpiration = defaultExpiration
	return &v
}

//设置清除过期时间的时间间隔，也就是每多少时间清除一次
func (this *defaultExpirationHelper) SetcleanupInterval(cleanupInterval time.Duration) *cleanupIntervalHelper {
	var v cleanupIntervalHelper
	v.ruleName = this.ruleName
	v.defaultExpiration = this.defaultExpiration
	v.cleanupInterval = cleanupInterval
	return &v
}

//设置默认时间段内，最多允许用户访问多少次
func (this *cleanupIntervalHelper) SetMaxVisitsNum(maxVisitsNum int) *maxVisitsNumHelper {
	var v maxVisitsNumHelper
	v.ruleName = this.ruleName
	v.defaultExpiration = this.defaultExpiration
	v.cleanupInterval = this.cleanupInterval
	v.maxVisitsNum = maxVisitsNum
	return &v
}

//设置在默认时间段内，预计最大在线用户数
func (this *maxVisitsNumHelper) SetMaximumNumberOfOnlineUsers(maximumNumberOfOnlineUsers int) *maximumNumberOfOnlineUsersHelper {
	var v maximumNumberOfOnlineUsersHelper
	v.ruleName = this.ruleName
	v.defaultExpiration = this.defaultExpiration
	v.cleanupInterval = this.cleanupInterval
	v.maxVisitsNum = this.maxVisitsNum
	v.maximumNumberOfOnlineUsers = maximumNumberOfOnlineUsers
	return &v
}

//
func (this *maximumNumberOfOnlineUsersHelper) CreateVisitercontrol() *Visitercontrol {
	return New(this.ruleName, this.defaultExpiration, this.cleanupInterval, this.maxVisitsNum, this.maximumNumberOfOnlineUsers)

}
