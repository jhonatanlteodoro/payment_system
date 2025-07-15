package ports

import "os"

type Worker interface {
	HasMessageToProcess() bool
	Process(shutdown chan os.Signal)
}
