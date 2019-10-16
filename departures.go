package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
)

func main() {
	// parse the command line arguments
	var (
		id                string
		filterMode        string
		filterDestination string
		filterLine        string
		width             int
		retries           int
		retryPause        time.Duration
		forceColor        bool
		min               int
		search            string
	)
	flag.StringVar(&id, "id", "900000100003", "ID of the stop")
	flag.StringVar(&filterMode, "filter-mode", "", "Filter the list for this mode of transporation (Comma separated)")
	flag.StringVar(&filterDestination, "filter-destination", "", "Filter the list for this destination (Comma separated)")
	flag.StringVar(&filterLine, "filter-line", "", "Filter the list for this line (Comma separated)")
	flag.IntVar(&width, "width", intEnv("WTF_WIDGET_WIDTH"), "Width of the output")
	flag.IntVar(&retries, "retries", 3, "Number of retries before giving up")
	flag.DurationVar(&retryPause, "retry-pause", time.Second, "Pause between retries")
	flag.IntVar(&min, "min", 60, "Number of minutes you want to see the departures for")
	flag.BoolVar(&forceColor, "force-color", false, "Use this flag to enforce color output even if the terminal does not report support")
	flag.StringVar(&search, "search", "", "Search for the stop name to get the stop ID")
	flag.Parse()

	// ensure valid retry values
	if retries < 0 {
		retries = 0
	}
	if retryPause < 0 {
		retryPause = 0
	}

	var err error

	if search != "" {
		var stations []station
		err = getJSON(&stations, "https://2.bvg.transport.rest/locations?query=%s&poi=false&addresses=false", search)
		if err != nil {
			fmt.Println("Could not query stations")
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Found %d station(s):\n", len(stations))
		for _, s := range stations {
			fmt.Printf("  %s - %s\n", s.ID, s.Name)
		}

		return
	}

	// set the color mode
	color.NoColor = color.NoColor && !forceColor

	// request the departures
	var deps []result
	for i := 0; i < retries+1; i++ {
		err = getJSON(&deps, "https://2.bvg.transport.rest/stations/%s/departures?duration=%d", id, min)
		if err == nil {
			break
		}
		time.Sleep(retryPause)
	}
	if err != nil {
		fmt.Println("Could not query departures")
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// initialize the filters
	fm := filterMap(filterMode)
	fd := filterMap(filterDestination)
	fl := filterMap(filterLine)

	// calculate the length of the columns
	lenName := 0
	lenDir := 0
	lenDep := 0
	from := time.Now().Add(-2 * time.Minute)
	until := time.Now().Add(time.Hour)
	for _, dep := range deps {
		if dep.When.Before(from) || dep.When.After(until) {
			continue
		}
		if fm != nil && !fm[dep.Line.Product] {
			continue
		}
		if fd != nil && !fd[dep.Direction] {
			continue
		}
		if fl != nil && !fl[dep.Line.Name] {
			continue
		}

		lenName = maxStringLen(dep.Line.Name, lenName)
		lenDir = maxStringLen(dep.Direction, lenDir)
		lenDep = maxStringLen(departureTime(dep), lenDep)
	}

	// adjust the column length
	if width > 0 {
		lenDir = width - lenName - lenDep - 2
		if lenDir < 1 {
			lenDir = 1
		}
	}

	// render the columns
	for _, dep := range deps {
		if dep.When.Before(from) || dep.When.After(until) {
			continue
		}
		if fm != nil && !fm[dep.Line.Product] {
			continue
		}
		if fd != nil && !fd[dep.Direction] {
			continue
		}
		if fl != nil && !fl[dep.Line.Name] {
			continue
		}

		departureColor := color.HiGreenString
		if dep.Delay > 0 {
			departureColor = color.RedString
		} else if dep.Delay < 0 {
			departureColor = color.YellowString
		}

		fmt.Println(
			color.WhiteString("%s", leftPad(dep.Line.Name, lenName)),
			rightPad(dep.Direction, lenDir),
			departureColor("%s", departureTime(dep)),
		)
	}
}

func getJSON(v interface{}, urlFormat string, values ...interface{}) error {
	resp, err := http.Get(fmt.Sprintf(urlFormat, values...))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	return d.Decode(v)
}

type result struct {
	TripID string `json:"tripId"`
	Stop   struct {
		Type     string `json:"type"`
		ID       string `json:"id"`
		Name     string `json:"name"`
		Location struct {
			Type      string  `json:"type"`
			ID        string  `json:"id"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
		Products struct {
			Suburban bool `json:"suburban"`
			Subway   bool `json:"subway"`
			Tram     bool `json:"tram"`
			Bus      bool `json:"bus"`
			Ferry    bool `json:"ferry"`
			Express  bool `json:"express"`
			Regional bool `json:"regional"`
		} `json:"products"`
	} `json:"stop"`
	When      time.Time `json:"when"`
	Direction string    `json:"direction"`
	Line      struct {
		Type     string `json:"type"`
		ID       string `json:"id"`
		FahrtNr  string `json:"fahrtNr"`
		Name     string `json:"name"`
		Public   bool   `json:"public"`
		Mode     string `json:"mode"`
		Product  string `json:"product"`
		Operator struct {
			Type string `json:"type"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"operator"`
		Symbol  string `json:"symbol"`
		Nr      int    `json:"nr"`
		Metro   bool   `json:"metro"`
		Express bool   `json:"express"`
		Night   bool   `json:"night"`
	} `json:"line"`
	Remarks []struct {
		Type string `json:"type"`
		Code string `json:"code"`
		Text string `json:"text"`
	} `json:"remarks"`
	Delay    int    `json:"delay"`
	Platform string `json:"platform"`
}

type station struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location struct {
		Type      string  `json:"type"`
		ID        string  `json:"id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Products struct {
		Suburban bool `json:"suburban"`
		Subway   bool `json:"subway"`
		Tram     bool `json:"tram"`
		Bus      bool `json:"bus"`
		Ferry    bool `json:"ferry"`
		Express  bool `json:"express"`
		Regional bool `json:"regional"`
	} `json:"products"`
}

func leftPad(s string, l int) string {
	r := []rune(s)
	if len(r) >= l {
		return string(r[:l])
	}

	return strings.Repeat(" ", l-len(r)) + string(r)
}

func rightPad(s string, l int) string {
	r := []rune(s)
	if len(r) >= l {
		return string(r[:l])
	}

	return string(r) + strings.Repeat(" ", l-len(r))
}

func departureTime(r result) string {
	if r.Delay == 0 {
		return r.When.Format("15:04")
	}
	return fmt.Sprintf("%s (%+d)", r.When.Format("15:04"), r.Delay/60)
}

func filterMap(filter string) map[string]bool {
	if filter == "" {
		return nil
	}

	fs := strings.Split(filter, ",")
	fm := make(map[string]bool, len(fs))
	for _, f := range fs {
		fm[strings.TrimSpace(f)] = true
	}
	return fm
}

func maxStringLen(s string, l int) int {
	c := utf8.RuneCountInString(s)
	if c > l {
		return c
	}
	return l
}

func intEnv(key string) int {
	i, _ := strconv.Atoi(os.Getenv(key))
	return i
}
