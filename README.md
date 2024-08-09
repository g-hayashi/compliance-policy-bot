# compliancepolicy notify bot
This application retrieves a list of devices and their owners from Intune,
checks against the specified compliance policy for each device,
and sends a direct message via Slack to the owners of devices that are not compliant.

## Setup

### Environment Variables or Google Cloud Secret Manager

Ensure the following environment variables are set on environment:

- `TENANT_ID` : Microsoft azure tenant
- `CLIENT_ID` : Microsoft azure app ID
- `CLIENT_SECRET`: Microsoft azure app secret
- `SLACK_TOKEN` : Bot token

or alternatively, these can be stored to and will attempt to retrive from Google Cloud Secret Manager.

- `projects/<project>/secrets/<name>/versions/latest`

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/g-hayashi/device-bot.git
    cd device-bot
    ```

2. Install the required Go packages:
    ```sh
    go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
    go get github.com/microsoftgraph/msgraph-sdk-go
    go get github.com/joho/godotenv
    go get github.com/slack-go/slack
    go get cloud.google.com/go/secretmanager/apiv1
    ```

3. Set the environment variables (optional):
    write credentials on .env

4. Run the application:
    ```sh
    go run app.go
    ```

## Required Permissions

### Intune and bot scopes

Ensure the azure application has the following permissions:

- `DeviceManagementApps.Read.All`
- `DeviceManagementManagedDevices.Read.All`
- `DeviceManagementConfiguration.Read.All`

Ensure the slack bot to have these permissions:

- `User.read`
- `users:read.email`
- `chat.write`
- `im.write`
