package cfg

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
	"os"
)

// Read parses a YAML configuration file into a given struct.
func Read[T any](ctx context.Context, path string, defaultValues T) (*T, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("path", path)
	log.V(1).Info("reading configuration file")

	// load the file
	f, err := os.Open(path)
	if err != nil {
		log.Error(err, "failed to open file for reading")
		return nil, err
	}
	defer f.Close()
	// decode the data
	var data T
	if err := yaml.NewDecoder(f).Decode(&data); err != nil {
		log.Error(err, "failed to parse YAML file")
		return nil, err
	}
	log.V(1).Info("successfully read configuration file")

	// merge with the defaults values
	if err := mergo.Merge(&defaultValues, &data, mergo.WithOverride); err != nil {
		log.Error(err, "failed to merge configuration with default values")
		return nil, err
	}
	return &defaultValues, nil
}
