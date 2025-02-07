package runners

import (
	"os"

	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/logger"
)

func importFile(in Request, ptr any) error {
	logger.Trace()

	path, err := NormalizePath(in)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	utils.UnmarshalData(b, ptr)

	return nil
}
