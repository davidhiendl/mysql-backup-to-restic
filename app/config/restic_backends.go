package config

type ResticBackends struct {
	// Local ResticBackendLocal `yaml:"local,omitempty"`
	// SFTP  ResticBackendSFTP  `yaml:"sftp,omitempty"` // unsupported due to difficult auth
	S3 ResticBackendS3 `yaml:"s3,omitempty"`
}

// type ResticBackendLocal struct {
// 	Path string `yaml:"path,omitempty"`
// }

type ResticBackendS3 struct {
	AccessKeyID     string `yaml:"accessKeyID,omitempty"`
	SecretAccessKey string `yaml:"secretAccessKey,omitempty"`
	Endpoint        string `yaml:"endpoint,omitempty"`
	Bucket          string `yaml:"bucket,omitempty"`
	Path            string `yaml:"path,omitempty"`
	ForceHTTP       bool   `yaml:"forceHTTP,omitempty"`
	Port            int    `yaml:"port,omitempty"`
}

func (b *ResticBackendS3) Enabled() bool {
	return len(b.AccessKeyID) > 0 && len(b.SecretAccessKey) > 0
}
