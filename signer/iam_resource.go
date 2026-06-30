package signer

import "strings"

func iamServiceAccountName(v string) string {
	v = strings.TrimSpace(v)
	if strings.HasPrefix(v, "projects/") {
		return v
	}
	return "projects/-/serviceAccounts/" + v
}

func iamDelegates(ds []string) []string {
	out := make([]string, 0, len(ds))
	for _, d := range ds {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}
		out = append(out, iamServiceAccountName(d))
	}
	return out
}
