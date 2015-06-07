package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"database/sql"
	"os"
	"strings"

	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/youtube/v3"
	_ "github.com/go-sql-driver/mysql"
)

var (
	maxResults = flag.Int64("max-results", 50, "Max Youtube results")
	except_words = [...]string{"歌ってみた", "踊ってみた"}
	query_words = [...]string{"花澤香菜", "花澤病"}
)


type YoutubeMovie struct {
	title string
	description string
	thumbnail string
}

func main() {
	flag.Parse()

	developerKey := os.Getenv("DEVELOPER_KEY")
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating Youtube client: %v", err)
	}

	for _, query := range query_words {

		nextPageToken := "hanazawabot_token"

		for nextPageToken != "" {
			if nextPageToken == "hanazawabot_token" {
				nextPageToken = ""
			}
			call := service.Search.List("id,snippet").
				Q(query).
				MaxResults(*maxResults).
				Type("video").
				PageToken(nextPageToken)

			response, err := call.Do()
			if err != nil {
				log.Fatalf("error making search API call: %v", err)
			}

			videos := make(map[string]YoutubeMovie)
			for _, item := range response.Items {
				switch item.Id.Kind {
				case "youtube#video":
					if except_check(item.Snippet.Title) && except_check(item.Snippet.Description) {
						videos[item.Id.VideoId] = YoutubeMovie{item.Snippet.Title, item.Snippet.Description, item.Snippet.Thumbnails.Default.Url}
					}
				}
			}
			nextPageToken = response.NextPageToken

			printIDs("videos", videos)

			db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
			if err != nil {
				panic(err.Error())
			}

			for id, youtube := range videos {

				_, err := db.Exec("insert into youtube_movies (title, movie_id, description, disabled, created_at) values (?, ?, ?, ?, ?)", youtube.title, id, youtube.description, 0, time.Now())
				if err != nil {
					fmt.Printf("mysql connect error: %v \n", err)
				}
			}

			db.Close()
		}
	}
}

func printIDs(sectionName string, matches map[string]YoutubeMovie) {
	fmt.Printf("%v:\n", sectionName)
	for id, youtube := range matches {
		fmt.Printf("[%v] %v : %v \n", id, youtube.title, youtube.description)
	}
	fmt.Printf("\n\n")
}

func except_check(word string) bool {
	for _, except := range except_words {
		if strings.Contains(word, except) {
			return false
		}
	}
	return true
}
