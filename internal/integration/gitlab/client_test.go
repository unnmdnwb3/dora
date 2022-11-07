package gitlab

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

func TestGetProjects(t *testing.T) {
	body := []Project{
		{
			CreatedAt: "2022-10-31T09:35:54.569Z",
			DefaultBranch: "main",
			NamespacedName: "example/backend",
			Url: "https://gitlab.com/example/frontend",
		},
		{
			CreatedAt: "2022-10-31T09:36:21.324Z",
			DefaultBranch: "main",
			NamespacedName: "example/frontend",
			Url: "https://gitlab.com/example/frontend",
		},
	}

	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json, _ := json.Marshal(body)
		w.Write(json)
		
	}))
	defer mock.Close()

	client := Client{
		Auth: "token",
		Uri: mock.URL,
	}

	projects, _ := client.GetProjects()
	require.Equal(t, *projects, body)
}