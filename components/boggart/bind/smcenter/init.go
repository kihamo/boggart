package smcenter

import (
	"net/url"

	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	link, _ := url.Parse("https://dom-24.net")
	boggart.RegisterBindType("smcenter:dom24", Type{
		Link:            link,
		BaseURL:         "/dgservicnew/",
		BillContentType: "application/pdf",
	})

	link, _ = url.Parse("https://www.rkcm.ru")
	boggart.RegisterBindType("smcenter:rkcm", Type{
		Link:            link,
		BaseURL:         "/muprcmytishi/",
		BillContentType: "image/png",
	})

	link, _ = url.Parse("https://api.sm-center.ru")
	boggart.RegisterBindType("smcenter:hashdomru", Type{
		Link:            link,
		BaseURL:         "/vremya/",
		BillContentType: "application/pdf",
	})

	boggart.RegisterBindType("smcenter:custom", Type{
		BaseURL:         "/",
		BillContentType: "application/pdf",
	})
}
