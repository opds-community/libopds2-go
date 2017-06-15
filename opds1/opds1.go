package opds1

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"
)

// Feed root element for acquisition or navigation feed
type Feed struct {
	ID           string    `xml:"id"`
	Title        string    `xml:"title"`
	Updated      time.Time `xml:"updated"`
	Entries      []Entry   `xml:"entry"`
	Links        []Link    `xml:"link"`
	TotalResults int       `xml:"totalResults"`
	ItemsPerPage int       `xml:"itemsPerPage"`
}

// Link link to different resources
type Link struct {
	Rel                 string                `xml:"rel,attr"`
	Href                string                `xml:"href,attr"`
	TypeLink            string                `xml:"type,attr"`
	Title               string                `xml:"title,attr"`
	FacetGroup          string                `xml:"facetGroup,attr"`
	Count               int                   `xml:"count,attr"`
	Price               Price                 `xml:"price"`
	IndirectAcquisition []IndirectAcquisition `xml:"indirectAcquisition"`
}

// Author represent the feed author or the entry author
type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

type Entry struct {
	Title      string     `xml:"title"`
	ID         string     `xml:"id"`
	Identifier string     `xml:"identifier"`
	Updated    *time.Time `xml:"updated"`
	Rights     string     `xml:"rights"`
	Publisher  string     `xml:"publisher"`
	Author     []Author   `xml:"author,omitempty"`
	Language   string     `xml:"language"`
	Issued     string     `xml:"issued"` // Check for format
	Published  *time.Time `xml:"published"`
	Category   []Category `xml:"category,omitempty"`
	Links      []Link     `xml:"link,omitempty"`
	Summary    Content    `xml:"summary"`
	Content    Content    `xml:"content"`
}

type Content struct {
	Content     string `xml:",cdata"`
	ContentType string `xml:"type,attr"`
}

type Category struct {
	Scheme string `xml:"scheme,attr"`
	Term   string `xml:"term,attr"`
	Label  string `xml:"label,attr"`
}

type Price struct {
	CurrencyCode string  `xml:"currencycode,attr"`
	Value        float64 `xml:",cdata"`
}

type IndirectAcquisition struct {
	TypeAcquisition     string                `xml:"type,attr"`
	IndirectAcquisition []IndirectAcquisition `xml:"indirectAcquisition"`
}

func ParseURL(url string) (*Feed, error) {
	var feed Feed

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, errReq := http.DefaultClient.Do(request)
	if errReq != nil {
		return nil, errReq
	}

	buff, errRead := ioutil.ReadAll(res.Body)
	if errRead != nil {
		return nil, errRead
	}

	xml.Unmarshal(buff, &feed)

	return &feed, nil
}
