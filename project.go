package bitbucket_v1

type Project struct {
	Key         string `json:"key"`
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
	Links       Links  `json:"links"`
	Client      *Client
}

func (project *Project) Repository(slug string) (*Repository, error) {
	repo, err := project.Client.GetRepository(project.Key, slug)
	if err != nil {
		return nil, err
	}
	repo.Project = project
	return repo, nil
}
