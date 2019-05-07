package low_level

import (
	"github.com/abbi-gaurav/go-learning-projects/ultimate-go-programming/ccp-in-go/conc-at-scale/error-handling/model"
	"os"
)

type LowLevelErr struct {
	error
}

func IsGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)

	if err != nil {
		return false, LowLevelErr{model.WrapError(err, err.Error())}
	}

	return info.Mode().Perm()&0100 == 0100, nil
}
