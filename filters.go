package main

type Filters struct {
	Users []string `json:"users"`
}

func (f *Filters) Filter(p *PullRequest) bool {
	if len(f.Users) == 0 {
		return false
	}
	for _, user := range f.Users {
		if user == p.Author || user == p.Assignee {
			return false
		}
	}
	return true
}
