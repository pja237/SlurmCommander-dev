package command

var (
	SacctJobCmdSwitches    = []string{"-n", "--json", "-j"}
	ScancelJobCmdSwitches  = []string{}
	SholdJobCmdSwitches    = []string{"hold"}
	SrequeueJobCmdSwitches = []string{"requeue"}
	SbatchCmdSwitches      = []string{}
	SacctmgrCmdSwitches    = []string{"list", "Association", "format=account", "-P", "-n"}
)
