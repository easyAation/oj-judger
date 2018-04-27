package store

import (
	"testing"

	"strings"

	"github.com/minio/minio-go"
	"github.com/open-fightcoder/oj-judger/common/g"
)

func TestMinio(t *testing.T) {
	g.LoadConfig("../../cfg/cfg.toml.debug")
	InitMinio()
	_, err := MinioClient.PutObject("ttt", "test.txt", strings.NewReader("aaaaaassssssdddd"), -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		panic(err)
	}
}
