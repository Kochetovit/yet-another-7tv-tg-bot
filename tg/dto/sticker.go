package dto

// https://core.telegram.org/bots/api#addstickertoset
type AddStickerToSetRequest struct {
	UserID  int64        `json:"user_id"`
	Name    string       `json:"name"`
	Sticker InputSticker `json:"sticker"`
}

// https://core.telegram.org/bots/api#createnewstickerset
type CreateNewStickerSet struct {
	UserID   int64          `json:"user_id"`
	Name     string         `json:"name"`
	Title    string         `json:"title"`
	Stickers []InputSticker `json:"stickers"`
}

// https://core.telegram.org/bots/api#inputsticker
type InputSticker struct {
	Sticker   string        `json:"sticker"`
	Format    StickerFormat `json:"format"`
	EmojiList []string      `json:"emoji_list"`
}

type StickerFormat string

const (
	StickerFormatStatic   StickerFormat = "static"
	StickerFormatAnimated StickerFormat = "animated"
	StickerFormatVideo    StickerFormat = "video"
)
