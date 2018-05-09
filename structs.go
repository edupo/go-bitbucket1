package bitbucket_v1

type PullRequest struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	FromRef     Ref    `json:"fromRef"`
	ToRef       Ref    `json:"toRef"`
	Repository  *Repository
}

type Ref struct {
	Id           string `json:"id"`
	DisplayId    string `json:"displayId"`
	LatestCommit string `json:"latestCommit"`
}

type User struct {
	Id           uint64 `json:"id"`
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
}

type Links struct {
	Self  []NamedLink `json:"self"`
	Clone []NamedLink `json:"clone"`
}

type NamedLink struct {
	Href string `json:"href"`
	Name string `json:"name"`
}
