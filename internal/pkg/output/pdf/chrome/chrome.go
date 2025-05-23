package chrome

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/go-playground/validator/v10"
	"github.com/seinshah/civic/internal/pkg/types"
)

type options struct {
	pageSize   types.PageSize
	pageMargin types.PageMargin
}

type Headless struct {
	config options
}

type Option func(*options)

var _ types.OutputGenerator = &Headless{}

func WithPageSize(size types.PageSize) Option {
	return func(o *options) {
		o.pageSize = size
	}
}

func WithPageMargin(margin types.PageMargin) Option {
	return func(o *options) {
		o.pageMargin = margin
	}
}

func NewHeadless(opts ...Option) *Headless {
	instanceOpts := options{
		pageSize: types.DefaultPageSize,
	}

	for _, opt := range opts {
		opt(&instanceOpts)
	}

	return &Headless{
		config: instanceOpts,
	}
}

func (h *Headless) Generate(ctx context.Context, content []byte) ([]byte, error) {
	if !h.config.pageSize.IsValid() {
		return nil, types.ErrInvalidPageSize
	}

	if err := validator.New().Struct(h.config.pageMargin); err != nil {
		return nil, errors.Join(types.ErrInvalidPageMargin, err)
	}

	newCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var result []byte

	printTask := chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(h.getLoadContentAction(string(content))),
		// Wait for the page to be fully loaded
		chromedp.WaitReady("body", chromedp.ByQuery),
		// Wait for all stylesheets to be loaded
		chromedp.ActionFunc(
			func(ctx context.Context) error {
				jsPoll := `
					Array.from(document.querySelectorAll('link[rel="stylesheet"]'))
						.every(link => link.sheet || link.onload === null)
				`

				return chromedp.Poll(
					jsPoll,
					nil,
					chromedp.WithPollingInterval(500*time.Millisecond), //nolint:mnd
					chromedp.WithPollingTimeout(10*time.Second),        //nolint:mnd
				).Do(ctx)
			},
		),
		// Give a tiny extra buffer for layout/paint
		chromedp.Sleep(500 * time.Millisecond), //nolint:mnd
		chromedp.ActionFunc(h.getPrintToPDFAction(&result)),
	}

	if err := chromedp.Run(newCtx, printTask); err != nil {
		return nil, err
	}

	return result, nil
}

func (h *Headless) getLoadContentAction(html string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		loadCtx, loadCancel := context.WithCancel(ctx)
		defer loadCancel()

		var wg sync.WaitGroup

		wg.Add(1)

		chromedp.ListenTarget(
			loadCtx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					loadCancel()
					wg.Done()
				}
			},
		)

		frameTree, err := page.GetFrameTree().Do(ctx)
		if err != nil {
			return err
		}

		if err = page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx); err != nil {
			return err
		}

		wg.Wait()

		return nil
	}
}

func (h *Headless) getPrintToPDFAction(output *[]byte) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		var err error

		if *output, _, err = page.PrintToPDF().
			WithDisplayHeaderFooter(false).
			WithPrintBackground(true).
			WithScale(1).
			WithPaperWidth(h.config.pageSize.GetWidthInch()).
			WithPaperHeight(h.config.pageSize.GetHeightInch()).
			WithMarginTop(h.config.pageMargin.Top).
			WithMarginRight(h.config.pageMargin.Right).
			WithMarginBottom(h.config.pageMargin.Bottom).
			WithMarginLeft(h.config.pageMargin.Left).
			Do(ctx); err != nil {
			return err
		}

		return nil
	}
}
