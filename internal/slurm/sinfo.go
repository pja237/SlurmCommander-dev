package slurm

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pja237/slurmcommander-dev/internal/openapi"
)

type SinfoJSON struct {
	Nodes []openapi.V0039Node
}

var gpuGresPattern = regexp.MustCompile(`^gpu\:([^\:]+)\:?(\d+)?`)

func ParseGRES(line string) *int {
	value := 0

	gres := strings.Split(line, ",")
	for _, g := range gres {
		if !strings.HasPrefix(g, "gpu:") {
			continue
		}

		matches := gpuGresPattern.FindStringSubmatch(g)
		if len(matches) == 3 {
			if matches[2] != "" {
				value, _ = strconv.Atoi(matches[2])
			} else {
				value, _ = strconv.Atoi(matches[1])
			}
		}
	}

	return &value
}
