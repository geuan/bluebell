package mysql

import (
	"bluebell/modules"
	"bluebell/settings"
	"testing"
)

// 单元测试就是自己与自己较劲的过程，习惯之后就可以提高开发效率

func init()  {
	dbCfg := settings.MySQLConfig{
		Host:         "106.54.119.58",
		User:         "root",
		Password:     "xc456789110",
		DB:           "bluebell",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}

	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}


}

func TestCreatePost(t *testing.T) {
	post := modules.Post{
		ID:          1,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a test",
	}

	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed,err:%v\n",err)
	}
	t.Logf("CreatePost insert record into mysql success")
}

