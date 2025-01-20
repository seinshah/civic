package cv

import (
	"context"
	"errors"
	"fmt"

	"github.com/seinshah/civic/internal/pkg/loader"
	"github.com/seinshah/civic/internal/pkg/types"
)

var ErrInvalidSchemaFormat = errors.New("schema file format does not match the schema")

func (h *Handler) parseSchemaFile(ctx context.Context) (*types.Schema, error) {
	confLoader, err := loader.NewGeneralLoader(h.schemaFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load the schema file (%s): %w", h.schemaFilePath, err)
	}

	content, err := confLoader.Load(ctx)
	if err != nil {
		return nil, err
	}

	data, err := types.NewSchema(content, h.schemaType)
	if err != nil {
		return nil, err
	}

	if err = data.IsValid(); err != nil {
		return nil, errors.Join(ErrInvalidSchemaFormat, err)
	}

	return data, nil
}
