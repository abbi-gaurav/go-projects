package intermediate

import (
	lowLevel "github.com/abbi-gaurav/go-learning-projects/ultimate-go-programming/ccp-in-go/conc-at-scale/error-handling/low-level"
	"github.com/abbi-gaurav/go-learning-projects/ultimate-go-programming/ccp-in-go/conc-at-scale/error-handling/model"
	"os/exec"
)

type InterimError struct {
	error
}

func RunJob(id string) error {
	const jonBinPath = "/bad/binary"
	isExecutable, err := lowLevel.IsGloballyExec(jonBinPath)
	if err != nil {
		return InterimError{model.WrapError(
			err,
			"cannot run job %s, binaries not available",
			id)}
	} else if !isExecutable {
		return InterimError{model.WrapError(
			nil,
			"cannot run job: %s, binaries not executable",
			id)}
	}

	return exec.Command(jonBinPath, "--id="+id).Run()
}
