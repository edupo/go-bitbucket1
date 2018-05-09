package bitbucket_v1

type Repository struct {
	Slug          string   `json:"slug"`
	Id            uint64   `json:"id"`
	Name          string   `json:"name"`
	ScmId         string   `json:"scmId"`
	State         string   `json:"state"`
	StatusMessage string   `json:"statusMessage"`
	Forkable      bool     `json:"forkable"`
	Project       *Project `json:"project"`
	Public        bool     `json:"public"`
	Links         Links    `json:"links"`
}

func (repository *Repository) GetPullRequests(opts *GetPullRequestsOpts) ([]*PullRequest, error) {
	opts.Project = repository.Project.Key
	opts.Repo = repository.Slug
	return repository.Project.Client.GetPullRequests(opts)
}
