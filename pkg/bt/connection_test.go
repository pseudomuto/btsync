package bt_test

import (
	"github.com/stretchr/testify/assert"

	"context"
	"testing"

	"github.com/pseudomuto/btsync/pkg/bt"
)

func TestNew(t *testing.T) {
	tests := []struct {
		msg    string
		config bt.ConnectionConfig
		err    string
	}{
		{
			"with a valid JWT",
			bt.ConnectionConfig{
				Project:  "project",
				Instance: "instance",
				JWT: []byte(`{
				"type": "service_account",
				"project_id": "proj",
				"private_key_id": "a4f549e019cec37d00875a7031c2669860e059a1",
				"private_key": "LOL NO PRIVATE KEY",
				"client_id": "116824582041636851187",
				"auth_uri": "https://accounts.google.com/o/oauth2/auth",
				"token_uri": "https://accounts.google.com/o/oauth2/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/proj.iam.gserviceaccount.com"
			}`),
			},
			"",
		},
		{
			"without a JWT",
			bt.ConnectionConfig{
				Project:  "project",
				Instance: "instance",
			},
			"",
		},
		{
			"with an invalid JWT",
			bt.ConnectionConfig{
				Project:  "project",
				Instance: "instance",
				JWT:      []byte("invalid JWT def"),
			},
			"invalid character 'i' looking for beginning of value",
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		conn, err := bt.New(ctx, test.config)

		if test.err == "" {
			assert.NotNil(t, conn, test.msg)
			assert.NoError(t, err, test.msg)
		} else {
			assert.Nil(t, conn, test.msg)
			assert.EqualError(t, err, test.err, test.msg)
		}
	}
}
