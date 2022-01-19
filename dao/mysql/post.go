package mysql

import (
	"bluebell/modules"
	"github.com/jmoiron/sqlx"
	"strings"
)

// CreatePost 创建帖子
func CreatePost(p *modules.Post) (err error) {
	sqlStr := `insert into post(post_id,title,content,author_id,community_id) values (?,?,?,?,?)`
	 _,err = db.Exec(sqlStr,p.ID,p.Title,p.Content,p.AuthorID,p.CommunityID)
	return
}

// GetPostByID 根据id 查询单个帖子数据
func GetPostByID(pid int64) (post *modules.Post,err error) {
	post = new(modules.Post)
	sqlStr := `select post_id,title,content,author_id,community_id from post where post_id=?`
	err = db.Get(post,sqlStr,pid)
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page,size int64)(posts []*modules.Post, err error)  {
	sqlStr := `select
	post_id,title,content,author_id,community_id 
	from post
	ORDER BY create_time
	DESC
	limit ?,?
	`
	posts = make([]*modules.Post,0,2)
	err = db.Select(&posts,sqlStr,(page-1)*size,size)
	return

}


// GetPostListByIDs 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*modules.Post,err error)  {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id,?)
	`
	// https://www.liwenzhou.com/posts/Go/sqlx/
	query, args, err := sqlx.In(sqlStr,ids,strings.Join(ids,","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)

	err = db.Select(&postList,query,args...) // !!!
	return
}