package r2

import (
	"net/http"
	"net/url"
	"testing"
)

type fakeResp struct{}

func (f fakeResp) Header() http.Header       { return http.Header{} }
func (f fakeResp) Write([]byte) (int, error) { return 0, nil }
func (f fakeResp) WriteHeader(int)           {}

var handlerId int
var pathVars Path

func f1(e *Env) {
	handlerId = 1
	pathVars = e.Path
}
func f2(e *Env) {
	handlerId = 2
	pathVars = e.Path
}
func f3(e *Env) {
	handlerId = 3
	pathVars = e.Path
}

type endpoint struct {
	method  string
	path    string
	handler Handler
}

// given a method and path, expect to reach funcId with pathVars
type question struct {
	method    string
	path      string
	handlerId int
	pathVar   Path
}

// List of Github API endpoints is
// Copyright (c) 2013 Julien Schmidt. All rights reserved.

var githubAPI = []*endpoint{
	// OAuth Authorizations
	{"GET", "/authorizations", f1},
	{"GET", "/authorizations/:id", f1},
	{"POST", "/authorizations", f1},
	{"PUT", "/authorizations/clients/:client_id", f1}, // should fail as :id and clients conflict
	{"PATCH", "/authorizations/:id", f1},
	{"DELETE", "/authorizations/:id", f1},
	{"GET", "/applications/:client_id/tokens/:access_token", f1},
	{"DELETE", "/applications/:client_id/tokens", f1},
	{"DELETE", "/applications/:client_id/tokens/:access_token", f1},

	// Activity
	{"GET", "/events", f1},
	{"GET", "/repos/:owner/:repo/events", f1},
	{"GET", "/networks/:owner/:repo/events", f1},
	{"GET", "/orgs/:org/events", f1},
	{"GET", "/users/:user/received_events", f1},
	{"GET", "/users/:user/received_events/public", f1},
	{"GET", "/users/:user/events", f1},
	{"GET", "/users/:user/events/public", f1},
	{"GET", "/users/:user/events/orgs/:org", f1},
	{"GET", "/feeds", f1},
	{"GET", "/notifications", f1},
	{"GET", "/repos/:owner/:repo/notifications", f1},
	{"PUT", "/notifications", f1},
	{"PUT", "/repos/:owner/:repo/notifications", f1},
	{"GET", "/notifications/threads/:id", f1},
	{"PATCH", "/notifications/threads/:id", f1},
	{"GET", "/notifications/threads/:id/subscription", f1},
	{"PUT", "/notifications/threads/:id/subscription", f1},
	{"DELETE", "/notifications/threads/:id/subscription", f1},
	{"GET", "/repos/:owner/:repo/stargazers", f1},
	{"GET", "/users/:user/starred", f1},
	{"GET", "/user/starred", f1},
	{"GET", "/user/starred/:owner/:repo", f1},
	{"PUT", "/user/starred/:owner/:repo", f1},
	{"DELETE", "/user/starred/:owner/:repo", f1},
	{"GET", "/repos/:owner/:repo/subscribers", f1},
	{"GET", "/users/:user/subscriptions", f1},
	{"GET", "/user/subscriptions", f1},
	{"GET", "/repos/:owner/:repo/subscription", f1},
	{"PUT", "/repos/:owner/:repo/subscription", f1},
	{"DELETE", "/repos/:owner/:repo/subscription", f1},
	{"GET", "/user/subscriptions/:owner/:repo", f1},
	{"PUT", "/user/subscriptions/:owner/:repo", f1},
	{"DELETE", "/user/subscriptions/:owner/:repo", f1},

	// Gists
	{"GET", "/users/:user/gists", f1},
	{"GET", "/gists", f1},
	{"GET", "/gists/public", f1},
	{"GET", "/gists/starred", f1},
	{"GET", "/gists/:id", f1},
	{"POST", "/gists", f1},
	{"PATCH", "/gists/:id", f1},
	{"PUT", "/gists/:id/star", f1},
	{"DELETE", "/gists/:id/star", f1},
	{"GET", "/gists/:id/star", f1},
	{"POST", "/gists/:id/forks", f1},
	{"DELETE", "/gists/:id", f1},

	// Git Data
	{"GET", "/repos/:owner/:repo/git/blobs/:sha", f1},
	{"POST", "/repos/:owner/:repo/git/blobs", f1},
	{"GET", "/repos/:owner/:repo/git/commits/:sha", f1},
	{"POST", "/repos/:owner/:repo/git/commits", f1},
	{"GET", "/repos/:owner/:repo/git/refs/*ref", f1},
	{"GET", "/repos/:owner/:repo/git/refs", f1},
	{"POST", "/repos/:owner/:repo/git/refs", f1},
	{"PATCH", "/repos/:owner/:repo/git/refs/*ref", f1},
	{"DELETE", "/repos/:owner/:repo/git/refs/*ref", f1},
	{"GET", "/repos/:owner/:repo/git/tags/:sha", f1},
	{"POST", "/repos/:owner/:repo/git/tags", f1},
	{"GET", "/repos/:owner/:repo/git/trees/:sha", f1},
	{"POST", "/repos/:owner/:repo/git/trees", f1},

	// Issues
	{"GET", "/issues", f1},
	{"GET", "/user/issues", f1},
	{"GET", "/orgs/:org/issues", f1},
	{"GET", "/repos/:owner/:repo/issues", f1},
	{"GET", "/repos/:owner/:repo/issues/:number", f1},
	{"POST", "/repos/:owner/:repo/issues", f1},
	{"PATCH", "/repos/:owner/:repo/issues/:number", f1},
	{"GET", "/repos/:owner/:repo/assignees", f1},
	{"GET", "/repos/:owner/:repo/assignees/:assignee", f1},
	{"GET", "/repos/:owner/:repo/issues/:number/comments", f1},
	{"GET", "/repos/:owner/:repo/issues/comments", f1},
	{"GET", "/repos/:owner/:repo/issues/comments/:id", f1},
	{"POST", "/repos/:owner/:repo/issues/:number/comments", f1},
	{"PATCH", "/repos/:owner/:repo/issues/comments/:id", f1},
	{"DELETE", "/repos/:owner/:repo/issues/comments/:id", f1},
	{"GET", "/repos/:owner/:repo/issues/:number/events", f1},
	{"GET", "/repos/:owner/:repo/issues/events", f1},
	{"GET", "/repos/:owner/:repo/issues/events/:id", f1},
	{"GET", "/repos/:owner/:repo/labels", f1},
	{"GET", "/repos/:owner/:repo/labels/:name", f1},
	{"POST", "/repos/:owner/:repo/labels", f1},
	{"PATCH", "/repos/:owner/:repo/labels/:name", f1},
	{"DELETE", "/repos/:owner/:repo/labels/:name", f1},
	{"GET", "/repos/:owner/:repo/issues/:number/labels", f1},
	{"POST", "/repos/:owner/:repo/issues/:number/labels", f1},
	{"DELETE", "/repos/:owner/:repo/issues/:number/labels/:name", f1},
	{"PUT", "/repos/:owner/:repo/issues/:number/labels", f1},
	{"DELETE", "/repos/:owner/:repo/issues/:number/labels", f1},
	{"GET", "/repos/:owner/:repo/milestones/:number/labels", f1},
	{"GET", "/repos/:owner/:repo/milestones", f1},
	{"GET", "/repos/:owner/:repo/milestones/:number", f1},
	{"POST", "/repos/:owner/:repo/milestones", f1},
	{"PATCH", "/repos/:owner/:repo/milestones/:number", f1},
	{"DELETE", "/repos/:owner/:repo/milestones/:number", f1},

	// Miscellaneous
	{"GET", "/emojis", f1},
	{"GET", "/gitignore/templates", f1},
	{"GET", "/gitignore/templates/:name", f1},
	{"POST", "/markdown", f1},
	{"POST", "/markdown/raw", f1},
	{"GET", "/meta", f1},
	{"GET", "/rate_limit", f1},

	// Organizations
	{"GET", "/users/:user/orgs", f1},
	{"GET", "/user/orgs", f1},
	{"GET", "/orgs/:org", f1},
	{"PATCH", "/orgs/:org", f1},
	{"GET", "/orgs/:org/members", f1},
	{"GET", "/orgs/:org/members/:user", f1},
	{"DELETE", "/orgs/:org/members/:user", f1},
	{"GET", "/orgs/:org/public_members", f1},
	{"GET", "/orgs/:org/public_members/:user", f1},
	{"PUT", "/orgs/:org/public_members/:user", f1},
	{"DELETE", "/orgs/:org/public_members/:user", f1},
	{"GET", "/orgs/:org/teams", f1},
	{"GET", "/teams/:id", f1},
	{"POST", "/orgs/:org/teams", f1},
	{"PATCH", "/teams/:id", f1},
	{"DELETE", "/teams/:id", f1},
	{"GET", "/teams/:id/members", f1},
	{"GET", "/teams/:id/members/:user", f1},
	{"PUT", "/teams/:id/members/:user", f1},
	{"DELETE", "/teams/:id/members/:user", f1},
	{"GET", "/teams/:id/repos", f1},
	{"GET", "/teams/:id/repos/:owner/:repo", f1},
	{"PUT", "/teams/:id/repos/:owner/:repo", f1},
	{"DELETE", "/teams/:id/repos/:owner/:repo", f1},
	{"GET", "/user/teams", f1},

	// Pull Requests
	{"GET", "/repos/:owner/:repo/pulls", f1},
	{"GET", "/repos/:owner/:repo/pulls/:number", f1},
	{"POST", "/repos/:owner/:repo/pulls", f1},
	{"PATCH", "/repos/:owner/:repo/pulls/:number", f1},
	{"GET", "/repos/:owner/:repo/pulls/:number/commits", f1},
	{"GET", "/repos/:owner/:repo/pulls/:number/files", f1},
	{"GET", "/repos/:owner/:repo/pulls/:number/merge", f1},
	{"PUT", "/repos/:owner/:repo/pulls/:number/merge", f1},
	{"GET", "/repos/:owner/:repo/pulls/:number/comments", f1},
	{"GET", "/repos/:owner/:repo/pulls/comments", f1},
	{"GET", "/repos/:owner/:repo/pulls/comments/:number", f1},
	{"PUT", "/repos/:owner/:repo/pulls/:number/comments", f1},
	{"PATCH", "/repos/:owner/:repo/pulls/comments/:number", f1},
	{"DELETE", "/repos/:owner/:repo/pulls/comments/:number", f1},

	// Repositories
	{"GET", "/user/repos", f1},
	{"GET", "/users/:user/repos", f1},
	{"GET", "/orgs/:org/repos", f1},
	{"GET", "/repositories", f1},
	{"POST", "/user/repos", f1},
	{"POST", "/orgs/:org/repos", f1},
	{"GET", "/repos/:owner/:repo", f1},
	{"PATCH", "/repos/:owner/:repo", f1},
	{"GET", "/repos/:owner/:repo/contributors", f1},
	{"GET", "/repos/:owner/:repo/languages", f1},
	{"GET", "/repos/:owner/:repo/teams", f1},
	{"GET", "/repos/:owner/:repo/tags", f1},
	{"GET", "/repos/:owner/:repo/branches", f1},
	{"GET", "/repos/:owner/:repo/branches/:branch", f1},
	{"DELETE", "/repos/:owner/:repo", f1},
	{"GET", "/repos/:owner/:repo/collaborators", f1},
	{"GET", "/repos/:owner/:repo/collaborators/:user", f1},
	{"PUT", "/repos/:owner/:repo/collaborators/:user", f1},
	{"DELETE", "/repos/:owner/:repo/collaborators/:user", f1},
	{"GET", "/repos/:owner/:repo/comments", f1},
	{"GET", "/repos/:owner/:repo/commits/:sha/comments", f1},
	{"POST", "/repos/:owner/:repo/commits/:sha/comments", f1},
	{"GET", "/repos/:owner/:repo/comments/:id", f1},
	{"PATCH", "/repos/:owner/:repo/comments/:id", f1},
	{"DELETE", "/repos/:owner/:repo/comments/:id", f1},
	{"GET", "/repos/:owner/:repo/commits", f1},
	{"GET", "/repos/:owner/:repo/commits/:sha", f1},
	{"GET", "/repos/:owner/:repo/readme", f1},
	{"GET", "/repos/:owner/:repo/contents/*path", f1},
	{"PUT", "/repos/:owner/:repo/contents/*path", f1},
	{"DELETE", "/repos/:owner/:repo/contents/*path", f1},
	{"GET", "/repos/:owner/:repo/:archive_format/:ref", f1},
	{"GET", "/repos/:owner/:repo/keys", f1},
	{"GET", "/repos/:owner/:repo/keys/:id", f1},
	{"POST", "/repos/:owner/:repo/keys", f1},
	{"PATCH", "/repos/:owner/:repo/keys/:id", f1},
	{"DELETE", "/repos/:owner/:repo/keys/:id", f1},
	{"GET", "/repos/:owner/:repo/downloads", f1},
	{"GET", "/repos/:owner/:repo/downloads/:id", f1},
	{"DELETE", "/repos/:owner/:repo/downloads/:id", f1},
	{"GET", "/repos/:owner/:repo/forks", f1},
	{"POST", "/repos/:owner/:repo/forks", f1},
	{"GET", "/repos/:owner/:repo/hooks", f1},
	{"GET", "/repos/:owner/:repo/hooks/:id", f1},
	{"POST", "/repos/:owner/:repo/hooks", f1},
	{"PATCH", "/repos/:owner/:repo/hooks/:id", f1},
	{"POST", "/repos/:owner/:repo/hooks/:id/tests", f1},
	{"DELETE", "/repos/:owner/:repo/hooks/:id", f1},
	{"POST", "/repos/:owner/:repo/merges", f1},
	{"GET", "/repos/:owner/:repo/releases", f1},
	{"GET", "/repos/:owner/:repo/releases/:id", f1},
	{"POST", "/repos/:owner/:repo/releases", f1},
	{"PATCH", "/repos/:owner/:repo/releases/:id", f1},
	{"DELETE", "/repos/:owner/:repo/releases/:id", f1},
	{"GET", "/repos/:owner/:repo/releases/:id/assets", f1},
	{"GET", "/repos/:owner/:repo/stats/contributors", f1},
	{"GET", "/repos/:owner/:repo/stats/commit_activity", f1},
	{"GET", "/repos/:owner/:repo/stats/code_frequency", f1},
	{"GET", "/repos/:owner/:repo/stats/participation", f1},
	{"GET", "/repos/:owner/:repo/stats/punch_card", f1},
	{"GET", "/repos/:owner/:repo/statuses/:ref", f1},
	{"POST", "/repos/:owner/:repo/statuses/:ref", f1},

	// Search
	{"GET", "/search/repositories", f1},
	{"GET", "/search/code", f1},
	{"GET", "/search/issues", f1},
	{"GET", "/search/users", f1},
	{"GET", "/legacy/issues/search/:owner/:repository/:state/:keyword", f1},
	{"GET", "/legacy/repos/search/:keyword", f1},
	{"GET", "/legacy/user/search/:keyword", f1},
	{"GET", "/legacy/user/email/:email", f1},

	// Users
	{"GET", "/users/:user", f1},
	{"GET", "/user", f1},
	{"PATCH", "/user", f1},
	{"GET", "/users", f1},
	{"GET", "/user/emails", f1},
	{"POST", "/user/emails", f1},
	{"DELETE", "/user/emails", f1},
	{"GET", "/users/:user/followers", f1},
	{"GET", "/user/followers", f1},
	{"GET", "/users/:user/following", f1},
	{"GET", "/user/following", f1},
	{"GET", "/user/following/:user", f1},
	{"GET", "/users/:user/following/:target_user", f1},
	{"PUT", "/user/following/:user", f1},
	{"DELETE", "/user/following/:user", f1},
	{"GET", "/users/:user/keys", f1},
	{"GET", "/user/keys", f1},
	{"GET", "/user/keys/:id", f1},
	{"POST", "/user/keys", f1},
	{"PATCH", "/user/keys/:id", f1},
	{"DELETE", "/user/keys/:id", f1},
}

var simpleAPI = []*endpoint{
	{"GET", "/", f1},
	{"GET", "/users", f1},
	{"GET", "/users/:user", f1},
	{"GET", "/users/:user/starred", f1},
	{"POST", "/users/:user/starred", f2},
	{"GET", "/users/topusers", f1},
	// {"GET", "/users/:invalid", f1},

	// {"GET", "/repos", f2},
	// {"GET", "/repos/best", f2},
	// {"GET", "/repos/:owner", f2},
	// {"GET", "/repos/:owner/:repo/stats", f1},

	{"GET", `/regex/integer/:n!\d+`, f3},
	{"GET", `/regex/lowercase/:word![a-z]+`, f3},
	{"PUT", `/regex/lowercase/:word![a-z]+`, f3},
}

var simpleApiQuestions = []*question{
	{"GET", "", 1, Path{}},
	{"GET", "/", 1, Path{}},
	{"POST", "/", 0, Path{}},
	{"GET", "/users", 1, Path{}},
	{"GET", "/users/", 1, Path{}},
	{"GET", "/users/dave", 1, Path{"user": "dave"}},
	{"GET", "/users/bill", 1, Path{"user": "bill"}},
	{"GET", "/users/bill/starred", 1, Path{"user": "bill"}},
	{"GET", "/users/topusers", 1, Path{}},
	{"DELETE", "/users/topusers/hey", 0, Path{}},

	{"GET", "/regex/integer/invalid", 0, Path{}},
	{"GET", "/regex/integer/999", 3, Path{}},
	{"GET", "/regex/lowercase/INVALID", 0, Path{}},
	{"GET", "/regex/lowercase/a", 3, Path{}},
	{"PUT", "/regex/lowercase/bb", 3, Path{}},
}

func makeRequest(method, path string) *http.Request {
	return &http.Request{Method: method, URL: &url.URL{Path: path}}
}

func apiTest(endpoints []*endpoint, questions []*question, t *testing.T) {
	r := NewRouter("")
	for _, ep := range endpoints {
		r.Route(ep.method, ep.path, ep.handler)
	}
	printTree("", "/", r.root, true)

	for _, ques := range questions {
		// reset pathVars and value of handlerId
		pathVars = make(Path)
		handlerId = 0

		req := makeRequest(ques.method, ques.path)
		r.ServeHTTP(fakeResp{}, req)

		if handlerId != ques.handlerId {
			t.Errorf("%v expected %v, got %v", ques.path, ques.handlerId, handlerId)
		}
		// todo: check pathVars valid
		// puts(pathVars)
	}

}

func TestSimple(t *testing.T) {
	apiTest(simpleAPI, simpleApiQuestions, t)

}

func TestGithub(t *testing.T) {
	apiTest(githubAPI, []*question{}, t)
}
