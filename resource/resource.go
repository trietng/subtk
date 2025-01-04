package resource

import (
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"io"
	"math"
	"os"
	"sync"
	"time"
	"trietng/subtk/common"
	"trietng/subtk/resource/errmsg"
)

// NOTE: no versioning for metadata at the moment
type Metadata struct {
	DataType       string
	Timestamp      int64  // in seconds
	UpdateInterval int64  // in seconds
	Checksum       uint32 // CRC-32 checksum of the resource file
}

const (
	// default update interval for resources (1 month)
	UPDATE_DEFAULT = 2592000 // 30 days
	// value to use when the update interval is not required
	UPDATE_NOT_REQUIRED = math.MaxInt64
	// value to use when the update interval is zero
	UPDATE_IMMEDIATE = 0
)

// organization of resources
// resources
// ├── resource_group_1
// │   ├── resource_1
// │   └── resource_2
// └── resource_group_2
// ...
// └── metadata.gob
// metadata.gob: a binary file that stores metadata of resources and also acts as a lock file
// resource_group: a directory that contains resources of the same type
// resource: a resource file
var (
	// metadata map: map[resource_group][resource]metadata
	metadata map[string]map[string]Metadata
	lock     sync.Mutex
)

func newMetadata[T any](data *T, updateInterval int64) *Metadata {
	return &Metadata{
		DataType:       fmt.Sprintf("%T", *data),
		Timestamp:      time.Now().Unix(),
		UpdateInterval: updateInterval,
	}
}

func init() {
	_, err := LoadMetadata()
	if err != nil {
		panic(err)
	}
}

func LoadMetadata() (common.CommonSignal, error) {
	lock.Lock()
	defer lock.Unlock()
	if metadata != nil {
		return common.SUBTK_DUPLICATED, nil
	}
	// create .subtk/resources directory if it does not exist
	home, _ := os.UserHomeDir()
	os.MkdirAll(fmt.Sprintf("%s/.subtk/resources", home), 0755)
	file, err := os.OpenFile(fmt.Sprintf("%s/.subtk/resources/metadata.gob", home), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return common.SUBTK_ERROR, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&metadata)
	if err == io.EOF {
		metadata = make(map[string]map[string]Metadata)
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(metadata)
		if err != nil {
			return common.SUBTK_ERROR, err
		} else {
			return common.SUBTK_INIT, nil
		}
	} else if err != nil {
		return common.SUBTK_ERROR, err
	}
	// integrity check
	for group, resources := range metadata {
		if gf, err := os.Stat(fmt.Sprintf("%s/.subtk/resources/%s", home, group)); err == nil {
			if !gf.IsDir() {
				return common.SUBTK_ERROR, errmsg.ErrGroupIntegrityCheckFailed(group)
			} else {
				for resource := range resources {
					if rf, err := os.Stat(fmt.Sprintf("%s/.subtk/resources/%s/%s.gob", home, group, resource)); err == nil {
						if rf.IsDir() {
							return common.SUBTK_ERROR, errmsg.ErrResourceIntegrityCheckFailed(group, resource)
						} else {
							// check if the checksum is correct
							bs, err := os.ReadFile(fmt.Sprintf("%s/.subtk/resources/%s/%s.gob", home, group, resource))
							if err != nil {
								return common.SUBTK_ERROR, errmsg.ErrResourceIntegrityCheckFailed(group, resource)
							}
							if crc32.ChecksumIEEE(bs) != metadata[group][resource].Checksum {
								return common.SUBTK_ERROR, errmsg.ErrResourceIntegrityCheckFailed(group, resource)
							}
						}
					} else {
						return common.SUBTK_ERROR, errmsg.ErrResourceIntegrityCheckFailed(group, resource)
					}
				}
			}
		} else {
			return common.SUBTK_ERROR, errmsg.ErrGroupIntegrityCheckFailed(group)
		}
	}
	return common.SUBTK_SUCCESS, nil
}

func SetResource[T any](group, resource string, data T, updateInterval int64) error {
	lock.Lock()
	defer lock.Unlock()
	if group == "" || group == "metadata" {
		return errmsg.ErrInvalidResourceGroupName
	}
	if resource == "" {
		return errmsg.ErrInvalidResourceName
	}
	if _, ok := metadata[group]; ok {
		metadata[group][resource] = *newMetadata(&data, updateInterval)
	} else {
		metadata[group] = make(map[string]Metadata)
		metadata[group][resource] = *newMetadata(&data, updateInterval)
	}
	home, _ := os.UserHomeDir()
	os.MkdirAll(fmt.Sprintf("%s/.subtk/resources/%s", home, group), 0755)
	file, err := os.OpenFile(fmt.Sprintf("%s/.subtk/resources/%s/%s.gob", home, group, resource), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	// close the file to flush the data to disk
	file.Close()
	// generate checksum for the resource file
	bs, err := os.ReadFile(fmt.Sprintf("%s/.subtk/resources/%s/%s.gob", home, group, resource))
	if err != nil {
		return err
	}
	tempMetadata := metadata[group][resource]
	tempMetadata.Checksum = crc32.ChecksumIEEE(bs)
	metadata[group][resource] = tempMetadata
	metadataFile, err := os.OpenFile(fmt.Sprintf("%s/.subtk/resources/metadata.gob", home), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer metadataFile.Close()
	encoder = gob.NewEncoder(metadataFile)
	err = encoder.Encode(metadata)
	if err != nil {
		return err
	}
	return nil
}

func GetResource[T any](group, resource string) (*T, common.ResourceSignal) {
	lock.Lock()
	defer lock.Unlock()
	if group != "" && group != "metadata" && resource != "" {
		if _, ok := metadata[group]; ok {
			if mr, ok := metadata[group][resource]; ok {
				if mr.DataType == fmt.Sprintf("%T", *new(T)) {
					home, _ := os.UserHomeDir()
					file, err := os.Open(fmt.Sprintf("%s/.subtk/resources/%s/%s.gob", home, group, resource))
					if err != nil {
						return nil, common.RESOURCE_ERROR
					}
					defer file.Close()
					decoder := gob.NewDecoder(file)
					var data T
					err = decoder.Decode(&data)
					if err != nil {
						return nil, common.RESOURCE_ERROR
					}
					// check if the resource is outdated
					if mr.UpdateInterval != UPDATE_NOT_REQUIRED && time.Now().Unix()-mr.Timestamp > mr.UpdateInterval {
						return &data, common.RESOURCE_OUTDATED
					}
					return &data, common.RESOURCE_OK
				}
			}
		}
	}
	return nil, common.RESOURCE_NOT_FOUND
}
