package openhab

// https://openhab.org/javadoc/latest/org/openhab/core/model/sitemap/sitemap/impl/switchimpl

type Sitemap struct {
	name     string
	label    string
	elements SitemapElements
}

type SitemapElements []*SitemapElement

type SitemapElement struct {
	icon       string
	label      string
	labelColor []string
	valueColor []string
	visibility []string
}

type SitemapElementFrame struct {
}

type SitemapElementDefault struct {
}
