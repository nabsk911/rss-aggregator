package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nabsk911/rss-aggregator/internal/db"
)

func startScrapping(db *db.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf("Scrapping started with concurrency %d and time between requests %v", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println(err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}

}
func scrapeFeed(database *db.Queries, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()

	_, err := database.MarkFeedAsFetched(context.Background(), feed.ID)

	if err != nil {
		log.Printf("Error marking feed %s as fetched: %v", feed.Url, err)
		return
	}

	feedData, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed %s: %v", feed.Url, err)
		return
	}

	for _, item := range feedData.Channel.Item {

		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}

		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing date %s: %v", item.PubDate, err)
			continue
		}
		_, err = database.CreatePost(context.Background(), db.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Description: description,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			PublishedAt: pubDate,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Error creating post for feed %s: %v", feed.Url, err)
			continue
		}
	}

	log.Printf("Feed %s collected %d items", feed.Name, len(feedData.Channel.Item))
}
