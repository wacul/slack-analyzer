package commands

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
)

func Word(slackToken string, targetChannels []string, targetWords []string, eachReport bool, allReportChannel string, oldest, latest DateTime) error {
	channelFilter := NewStringFilter(targetChannels)

	api := slack.New(slackToken)
	channels, err := api.GetChannels(true)
	if err != nil {
		return fmt.Errorf("failed to get channels: %s", err)
	}

	var allWords StringCounter
	allReportChannelID := ""
	for _, c := range channels {
		if c.Name == allReportChannel {
			allReportChannelID = c.ID
		}

		if !channelFilter.Match(c.Name) {
			continue
		}
		history, err := api.GetChannelHistory(c.ID, slack.HistoryParameters{
			Oldest: fmt.Sprintf("%d.000000", oldest.Time().Unix()),
			Latest: fmt.Sprintf("%d.000000", latest.Time().Unix()),
			Count:  1000,
		})
		if err != nil {
			return fmt.Errorf("failed to get channel %s history: %s", c.Name, err)
		}
		var eachWords StringCounter
		log.Printf("%d messages found in %s", len(history.Messages), c.Name)
		for _, m := range history.Messages {
			for _, w := range targetWords {
				cnt := strings.Count(m.Text, w)
				eachWords.Add(w, cnt)
				allWords.Add(w, cnt)
			}
		}
		if eachReport {
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "Words in this channel from %s to %s\n", oldest.Time().Format("01-02 15:04"), latest.Time().Format("01-02 15:04"))
			fmt.Fprintln(&buf)
			eachWords.Fprint(&buf)
			if _, _, err := api.PostMessage(c.ID, slack.MsgOptionText(buf.String(), false)); err != nil {
				return fmt.Errorf("failed to post result to %s: %s", c.Name, err)
			}
		}
	}
	if allReportChannel != "" {
		if allReportChannelID == "" {
			return fmt.Errorf("failed to find public channel %s", allReportChannel)
		}
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Words in all public channel from %s to %s\n", oldest.Time().Format("01-02 15:04"), latest.Time().Format("01-02 15:04"))
		fmt.Fprintln(&buf)
		allWords.Fprint(&buf)
		if _, _, err := api.PostMessage(allReportChannelID, slack.MsgOptionText(buf.String(), false)); err != nil {
			return fmt.Errorf("failed to post all channel report to %s: %s", allReportChannel, err)
		}
	}
	return nil
}
