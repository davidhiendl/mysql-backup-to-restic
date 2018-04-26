package restic

import (
	"path/filepath"
	"strings"
	"fmt"
	"os"
	"github.com/sirupsen/logrus"
)

const (
	RESTIC_REPOSITORY = "RESTIC_REPOSITORY"
	RESTIC_PASSWORD   = "RESTIC_PASSWORD"
	TMPDIR            = "TMPDIR"

	AWS_ACCESS_KEY_ID     = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"
)

func (rt *Restic) configure() (error) {

	cfg := rt.cfg

	rt.sh.SetEnv(RESTIC_PASSWORD, cfg.Restic.Password)

	tmpDir := filepath.Join(rt.cfg.Common.ScratchDir, "restic-tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	rt.sh.SetEnv(TMPDIR, tmpDir)

	if cfg.Restic.Backends.S3.Enabled() {
		backend := &cfg.Restic.Backends.S3

		prefix := strings.TrimPrefix(filepath.Join(backend.Bucket, backend.Path), "/")
		repo := fmt.Sprintf("s3:%s/%s", backend.Endpoint, prefix)

		rt.sh.SetEnv(RESTIC_REPOSITORY, repo)
		rt.sh.SetEnv(AWS_ACCESS_KEY_ID, backend.AccessKeyID)
		rt.sh.SetEnv(AWS_SECRET_ACCESS_KEY, backend.SecretAccessKey)
	} else {
		logrus.Fatalf("no backend configured for Restic")
	}

	return nil
}

func (rt *Restic) DumpEnv() {
	out, err := rt.sh.Command("env").Output()
	if err != nil {
		logrus.Fatalf("error trying to debug environment: %+v", err)
	}
	logrus.Debugf("ENV:\n%s", string(out))
}
