package auth

import (
	"fmt"
	"regexp"
)

type Permissions struct {
	Enabled bool
	users   map[string][]string
	groups  map[string][]permission
}

type permission struct {
	resource *regexp.Regexp
	action   *regexp.Regexp
}

func NewPermissions(enabled bool) (*Permissions, error) {
	permissions := Permissions{
		Enabled: enabled,
		users:   map[string][]string{},
		groups:  map[string][]permission{},
	}

	return &permissions, nil
}

func (p *Permissions) Validate(clientID, resource, action string) error {
	groups, ok := p.users[clientID]
	if !ok {
		return fmt.Errorf("%s: Not a member of any groups", clientID)
	}

	for _, g := range groups {
		if permissions, ok := p.groups[g]; ok {
			for _, q := range permissions {
				if q.resource.MatchString(resource) && q.action.MatchString(action) {
					fmt.Printf(">>> GOTCHA: %v %s %s\n", q, resource, action)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("%s: Not authorised for %s:%s", clientID, resource, action)
}
