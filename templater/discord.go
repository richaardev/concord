package templater

import (
	"github.com/disgoorg/disgo/discord"
)

func DiscordUserPlaceholders(user discord.User) map[string]interface{} {
	return map[string]interface{}{
		"id":                 user.ID,
		"mention":            user.Mention(),
		"username":           user.Username,
		"global_name":        user.GlobalName,
		"avatar_url":         user.EffectiveAvatarURL(),
		"default_avatar_url": user.DefaultAvatarURL(),
		"discriminator":      user.Discriminator,
		"created_at":         user.CreatedAt().Format("02 Jan 06 15:04"),
		"tag":                user.Tag(),
	}
}

func DiscordMemberPlaceholders(member discord.Member) map[string]interface{} {
	return map[string]interface{}{
		"id":                 member.User.ID,
		"username":           member.User.Username,
		"mention":            member.User.Mention(),
		"global_name":        member.User.GlobalName,
		"nickname":           member.Nick,
		"avatar_url":         member.User.EffectiveAvatarURL(),
		"default_avatar_url": member.User.DefaultAvatarURL(),
		"discriminator":      member.User.Discriminator,
		"tag":                member.User.Tag(),
		"joined_at":          member.JoinedAt.Format("02 Jan 06 15:04"),
		"created_at":         member.User.CreatedAt().Format("02 Jan 06 15:04"),
	}
}

func DiscordRolePlaceholders(role discord.Role) map[string]interface{} {
	return map[string]interface{}{
		"id":          role.ID,
		"name":        role.Name,
		"description": role.Description,
		"mention":     role.Mention(),
		"icon_url":    role.IconURL(),
	}
}

func DiscordGuildPlaceholders(guild discord.Guild) map[string]interface{} {
	return map[string]interface{}{
		"id":                            guild.ID,
		"name":                          guild.Name,
		"icon_url":                      guild.IconURL(),
		"splash_url":                    guild.SplashURL(),
		"discovery_splash_url":          guild.DiscoverySplashURL(),
		"owner_id":                      guild.OwnerID,
		"afk_channel_id":                guild.AfkChannelID,
		"afk_timeout":                   guild.AfkTimeout,
		"widget_enabled":                guild.WidgetEnabled,
		"widget_channel_id":             guild.WidgetChannelID,
		"verification_level":            guild.VerificationLevel,
		"default_message_notifications": guild.DefaultMessageNotifications,
		"explicit_content_filter":       guild.ExplicitContentFilter,
		"features":                      guild.Features,
		"mfa_level":                     guild.MFALevel,
		"application_id":                guild.ApplicationID,
		"system_channel_id":             guild.SystemChannelID,
		"system_channel_flags":          guild.SystemChannelFlags,
		"rules_channel_id":              guild.RulesChannelID,
		"member_count":                  guild.MemberCount,
		"max_presences":                 guild.MaxPresences,
		"max_members":                   guild.MaxMembers,
		"vanity_url_code":               guild.VanityURLCode,
		"description":                   guild.Description,
		"banner_url":                    guild.BannerURL(),
		"premium_tier":                  guild.PremiumTier,
		"premium_subscription_count":    guild.PremiumSubscriptionCount,
		"preferred_locale":              guild.PreferredLocale,
		"public_updates_channel_id":     guild.PublicUpdatesChannelID,
		"max_video_channel_users":       guild.MaxVideoChannelUsers,
		"max_stage_video_channel_users": guild.MaxStageVideoChannelUsers,
		"welcome_screen":                guild.WelcomeScreen,
		"nsfw_level":                    guild.NSFWLevel,
		"premium_progress_bar_enabled":  guild.PremiumProgressBarEnabled,
		"joined_at":                     guild.JoinedAt.Format("02 Jan 06 15:04"),
		"safety_alerts_channel_id":      guild.SafetyAlertsChannelID,
		"incidents_data":                guild.IncidentsData,
		"approximate_member_count":      guild.ApproximateMemberCount,
		"approximate_presence_count":    guild.ApproximatePresenceCount,
	}
}
