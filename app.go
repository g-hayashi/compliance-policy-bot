package main

import (
	"fmt"
	"os"
	"strings"

	"main/tools"

	"github.com/joho/godotenv"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/slack-go/slack"
)

type Policy struct {
	ID   string
	Desc string
	Name string
}

const (
	COMPLIANCE_POLICY_PREFIX = "test"
	MESSAGE_TITLE            = "このボットは端末 :computer: をチェックし、推奨を提案する社内ツールです。各推奨事項を確認してください。\n\n Hello, This internal bot checks the device and provides our recommendations. Please check each items"
	MESSAGE_FOOTER           = "対応方法など詳細はこちらをご確認ください (https://google.com)\n\nCheck the details here (https://google.com)"
)

func main() {
	graphHelper := tools.NewGraphHelper()
	slackHelper := tools.NewSlackHelper()

	err := Initialize(graphHelper, slackHelper)
	if err != nil {
		fmt.Println("Error Initializing:", err)
		return
	}
	policies, err := GetCompliancePolicies(graphHelper)
	if err != nil {
		fmt.Println("Error getting Policies:", err)
		return
	}
	messages, err := CheckDevices(graphHelper, policies)
	if err != nil {
		fmt.Println("Error checking devices:", err)
		return
	}

	// Send messages to slack
	err = SendMessages(messages, slackHelper)
	if err != nil {
		fmt.Println("Error Sending Slack Message:", err)
		return
	}
}

func Initialize(g *tools.GraphHelper, s *tools.SlackHelper) error {
	godotenv.Load(".env")
	err := godotenv.Load()
	var (
		tenantId, clientId, clientSecret, slackToken string
	)
	if err == nil {
		tenantId = os.Getenv("TENANT_ID")
		clientId = os.Getenv("CLIENT_ID")
		clientSecret = os.Getenv("CLIENT_SECRET")
		slackToken = os.Getenv("SLACK_TOKEN")
	} else {
		var err error
		prefix := "projects/YOUR_PROJECT_ID/secrets/"
		tenantId, err = tools.GetSecret(prefix + "tenantId/versions/latest")
		if err != nil {
			return err
		}
		clientId, err = tools.GetSecret(prefix + "clientId/versions/latest")
		if err != nil {
			return err
		}
		clientSecret, err = tools.GetSecret(prefix + "clientSecret/versions/latest")
		if err != nil {
			return err
		}
		slackToken, err = tools.GetSecret(prefix + "slackToken/versions/latest")
		if err != nil {
			return err
		}
	}
	err = g.InitializeGraph(tenantId, clientId, clientSecret)
	if err != nil {
		return err
	}

	err = s.InitializeSlack(slackToken)
	if err != nil {
		return err
	}

	return nil
}

func GetCompliancePolicies(g *tools.GraphHelper) ([]Policy, error) {
	result, err := g.GetCompliancePolicies()
	if err != nil {
		return nil, err
	}
	var policies []Policy
	for _, policy := range result.GetValue() {
		policies = append(policies, Policy{
			ID:   *policy.GetId(),
			Desc: *policy.GetDescription(),
			Name: *policy.GetDisplayName(),
		})
	}
	return policies, nil
}

func CheckDevices(g *tools.GraphHelper, policies []Policy) (map[string][]string, error) {
	messages := make(map[string][]string)

	for _, policy := range policies {
		// コンプライアンスポリシーに対応するデバイス一覧を取得
		result, err := g.GetDevicesWithCompliancePolicy(policy.ID)
		if err != nil {
			return nil, err
		}
		for _, device := range result.GetValue() {
			// デバイスのコンプライアンス準拠状況をチェック
			if *device.GetStatus() != models.COMPLIANT_COMPLIANCESTATUS {
				continue
			}
			// メッセージの生成
			message := fmt.Sprintf("1. %s \n\n Device Name %s, Model %s\n",
				policy.Desc, *device.GetDeviceDisplayName(), *device.GetDeviceModel())
			owner := *device.GetUserName()
			messages[owner] = append(messages[owner], message)
		}
	}
	return messages, nil
}

func SendMessages(messages map[string][]string, s *tools.SlackHelper) error {
	for email, message := range messages {
		slack_message := slack.MsgOptionBlocks(
			// ------- ここからSlackMessageBlocks -------------
			&slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: MESSAGE_TITLE,
				},
			},
			slack.NewDividerBlock(),
			&slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: strings.Join(message, ","),
				},
			},
			slack.NewDividerBlock(),
			&slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: MESSAGE_FOOTER,
				},
			},
			// ------- ここまでSlackMessageBlocks -------------
		)
		// Send Message
		err := s.SendDM(email, slack_message)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Example of Device Graph API
type Device struct {
	ID           string
	Name         string
	OwnerId      string
	OwnerEmail   string
	OwnerType    models.ManagedDeviceOwnerType
	Serial       string
	Manufacturer string
	isEncryped   bool
}
	func GetDevices(g *tools.GraphHelper) ([]Device, error) {
	result, err := g.GetDevices()
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}
	var devices []Device
	for _, device := range result.GetValue() {
		devices = append(devices, Device{
			ID:           *device.GetId(),
			Name:         *device.GetDeviceName(),
			OwnerId:      *device.GetUserId(),
			OwnerEmail:   *device.GetEmailAddress(),
			OwnerType:    *device.GetManagedDeviceOwnerType(),
			Serial:       *device.GetSerialNumber(),
			Manufacturer: *device.GetManufacturer(),
			isEncryped:   *device.GetIsEncrypted(),
		})
	}
	return devices, nil
}
*/
