package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/nlopes/slack"
)

// nolint
var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	var (
		targetChannels []string
		targetEmoji    []string
		slackToken     string
		eachReport     bool
		allReport      string
	)
	app := kingpin.New("slack-analyzer", "Analyzes messages in slack").Author("wacul").Version(version)
	app.Flag("channel", "Target channels (not specifing, all public channels)").StringsVar(&targetChannels)
	app.Flag("emoji", "Target emoji (not specifing, all emoji)").StringsVar(&targetEmoji)
	app.Flag("slack-token", "Slack App OAuth App Token").Envar("SLACK_TOKEN").Required().StringVar(&slackToken)
	app.Flag("each-channel-report", "Show each channel report").Default("false").BoolVar(&eachReport)
	app.Flag("all-channel-report", "A channel name to post all public channel report").StringVar(&allReport)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	isTargetChannel := func(s string) bool {
		return true
	}
	if len(targetChannels) > 0 {
		m := map[string]struct{}{}
		for _, c := range targetChannels {
			m[c] = struct{}{}
		}
		isTargetChannel = func(name string) bool {
			_, ok := m[name]
			return ok
		}
	}

	isTargetEmoji := func(s string) bool {
		return true
	}
	if len(targetEmoji) > 0 {
		m := map[string]struct{}{}
		for _, c := range targetEmoji {
			m[c] = struct{}{}
		}
		isTargetEmoji = func(name string) bool {
			_, ok := m[name]
			return ok
		}
	}

	api := slack.New(slackToken)
	channels, err := api.GetChannels(true)
	if err != nil {
		log.Fatalf("failed to get channels: %s", err)
	}

	latest := time.Now()
	oldest := time.Now().Add(-7 * 24 * time.Hour)
	allReactions := map[string]int{}
	for _, c := range channels {
		if !isTargetChannel(c.Name) {
			continue
		}
		history, err := api.GetChannelHistory(c.ID, slack.HistoryParameters{
			Oldest: fmt.Sprintf("%d.000000", oldest.Unix()),
			Latest: fmt.Sprintf("%d.000000", latest.Unix()),
			Count:  1000,
		})
		if err != nil {
			log.Fatalf("failed to get channel %s history: %s", c.Name, err)
		}
		eachReactions := map[string]int{}
		log.Printf("%d messages found in %s", len(history.Messages), c.Name)
		for _, m := range history.Messages {
			for _, r := range m.Reactions {
				if !isTargetEmoji(r.Name) {
					continue
				}
				eachReactions[r.Name] += r.Count
				allReactions[r.Name] += r.Count
			}
		}
		if eachReport {
			var index int
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "Emoji reactions in this channel from %s to %s\n", oldest.Format("01-02 15:04"), latest.Format("01-02 15:04"))
			fmt.Fprintln(&buf)
			for name, count := range eachReactions {
				fmt.Fprintf(&buf, ":%s: : %d回", name, count)
				if (index-3)%4 == 0 {
					fmt.Fprintln(&buf)
				} else {
					fmt.Fprint(&buf, "  ")
				}
				index++
			}
			if _, _, err := api.PostMessage(c.ID, slack.MsgOptionText(buf.String(), false)); err != nil {
				log.Fatalf("failed to post result to %s: %s", c.Name, err)
			}
		}
	}
	if allReport != "" {
		var channelID string
		for _, c := range channels {
			if c.Name == allReport {
				channelID = c.ID
			}
		}
		if channelID == "" {
			log.Fatalf("failed to find public channel %s", allReport)
		}
		var index int
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Emoji reactions in all public channel from %s to %s\n", oldest.Format("01-02 15:04"), latest.Format("01-02 15:04"))
		fmt.Fprintln(&buf)
		for name, count := range allReactions {
			fmt.Fprintf(&buf, ":%s: : %d回", name, count)
			if (index-3)%4 == 0 {
				fmt.Fprintln(&buf)
			} else {
				fmt.Fprint(&buf, "  ")
			}
			index++
		}
		if _, _, err := api.PostMessage(channelID, slack.MsgOptionText(buf.String(), false)); err != nil {
			log.Fatalf("failed to post all channel report to %s: %s", allReport, err)
		}
	}
}
