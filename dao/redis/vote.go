package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote float64 = 432
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated = errors.New("不允许重复投票")

)

func CreatePost(postID,communityID int64) error  {
	pipeline := client.TxPipeline()
	// 帖子时间
	 pipeline.ZAdd(getRedisKey(keyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(keyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 补充：把贴 id 加入到社区的 set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey,postID)
	_, err := pipeline.Exec()
	return  err

}

func VoteForPost(userID, postID string, value float64) error {
	// 1、判断投票的限制
	postTime := client.ZScore(getRedisKey(keyPostTimeZSet),postID).Val()
	if float64(time.Now().Unix()) - postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2和3需要放到一个pipeline事务中操作
	// 2、更新帖子的分数
	// 先查当前用户给当前的帖子的投票记录
	ov := client.ZScore(getRedisKey(keyPostVotedZSetPF+postID),userID).Val()

	// 更新： 如果这一次投票的值和之前保持的值一致，就提示不允许重复投票
	if value == ov {
		return  ErrVoteRepeated
	}

	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}

	diff := math.Abs(ov-value)  // 计算两次投票的差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(keyPostScoreZSet),op * diff * scorePerVote,postID)
	// 3、 记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(keyPostVotedZSetPF+postID), postID)
	} else {

		pipeline.ZAdd(getRedisKey(keyPostVotedZSetPF+postID),redis.Z{
			Score:  value,  // 赞成票还是反对票
			Member: nil,
		})
	}
	_, err := pipeline.Exec()
	return  err
}