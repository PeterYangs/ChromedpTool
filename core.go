package ChromedpTool

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"time"
)

type ChromedpTool struct {
	debug   bool
	timeout time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
}

func New() *ChromedpTool {

	return &ChromedpTool{
		debug:   false,
		timeout: 5 * time.Minute,
	}
}

func (c *ChromedpTool) Debug() *ChromedpTool {

	c.debug = true

	return c
}

func (c *ChromedpTool) initChrome() (context.Context, context.CancelFunc) {

	var opts []chromedp.ExecAllocatorOption

	if c.debug {
		opts = append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.WindowSize(1920, 1080),
		)
	} else {
		opts = append(chromedp.DefaultExecAllocatorOptions[:],

			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.WindowSize(1920, 1080),
		)
	}

	cc, _ := context.WithTimeout(context.Background(), c.timeout)

	allocCtx, cancel := chromedp.NewExecAllocator(cc, opts...)

	//create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
	)

	c.ctx = ctx

	c.cancel = cancel

	return ctx, cancel
}

func (c *ChromedpTool) SetCookie(name, value, domain, path string, httpOnly, secure bool) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		err := network.SetCookie(name, value).
			WithExpires(&expr).
			WithDomain(domain).
			WithPath(path).
			WithHTTPOnly(httpOnly).
			WithSecure(secure).
			Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (c *ChromedpTool) WaitElementExist(ctx context.Context, selector string, timeout time.Duration) bool {
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := chromedp.Run(ctxTimeout,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
	)

	return err == nil
}
