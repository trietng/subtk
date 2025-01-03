package resource

import (
	"fmt"
	"os"
	"sync"
	"trietng/subtk/common"
)

// organization of resources
// resources
// ├── resource_group_1
// │   ├── resource_1
// │   └── resource_2
// └── resource_group_2
// ...
// └── metadata
// metadata: a binary file that stores metadata of resources
// resource_group_<n>: a directory that contains resources of the same type
// resource: a resource file
var (
	// metadata map: map[resource_group][resource]timestamp
	metadata map[string]map[string]int64
	lock     sync.Mutex
)

func LoadExternalMetadata() (common.CommonSignal, error) {
	lock.Lock()
	defer lock.Unlock()
	if metadata != nil {
		return common.SUBTK_DUPLICATED, nil
	}
	// create .subtk/resources directory if it does not exist
	home, _ := os.UserHomeDir()
	os.MkdirAll(fmt.Sprintf("%s/.subtk/resources", home), 0755)
	return common.SUBTK_INIT, nil
}