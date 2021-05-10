package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
)

const (
	omdbURL = "https://www.omdbapi.com/?"
	usage   = `Usage:
omdb search terminator
omdb search  -type movie terminator
omdb search  -type series terminator
omdb title -t terminator
omdb title terminator
omdb title -type series terminator
omdb title -id tt0851851`
)

var (
	apikey = os.Getenv("OMDBAPIKEY")
)

type Query struct {
	Vals url.Values
}

func NewQuery() *Query {
	q := new(Query)
	q.Vals = make(url.Values)
	q.Vals.Add("apikey", apikey)
	q.Vals.Add("page", "1")
	switch {
	case help.Parsed():

	case title.Parsed():
		if *tomdb != "" {
			q.Vals.Add("type", *tomdb)
		}
		switch {
		case *tt:
			q.Vals.Add("t", url.QueryEscape(strings.Join(title.Args(), " ")))
		case *tid:
			q.Vals.Add("i", url.QueryEscape(strings.Join(title.Args(), " ")))
		default:
			q.Vals.Add("t", url.QueryEscape(strings.Join(title.Args(), " ")))
		}
	case search.Parsed():
		if *somdb != "" {
			q.Vals.Add("type", *somdb)
		}
		q.Vals.Add("s", url.QueryEscape(strings.Join(search.Args(), " ")))
	}
	return q
}

func (q *Query) Next() {
	page, _ := strconv.Atoi(q.Vals.Get("page"))
	page++
	q.Vals.Del("page")
	q.Vals.Add("page", strconv.Itoa(page))
}

type Movie struct {
	Title     string
	Year      string
	Genre     string
	Runtime   string
	Ratings   []Ratings `json:"Ratings,omitempty"`
	Poster    string
	Website   string
	BoxOffice string
	FileName  string `json:"SaveFileName,omitempty"`
	Type      string
	IMDBID    string `json:"imdbID"`
}

type Ratings struct {
	Source string
	Value  string
}

func (m *Movie) String() string {
	data, _ := json.MarshalIndent(m, "", "\t")
	return string(data)
}

func (m *Movie) HasPoster() bool {
	return m.Poster != "" && m.Poster != "N/A"
}

func (m *Movie) HasTitle() bool {
	return m.Title != ""
}

func (m *Movie) TitleSplit() []string {
	return strings.Split(m.Title, " ")
}

func (m *Movie) DnldPoster() error {
	resp, err := http.Get(m.Poster)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d", resp.StatusCode)
	}
	repl := strings.NewReplacer(" ", "", "/", "", ":", "", "?", "")
	fname := repl.Replace(m.Title) + "--" + m.IMDBID
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	io.Copy(f, resp.Body)
	m.FileName = fname
	return nil
}

func (m *Movie) PrintPosterSave(print bool) error {
	if m.HasTitle() && print {
		fmt.Println(m)
	}
	if m.HasPoster() {
		if err := m.DnldPoster(); err != nil {
			return err
		}
	}
	return nil
}

type Movies struct {
	Search   []*Movie
	Response string
}

func (m *Movies) String() string {
	data, _ := json.MarshalIndent(m, "", "  ")
	return string(data)
}

func (m *Movies) PrintPosterSave() error {
	for _, movie := range m.Search {
		if err := movie.PrintPosterSave(false); err != nil {
			return err
		}
		log.Printf("Processing %v", movie.Title)
	}
	return nil
}

func (m *Movies) GenHTMLPage() error {
	report := template.Must(template.New("moviegallery").
		Funcs(template.FuncMap{"hasPoster": hasPoster}).
		Parse(templ))
	f, err := os.Create("films.html")
	if err != nil {
		return err
	}
	defer f.Close()
	if err := report.Execute(f, m); err != nil {
		return err
	}
	return nil
}

func SearchMovie() (*Movie, error) {
	resp, err := http.Get(omdbURL + NewQuery().Vals.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	m := new(Movie)
	if err := json.NewDecoder(resp.Body).Decode(m); err != nil {
		return nil, err
	}
	return m, nil
}

func Search() (*Movies, error) {
	m := new(Movies)
	q := NewQuery()
	for {
		n := new(Movies)
		resp, err := http.Get(omdbURL + q.Vals.Encode())
		if err != nil {
			return m, err
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&n); err != nil {
			return m, err
		}
		if n.Response == "False" {
			break
		}
		m.Search = append(m.Search, n.Search...)
		q.Next()
	}
	return m, nil
}
