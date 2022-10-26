package notify

import (
	"os"

	"github.com/slack-go/slack"
	"gopkg.in/yaml.v2"
)

type SlackNotify struct {
	config slackConfig
	c      *slack.Client
}

type slackConfig struct {
	Token     string `yaml:"token"`
	ChannelID string `yaml:"channelID"`
}

func NewSlackNotify(configPath string) (SlackNotify, error) {
	sn := SlackNotify{}
	buf, err := os.ReadFile(configPath)
	if err != nil {
		return sn, err
	}
	sn.config = slackConfig{}
	err = yaml.Unmarshal(buf, &sn.config)
	if err != nil {
		return sn, err
	}
	sn.c = slack.New(sn.config.Token)
	return sn, nil
}

func (sn *SlackNotify) SendSuccessEvent(title, message string) error {
	err := sn.send(title, message, "good")
	return err
}

func (sn *SlackNotify) SendFailEvent(title, message string) error {
	err := sn.send(title, message, "danger")
	return err
}

func (sn *SlackNotify) send(title, message, color string) error {
	var attachment slack.Attachment
	if message != "" {
		attachment = slack.Attachment{
			Title:   title,
			Pretext: "k8s Event",
			Text:    message,
			Color:   color,
		}
	} else {
		attachment = slack.Attachment{
			Title:   message,
			Pretext: "k8s Event",
			Color:   color,
		}
	}

	_, _, err := sn.c.PostMessage(
		sn.config.ChannelID,
		slack.MsgOptionAttachments(attachment),
	)
	return err
}
