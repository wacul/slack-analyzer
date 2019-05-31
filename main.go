package main

import (
	"log"
	"os"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/slack-stamps/commands"
)

// nolint
var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	app := kingpin.New("slack-analyzer", "Analyzes messages in slack").Author("wacul").Version(version)

	cmds := map[string]func() error{}
	for _, f := range []func(*kingpin.Application) (string, func() error){
		reaction,
		find,
		word,
	} {
		key, run := f(app)
		cmds[key] = run
	}
	if err := cmds[kingpin.MustParse(app.Parse(os.Args[1:]))](); err != nil {
		log.Fatalf("error: %s", err)
	}
}

func reaction(app *kingpin.Application) (string, func() error) {
	var (
		slackToken       string
		targetChannels   []string
		targetEmoji      []string
		eachReport       bool
		allReportChannel string

		oldest commands.DateTime
		latest commands.DateTime
	)

	cmd := app.Command("reaction", "Analyze reactions")
	cmd.Flag("slack-token", "Slack App OAuth App Token").Envar("SLACK_TOKEN").Required().StringVar(&slackToken)
	cmd.Flag("channel", "Target channels (not specifing, all public channels)").StringsVar(&targetChannels)
	cmd.Flag("emoji", "Target emoji (not specifing, all emoji)").StringsVar(&targetEmoji)
	cmd.Flag("each-channel-report", "Show each channel report").Default("false").BoolVar(&eachReport)
	cmd.Flag("all-channel-report", "A channel name to post all public channel report").StringVar(&allReportChannel)
	cmd.Flag("oldest", "A time of the oldest message that should be analyzed").Default(time.Now().Add(-7 * 24 * time.Hour).Format(commands.DateFormat)).SetValue(&oldest)
	cmd.Flag("latest", "A time of the latest message that should be analyzed").Default(time.Now().Format(commands.DateFormat)).SetValue(&latest)

	return cmd.FullCommand(), func() error {
		return commands.Reaction(slackToken, targetChannels, targetEmoji, eachReport, allReportChannel, oldest, latest)
	}
}

func find(app *kingpin.Application) (string, func() error) {
	var (
		slackToken       string
		targetChannels   []string
		targetWord       []string
		eachReport       bool
		allReportChannel string

		oldest commands.DateTime
		latest commands.DateTime
	)
	cmd := app.Command("word", "Analyze words")
	cmd.Flag("slack-token", "Slack App OAuth App Token").Envar("SLACK_TOKEN").Required().StringVar(&slackToken)
	cmd.Flag("channel", "Target channels (not specifing, all public channels)").StringsVar(&targetChannels)
	cmd.Flag("word", "Target word").Required().StringsVar(&targetWord)
	cmd.Flag("each-channel-report", "Show each channel report").Default("false").BoolVar(&eachReport)
	cmd.Flag("all-channel-report", "A channel name to post all public channel report").StringVar(&allReportChannel)
	cmd.Flag("oldest", "A time of the oldest message that should be analyzed").Default(time.Now().Add(-7 * 24 * time.Hour).Format(commands.DateFormat)).SetValue(&oldest)
	cmd.Flag("latest", "A time of the latest message that should be analyzed").Default(time.Now().Format(commands.DateFormat)).SetValue(&latest)

	return cmd.FullCommand(), func() error {
		return commands.Find(slackToken, targetChannels, targetWord, eachReport, allReportChannel, oldest, latest)
	}
}

func word(app *kingpin.Application) (string, func() error) {
	var (
		slackToken       string
		targetChannels   []string
		targetWord       []string
		eachReport       bool
		allReportChannel string

		oldest commands.DateTime
		latest commands.DateTime
	)
	cmd := app.Command("word", "Analyze words")
	cmd.Flag("slack-token", "Slack App OAuth App Token").Envar("SLACK_TOKEN").Required().StringVar(&slackToken)
	cmd.Flag("channel", "Target channels (not specifing, all public channels)").StringsVar(&targetChannels)
	cmd.Flag("word", "Target word").Required().StringsVar(&targetWord)
	cmd.Flag("each-channel-report", "Show each channel report").Default("false").BoolVar(&eachReport)
	cmd.Flag("all-channel-report", "A channel name to post all public channel report").StringVar(&allReportChannel)
	cmd.Flag("oldest", "A time of the oldest message that should be analyzed").Default(time.Now().Add(-7 * 24 * time.Hour).Format(commands.DateFormat)).SetValue(&oldest)
	cmd.Flag("latest", "A time of the latest message that should be analyzed").Default(time.Now().Format(commands.DateFormat)).SetValue(&latest)

	return cmd.FullCommand(), func() error {
		return commands.Word(slackToken, targetChannels, targetWord, eachReport, allReportChannel, oldest, latest)
	}
}
