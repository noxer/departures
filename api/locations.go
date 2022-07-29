package api

import "context"

type Location struct {
	Type     string
	ID       string
	Name     string
	Location struct {
		Type      string
		ID        string
		Latitude  float64
		Longitude float64
	}
	Products    map[string]bool
	StationDHID string
}

type locationsQuery struct {
	c            *Client
	query        string
	fuzzy        bool
	results      int
	stops        bool
	addresses    bool
	poi          bool
	linesOfStops bool
	language     string
}

func (c *Client) Locations(query string) *locationsQuery {
	return &locationsQuery{
		c:            c,
		query:        query,
		fuzzy:        true,
		results:      10,
		stops:        true,
		addresses:    true,
		poi:          true,
		linesOfStops: false,
		language:     "en",
	}
}

func (q *locationsQuery) Do(ctx context.Context) ([]Location, error) {
	const u = "/locations?query=%s&fuzzy=%t&results=%d&stops=%t&addresses=%t&poi=%t&linesOfStops=%t&language=%s&pretty=false"

	var locs []Location
	err := q.c.getJSON(ctx, &locs, u, q.query, q.fuzzy, q.results, q.stops, q.addresses, q.poi, q.linesOfStops, q.language)
	if err != nil {
		return nil, err
	}

	return locs, nil
}

func (q *locationsQuery) Fuzzy(f bool) *locationsQuery {
	q.fuzzy = f
	return q
}

func (q *locationsQuery) Results(r int) *locationsQuery {
	q.results = r
	return q
}

func (q *locationsQuery) Stops(s bool) *locationsQuery {
	q.stops = s
	return q
}

func (q *locationsQuery) Addresses(a bool) *locationsQuery {
	q.addresses = a
	return q
}

func (q *locationsQuery) POI(p bool) *locationsQuery {
	q.poi = p
	return q
}

func (q *locationsQuery) LinesOfStops(l bool) *locationsQuery {
	q.linesOfStops = l
	return q
}

func (q *locationsQuery) Language(l string) *locationsQuery {
	q.language = l
	return q
}
