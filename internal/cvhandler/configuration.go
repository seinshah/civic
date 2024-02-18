package cvhandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/seinshah/cvci/internal/pkg/loader"
	"github.com/seinshah/cvci/internal/pkg/types"
	"gopkg.in/yaml.v3"
)

type ConfigurationConfig struct {
	// Loader is the loader that is ready to load the configuration.
	Loader loader.Loader `validate:"required"`
}

type Configuration struct {
	config  ConfigurationConfig
	content []byte
	data    *types.Schema
}

var (
	ErrInvalidContent             = errors.New("configuration format is not valid")
	ErrInvalidConfigurationFormat = errors.New("configuration format does not match the schema")
)

func NewConfiguration(ctx context.Context, config ConfigurationConfig) (*Configuration, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("invalid initialization of configuration engine: %w", err)
	}

	content, err := config.Loader.Load(ctx)
	if err != nil {
		return nil, err
	}

	var data *types.Schema

	if err = defaults.Set(data); err != nil {
		return nil, errors.Join(ErrInvalidContent, err)
	}

	if err = yaml.Unmarshal(content, data); err != nil {
		return nil, errors.Join(ErrInvalidContent, err)
	}

	return &Configuration{
		config:  config,
		content: content,
		data:    data,
	}, nil
}

func (c *Configuration) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(c.data); err != nil {
		return errors.Join(ErrInvalidConfigurationFormat, err)
	}

	return nil
}

func (c *Configuration) Data() *types.Schema {
	return c.data
}
