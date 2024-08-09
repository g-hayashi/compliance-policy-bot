package tools

import (
	"fmt"

	"github.com/slack-go/slack"
)

type SlackHelper struct {
	slackClient *slack.Client
}

func NewSlackHelper() *SlackHelper {
	s := &SlackHelper{}
	return s
}

func (s *SlackHelper) InitializeSlack(token string) error {
	s.slackClient = slack.New(token)
	return nil
}

func (s *SlackHelper) GetUserIdFromEmail(email string) (*slack.User, error) {
	user, err := s.slackClient.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (s *SlackHelper) SendDM(email string, message slack.MsgOption) error {
	user, err := s.GetUserIdFromEmail(email)
	if err != nil {
		return fmt.Errorf("getting slack ID from email: %v", err)
	}
	channel, _, _, err := s.slackClient.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{user.ID},
	})
	if err != nil {
		return fmt.Errorf("opening slack conversation: %v", err)
	}
	_, _, err = s.slackClient.PostMessage(
		channel.ID,
		message,
	)
	if err != nil {
		return fmt.Errorf("sending slack DM: %v", err)
	}
	return nil
}
