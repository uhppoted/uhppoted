package auth

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"uhppoted/kvs"
)

type Permissions struct {
	Enabled bool
	users   *kvs.KeyValueStore
	groups  *kvs.KeyValueStore
}

type permission struct {
	resource *regexp.Regexp
	action   *regexp.Regexp
}

func (p permission) String() string {
	return fmt.Sprintf("resource:`%s` action:`%s`", p.resource, p.action)
}

func NewPermissions(enabled bool, users, groups string, logger *log.Logger) (*Permissions, error) {
	separator := regexp.MustCompile(`\s*,\s*`)

	u := func(value string) (interface{}, error) {
		return separator.Split(value, -1), nil
	}

	g := func(value string) (interface{}, error) {
		permissions := []permission{}
		re := regexp.MustCompile(`(.*?):(.*)`)
		tokens := separator.Split(value, -1)
		for _, s := range tokens {
			if match := re.FindStringSubmatch(s); len(match) == 3 {
				resource, err := regexp.Compile("^" + strings.ReplaceAll(match[1], "*", ".*") + "$")
				if err != nil {
					return permissions, err
				}

				action, err := regexp.Compile("^" + strings.ReplaceAll(match[2], "*", ".*") + "$")
				if err != nil {
					return permissions, err
				}

				permissions = append(permissions, permission{
					resource: resource,
					action:   action,
				})
			}
		}

		return permissions, nil
	}

	permissions := Permissions{
		Enabled: enabled,
		users:   kvs.NewKeyValueStore("permissions:users", u, logger),
		groups:  kvs.NewKeyValueStore("permissions:groups", g, logger),
	}

	if enabled {
		err := permissions.users.LoadFromFile(users)
		if err != nil {
			return nil, err
		}

		err = permissions.groups.LoadFromFile(groups)
		if err != nil {
			return nil, err
		}

		permissions.users.Watch(users, logger)
		permissions.groups.Watch(groups, logger)
	}

	return &permissions, nil
}

func (p *Permissions) Validate(clientID, resource, action string) error {
	groups, ok := p.users.Get(clientID)
	if !ok {
		return fmt.Errorf("%s: Not a member of any groups", clientID)
	}

	for _, g := range groups.([]string) {
		if permissions, ok := p.groups.Get(g); ok {
			for _, q := range permissions.([]permission) {
				if q.resource.MatchString(resource) && q.action.MatchString(action) {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("%s: Not authorised for %s:%s", clientID, resource, action)
}
