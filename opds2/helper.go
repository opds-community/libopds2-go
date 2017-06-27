package opds2

// AddLink add a new link in feed information
// at minimum the self link
func (feed *Feed) AddLink(href string, rel string, typeLink string, templated bool) {
	var l Link

	l.Href = href
	l.Rel = append(l.Rel, rel)
	l.TypeLink = typeLink
	if templated == true {
		l.Templated = true
	}

	feed.Links = append(feed.Links, l)
}

// AddImage add a image link to Publication
func (publication *Publication) AddImage(href string, typeImage string, height int, width int) {
	var i Link

	i.Href = href
	i.TypeLink = typeImage
	if height > 0 {
		i.Height = height
	}
	if width > 0 {
		i.Width = width
	}

	publication.Images = append(publication.Images, i)
}

// AddLink add a new link to the publication
func (publication *Publication) AddLink(href string, typeLink string, rel string, title string) {
	var l Link

	l.Href = href
	l.TypeLink = typeLink
	if rel != "" {
		l.Rel = append(l.Rel, rel)
	}
	if title != "" {
		l.Title = title
	}

	publication.Links = append(publication.Links, l)
}

// AddAuthor add author to publication with all parameters mostly optional
func (publication *Publication) AddAuthor(name string, identifier string, sortAs string, href string, typeLink string) {
	var c Contributor
	var l Link

	c.Name.SingleString = name
	if identifier != "" {
		c.Identifier = identifier
	}
	if sortAs != "" {
		c.SortAs = sortAs
	}
	if href != "" {
		l.Href = href
	}
	if typeLink != "" {
		l.TypeLink = typeLink
	}

	if l.Href != "" {
		c.Links = append(c.Links, l)
	}

	publication.Metadata.Author = append(publication.Metadata.Author, c)
}

// AddSerie add serie to publication
func (publication *Publication) AddSerie(name string, position float32, href string, typeLink string) {
	var c Collection
	var l Link

	c.Name = name
	c.Position = position

	if publication.Metadata.BelongsTo == nil {
		publication.Metadata.BelongsTo = &BelongsTo{}
	}

	if typeLink != "" {
		l.TypeLink = typeLink
	}

	if l.Href != "" {
		c.Links = append(c.Links, l)
	}

	publication.Metadata.BelongsTo.Series = append(publication.Metadata.BelongsTo.Series, c)
}

// AddPublisher add publisher to publication
func (publication *Publication) AddPublisher(name string, href string, typeLink string) {
	var c Contributor
	var l Link

	c.Name.SingleString = name

	if typeLink != "" {
		l.TypeLink = typeLink
	}

	if l.Href != "" {
		c.Links = append(c.Links, l)
	}

	publication.Metadata.Publisher = append(publication.Metadata.Publisher, c)
}

// AddNavigation add navigation element in feed
func (feed *Feed) AddNavigation(title string, href string, rel string, typeLink string) {
	var l Link

	l.Href = href
	l.TypeLink = typeLink
	l.Rel = append(l.Rel, rel)
	if title != "" {
		l.Title = title
	}

	feed.Navigation = append(feed.Navigation, l)
}

// AddPagination add pagination and link information in feed
func (feed *Feed) AddPagination(numberItems int, itemsPerPage int, currentPage int, nextLink string, prevLink string, firstLink string, lastLink string) {

	feed.Metadata.CurrentPage = currentPage
	feed.Metadata.ItemsPerPage = itemsPerPage
	feed.Metadata.NumberOfItems = numberItems

	if nextLink != "" {
		feed.AddLink(nextLink, "next", "application/opds+json", false)
	}
	if prevLink != "" {
		feed.AddLink(prevLink, "previous", "application/opds+json", false)
	}
	if firstLink != "" {
		feed.AddLink(firstLink, "first", "application/opds+json", false)
	}
	if lastLink != "" {
		feed.AddLink(lastLink, "last", "application/opds+json", false)
	}
}

// AddFacet add link to facet handler multiple add
func (feed *Feed) AddFacet(link Link, group string) {
	var facet Facet

	for i, f := range feed.Facets {
		if f.Metadata.Title == group {
			feed.Facets[i].Links = append(feed.Facets[i].Links, link)
			return
		}
	}

	facet.Metadata.Title = group
	facet.Links = append(facet.Links, link)
	feed.Facets = append(feed.Facets, facet)
}

// AddPublicationInGroup smart adding of publication in Group
func (feed *Feed) AddPublicationInGroup(publication Publication, collLink Link) {
	var group Group

	for i, g := range feed.Groups {
		for _, l := range g.Links {
			if l.Href == collLink.Href {
				feed.Groups[i].Publications = append(feed.Groups[i].Publications, publication)
				return
			}
		}
	}

	group.Metadata.Title = collLink.Title
	group.Publications = append(group.Publications, publication)
	group.Links = append(group.Links, Link{Rel: []string{"self"}, Title: collLink.Title, Href: collLink.Href})
	feed.Groups = append(feed.Groups, group)
}

// AddNavigationInGroup add a navigation link to Group
func (feed *Feed) AddNavigationInGroup(link Link, collLink Link) {
	var group Group

	for i, g := range feed.Groups {
		for _, l := range g.Links {
			if l.Href == collLink.Href {
				feed.Groups[i].Navigation = append(feed.Groups[i].Navigation, link)
				return
			}
		}
	}

	group.Metadata.Title = collLink.Title
	group.Navigation = append(group.Navigation, link)
	group.Links = append(group.Links, Link{Rel: []string{"self"}, Title: collLink.Title, Href: collLink.Href})
	feed.Groups = append(feed.Groups, group)
}
