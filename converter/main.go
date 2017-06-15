package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/opds-community/libopds2-go/opds1"
	"github.com/opds-community/libopds2-go/opds2"
)

func main() {

	feed, err := opds1.ParseURL(os.Args[1])
	if err != nil {
		fmt.Println(err)
	} else {
		opds2feed := fillOPDS2(feed, os.Args[1])
		j, _ := JSONMarshal(opds2feed, true)
		var identJSON bytes.Buffer

		json.Indent(&identJSON, j, "", " ")
		fmt.Println(identJSON.String())
	}

}

func fillOPDS2(feed *opds1.Feed, url string) opds2.Feed {
	var opds2feed opds2.Feed

	// If acquisition link check if rel='collection' than mean it is a group, if there no rel it is a publication

	opds2feed.Metadata.Title = feed.Title
	opds2feed.Metadata.Modified = &feed.Updated
	if feed.TotalResults != 0 {
		opds2feed.Metadata.NumberOfItems = feed.TotalResults
	}
	if feed.ItemsPerPage != 0 {
		opds2feed.Metadata.ItemsPerPage = feed.ItemsPerPage
	}

	for _, entry := range feed.Entries {
		// Get all entry, if entry has a acquisition puts in publication else it is a navigation link put in in links objetcs to
		isAnNavigation := true
		collLink := opds2.Link{}

		for _, l := range entry.Links {
			if strings.Contains(l.Rel, "http://opds-spec.org/acquisition") {
				isAnNavigation = false
			}
			if l.Rel == "collection" || l.Rel == "http://opds-spec.org/group" {
				collLink.Rel = []string{"collection"}
				collLink.Href = l.Href
				collLink.Title = l.Title
			}
		}

		if isAnNavigation == false {
			p := opds2.Publication{}
			p.Metadata.Title.SingleString = entry.Title
			if entry.Identifier != "" {
				p.Metadata.Identifier = entry.Identifier
			} else {
				p.Metadata.Identifier = entry.ID
			}
			p.Metadata.Language = []string{entry.Language}
			p.Metadata.Modified = entry.Updated
			p.Metadata.PublicationDate = entry.Published
			p.Metadata.Rights = entry.Rights
			if entry.Publisher != "" {
				c := opds2.Contributor{}
				c.Name.SingleString = entry.Publisher
				p.Metadata.Publisher = append(p.Metadata.Publisher, c)
			}

			for _, cat := range entry.Category {
				p.Metadata.Subject = append(p.Metadata.Subject, opds2.Subject{Code: cat.Term, Name: cat.Label, Scheme: cat.Scheme})
			}

			for _, aut := range entry.Author {
				cont := opds2.Contributor{}
				cont.Name.SingleString = aut.Name
				cont.Identifier = aut.URI
				p.Metadata.Author = append(p.Metadata.Author, cont)
			}

			// for html resource like description, atom:summary go to description
			// if atom:content use it in description else use summary
			if entry.Content.Content != "" {
				p.Metadata.Description = entry.Content.Content
			} else if entry.Summary.Content != "" {
				p.Metadata.Description = entry.Summary.Content
			}

			for _, link := range entry.Links {
				l := opds2.Link{}
				l.Href = link.Href
				l.TypeLink = link.TypeLink
				l.Rel = []string{link.Rel}
				l.Title = link.Title

				if len(link.IndirectAcquisition) > 0 {
					if l.Properties == nil {
						l.Properties = &opds2.Properties{}
					}

					for _, ia := range link.IndirectAcquisition {
						ind := opds2.IndirectAcquisition{}
						ind.TypeAcquisition = ia.TypeAcquisition
						if len(ia.IndirectAcquisition) > 0 {
							for _, iac := range ia.IndirectAcquisition {
								cia := opds2.IndirectAcquisition{}
								cia.TypeAcquisition = iac.TypeAcquisition
								ind.Child = append(ind.Child, cia)
							}
						}
						l.Properties.IndirectAcquisition = append(l.Properties.IndirectAcquisition, ind)
					}
				}

				if link.Price.CurrencyCode != "" {
					if l.Properties == nil {
						l.Properties = &opds2.Properties{}
					}
					l.Properties.Price = &opds2.Price{}
					l.Properties.Price.Currency = link.Price.CurrencyCode
					l.Properties.Price.Value = link.Price.Value
				}

				if link.Rel == "collection" || link.Rel == "http://opds-spec.org/group" {
				} else if link.Rel == "http://opds-spec.org/image" || link.Rel == "http://opds-spec.org/image/thumbnail" {
					p.Images = append(p.Images, l)
				} else {
					p.Links = append(p.Links, l)
				}
			}

			if collLink.Href != "" {
				opds2feed.AddPublicationInGroup(p, collLink)
			} else {
				opds2feed.Publications = append(opds2feed.Publications, p)
			}
		} else {
			linkNav := opds2.Link{}
			linkNav.Title = entry.Title
			linkNav.Rel = []string{entry.Links[0].Rel}
			linkNav.TypeLink = entry.Links[0].TypeLink
			linkNav.Href = entry.Links[0].Href

			if collLink.Href != "" {
				opds2feed.AddNavigationInGroup(linkNav, collLink)
			} else {
				opds2feed.Navigation = append(opds2feed.Navigation, linkNav)
			}
		}
	}

	for _, l := range feed.Links {
		linkFeed := opds2.Link{}
		linkFeed.Href = l.Href
		linkFeed.Rel = []string{l.Rel}
		linkFeed.TypeLink = l.TypeLink
		linkFeed.Title = l.Title

		if l.Rel == "http://opds-spec.org/facet" {
			linkFeed.Properties = &opds2.Properties{NumberOfItems: l.Count}
			opds2feed.AddFacet(linkFeed, l.FacetGroup)
		} else {
			opds2feed.Links = append(opds2feed.Links, linkFeed)
		}
	}

	return opds2feed
}

// JSONMarshal override marshalling function to fix some encoding
func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}
