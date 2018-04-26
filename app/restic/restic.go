package restic

import (
	"strconv"
	shell "github.com/codeskyblue/go-sh"
	"github.com/davidhiendl/mysql-backup-to-restic/app/config"
	"path/filepath"
	"strings"
	"time"
	"errors"
	"github.com/sirupsen/logrus"
	"os/exec"
)

const (
	RESTIC_BINARY     = "restic"
	FLAG_KEEP_LAST    = "--keep-last"
	FLAG_KEEP_HOURLY  = "--keep-hourly"
	FLAG_KEEP_DAILY   = "--keep-daily"
	FLAG_KEEP_WEEKLY  = "--keep-weekly"
	FLAG_KEEP_MONTHLY = "--keep-monthly"
	FLAG_KEEP_YEARLY  = "--keep-yearly"
	FLAG_KEEP_TAG     = "--keep-tag"
)

type Restic struct {
	cfg      *config.Config
	sh       *shell.Session
	hostname string
	exe      string
}

func New(cfg *config.Config) (*Restic) {
	w := Restic{
		sh:  shell.NewSession(),
		cfg: cfg,
	}

	// find binary
	binary, err := exec.LookPath(RESTIC_BINARY)
	if err != nil {
		logrus.Fatalf("failed to find restic binary: %v, err: %v \n", RESTIC_BINARY, err)
	}
	w.exe = binary

	w.sh.ShowCMD = true
	w.configure()
	w.DumpEnv()

	return &w
}

type Snapshot struct {
	ID       string    `json:"id"`
	Time     time.Time `json:"time"`
	Tree     string    `json:"tree"`
	Paths    []string  `json:"paths"`
	Hostname string    `json:"hostname"`
	Username string    `json:"username"`
	UID      int       `json:"uid"`
	Gid      int       `json:"gid"`
	Tags     []string  `json:"tags"`
}

func (rt *Restic) InitRepositoryIfAbsent() error {
	// TODO find a way to silence the error if repository is not yet initialized
	args := rt.appendCacheDirFlag([]interface{}{"snapshots", "--json"})
	err := rt.run(args)

	if err != nil {
		args = rt.appendCacheDirFlag([]interface{}{"init"})
		return rt.run(args)
	}

	return nil
}

func (rt *Restic) Backup(path string, tags []string) error {
	args := []interface{}{"backup", path, "--force"}

	// set hostname if overridden
	if rt.cfg.Restic.Hostname != "" {
		args = append(args, "--hostname")
		args = append(args, rt.cfg.Restic.Hostname)
	}

	// add tags if any
	for _, tag := range tags {
		args = append(args, "--tag")
		args = append(args, tag)
	}

	args = rt.appendCacheDirFlag(args)

	return rt.run(args)
}

func (rt *Restic) Forget(retentionPolicy *config.RetentionPolicy) error {
	// Get retentionPolicy for fileGroup, ignore if not found

	args := []interface{}{"forget"}

	if retentionPolicy.KeepLast > 0 {
		args = append(args, string(FLAG_KEEP_LAST))
		args = append(args, strconv.Itoa(retentionPolicy.KeepLast))
	}

	if retentionPolicy.KeepHourly > 0 {
		args = append(args, string(FLAG_KEEP_HOURLY))
		args = append(args, strconv.Itoa(retentionPolicy.KeepHourly))
	}

	if retentionPolicy.KeepDaily > 0 {
		args = append(args, string(FLAG_KEEP_DAILY))
		args = append(args, strconv.Itoa(retentionPolicy.KeepDaily))
	}

	if retentionPolicy.KeepWeekly > 0 {
		args = append(args, string(FLAG_KEEP_WEEKLY))
		args = append(args, strconv.Itoa(retentionPolicy.KeepWeekly))
	}

	if retentionPolicy.KeepMonthly > 0 {
		args = append(args, string(FLAG_KEEP_MONTHLY))
		args = append(args, strconv.Itoa(retentionPolicy.KeepMonthly))
	}

	if retentionPolicy.KeepYearly > 0 {
		args = append(args, string(FLAG_KEEP_YEARLY))
		args = append(args, strconv.Itoa(retentionPolicy.KeepYearly))
	}

	if len(retentionPolicy.KeepTags) > 0 {
		for _, tag := range retentionPolicy.KeepTags {
			args = append(args, string(FLAG_KEEP_TAG))
			args = append(args, tag)
		}
	}

	if retentionPolicy.Prune {
		args = append(args, "--prune")
	}

	if retentionPolicy.DryRun {
		args = append(args, "--dry-run")
	}

	if len(args) > 1 {
		args = rt.appendCacheDirFlag(args)
		return rt.run(args)
	}
	return nil
}

func (rt *Restic) Check() error {
	args := rt.appendCacheDirFlag([]interface{}{"check"})
	return rt.run(args)
}

func (rt *Restic) appendCacheDirFlag(args []interface{}) []interface{} {
	if rt.cfg.Restic.CacheEnable {
		cacheDir := filepath.Join(rt.cfg.Common.ScratchDir, "restic-cache")
		return append(args, "--cache-dir", cacheDir)
	}
	return append(args, "--no-cache")
}

func (rt *Restic) run(args []interface{}) error {
	out, err := rt.sh.Command(rt.exe, args...).CombinedOutput()
	if err != nil {
		logrus.Errorf("Error running command '%s %s' output:\n%s\n", rt.exe, args, string(out))
		parts := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
		if len(parts) > 1 {
			parts = parts[len(parts)-1:]
			return errors.New(parts[0])
		}
	}
	return err
}
