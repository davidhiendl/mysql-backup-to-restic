package config

import "errors"

func (cfg *Config) Validate() []error {
	errs := make([]error, 0)

	if len(cfg.Restic.Password) <= 0 {
		errs = append(errs, errors.New("missing restic password"))
	}

	// TODO much more validation rules to fail immediately

	return errs
}
