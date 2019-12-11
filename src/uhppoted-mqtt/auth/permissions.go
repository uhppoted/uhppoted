package auth

import (
	"fmt"
	"regexp"
	"strings"
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

func NewPermissions(enabled bool, users, groups string) (*Permissions, error) {
	permissions := Permissions{
		Enabled: enabled,
		users:   map[string][]string{},
		groups:  map[string][]permission{},
	}

	if enabled {
		u, err := getUsers(users)
		if err != nil {
			return nil, err
		}

		g, err := getGroups(groups)
		if err != nil {
			return nil, err
		}

		permissions.users = u
		permissions.groups = g
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
					return nil
				}
			}
		}
	}

	return fmt.Errorf("%s: Not authorised for %s:%s", clientID, resource, action)
}

func getUsers(path string) (map[string][]string, error) {
	users := map[string][]string{}
	separator := regexp.MustCompile(`\s*,\s*`)
	err := load(path, func(key, value string) error {
		users[key] = separator.Split(value, -1)
		return nil
	})

	return users, err
}

func getGroups(path string) (map[string][]permission, error) {
	groups := map[string][]permission{}
	separator := regexp.MustCompile(`\s*,\s*`)
	re := regexp.MustCompile(`(.*?):(.*)`)
	err := load(path, func(key, value string) error {
		list := separator.Split(value, -1)
		for _, s := range list {
			match := re.FindStringSubmatch(s)
			if len(match) == 3 {
				resource, err := regexp.Compile("^" + strings.ReplaceAll(match[1], "*", ".*") + "$")
				if err != nil {
					return err
				}

				action, err := regexp.Compile("^" + strings.ReplaceAll(match[2], "*", ".*") + "$")
				if err != nil {
					return err
				}

				groups[key] = append(groups[key], permission{
					resource: resource,
					action:   action,
				})
			}
		}
		return nil
	})

	return groups, err
}
