package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type fakeSecretsClient struct {
	output   *secretsmanager.GetSecretValueOutput
	err      error
	secretID string
}

func (f *fakeSecretsClient) GetSecretValue(_ context.Context, input *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	f.secretID = *input.SecretId
	return f.output, f.err
}

func TestSecretValueSDKv2(t *testing.T) {
	text := `{"secrets": ["database"]}`
	for _, tc := range []struct {
		name   string
		output *secretsmanager.GetSecretValueOutput
	}{
		{"string", &secretsmanager.GetSecretValueOutput{SecretString: &text}},
		// SDK v2 already decodes the wire-format base64 into these bytes.
		{"binary", &secretsmanager.GetSecretValueOutput{SecretBinary: []byte(text)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &fakeSecretsClient{output: tc.output}
			manager := &Config{client: client}
			value, err := manager.getSecretValue("application")
			if err != nil || value == nil || *value != text || client.secretID != "application" {
				t.Fatalf("unexpected secret: value=%v error=%v id=%q", value, err, client.secretID)
			}
		})
	}
}

func TestSecretValuePropagatesError(t *testing.T) {
	want := errors.New("access denied")
	manager := &Config{client: &fakeSecretsClient{err: want}}
	value, err := manager.getSecretValue("application")
	if value != nil || !errors.Is(err, want) {
		t.Fatalf("expected original AWS error, got %v, %v", value, err)
	}
}
