package config

var Defaults = Config{
	Common{
		ScratchDir: "/tmp/mysqldump-to-restic",
	},
	MySQL{
		Port:     3306,
		Host:     "localhost",
		Username: "root",
	},
	Dump{
		CompressWithGz: true,
	},
	Databases{
		ExcludeSystem: true,
	},
	Restic{
		Backends: ResticBackends{
			S3: ResticBackendS3{
				Endpoint:  "s3.amazonaws.com",
				Bucket:    "restic",
				Path:      "restic",
				ForceHTTP: false,
				Port:      443,
			},
		},
	},
	RetentionPolicy{
		Prune:  true,
		Check:  true,
		DryRun: false,
	},
}
