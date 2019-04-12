package commands

import (
	"bytes"
	"fmt"
	"log"

	"github.com/nlopes/slack"
)

func Raw(slackToken string, targetChannel string, oldest, latest DateTime) error {
	channelFilter := NewStringFilter([]string{targetChannel})

	api := slack.New(slackToken)
	channels, err := api.GetChannels(true)
	if err != nil {
		return fmt.Errorf("failed to get channels: %s", err)
	}

	var allReactions StringCounter
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
		var eachReactions StringCounter
		log.Printf("%d messages found in %s", len(history.Messages), c.Name)
		for _, m := range history.Messages {
			for _, r := range m.Reactions {
				if !emojiFilter.Match(r.Name) {
					continue
				}
				eachReactions.Add(":"+r.Name+":", r.Count)
				allReactions.Add(":"+r.Name+":", r.Count)
			}
		}
		if eachReport {
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "Emoji reactions in this channel from %s to %s\n", oldest.Time().Format("01-02 15:04"), latest.Time().Format("01-02 15:04"))
			fmt.Fprintln(&buf)
			eachReactions.Fprint(&buf)
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
		fmt.Fprintf(&buf, "Emoji reactions in all public channel from %s to %s\n", oldest.Time().Format("01-02 15:04"), latest.Time().Format("01-02 15:04"))
		fmt.Fprintln(&buf)
		allReactions.Fprint(&buf)
		if _, _, err := api.PostMessage(allReportChannelID, slack.MsgOptionText(buf.String(), false)); err != nil {
			return fmt.Errorf("failed to post all channel report to %s: %s", allReportChannel, err)
		}
	}
	return nil
}
