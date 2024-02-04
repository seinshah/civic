package configuration

import (
	"context"
	"errors"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/seinshah/cvci/internal/loader"
	"github.com/seinshah/cvci/internal/types"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// Loader is the loader that is ready to load the configuration.
	Loader loader.Loader `validate:"required"`
}

type Engine struct {
	config  Config
	content []byte
	data    *types.Schema
}

var (
	ErrInvalidContent             = errors.New("configuration format is not valid")
	ErrInvalidConfigurationFormat = errors.New("configuration format does not match the schema")
)

func NewEngine(ctx context.Context, config Config) (*Engine, error) {
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

	return &Engine{
		config:  config,
		content: content,
		data:    data,
	}, nil
}

func (e *Engine) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(e.data); err != nil {
		return errors.Join(ErrInvalidConfigurationFormat, err)
	}

	return nil
}

func (e *Engine) SchemaData() *types.Schema {
	return e.data
}
