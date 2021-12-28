/**
 * @Author: koulei
 * @Description: TODO
 * @File:  redis_factory
 * @Version: 1.0.0
 * @Date: 2021/5/11 16:38
 */

package redis_factory

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var (
	clusterClient *redis.ClusterClient
	singleClient  *redis.Client
	mode          int
)

// CreateRedis 初始化redis，自动识别是单机和集群模式
func CreateRedis(nodes []string) *RedisClient {
	once := sync.Once{}
	once.Do(func() {
		if len(nodes) == 0 {
			log.Errorf("redis未配置,请重试...")
			os.Exit(1)
		} else if len(nodes) >= 6 {
			mode = 6
			clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: nodes,
				// 超时
				DialTimeout:  5 * time.Second, // 连接建立超时时间，默认5秒
				ReadTimeout:  3 * time.Second, // 读超时，默认3秒， -1表示取消读超时
				WriteTimeout: 3 * time.Second, // 写超时，默认等于读超时
				PoolTimeout:  4 * time.Second, // 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒
				// 闲置连接检查包括IdleTimeout，MaxConnAge
				IdleCheckFrequency: 60 * time.Second, // 闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理
				IdleTimeout:        time.Minute,      // 闲置超时，默认5分钟，-1表示取消闲置超时检查
				MaxConnAge:         0 * time.Second,  // 连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
				// 命令执行失败时的重试策略
				MaxRedirects:    0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
				MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
				MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

				PoolSize:     16,
				MinIdleConns: 10,
			})
		} else {
			mode = 1
			singleClient = redis.NewClient(&redis.Options{
				Addr:         nodes[0],
				PoolSize:     16,
				MinIdleConns: 10,
				DialTimeout:  5 * time.Millisecond,
			})
		}

		if mode == 1 {
			_, err := singleClient.Ping(context.Background()).Result()
			if err != nil {
				log.Errorf(err.Error())
			}
		}
		if mode == 6 {
			_, err := clusterClient.Ping(context.Background()).Result()
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	})
	return &RedisClient{}
}

type RedisClient struct {
}

// GetConn 获取redis连接
func (rc *RedisClient) GetConn() redis.Cmdable {
	switch mode {
	case 1:
		return rc.singleClient()
	case 6:
		return rc.clusterClient()
	default:
		log.Warnf("获取redis client错误: 无法识别的客户端模式\n")
		return nil
	}
}

func (rc *RedisClient) clusterClient() *redis.ClusterClient {
	return clusterClient
}

func (rc *RedisClient) singleClient() *redis.Client {
	return singleClient
}
