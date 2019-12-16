package auth

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"uhppoted/kvs"
)

type Permissions struct {
	Enabled bool
	users   *kvs.KeyValueStore
	groups  struct {
		groups   map[string][]permission
		filepath string
		guard    sync.Mutex
	}
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
	permissions := Permissions{
		Enabled: enabled,

		users: kvs.NewKeyValueStore(
			"permissions:users",
			func(value string) (interface{}, error) {
				return separator.Split(value, -1), nil
			}),

		groups: struct {
			groups   map[string][]permission
			filepath string
			guard    sync.Mutex
		}{
			groups:   map[string][]permission{},
			filepath: groups,
			guard:    sync.Mutex{},
		},
	}

	if enabled {
		err := permissions.users.LoadFromFile(users)
		if err != nil {
			return nil, err
		}

		g, err := getGroups(groups)
		if err != nil {
			return nil, err
		}

		permissions.groups.groups = g

		q := func() error {
			return permissions.reloadGroups(logger)
		}

		permissions.users.Watch(users, logger)
		watch(groups, q, logger)
	}

	return &permissions, nil
}

func (p *Permissions) Validate(clientID, resource, action string) error {
	groups, ok := p.users.Get(clientID)
	if !ok {
		return fmt.Errorf("%s: Not a member of any groups", clientID)
	}

	p.groups.guard.Lock()
	defer p.groups.guard.Unlock()

	for _, g := range groups.([]string) {
		if permissions, ok := p.groups.groups[g]; ok {
			for _, q := range permissions {
				if q.resource.MatchString(resource) && q.action.MatchString(action) {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("%s: Not authorised for %s:%s", clientID, resource, action)
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

func (p *Permissions) reloadGroups(log *log.Logger) error {
	groups, err := getGroups(p.groups.filepath)
	if err != nil {
		return err
	}

	p.groups.guard.Lock()
	defer p.groups.guard.Unlock()

	if !reflect.DeepEqual(groups, p.groups.groups) {
		for k, v := range groups {
			p.groups.groups[k] = v
		}

		for k, _ := range p.groups.groups {
			if _, ok := groups[k]; !ok {
				delete(p.groups.groups, k)
			}
		}

		log.Printf("WARN  Updated permissions:groups from '%s'", p.groups.filepath)
	}

	return nil
}
