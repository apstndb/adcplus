package signer

import (
	"reflect"
	"testing"
)

func TestIamServiceAccountName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "email is canonicalized",
			input: "sa@project.iam.gserviceaccount.com",
			want:  "projects/-/serviceAccounts/sa@project.iam.gserviceaccount.com",
		},
		{
			name:  "already canonical",
			input: "projects/-/serviceAccounts/sa@project.iam.gserviceaccount.com",
			want:  "projects/-/serviceAccounts/sa@project.iam.gserviceaccount.com",
		},
		{
			name:  "trims whitespace",
			input: "  sa@project.iam.gserviceaccount.com  ",
			want:  "projects/-/serviceAccounts/sa@project.iam.gserviceaccount.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iamServiceAccountName(tt.input); got != tt.want {
				t.Errorf("iamServiceAccountName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIamDelegates(t *testing.T) {
	got := iamDelegates([]string{
		"d1@p.iam.gserviceaccount.com",
		"projects/-/serviceAccounts/d2@p.iam.gserviceaccount.com",
		"   ",
	})
	want := []string{
		"projects/-/serviceAccounts/d1@p.iam.gserviceaccount.com",
		"projects/-/serviceAccounts/d2@p.iam.gserviceaccount.com",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("iamDelegates() = %v, want %v", got, want)
	}
}
