package runners

import (
	"os"

	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/logger"
)

// importFile will take a Request object and attempt to load the data from disk
// and unmarshal it into `ptr`
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
