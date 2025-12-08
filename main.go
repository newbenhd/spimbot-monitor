package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var standingHTML string
	var gamesHTML string

	err := chromedp.Run(ctx, chromedp.Navigate("https://spimarena.mrsheep.dev/#/"), chromedp.WaitVisible(`#standings`, chromedp.ByQuery), chromedp.OuterHTML(`#standings`, &standingHTML, chromedp.ByQuery), chromedp.OuterHTML(`#games`, &gamesHTML, chromedp.ByQuery))
	if err != nil {
		log.Fatal(err)
	}

	// prep
	standingDoc, err := goquery.NewDocumentFromReader(strings.NewReader(standingHTML))
	if err != nil {
		log.Fatal(err)
	}
	gamesDoc, err := goquery.NewDocumentFromReader(strings.NewReader(gamesHTML))
	if err != nil {
		log.Fatal(err)
	}

	ptCache := make(map[string]int)
	gameCache := make(map[string]int)
	filterCache := make(map[string]bool)

	standingDoc.Find("#standings .tableContainer .team-name").Each(func(i int, tr *goquery.Selection) {
		tr.Find("a").Each(func(j int, td *goquery.Selection) {
			gameCache[td.Text()] = 0
			ptCache[td.Text()] = 0
		})
	})

	gamesDoc.Find("#games tr").Each(func(j int, tr *goquery.Selection) {
		name := tr.Find(".match-name a").Text()
		score := tr.Find(".scores").Text()

		if _, ok := filterCache[name]; ok {
			return
		}
		filterCache[name] = true

		result := strings.Split(score, "-")
		a, err := strconv.Atoi(strings.TrimSpace(result[0]))
		if err != nil {
			log.Fatal(err)
		}
		b, err := strconv.Atoi(strings.TrimSpace(result[1]))
		if err != nil {
			log.Fatal(err)
		}

		players := strings.Split(name, " ")
		playerA := strings.TrimSpace(players[0])
		playerB := strings.TrimSpace(players[1])

		gameCache[playerA] += 1
		gameCache[playerB] += 1

		if a > b {
			ptCache[playerA] += 1
		} else {
			ptCache[playerB] += 1
		}
	})

	type Outcome struct {
		Key   string
		Value float32
	}
	output := []Outcome{}
	for k, v := range ptCache {
		output = append(output, Outcome{Key: k, Value: (float32(v) / float32(gameCache[k])) * 100})
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i].Value > output[j].Value
	})

	// query
	for _, outcome := range output {
		fmt.Printf("%s: %.f%%\n", outcome.Key, outcome.Value)
	}
}
