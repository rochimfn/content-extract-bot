// THIS FILE IS AUTOGENERATED. DO NOT EDIT.
// Regen by running 'go generate' in the repo root.

package gotgbot

// The consts listed below represent all the update types that can be requested from telegram.
const (
	UpdateTypeMessage                 = "message"
	UpdateTypeEditedMessage           = "edited_message"
	UpdateTypeChannelPost             = "channel_post"
	UpdateTypeEditedChannelPost       = "edited_channel_post"
	UpdateTypeBusinessConnection      = "business_connection"
	UpdateTypeBusinessMessage         = "business_message"
	UpdateTypeEditedBusinessMessage   = "edited_business_message"
	UpdateTypeDeletedBusinessMessages = "deleted_business_messages"
	UpdateTypeMessageReaction         = "message_reaction"
	UpdateTypeMessageReactionCount    = "message_reaction_count"
	UpdateTypeInlineQuery             = "inline_query"
	UpdateTypeChosenInlineResult      = "chosen_inline_result"
	UpdateTypeCallbackQuery           = "callback_query"
	UpdateTypeShippingQuery           = "shipping_query"
	UpdateTypePreCheckoutQuery        = "pre_checkout_query"
	UpdateTypePoll                    = "poll"
	UpdateTypePollAnswer              = "poll_answer"
	UpdateTypeMyChatMember            = "my_chat_member"
	UpdateTypeChatMember              = "chat_member"
	UpdateTypeChatJoinRequest         = "chat_join_request"
	UpdateTypeChatBoost               = "chat_boost"
	UpdateTypeRemovedChatBoost        = "removed_chat_boost"
)

// GetType is a helper method to easily identify the type of update that is being received.
func (u Update) GetType() string {
	switch {
	case u.Message != nil:
		return UpdateTypeMessage

	case u.EditedMessage != nil:
		return UpdateTypeEditedMessage

	case u.ChannelPost != nil:
		return UpdateTypeChannelPost

	case u.EditedChannelPost != nil:
		return UpdateTypeEditedChannelPost

	case u.BusinessConnection != nil:
		return UpdateTypeBusinessConnection

	case u.BusinessMessage != nil:
		return UpdateTypeBusinessMessage

	case u.EditedBusinessMessage != nil:
		return UpdateTypeEditedBusinessMessage

	case u.DeletedBusinessMessages != nil:
		return UpdateTypeDeletedBusinessMessages

	case u.MessageReaction != nil:
		return UpdateTypeMessageReaction

	case u.MessageReactionCount != nil:
		return UpdateTypeMessageReactionCount

	case u.InlineQuery != nil:
		return UpdateTypeInlineQuery

	case u.ChosenInlineResult != nil:
		return UpdateTypeChosenInlineResult

	case u.CallbackQuery != nil:
		return UpdateTypeCallbackQuery

	case u.ShippingQuery != nil:
		return UpdateTypeShippingQuery

	case u.PreCheckoutQuery != nil:
		return UpdateTypePreCheckoutQuery

	case u.Poll != nil:
		return UpdateTypePoll

	case u.PollAnswer != nil:
		return UpdateTypePollAnswer

	case u.MyChatMember != nil:
		return UpdateTypeMyChatMember

	case u.ChatMember != nil:
		return UpdateTypeChatMember

	case u.ChatJoinRequest != nil:
		return UpdateTypeChatJoinRequest

	case u.ChatBoost != nil:
		return UpdateTypeChatBoost

	case u.RemovedChatBoost != nil:
		return UpdateTypeRemovedChatBoost

	default:
		return "unknown"
	}
}

// The consts listed below represent all the parse_mode options that can be sent to telegram.
const (
	ParseModeHTML       = "HTML"
	ParseModeMarkdownV2 = "MarkdownV2"
	ParseModeMarkdown   = "Markdown"
	ParseModeNone       = ""
)

// The consts listed below represent all the chat action options that can be sent to telegram.
const (
	ChatActionTyping          = "typing"
	ChatActionUploadPhoto     = "upload_photo"
	ChatActionRecordVideo     = "record_video"
	ChatActionUploadVideo     = "upload_video"
	ChatActionRecordVoice     = "record_voice"
	ChatActionUploadVoice     = "upload_voice"
	ChatActionUploadDocument  = "upload_document"
	ChatActionChooseSticker   = "choose_sticker"
	ChatActionFindLocation    = "find_location"
	ChatActionRecordVideoNote = "record_video_note"
	ChatActionUploadVideoNote = "upload_video_note"
)

// The consts listed below represent all the sticker types that can be obtained from telegram.
const (
	StickerTypeRegular     = "regular"
	StickerTypeMask        = "mask"
	StickerTypeCustomEmoji = "custom_emoji"
)

// The consts listed below represent all the chat types that can be obtained from telegram.
const (
	ChatTypePrivate    = "private"
	ChatTypeGroup      = "group"
	ChatTypeSupergroup = "supergroup"
	ChatTypeChannel    = "channel"
)
