/**
 * @Author: koulei
 * @Description:
 * @File: mysql_test
 * @Version: 1.0.0
 * @Date: 2021/8/17 12:06
 */

package mysql

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

type OrderCreateModel struct {
	gorm.Model   `json:"-"`
	UserID       string  `json:"userid" gorm:"type:string;size:19;index;comment:学生ID"`
	ParentID     string  `json:"parentid" gorm:"type:string;size:20;comment:家长ID"`
	Username     string  `json:"username" gorm:"type:string;size:20;comment:学生名字"`
	Mobile       string  `json:"mobile" gorm:"type:string;size:20;comment:家长手机号码"`
	UserToken    string  `json:"user_token" gorm:"type:string;size:50;comment:用户令牌"`
	OrderToken   string  `json:"order_token" gorm:"type:string;size:50;comment:订购标识令牌"`
	CardCode     string  `json:"card_code" gorm:"type:string;size:20;comment:学生卡ID"`
	Role         string  `json:"role" gorm:"type:string;size:10;comment:角色,1学生 2家长 3教师 4管理员"`
	OpType       int     `json:"optype" gorm:"type:string;size:5;comment:操作类型,订购为0,变更为1"`
	AppLifecycle int     `json:"applifecycle" gorm:"type:string;size:2;comment:0体验期,1正式期"`
	ClassID      int     `json:"class_id" gorm:"type:string;size:10;comment:班级ID"`
	Classname    string  `json:"class_name" gorm:"type:string;size:100;comment:班级名称"`
	GradeID      string  `json:"grade_id" gorm:"type:string;size:10;comment:年级ID"`
	GradeName    string  `json:"grade_name" gorm:"type:string;size:100;comment:年级名称"`
	SchoolID     string  `json:"school_id" gorm:"type:string;size:10;comment:学校ID"`
	SchoolName   string  `json:"school_name" gorm:"type:string;size:100;comment:学校名称"`
	CountyID     string  `json:"countyid" gorm:"type:string;size:10;comment:地区编号"`
	CityID       string  `json:"cityid" gorm:"type:string;size:6;comment:城市编号"`
	AppID        string  `json:"appid" gorm:"type:string;size:10;comment:应用ID"`
	AppName      string  `json:"appname" gorm:"type:string;size:30;comment:应用名称"`
	ServiceID    string  `json:"serviceid" gorm:"type:string;size:22;comment:服务ID"`
	ServiceName  string  `json:"servicename" gorm:"type:string;size:32;comment:服务名称"`
	Fee          float64 `json:"fee" gorm:"type:string;size:4;comment:实际支付费用"`
	FeeType      string  `json:"feetype" gorm:"type:string;size:2;comment:0为一次性,1为包月，2为包年"`
	BeginTime    string  `json:"begintime" gorm:"type:time;size:14;default:NULL;comment:订单生效时间"`
	CreateTime   string  `json:"createtime" gorm:"type:time;size:14;default:NULL;comment:订单创建时间"`
	EndTime      string  `json:"endtime" gorm:"type:time;size:14;default:NULL;comment:订单到期时间"`
}

var sess = dbc.Session(&gorm.Session{SkipDefaultTransaction: true})

func Select() {
	var order int64
	sess.Table("order_create_model").Select("id").Where("id = ?", 10).Find(&order)
}
func SelectTrans() {
	var order int64
	dbc.Table("order_create_model").Select("id").Where("id = ?", 10).Find(&order)
}

func Update() {

}

func Insert() {

}

var dbc = CreateMysql("localhost", "root", "123456", "scmebap", "info").GetConn()

func TestSelect(t *testing.T) {
	var order OrderCreateModel
	dbc.Model(OrderCreateModel{}).Where("id = ?", 1).First(&order)
	fmt.Println(order)
}

func BenchmarkSelect(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Select()
		}
	})
}
