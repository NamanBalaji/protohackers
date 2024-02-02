package chat

import (
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"regexp"
	"sort"
	"sync"
)

var (
	RgxUsername           = regexp.MustCompile(`^[a-zA-Z0-9]{1,}$`)
	ErrDuplicatedUsername = errors.New("duplicated user name")
	ErrInvalidUsername    = errors.New("invalid username")
)

type Chat struct {
	sync.Mutex
	users map[string]io.Writer
}

func New() *Chat {
	return &Chat{users: make(map[string]io.Writer)}
}

func (c *Chat) AddUser(username string, conn io.Writer) error {
	err := c.validateUserName(username)
	if err != nil {
		return err
	}

	c.Lock()
	c.users[username] = conn
	c.Unlock()

	c.Broadcast(username, fmt.Sprintf("* %s has entered the room\n", username))
	fmt.Fprintf(conn, "* the room contains: %s\n", c.ListUsers(username))

	return nil
}

func (c *Chat) RemoveUser(username string) {
	c.Lock()
	delete(c.users, username)
	c.Unlock()
	c.Broadcast(username, fmt.Sprintf("* %s has left the room\n", username))
}

func (c *Chat) Broadcast(sender, message string) {
	for user, conn := range c.users {
		if sender != user {
			if _, err := fmt.Fprint(conn, message); err != nil {
				log.Printf("Failed to broadcast message to other users: %s", err)
			}
		}
	}
}

func (c *Chat) ListUsers(except string) []string {
	users := maps.Clone(c.users)
	delete(users, except)

	var names []string
	for n, _ := range users {
		names = append(names, n)
	}
	sort.Strings(names)

	return names
}

func (c *Chat) SendMessage(username, message string) {
	c.Broadcast(username, fmt.Sprintf("[%s] %s\n", username, message))
}

func (c *Chat) validateUserName(username string) error {
	if _, ok := c.users[username]; ok {
		return ErrDuplicatedUsername
	}

	if !RgxUsername.MatchString(username) {
		return ErrInvalidUsername
	}

	return nil
}
