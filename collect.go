package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	rss "github.com/ungerik/go-rss"
)

// dayrate is a result of parsing rss item entry,
// Rates is a map of Currency as a key and Rate as a value.
type dayrate struct {
	Day   time.Time
	Rates map[string]float64
}

func collect(c *cli.Context) error {
	db, err := sqlx.Connect("mysql", c.GlobalString("database"))
	if err != nil {
		return err
	}

	log.Println("Fetching RSS feed from", c.String("source"))
	channel, err := rss.Read(c.String("source"))
	if err != nil {
		return err
	}

	for _, item := range channel.Item {
		day, err := parseDay(item)
		if err != nil {
			return err
		}
		tx := db.MustBegin()
		for cur, rate := range day.Rates {
			tx.MustExec("INSERT IGNORE INTO rates (day, currency, rate) VALUES (?, ?, ?)", day.Day, cur, rate)
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	if err := db.Close(); err != nil {
		return err
	}
	log.Println("Finished processing", len(channel.Item), "days")
	return nil
}

func parseDay(item rss.Item) (*dayrate, error) {
	date, err := item.PubDate.Parse()
	if err != nil {
		return nil, err
	}
	rates, err := parseRates(item.Description)
	if err != nil {
		return nil, err
	}
	return &dayrate{Day: date, Rates: rates}, nil
}

func parseRates(s string) (map[string]float64, error) {
	things := strings.Split(strings.TrimSpace(s), " ")

	if len(things)%2 != 0 {
		// Odd items, something went horribly wrong with the data???
		return nil, errors.New("received odd data, something is horribly wrong, check source for validity")
	}

	rates := make(map[string]float64, len(things)/2)

	for i := 0; i < len(things)/2; i++ {
		f, err := strconv.ParseFloat(things[i*2+1], 64)
		if err != nil {
			return nil, err
		}
		rates[things[i*2]] = f
	}
	return rates, nil
}
