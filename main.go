package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	title  = flag.NewFlagSet("title", flag.ExitOnError)
	tomdb  = title.String("type", "", "omdb movie, series, episode, game")
	tt     = title.Bool("t", false, "search for movie title")
	tid    = title.Bool("id", false, "search for movie id")
	search = flag.NewFlagSet("search", flag.ExitOnError)
	somdb  = search.String("type", "", "omdb movie, series, episode")
	help   = flag.NewFlagSet("help", flag.ExitOnError)
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected help, title or search subcommands")
	}
	switch os.Args[1] {
	case "help":
		fmt.Println(usage)
	case "title":
		title.Parse(os.Args[2:])
		movie, err := SearchMovie()
		if err != nil {
			log.Fatal(err)
		}
		if err := movie.PrintPosterSave(true); err != nil {
			log.Fatal(err)
		}
	case "search":
		search.Parse(os.Args[2:])
		movies, err := Search()
		if err != nil {
			log.Fatal(err)
		}
		if err := movies.PrintPosterSave(); err != nil {
			log.Fatal(err)
		}
		if err := movies.GenHTMLPage(); err != nil {
			log.Fatal(err)
		}
	}

}
