package neighbors

// Neighbor contains info about a neighboring server
type Neighbor struct {
	URL string
}

// Neighbors the configuration of neighboring
// servers
type Neighbors []Neighbor

// New creates new neighbors from the urls
func New(URLs ...string) Neighbors {
	L := len(URLs)

	NS := make(Neighbors, L)

	for i, URL := range URLs {
		NS[i] = Neighbor{
			URL: URL,
		}
	}

	return NS
}
