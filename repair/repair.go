package repair

import (
	"fmt"
	"os"
	"sync"
)

var (
	lock sync.Mutex
)

func ResetResources() {
	lock.Lock()
	defer lock.Unlock()
	home, _ := os.UserHomeDir()
	// delete all resources
	err := os.RemoveAll(fmt.Sprintf("%s/.subtk/resources", home))
	if err != nil {
		panic(err)
	}
}