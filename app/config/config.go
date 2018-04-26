package config

type Config struct {
	Common          Common          `yaml:"common"`
	MySQL           MySQL           `yaml:"mysql"`
	Dump            Dump            `yaml:"dump"`
	Databases       Databases       `yaml:"databases"`
	Restic          Restic          `yaml:"restic"`
	RetentionPolicy RetentionPolicy `yaml:"retentionPolicy"`
}

type Common struct {
	ScratchDir string `yaml:"scratchDir,omitempty"`
}

type MySQL struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
}

type Dump struct {
	CompressWithGz bool `yaml:"compressWithGz,omitempty"`
}

type Databases struct {
	ExcludeSystem bool     `yaml:"excludeSystem,omitempty"`
	Exclude       []string `yaml:"exclude,omitempty"`
	Include       []string `yaml:"include,omitempty"`
}

type Restic struct {
	Hostname    string         `yaml:"hostname,omitempty"`
	Password    string         `yaml:"password,omitempty"`
	Backends    ResticBackends `yaml:"backends,omitempty"`
	CacheEnable bool           `yaml:"cacheEnable,omitempty"`
}
