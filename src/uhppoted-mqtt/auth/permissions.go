package auth

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

type Permissions struct {
	Enabled bool
	users   struct {
		users    map[string][]string
		filepath string
		guard    sync.Mutex
	}
	groups struct {
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
	permissions := Permissions{
		Enabled: enabled,
		users: struct {
			users    map[string][]string
			filepath string
			guard    sync.Mutex
		}{
			users:    map[string][]string{},
			filepath: users,
			guard:    sync.Mutex{},
		},
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
		u, err := getUsers(users)
		if err != nil {
			return nil, err
		}

		g, err := getGroups(groups)
		if err != nil {
			return nil, err
		}

		permissions.users.users = u
		permissions.groups.groups = g

		p := func() error {
			return permissions.reloadUsers(logger)
		}

		q := func() error {
			return permissions.reloadGroups(logger)
		}

		watch(users, p, logger)
		watch(groups, q, logger)
	}

	return &permissions, nil
}

func (p *Permissions) Validate(clientID, resource, action string) error {
	p.users.guard.Lock()
	defer p.users.guard.Unlock()

	groups, ok := p.users.users[clientID]
	if !ok {
		return fmt.Errorf("%s: Not a member of any groups", clientID)
	}

	p.groups.guard.Lock()
	defer p.groups.guard.Unlock()

	for _, g := range groups {
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

func (p *Permissions) reloadUsers(log *log.Logger) error {
	users, err := getUsers(p.users.filepath)
	if err != nil {
		return err
	}

	p.users.guard.Lock()
	defer p.users.guard.Unlock()

	if !reflect.DeepEqual(users, p.users.users) {
		for k, v := range users {
			p.users.users[k] = v
		}

		for k, _ := range p.users.users {
			if _, ok := users[k]; !ok {
				delete(p.users.users, k)
			}
		}

		log.Printf("WARN  Updated permissions:users from '%s'", p.users.filepath)
	}

	return nil
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
