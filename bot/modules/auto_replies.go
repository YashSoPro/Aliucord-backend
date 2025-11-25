package modules

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func init() {
	modules = append(modules, initAutoReplies)
}

func r(regex string) *regexp.Regexp {
	return regexp.MustCompile("(?i)" + regex)
}

var lastReplyCache = map[discord.UserID]string{}

func reply(msg *gateway.MessageCreateEvent, reply string) {
	if lastReplyCache[msg.Author.ID] == reply {
		return
	}
	_, err := s.SendTextReply(msg.ChannelID, reply, msg.ID)
	logger.LogIfErr(err)
	lastReplyCache[msg.Author.ID] = reply
}

const (
	JustAsk       = "https://dontasktoask.com/"
	CheckThePins  = "<a:checkpins:859804429536198676>"
	MentionHelp   = "Rule 9: Don't dm or mention for support"
	ElaborateHelp = "We can't help you if you don't tell us your issue. "

	// ✅ NEW GUIDES (YutaPlug)
	BeginnerGuide = "https://yutaplug.github.io/Aliucord/beginner"
	ThemerGuide   = "https://yutaplug.github.io/Aliucord/themer"
	ForksGuide    = "https://yutaplug.github.io/Aliucord/forks"
	SoundsGuide   = "https://yutaplug.github.io/Aliucord/sounds"
	UserPFPBG     = "https://yutaplug.github.io/Aliucord/userpfpbg"
	FindPlugins   = "https://yutaplug.github.io/Aliucord/findplugins"
	NewUIGuide    = "https://yutaplug.github.io/Aliucord/newui"
	Backports     = "https://yutaplug.github.io/Aliucord/backports"
	Changelog     = "https://yutaplug.github.io/Aliucord/changelog"
	TokenGuide    = "https://yutaplug.github.io/Aliucord/token"
	OldUIGuide    = "https://yutaplug.github.io/Aliucord/oldui"

	FullTransparency = "1. Are you using a theme that requires full transparency? If the answer is no, then that's the problem. Normally in the description says what transparency you need to use. 2. Are you using a custom ROM? If the answer is yes, then we can't do nothing about it."
	AliuCrash        = "Send crashlogs (check Crashes in Settings, and copy the most recent), if there aren't any crashlogs, then we can't do nothing about it."
	FreeNitro        = "Not possible. Nitrospoof exists for \"free\" emotes, UserBG exists for using a custom banner, UserPFP exists for using a custom profile picture."
	Usage            = "Go to the plugin's repository and read the readme. Chances are the dev added a description."
	BetterInternet   = "This happens when you have an old/misbehaving router. Use mobile data or maybe a VPN."
	PluginDownloader = "PluginDownloader is now a part of Aliucord. If download is missing, update Aliucord."
)

func initAutoReplies() {
	cfg := config.AutoReplyConfig
	if !cfg.Enabled {
		return
	}

	PRD := fmt.Sprintf("%s 👉 <#%s>", CheckThePins, cfg.PRD)
	FindPluginGuide := fmt.Sprintf("Search in <#%s> and <#%s>. If it doesn't exist, check %s",
		cfg.PluginsList, cfg.NewPlugins, FindPlugins)

	autoRepliesString := map[string]string{
		"a plugin to":           PRD,
		"can you make":          PRD,
		"how do i use":          Usage,
		"free nitro":            FreeNitro,
		"handshake exception":   BetterInternet,
		"connection terminated": BetterInternet,
	}

	autoRepliesRegex := map[*regexp.Regexp]string{
		r("^(?:i need )?help(?: me)?$"):                            ElaborateHelp,
		r("<@!?\\d{2,19}> help"):                           MentionHelp,
		r("help <@!?\\d{2,19}>"):                           MentionHelp,
		r("animated (profile|avatar|pfp)"):                 FreeNitro,
		r("is there a plugin.+"):                           FindPluginGuide,
		r("^where(?: i)?'?s(?: the )?.+ plugin$"):          FindPluginGuide,
		r("^can (?:someone|anybody|anyone|you) help\??$"): JustAsk,

		// ✅ UPDATED GUIDE RESPONSES
		r("how (?:to|do i|do you) install aliucord"): BeginnerGuide,
		r("how (?:to|do i|do you) use themer"):         ThemerGuide,
		r("how (?:to|do i|do you) add sounds"):         SoundsGuide,
		r("how (?:to|do i|do you) change pfp|bg"):      UserPFPBG,
		r("new ui"):                                     NewUIGuide,
		r("old ui"):                                     OldUIGuide,
		r("forks"):                                      ForksGuide,
		r("backport"):                                   Backports,
		r("changelog"):                                  Changelog,
		r("token"):                                      TokenGuide,

		r("aliucord (crashed|keeps crashing|crash|crashes)"): AliuCrash,
		r("full transparency not work"):                   FullTransparency,
	}

	s.AddHandler(func(msg *gateway.MessageCreateEvent) {
		if msg.Member == nil || len(msg.Attachments) > 0 || (msg.ReferencedMessage != nil && msg.ReferencedMessage.Author.ID == msg.Author.ID) || msg.Author.Bot || strings.HasPrefix(msg.Content, "Quick Aliucord ") {
			return
		}

		c, err := s.Channel(msg.ChannelID)
		if err == nil {
			if c.ID != cfg.PRD && c.ParentID != cfg.SupportCategory {
				return
			}
		}

		for _, role := range msg.Member.RoleIDs {
			if slices.Contains(cfg.IgnoredRoles, role) {
				return
			}
		}

		for regex, value := range autoRepliesRegex {
			if regex.MatchString(msg.Content) {
				reply(msg, value)
				return
			}
		}

		content := strings.ToLower(msg.Content)
		for trigger, value := range autoRepliesString {
			if strings.Contains(content, trigger) {
				reply(msg, value)
				return
			}
		}
	})
}
