package types

import "errors"

//go:generate go tool go-enum --nocase --names

var ErrInvalidPageMargin = errors.New("invalid page margin")

// PageSize is the type for defining the page size for the PDF.
// ENUM(A4, B4, A, Arch-A, Letter).
type PageSize string

// PageMargin is the type for defining the page margin for the PDF.
type PageMargin struct {
	Top    float64 `json:"top,omitempty"    validate:"gte=0,lt=3" yaml:"top"`
	Right  float64 `json:"right,omitempty"  validate:"gte=0,lt=3" yaml:"right"`
	Bottom float64 `json:"bottom,omitempty" validate:"gte=0,lt=3" yaml:"bottom"`
	Left   float64 `json:"left,omitempty"   validate:"gte=0,lt=3" yaml:"left"`
}

const (
	pageSizeA4Width  = 8.27
	pageSizeA4Height = 11.69

	pageSizeB4Width  = 10.12
	pageSizeB4Height = 14.33

	pageSizeAWidth  = 8.5
	pageSizeAHeight = 11

	pageSizeArchAWidth  = 9
	pageSizeArchAHeight = 12
)

// GetWidthInch returns the width of the page in inch.
func (p PageSize) GetWidthInch() float64 {
	switch p {
	case PageSizeA4:
		return pageSizeA4Width
	case PageSizeB4:
		return pageSizeB4Width
	case PageSizeA, PageSizeLetter:
		return pageSizeAWidth
	case PageSizeArchA:
		return pageSizeArchAWidth
	}

	return DefaultPageSize.GetWidthInch()
}

// GetHeightInch returns the height of the page in inch.
func (p PageSize) GetHeightInch() float64 {
	switch p {
	case PageSizeA4:
		return pageSizeA4Height
	case PageSizeB4:
		return pageSizeB4Height
	case PageSizeA, PageSizeLetter:
		return pageSizeAHeight
	case PageSizeArchA:
		return pageSizeArchAHeight
	}

	return DefaultPageSize.GetHeightInch()
}
