package tableau

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient_WithUsernamePassword(t *testing.T) {
	t.Parallel()
	config := Config{
		Name:       "test-tableau",
		Host:       "tableau.example.com",
		Username:   "user",
		Password:   "pass",
		SiteID:     "site123",
		APIVersion: "3.4",
	}

	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, config, client.config)
}

func TestNewClient_WithPAT(t *testing.T) {
	t.Parallel()
	config := Config{
		Name:                      "test-tableau",
		Host:                      "tableau.example.com",
		PersonalAccessTokenName:   "my-token",
		PersonalAccessTokenSecret: "my-secret",
		SiteID:                    "site123",
		APIVersion:                "3.4",
	}

	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, config, client.config)
}

func TestNewClient_MissingRequiredFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{
			name: "missing host",
			config: Config{
				Username: "user",
				Password: "pass",
				SiteID:   "site123",
			},
			wantErr: "host is required for Tableau connection",
		},
		{
			name: "missing site_id",
			config: Config{
				Host:     "tableau.example.com",
				Username: "user",
				Password: "pass",
			},
			wantErr: "site_id is required for Tableau connection",
		},
		{
			name: "missing both auth methods",
			config: Config{
				Host:   "tableau.example.com",
				SiteID: "site123",
			},
			wantErr: "either personal access token (name and secret) or username and password are required for Tableau connection",
		},
		{
			name: "incomplete PAT (missing secret)",
			config: Config{
				Host:                    "tableau.example.com",
				SiteID:                  "site123",
				PersonalAccessTokenName: "my-token",
			},
			wantErr: "either personal access token (name and secret) or username and password are required for Tableau connection",
		},
		{
			name: "incomplete username/password (missing password)",
			config: Config{
				Host:     "tableau.example.com",
				SiteID:   "site123",
				Username: "user",
			},
			wantErr: "either personal access token (name and secret) or username and password are required for Tableau connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, err := NewClient(tt.config)
			require.Error(t, err)
			require.Nil(t, client)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestNewClient_DefaultAPIVersion(t *testing.T) {
	t.Parallel()
	config := Config{
		Host:     "tableau.example.com",
		Username: "user",
		Password: "pass",
		SiteID:   "site123",
		// APIVersion not set
	}

	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, "3.4", client.config.APIVersion)
}
