package bot

import (
	telemux "github.com/and3rson/telemux/v2"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Filter      telemux.FilterFunc
	Action      Action
}

type Commands []Command

var (
	lastDialogPart string
)

func SetLastDialogPart(command string) {
	lastDialogPart = command
}

func GetLastDialogPart() string {
	return lastDialogPart
}

func FilterDefault(u *telemux.Update, name string) bool {
	if u.Message != nil {
		if strings.HasPrefix(u.Message.Text, "/"+name) {
			return true
		}
	}
	return false
}
