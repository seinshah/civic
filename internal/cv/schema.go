package cv

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

var (
	ErrInvalidContent             = errors.New("configuration format is not valid")
	ErrInvalidConfigurationFormat = errors.New("configuration format does not match the schema")
)

func (h *Handler) parseSchemaFile(ctx context.Context) (*types.Schema, error) {
	confLoader, err := loader.NewGeneralLoader(h.schemaFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load cvci schema file (%s): %w", h.schemaFilePath, err)
	}

	content, err := confLoader.Load(ctx)
	if err != nil {
		return nil, err
	}

	data := types.Schema{}

	if err = defaults.Set(&data); err != nil {
		return nil, errors.Join(ErrInvalidContent, err)
	}

	if err = yaml.Unmarshal(content, &data); err != nil {
		return nil, errors.Join(ErrInvalidContent, err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err = validate.Struct(&data); err != nil {
		return nil, errors.Join(ErrInvalidConfigurationFormat, err)
	}

	return &data, nil
}
