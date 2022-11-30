package defaults

import "log"

const (
	TickMin = 3 // minimal time in seconds between that can be set in config file. If not set or less then, Set to this value.

	// Global
	AppName     = "scom"
	SiteConfDir = "/etc/" + AppName + "/"

	// Job from templates tab
	TemplatesDir        = SiteConfDir + "templates"
	TemplatesSuffix     = ".sbatch"
	TemplatesDescSuffix = ".desc"

	// logging
	LogFlag = log.Lshortfile | log.Lmicroseconds

	// SlurmCommander
	SCAppName      = AppName
	SCLogFile      = "scdebug.log"
	SCConfFileName = "scom.conf"
	SCSiteConfFile = SiteConfDir + SCConfFileName

	// SlurmCommanderCache
	SccAppName      = "sccache"
	SccLogFile      = "sccdebug.log"
	SccConfFileName = "scc.conf"
	SccSiteConfFile = SiteConfDir + SccConfFileName
	SccPrefix       = "/bin"

	SccRefreshT = 5
	SccPort     = 10237
)

var (
	// default paths
	BinPaths = map[string]string{
		"sacct":    "/bin/sacct",
		"sstat":    "/bin/sstat",
		"sinfo":    "/bin/sinfo",
		"squeue":   "/bin/squeue",
		"sbatch":   "/bin/sbatch",
		"scancel":  "/bin/scancel",
		"scontrol": "/bin/scontrol",
		"sacctmgr": "/bin/sacctmgr",
	}

	SqueueCmdSwitches = []string{"-a", "--json"}
)
