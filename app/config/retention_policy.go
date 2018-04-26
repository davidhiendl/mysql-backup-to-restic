package config

type RetentionPolicy struct {
	KeepLast    int      `yaml:"keepLast,omitempty"`
	KeepHourly  int      `yaml:"keepHourly,omitempty"`
	KeepDaily   int      `yaml:"keepDaily,omitempty"`
	KeepWeekly  int      `yaml:"keepWeekly,omitempty"`
	KeepMonthly int      `yaml:"keepMonthly,omitempty"`
	KeepYearly  int      `yaml:"keepYearly,omitempty"`
	KeepTags    []string `yaml:"keepTags,omitempty"`
	Prune       bool     `yaml:"prune,omitempty"`
	Check       bool     `yaml:"check,omitempty"`
	DryRun      bool     `yaml:"dryRun,omitempty"`
}

// check if any keep policy is set
func (rp *RetentionPolicy) HasKeepDirective() bool {

	if rp.KeepLast > 0 {
		return true
	}

	if rp.KeepHourly > 0 {
		return true
	}

	if rp.KeepDaily > 0 {
		return true
	}

	if rp.KeepWeekly > 0 {
		return true
	}

	if rp.KeepMonthly > 0 {
		return true
	}

	if rp.KeepYearly > 0 {
		return true
	}

	if len(rp.KeepTags) > 0 {
		return true
	}

	return false
}
