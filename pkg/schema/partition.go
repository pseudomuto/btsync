package schema

import (
	yaml "gopkg.in/yaml.v2"

	"context"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

// A Descriptor is an object with a name and description.
type Descriptor struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// A Partition represents a set of table definitions
type Partition struct {
	Name   string   `yaml:"name"`
	Tables []*Table `yaml:"tables"`
}

// A Table represents a CBT table and the set of column families it defines
type Table struct {
	Descriptor     `yaml:",inline"`
	ColumnFamilies []*ColumnFamily `yaml:"columnFamilies"`
}

// A ColumnFamily represents a CBT column family including any retention policies that should be applied.
type ColumnFamily struct {
	Descriptor `yaml:",inline"`
}

// A RetentionPolicy describes the policy that should be applied to a column family.
type RetentionPolicy struct {
	Versions int           `yaml:"versions"`
	TTL      time.Duration `yaml:"ttl"`
}

// A PartitionResult decribes the result of parsing a partition file.
type PartitionResult struct {
	Partition Partition
	Err       error
}

// ParseDirectory parses each yaml file in the specified directory into a partition object.
func ParseDirectory(ctx context.Context, dir string) <-chan *PartitionResult {
	files, _ := filepath.Glob(filepath.Join(dir, "*.yml"))
	stream := make(chan *PartitionResult, len(files))

	go func() {
		defer close(stream)

		for _, pFile := range files {
			select {
			case <-ctx.Done():
				return
			case stream <- parsePartition(pFile):
			}
		}
	}()

	return stream
}

func parsePartition(path string) *PartitionResult {
	partName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	pr := PartitionResult{Partition: Partition{Name: partName}}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		pr.Err = err
		return &pr
	}

	pr.Err = yaml.Unmarshal(data, &pr.Partition)
	return &pr
}
