package store

import (
	"sync"

	"log"

	"fmt"

	minio "github.com/minio/minio-go"
	"github.com/open-fightcoder/oj-judger/common/g"
)

var MinioClient *minio.Client
var onceExec sync.Once

func InitMinio() {
	onceExec.Do(func() {
		cfg := g.Conf()
		fmt.Println(cfg)
		var err error
		MinioClient, err = minio.New(cfg.Minio.Endpoint, cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, cfg.Minio.Secure)
		if err != nil {
			log.Fatalln("fail to connect minio", g.Conf().Minio.Endpoint, err)
		}
	})
}
