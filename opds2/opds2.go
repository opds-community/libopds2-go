// Package opds2 provide parsing and generation method for an OPDS2 feed
// https://github.com/opds-community/opds-revision/blob/master/opds-2.0.md
package opds2

import (
	"encoding/json"
	"time"
)

// Feed is a collection as defined in Readium Web Publication Manifest
type Feed struct {
	Context      []string      `json:"@context,omitempty"`
	Metadata     Metadata      `json:"metadata"`
	Links        []Link        `json:"links"`
	Facets       []Facet       `json:"facets,omitempty"`
	Groups       []Group       `json:"groups,omitempty"`
	Publications []Publication `json:"publications,omitempty"`
	Navigation   []Link        `json:"navigation,omitempty"`
}

// Publication is a collection for a given publication
type Publication struct {
	Metadata PublicationMetadata `json:"metadata"`
	Links    []Link              `json:"links"`
	Images   []Link              `json:"images"`
}

// Metadata has a limited subset of metadata compared to a publication
type Metadata struct {
	RDFType       string     `json:"@type,omitempty"`
	Title         string     `json:"title"`
	NumberOfItems int        `json:"numberOfItems,omitempty"`
	ItemsPerPage  int        `json:"itemsPerPage,omitempty"`
	CurrentPage   int        `json:"currentPage,omitempty"`
	Modified      *time.Time `json:"modified,omitempty"`
}

// Facet is a collection that contains a facet group
type Facet struct {
	Metadata Metadata `json:"metadata"`
	Links    []Link   `json:"links"`
}

// Group is a group collection that must contain publications
type Group struct {
	Metadata     Metadata      `json:"metadata"`
	Links        []Link        `json:"links,omitempty"`
	Publications []Publication `json:"publications,omitempty"`
	Navigation   []Link        `json:"navigation,omitempty"`
}

// Link object used in collections and links
type Link struct {
	Href       string        `json:"href"`
	TypeLink   string        `json:"type,omitempty"`
	Rel        StringOrArray `json:"rel,omitempty"`
	Height     int           `json:"height,omitempty"`
	Width      int           `json:"width,omitempty"`
	Title      string        `json:"title,omitempty"`
	Properties *Properties   `json:"properties,omitempty"`
	Duration   string        `json:"duration,omitempty"`
	Templated  bool          `json:"templated,omitempty"`
	Children   []Link        `json:"children,omitempty"`
	Bitrate    int           `json:"bitrate,omitempty"`
}

// Properties object use to link properties
// Use also in Rendition for fxl
type Properties struct {
	NumberOfItems       int                   `json:"numberOfItems,omitempty"`
	Price               *Price                `json:"price,omitempty"`
	IndirectAcquisition []IndirectAcquisition `json:"indirectAcquisition,omitempty"`
}

// IndirectAcquisition store
type IndirectAcquisition struct {
	TypeAcquisition string                `json:"type"`
	Child           []IndirectAcquisition `json:"child,omitempty"`
}

// Price price information
type Price struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

// PublicationMetadata for the default context in WebPub
type PublicationMetadata struct {
	RDFType         string        `json:"@type,omitempty"` //Defaults to schema.org for EBook
	Title           MultiLanguage `json:"title"`
	Identifier      string        `json:"identifier"`
	Author          []Contributor `json:"author,omitempty"`
	Translator      []Contributor `json:"translator,omitempty"`
	Editor          []Contributor `json:"editor,omitempty"`
	Artist          []Contributor `json:"artist,omitempty"`
	Illustrator     []Contributor `json:"illustrator,omitempty"`
	Letterer        []Contributor `json:"letterer,omitempty"`
	Penciler        []Contributor `json:"penciler,omitempty"`
	Colorist        []Contributor `json:"colorist,omitempty"`
	Inker           []Contributor `json:"inker,omitempty"`
	Narrator        []Contributor `json:"narrator,omitempty"`
	Contributor     []Contributor `json:"contributor,omitempty"`
	Publisher       []Contributor `json:"publisher,omitempty"`
	Imprint         []Contributor `json:"imprint,omitempty"`
	Language        StringOrArray `json:"language,omitempty"`
	Modified        *time.Time    `json:"modified,omitempty"`
	PublicationDate *time.Time    `json:"published,omitempty"`
	Description     string        `json:"description,omitempty"`
	Source          string        `json:"source,omitempty"`
	Rights          string        `json:"rights,omitempty"`
	Subject         []Subject     `json:"subject,omitempty"`
	BelongsTo       *BelongsTo    `json:"belongs_to,omitempty"`
	Duration        int           `json:"duration,omitempty"`
}

// Contributor construct used internally for all contributors
type Contributor struct {
	Name       MultiLanguage `json:"name,omitempty"`
	SortAs     string        `json:"sort_as,omitempty"`
	Identifier string        `json:"identifier,omitempty"`
	Role       string        `json:"role,omitempty"`
	Links      []Link        `json:"links,omitempty"`
}

// Subject as based on EPUB 3.1 and WebPub
type Subject struct {
	Name   string `json:"name"`
	SortAs string `json:"sort_as,omitempty"`
	Scheme string `json:"scheme,omitempty"`
	Code   string `json:"code,omitempty"`
}

// BelongsTo is a list of collections/series that a publication belongs to
type BelongsTo struct {
	Series     []Collection `json:"series,omitempty"`
	Collection []Collection `json:"collection,omitempty"`
}

// Collection construct used for collection/serie metadata
type Collection struct {
	Name       string  `json:"name"`
	SortAs     string  `json:"sort_as,omitempty"`
	Identifier string  `json:"identifier,omitempty"`
	Position   float32 `json:"position,omitempty"`
	Links      []Link  `json:"links,omitempty"`
}

// MultiLanguage store the a basic string when we only have one lang
// Store in a hash by language for multiple string representation
type MultiLanguage struct {
	SingleString string
	MultiString  map[string]string
}

// StringOrArray is a array of string redifine for overriding json
// marshalling and unmarshalling
type StringOrArray []string

// MarshalJSON overwrite json marshalling for MultiLanguage
// when we have an entry in the Multi fields we use it
// otherwise we use the single string
func (m MultiLanguage) MarshalJSON() ([]byte, error) {
	if len(m.MultiString) > 0 {
		return json.Marshal(m.MultiString)
	}
	return json.Marshal(m.SingleString)
}

func (m MultiLanguage) String() string {
	if len(m.MultiString) > 0 {
		for _, s := range m.MultiString {
			return s
		}
	}
	return m.SingleString
}

// MarshalJSON overwrite json marshalling for handling string or array
func (r StringOrArray) MarshalJSON() ([]byte, error) {
	if len(r) == 1 {
		return json.Marshal(r[0])
	}
	return json.Marshal(r)
}

func (publication *Publication) findFirstLinkByRel(rel string) Link {

	for _, l := range publication.Links {
		for _, r := range l.Rel {
			if r == rel {
				return l
			}
		}
	}

	return Link{}
}

// New create a new feed structure
func New(title string) Feed {
	var feed Feed

	feed.Metadata.Title = title
	t := time.Now()
	feed.Metadata.Modified = &t

	return feed
}
