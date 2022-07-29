package api

import (
	"context"
	"strconv"
)

type StopNearby struct {
	Type     string
	ID       string
	Name     string
	Location struct {
		Type      string
		ID        string
		Latitude  float64
		Longitude float64
	}
	Products map[string]bool
	Distance int
}

type stopsNearbyQuery struct {
	c            *Client
	latitude     float64
	longitude    float64
	results      int
	distance     int
	stops        bool
	poi          bool
	linesOfStops bool
	language     string
}

func (c *Client) StopsNearby(latitude, longitude float64) *stopsNearbyQuery {
	return &stopsNearbyQuery{
		c:            c,
		latitude:     latitude,
		longitude:    longitude,
		results:      8,
		distance:     -1,
		stops:        true,
		poi:          false,
		linesOfStops: false,
		language:     "en",
	}
}

func (q *stopsNearbyQuery) Do(ctx context.Context) ([]StopNearby, error) {
	const u = "/stops/nearby?latitude=%f&longitude=%f&results=%d&distance=%s&stops=%t&poi=%t&linesOfStops=%t&language=%s&pretty=false"

	distance := ""
	if q.distance > 0 {
		distance = strconv.Itoa(q.distance)
	}

	var stops []StopNearby
	err := q.c.getJSON(ctx, &stops, u, q.latitude, q.longitude, q.results, distance, q.stops, q.poi, q.linesOfStops, q.language)
	if err != nil {
		return nil, err
	}

	return stops, nil
}

func (q *stopsNearbyQuery) Results(r int) *stopsNearbyQuery {
	q.results = r
	return q
}

func (q *stopsNearbyQuery) Distance(d int) *stopsNearbyQuery {
	q.distance = d
	return q
}

func (q *stopsNearbyQuery) Stops(s bool) *stopsNearbyQuery {
	q.stops = s
	return q
}

func (q *stopsNearbyQuery) POI(p bool) *stopsNearbyQuery {
	q.poi = p
	return q
}

func (q *stopsNearbyQuery) LinesOfStops(l bool) *stopsNearbyQuery {
	q.linesOfStops = l
	return q
}

func (q *stopsNearbyQuery) Language(l string) *stopsNearbyQuery {
	q.language = l
	return q
}
