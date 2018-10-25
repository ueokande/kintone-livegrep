package lgrunner

import "testing"

func TestParseGitHubURL(t *testing.T) {
	cases := []struct {
		url string

		origin string
		owner  string
		repo   string
		e      bool
	}{
		{"https://github.com/kubernetes/client-go", "github.com", "kubernetes", "client-go", false},
	}

	for _, c := range cases {
		origin, owner, repo, err := ParseGitHubURL(c.url)
		if c.e && err == nil {
			t.Error("unexpected err == nil")
			continue
		}
		if !c.e && err != nil {
			t.Error("unexpected err", err)
			continue
		}
		if origin != c.origin || owner != c.owner || repo != c.repo {
			t.Errorf("unpexpected url: origin=%s, owner=%s, repo=%s", origin, owner, repo)
		}

	}
}
