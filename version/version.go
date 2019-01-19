package version

import "fmt"

var (
	Version, GitCommit string
)

func BuildVersion() string {
	if len(Version) == 0 {
		Version = "dev"
	}
	return fmt.Sprintf("{\"version\":\"%s\", \"git_commit\":\"%s\"}", Version, GitCommit)
}
