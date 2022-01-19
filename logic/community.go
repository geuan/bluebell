package logic

import (
	"bluebell/dao/mysql"
	"bluebell/modules"
)

func GetCommunityList() ([]*modules.Community,error) {
	//查询数据库 查找到所有的 community 并返回
	return mysql.GetCommunityList()
}


func GetCommunityDetail(id int64) (*modules.CommunityDetail,error) {
	return mysql.GetCommunityDetailByID(id)
}