package opds2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// ParseURL parse the opds2 feed from an url
func ParseURL(url string) (*Feed, error) {

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

	feed, errParse := ParseBuffer(buff)
	if errParse != nil {
		return &Feed{}, errParse
	}

	return feed, nil
}

// ParseFile parse opds2 from a file on filesystem
func ParseFile(filePath string) (*Feed, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return &Feed{}, err
	}
	buff, errRead := ioutil.ReadAll(f)
	if err != nil {
		return &Feed{}, errRead
	}

	feed, errParse := ParseBuffer(buff)
	if errParse != nil {
		return &Feed{}, errParse
	}

	return feed, nil
}

// ParseBuffer parse opds2 feed from a buffer of byte usually get
// from a file or url
func ParseBuffer(buff []byte) (*Feed, error) {
	var feed Feed

	errParse := json.Unmarshal(buff, &feed)

	if errParse != nil {
		fmt.Println(errParse)
	}

	return &feed, nil
}

// UnmarshalJSON make all unmarshalling by hand to handle all case
func (feed *Feed) UnmarshalJSON(data []byte) error {
	var info map[string]interface{}

	json.Unmarshal(data, &info)

	for k, v := range info {
		switch k {
		case "@context":
			switch v.(type) {
			case string:
				feed.Context = append(feed.Context, v.(string))
			case []string:
				feed.Context = v.([]string)
			}
		case "metadata":
			parseMetadata(&feed.Metadata, v)
		case "links":
			parseLinks(feed, v)
		case "facets":
			parseFacets(feed, v)
		case "publications":
			parsePublications(feed, v)
		case "navigation":
			parseNavigation(feed, v)
		case "groups":
			parseGroups(feed, v)
		}
	}

	return nil
}

func parseMetadata(m *Metadata, data interface{}) {

	info := data.(map[string]interface{})
	for k, v := range info {
		switch k {
		case "title":
			m.Title = v.(string)
		case "numberOfItems":
			m.NumberOfItems = int(v.(float64))
		case "itemsPerPage":
			m.ItemsPerPage = int(v.(float64))
		case "modified":
			t, err := time.Parse(time.RFC3339, v.(string))
			if err == nil {
				m.Modified = &t
			}
		case "type":
			m.RDFType = v.(string)
		case "currentPage":
			m.CurrentPage = int(v.(float64))
		}
	}
}

func parseLinks(feed *Feed, data interface{}) {
	infoA := data.([]interface{})
	for _, vA := range infoA {
		l := parseLink(vA)
		feed.Links = append(feed.Links, l)
	}
}

func parseLink(data interface{}) Link {
	info := data.(map[string]interface{})
	l := Link{}
	for k, v := range info {
		switch k {
		case "title":
			l.Title = v.(string)
		case "href":
			l.Href = v.(string)
		case "type":
			l.TypeLink = v.(string)
		case "rel":
			switch v.(type) {
			case string:
				l.Rel = append(l.Rel, v.(string))
			case []string:
				l.Rel = v.([]string)
			}
		case "height":
			l.Height = int(v.(float64))
		case "width":
			l.Width = int(v.(float64))
		case "bitrate":
			l.Bitrate = int(v.(float64))
		case "duration":
			l.Duration = strconv.FormatFloat(v.(float64), 'f', -1, 64)
		case "templated":
			l.Templated = v.(bool)
		case "properties":
			p := Properties{}
			infoProp := v.(map[string]interface{})
			for kp, vp := range infoProp {
				switch kp {
				case "numberOfItems":
					p.NumberOfItems = int(vp.(float64))
				case "indirectAcquisition":
					infoIndir := vp.([]interface{})
					for _, in := range infoIndir {
						indir := parseIndirectAcquisition(in)
						p.IndirectAcquisition = append(p.IndirectAcquisition, indir)
					}
				case "price":
					pr := Price{}
					infoPrice := vp.(map[string]interface{})
					for kpr, vpr := range infoPrice {
						switch kpr {
						case "currency":
							pr.Currency = vpr.(string)
						case "value":
							pr.Value = vpr.(float64)
						}
					}
					p.Price = &pr
				}
			}
			l.Properties = &p
		case "children":
			lc := parseLink(v)
			l.Children = append(l.Children, lc)
		}
	}

	return l
}

func parseIndirectAcquisition(data interface{}) IndirectAcquisition {
	var i IndirectAcquisition

	info := data.(map[string]interface{})
	for k, v := range info {
		switch k {
		case "type":
			i.TypeAcquisition = v.(string)
		case "child":
			infoA := v.([]interface{})
			for _, in := range infoA {
				indirect := parseIndirectAcquisition(in)
				i.Child = append(i.Child, indirect)
			}
		}
	}

	return i
}

func parseFacets(feed *Feed, data interface{}) {
	info := data.([]interface{})
	f := Facet{}
	for _, fa := range info {
		infoA := fa.(map[string]interface{})
		for k, v := range infoA {
			switch k {
			case "metadata":
				parseMetadata(&f.Metadata, v)
			case "links":
				infoAL := v.([]interface{})
				for _, vA := range infoAL {
					l := parseLink(vA)
					f.Links = append(f.Links, l)
				}
			}
		}
		feed.Facets = append(feed.Facets, f)
	}
}

func parseGroups(feed *Feed, data interface{}) {
	info := data.([]interface{})
	for _, ga := range info {
		g := Group{}
		infoA := ga.(map[string]interface{})
		for k, v := range infoA {
			switch k {
			case "metadata":
				parseMetadata(&g.Metadata, v)
			case "links":
				infoAL := v.([]interface{})
				for _, vA := range infoAL {
					l := parseLink(vA)
					g.Links = append(g.Links, l)
				}
			case "navigation":
				infoAN := v.([]interface{})
				for _, vAN := range infoAN {
					l := parseLink(vAN)
					g.Navigation = append(g.Navigation, l)
				}
			case "publications":
				infoP := v.([]interface{})
				for _, vP := range infoP {
					p := parsePublication(vP)
					g.Publications = append(g.Publications, p)
				}
			}
		}
		feed.Groups = append(feed.Groups, g)
	}
}

func parsePublications(feed *Feed, data interface{}) {
	info := data.([]interface{})
	for _, fa := range info {
		p := parsePublication(fa)
		feed.Publications = append(feed.Publications, p)
	}
}

func parsePublication(data interface{}) Publication {
	var p Publication

	infoA := data.(map[string]interface{})
	for k, v := range infoA {
		switch k {
		case "metadata":
			parsePublicationMetadata(&p.Metadata, v)
		case "links":
			infoAL := v.([]interface{})
			for _, vA := range infoAL {
				l := parseLink(vA)
				p.Links = append(p.Links, l)
			}
		case "images":
			infoAL := v.([]interface{})
			for _, vA := range infoAL {
				l := parseLink(vA)
				p.Images = append(p.Images, l)
			}
		}
	}

	return p
}

func parsePublicationMetadata(metadata *PublicationMetadata, data interface{}) {
	info := data.(map[string]interface{})
	for k, v := range info {
		switch k {
		case "title": // handle multistring
			metadata.Title.SingleString = v.(string)
		case "identifier":
			metadata.Identifier = v.(string)
		case "@type":
			metadata.RDFType = v.(string)
		case "modified":
			t, err := time.Parse(time.RFC3339, v.(string))
			if err == nil {
				metadata.Modified = &t
			}
		case "type":
			metadata.RDFType = v.(string)
		case "author":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Author = append(metadata.Author, cont)
			}
		case "translator":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Translator = append(metadata.Translator, cont)
			}
		case "editor":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Editor = append(metadata.Editor, cont)
			}
		case "artist":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Artist = append(metadata.Artist, cont)
			}
		case "illustrator":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Illustrator = append(metadata.Illustrator, cont)
			}
		case "letterer":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Letterer = append(metadata.Letterer, cont)
			}
		case "penciler":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Penciler = append(metadata.Penciler, cont)
			}
		case "colorist":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Colorist = append(metadata.Colorist, cont)
			}
		case "inker":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Inker = append(metadata.Inker, cont)
			}
		case "narrator":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Narrator = append(metadata.Narrator, cont)
			}
		case "contributor":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Contributor = append(metadata.Contributor, cont)
			}
		case "publisher":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Publisher = append(metadata.Publisher, cont)
			}
		case "imprint":
			c := parseContributors(v)
			for _, cont := range c {
				metadata.Imprint = append(metadata.Imprint, cont)
			}
		case "language":
		case "published":
			t, err := time.Parse(time.RFC3339, v.(string))
			if err == nil {
				metadata.PublicationDate = &t
			}
		case "description":
			metadata.Description = v.(string)
		case "source":
			metadata.Source = v.(string)
		case "rights":
			metadata.Rights = v.(string)
		case "subject":
			infoS := v.([]interface{})
			for _, sub := range infoS {
				s := Subject{}
				subject := sub.(map[string]interface{})
				for ks, vs := range subject {
					switch ks {
					case "name":
						s.Name = vs.(string)
					case "sort_as":
						s.SortAs = vs.(string)
					case "scheme":
						s.Scheme = vs.(string)
					case "code":
						s.Code = vs.(string)
					}
				}
				metadata.Subject = append(metadata.Subject, s)
			}
		case "belongs_to":
			belong := BelongsTo{}
			infoB := v.(map[string]interface{})
			for kb, vb := range infoB {
				switch kb {
				case "series":
					switch vb.(type) {
					case string:
						belong.Series = append(belong.Series, Collection{Name: vb.(string)})
					case []interface{}:
						for _, colls := range vb.([]interface{}) {
							coll := parseCollection(colls)
							belong.Series = append(belong.Series, coll)
						}
					case interface{}:
						coll := parseCollection(vb)
						belong.Series = append(belong.Series, coll)
					}
				case "collection":
					switch vb.(type) {
					case string:
						belong.Collection = append(belong.Collection, Collection{Name: vb.(string)})
					case []interface{}:
						for _, colls := range vb.([]interface{}) {
							coll := parseCollection(colls)
							belong.Collection = append(belong.Collection, coll)
						}
					case interface{}:
						coll := parseCollection(vb)
						belong.Collection = append(belong.Collection, coll)
					}
				}
			}
			metadata.BelongsTo = &belong
		case "duration":
			metadata.Duration = int(v.(float64))
		}
	}
}

func parseCollection(data interface{}) Collection {
	var collection Collection

	info := data.(map[string]interface{})
	for k, v := range info {
		switch k {
		case "name":
			collection.Name = v.(string)
		case "sort_as":
			collection.SortAs = v.(string)
		case "identifier":
			collection.Identifier = v.(string)
		case "position":
			collection.Position = float32(v.(float64))
		case "links":
			infoL := v.([]interface{})
			for _, l := range infoL {
				link := parseLink(l)
				collection.Links = append(collection.Links, link)
			}
		}
	}

	return collection
}

func parseContributors(data interface{}) []Contributor {
	var c []Contributor

	switch data.(type) {
	case string:
		cont := Contributor{}
		cont.Name.SingleString = data.(string)
		c = append(c, cont)
	case []interface{}:
		infoA := data.([]interface{})
		for _, i := range infoA {
			cont := parseContributor(i)
			c = append(c, cont)
		}
	case interface{}:
		cont := parseContributor(data)
		c = append(c, cont)
	}
	return c
}

func parseContributor(data interface{}) Contributor {
	var c Contributor

	switch data.(type) {
	case string:
		c.Name.SingleString = data.(string)
	default:
		info := data.(map[string]interface{})
		for k, v := range info {
			switch k {
			case "name":
				switch v.(type) {
				case string:
					c.Name.SingleString = v.(string)
				case map[string]interface{}:
					infoN := v.(map[string]interface{})
					c.Name.MultiString = make(map[string]string)
					for kn, vn := range infoN {
						c.Name.MultiString[kn] = vn.(string)
					}
				}
			case "identifier":
				c.Identifier = v.(string)
			case "sort_as":
				c.SortAs = v.(string)
			case "role":
				c.Role = v.(string)
			case "links":
				l := parseLink(v)
				c.Links = append(c.Links, l)
			}
		}
	}

	return c
}

func parseNavigation(feed *Feed, data interface{}) {
	infoA := data.([]interface{})
	for _, vA := range infoA {
		l := parseLink(vA)
		feed.Navigation = append(feed.Navigation, l)
	}
}

// UnmarshalJSON overwrite json unmarshalling for Rel for handling
// when we have a array of a string
// func (r *StringOrArray) UnmarshalJSON(data []byte) error {
// 	var relAr []string
//
// 	if data[0] == '[' {
// 		err := json.Unmarshal(data, &relAr)
// 		if err != nil {
// 			return err
// 		}
// 		for _, ra := range relAr {
// 			*r = append(*r, ra)
// 		}
// 	} else {
// 		*r = append(*r, string(data))
// 	}
//
// 	return nil
// }

// UnmarshalJSON overwrite json unmarshalling for MultiLanguage
// when we have an entry in the Multi fields we use it
// otherwise we use the single string
// func (m *MultiLanguage) UnmarshalJSON(data []byte) error {
// 	var mParse map[string]string
//
// 	if data[0] == '{' {
// 		json.Unmarshal(data, &mParse)
// 		m.MultiString = mParse
// 	} else {
// 		m.SingleString = string(data)
// 	}
//
// 	return nil
// }
