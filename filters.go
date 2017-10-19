package main

type Filters struct {
	users []string
}

func (f *Filters) Filter(p *PullRequest) bool {
	if len(f.users) == 0 {
		return false
	}
	for _, user := range f.users {
		if user == p.Author || user == p.Assignee {
			return false
		}
	}
	return true
}

func (f *Filters) SetUsers(values []string) {
	f.users = values
}
