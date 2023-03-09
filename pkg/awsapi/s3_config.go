package awsapi

import (
	"os"

	"github.com/spf13/pflag"

	"zq-xu/warehouse-admin/pkg/config"
	"zq-xu/warehouse-admin/pkg/log"
)

const (
	awsS3ConfigName = "AWSS3Config"

	AWSS3ACLEnv = "AWSS3ACL"
)

var (
	S3Cfg = &S3Config{}
)

type S3Config struct {
	// public read for storagegateway.ObjectACLPublicRead
	ACL string
}

func InitS3Config() {
	config.RegisterCfg(awsS3ConfigName, S3Cfg)
}

// AddFlags adds flags for S3Config
func (sc *S3Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&sc.ACL, "aws-s3-acl", os.Getenv(AWSS3ACLEnv), "the ACL.")
}

func (sc *S3Config) Revise() {
	log.Logger.Infof("AWS S3Config is ACL: %v", sc.ACL)
}
