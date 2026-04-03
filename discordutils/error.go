package discordutils

import (
	"errors"

	"github.com/disgoorg/disgo/rest"
)

func ParseDiscordError(err error) {
	var restError rest.Error
	if errors.As(err, &restError) {
		println(string(restError.RsBody))
	}
}
