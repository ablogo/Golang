package webscraping

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"product-tracking/kafka"
	"product-tracking/service"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type DevTools struct {
	ProductSvc *service.ProductService
	KafkaSvc   *kafka.KafkaSvc
}

func check(err error, product_id int) {
	var evalErr *rod.EvalError
	if errors.Is(err, context.DeadlineExceeded) {
		slog.Error("Price not found, timeout.", "product_id", product_id)
	} else if errors.As(err, &evalErr) {
		slog.Error(evalErr.Error(), "product_id", product_id)
	} else if err != nil {
		slog.Error("Price not found.", "description", err.Error(), "product_id", product_id)
	}
}

func (d *DevTools) Execute() {
	products := d.ProductSvc.GetProducts(false)

	if products == nil || len(*products) == 0 {
		slog.Info("No products found on the db.")
		return
	}

	slog.Info("Opening Browser..")
	u := launcher.NewUserMode().
		Set("headless", "new").
		Set("no-sandbox").
		Set("disable-gpu").
		Set("window-size", "1920,1080").
		Set("enable-unsafe-swiftshader").
		Set("disable-blink-features", "AutomationControlled").
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
		MustLaunch() // local dev
		//Bin("/usr/bin/google-chrome-stable").MustLaunch() // prod deployment, set chrome path on the server.

	browser := rod.New().ControlURL(u).MustConnect().SlowMotion(time.Second * 2).NoDefaultDevice()
	defer browser.Close()

	page, err := browser.Page(proto.TargetCreateTarget{URL: d.KafkaSvc.Settings.INITIAL_WEBPAGE})
	if err != nil {
		slog.Error("Fail to start broswer.", "description ", err.Error())
	}

	err = page.WaitLoad()
	if err != nil {
		slog.Error("Fail to load intial page.", "description ", err.Error())
	}

	for _, v := range *products {
		price := "0"

		page.Navigate(v.Url)
		if err != nil {
			slog.Error("Product page not found.", "product_id ", v.ID)
		}

		err = page.WaitLoad()
		if err != nil {
			slog.Error("Product page loading problems.", "product_id ", v.ID)
		}

		tags := strings.Split(v.TagId, ";")
		for _, tag := range tags {

			elem, err := page.Timeout(time.Second * 15).Element(tag)
			if err != nil {
				check(err, v.ID)
			} else {
				price, err = elem.Text()
				slog.Info("HTML element text. " + price)
			}

			clean_price, err := d.ProductSvc.RegisterPrice(v.ID, price)
			if clean_price > 0 {

				isLowestPrice := d.ProductSvc.IsLowestPrice(v.ID, clean_price)
				if isLowestPrice {
					msg := "Lower price detected, product id. " + strconv.Itoa(v.ID)
					d.KafkaSvc.SendMessage(msg, "prices", "notifications")
					slog.Info(msg)
				}
			}

			time.Sleep(time.Second * time.Duration(rand.Intn(20)+6))
		}
	}

	slog.Info("Browser closed.")
}
