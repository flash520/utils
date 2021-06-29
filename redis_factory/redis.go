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
				Addrs:        nodes,
				DialTimeout:  10 * time.Second,
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,

				MaxRedirects: 8,

				PoolSize:           10,
				PoolTimeout:        30 * time.Second,
				IdleTimeout:        time.Minute,
				IdleCheckFrequency: 100 * time.Millisecond,
			})
		} else {
			mode = 1
			singleClient = redis.NewClient(&redis.Options{
				Addr:        nodes[0],
				DialTimeout: 5 * time.Millisecond,
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
