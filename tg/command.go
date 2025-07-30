package tg

import "errors"

const helpMessage = "Commands:\n" +
	"/start - Create sticker pack\n" +
	"/add [png|gif] <url> [-S | -M | -L] [-e <emote>]"

type command interface {
	parseFlags(parts []string) error
	exec(t *Bot) (string, error)
}

type startCommand struct{}

func (c *startCommand) parseFlags(_ []string) error {
	return nil
}

func (c *startCommand) exec(t *Bot) (string, error) {
	if err := t.CreateStickerSet(t.cfg.StickerPackName + "_by_" + t.cfg.BotName); err != nil {
		return "", err
	}

	return "Pack was created", nil
}

// divide it into png and gif commands
type addCommand struct {
	url         string
	stickerType StickerType
	opts        []addStickerOptionsFn
}

func (c *addCommand) parseFlags(parts []string) error {
	if len(parts) < 3 {
		return errors.New("pass sticker type and url to 7tv emote")
	}

	switch parts[1] {
	case "png":
		c.stickerType = stickerTypePNG
	case "gif":
		c.stickerType = stickerTypeGIF
	default:
		return errors.New("invalid sticker type")
	}

	c.url = parts[2]

	for i := 3; i < len(parts); {
		switch parts[i] {
		case "-S":
			c.opts = append(c.opts, WithSize(SizeS))
		case "-M":
			c.opts = append(c.opts, WithSize(SizeM))
		case "-L":
			c.opts = append(c.opts, WithSize(SizeL))
		case "-e":
			i++

			if i >= len(parts) {
				return errors.New("pass emote")
			}

			c.opts = append(c.opts, WithEmote(parts[i]))
		}
	}

	return nil
}

func (c *addCommand) exec(t *Bot) (string, error) {
	if err := t.AddStickerToSet(c.url, c.stickerType, c.opts...); err != nil {
		return "", err
	}

	return "Sticker was added", nil
}

type helpCommand struct{}

func (c *helpCommand) parseFlags(_ []string) error {
	return nil
}

func (c *helpCommand) exec(t *Bot) (string, error) {
	return helpMessage, nil
}
