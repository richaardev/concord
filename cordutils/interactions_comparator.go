package cordutils

import (
	"log/slog"
	"reflect"
	"slices"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/omit"
)

type commandIdentity struct {
	name    string
	cmdType discord.ApplicationCommandType
}

func CheckIfCommandsIsEqual(oldCommands []discord.ApplicationCommand, desiredCommands []discord.ApplicationCommandCreate) bool {
	if len(oldCommands) != len(desiredCommands) {
		logComparison("command_count_mismatch", map[string]any{
			"existing": len(oldCommands),
			"desired":  len(desiredCommands),
		})
		return false
	}

	existing := make(map[commandIdentity]discord.ApplicationCommand, len(oldCommands))
	for _, cmd := range oldCommands {
		key := commandIdentity{name: cmd.Name(), cmdType: cmd.Type()}
		existing[key] = cmd
	}

	for _, desired := range desiredCommands {
		key := commandIdentity{name: desired.CommandName(), cmdType: desired.Type()}
		existingCmd, ok := existing[key]
		if !ok {
			logComparison("command_not_found", map[string]any{
				"name": desired.CommandName(),
				"type": desired.Type(),
			})
			return false
		}

		if !CompareApplicationCommand(existingCmd, desired) {
			logComparison("command_payload_mismatch", map[string]any{
				"name": desired.CommandName(),
				"type": desired.Type(),
			})
			return false
		}

		delete(existing, key)
	}

	if len(existing) > 0 {
		for _, cmd := range existing {
			logComparison("unexpected_command_found", map[string]any{
				"name": cmd.Name(),
				"type": cmd.Type(),
			})
			break
		}
		return false
	}

	return true
}

func CompareApplicationCommand(existing discord.ApplicationCommand, desired discord.ApplicationCommandCreate) bool {
	switch oldCmd := existing.(type) {
	case discord.SlashCommand:
		desiredCmd, ok := desired.(discord.SlashCommandCreate)
		if !ok {
			logComparison("type_cast_failure", map[string]any{"expected": "SlashCommandCreate", "got": desired})
			return false
		}
		return CompareSlashCommand(oldCmd, desiredCmd)
	case discord.UserCommand:
		desiredCmd, ok := desired.(discord.UserCommandCreate)
		if !ok {
			logComparison("type_cast_failure", map[string]any{"expected": "UserCommandCreate", "got": desired})
			return false
		}
		return CompareUserCommand(oldCmd, desiredCmd)
	case discord.MessageCommand:
		desiredCmd, ok := desired.(discord.MessageCommandCreate)
		if !ok {
			logComparison("type_cast_failure", map[string]any{"expected": "MessageCommandCreate", "got": desired})
			return false
		}
		return CompareMessageCommand(oldCmd, desiredCmd)
	case discord.EntryPointCommand:
		desiredCmd, ok := desired.(discord.EntryPointCommandCreate)
		if !ok {
			logComparison("type_cast_failure", map[string]any{"expected": "EntryPointCommandCreate", "got": desired})
			return false
		}
		return CompareEntryPointCommand(oldCmd, desiredCmd)
	default:
		logComparison("unsupported_command_type", map[string]any{"type": existing.Type()})
		return false
	}
}

func CompareSlashCommand(existing discord.SlashCommand, desired discord.SlashCommandCreate) bool {
	if existing.Name() != desired.Name || existing.Description != desired.Description {
		logComparison("slash_metadata_mismatch", map[string]any{
			"existing_name":        existing.Name(),
			"desired_name":         desired.Name,
			"existing_description": existing.Description,
			"desired_description":  desired.Description,
		})
		return false
	}

	if !compareLocaleMap(existing.NameLocalizations(), desired.NameLocalizations) {
		logComparison("slash_name_locale_mismatch", map[string]any{"existing": existing.NameLocalizations(), "desired": desired.NameLocalizations})
		return false
	}

	if !compareLocaleMap(existing.DescriptionLocalizations, desired.DescriptionLocalizations) {
		logComparison("slash_description_locale_mismatch", map[string]any{"existing": existing.DescriptionLocalizations, "desired": desired.DescriptionLocalizations})
		return false
	}

	if !compareNullablePermissions(existing.DefaultMemberPermissions(), desired.DefaultMemberPermissions) {
		logComparison("slash_permissions_mismatch", map[string]any{"existing": existing.DefaultMemberPermissions(), "desired": desired.DefaultMemberPermissions})
		return false
	}

	// if !compareOptionalBool(existing.DMPermission(), desired.DMPermission, true) {
	// 	logComparison("slash_dm_permission_mismatch", map[string]any{"existing": existing.DMPermission(), "desired": desired.DMPermission})
	// 	return false
	// }

	if !compareOptionalBool(existing.NSFW(), desired.NSFW, false) {
		logComparison("slash_nsfw_mismatch", map[string]any{"existing": existing.NSFW(), "desired": desired.NSFW})
		return false
	}

	if !sameElements(existing.IntegrationTypes(), desired.IntegrationTypes) && len(desired.IntegrationTypes) > 0 {
		logComparison("slash_integration_types_mismatch", map[string]any{"existing": existing.IntegrationTypes(), "desired": desired.IntegrationTypes})
		return false
	}

	if !sameElements(existing.Contexts(), desired.Contexts) {
		logComparison("slash_contexts_mismatch", map[string]any{"existing": existing.Contexts(), "desired": desired.Contexts})
		return false
	}

	return compareOptions(existing.Options, desired.Options)
}

func CompareUserCommand(existing discord.UserCommand, desired discord.UserCommandCreate) bool {
	return compareSharedContextCommandFields(
		existing.Name(),
		desired.Name,
		existing.NameLocalizations(),
		desired.NameLocalizations,
		existing.DefaultMemberPermissions(),
		desired.DefaultMemberPermissions,
		existing.NSFW(),
		desired.NSFW,
		existing.IntegrationTypes(),
		desired.IntegrationTypes,
		existing.Contexts(),
		desired.Contexts,
	)
}

func CompareMessageCommand(existing discord.MessageCommand, desired discord.MessageCommandCreate) bool {
	return compareSharedContextCommandFields(
		existing.Name(),
		desired.Name,
		existing.NameLocalizations(),
		desired.NameLocalizations,
		existing.DefaultMemberPermissions(),
		desired.DefaultMemberPermissions,
		existing.NSFW(),
		desired.NSFW,
		existing.IntegrationTypes(),
		desired.IntegrationTypes,
		existing.Contexts(),
		desired.Contexts,
	)
}

func CompareEntryPointCommand(existing discord.EntryPointCommand, desired discord.EntryPointCommandCreate) bool {
	if !compareSharedContextCommandFields(
		existing.Name(),
		desired.Name,
		existing.NameLocalizations(),
		desired.NameLocalizations,
		existing.DefaultMemberPermissions(),
		desired.DefaultMemberPermissions,
		existing.NSFW(),
		desired.NSFW,
		existing.IntegrationTypes(),
		desired.IntegrationTypes,
		existing.Contexts(),
		desired.Contexts,
	) {
		return false
	}

	if existing.Handler != desired.Handler {
		logComparison("entrypoint_handler_mismatch", map[string]any{"existing": existing.Handler, "desired": desired.Handler})
		return false
	}

	return true
}

func compareSharedContextCommandFields(
	existingName, desiredName string,
	existingLocales, desiredLocales map[discord.Locale]string,
	existingPermissions discord.Permissions,
	desiredPermissions omit.Omit[*discord.Permissions],
	existingNSFW bool,
	desiredNSFW *bool,
	existingIntegrations []discord.ApplicationIntegrationType,
	desiredIntegrations []discord.ApplicationIntegrationType,
	existingContexts []discord.InteractionContextType,
	desiredContexts []discord.InteractionContextType,
) bool {
	if existingName != desiredName {
		logComparison("context_name_mismatch", map[string]any{"existing": existingName, "desired": desiredName})
		return false
	}

	if !compareLocaleMap(existingLocales, desiredLocales) {
		logComparison("context_locale_mismatch", map[string]any{"existing": existingLocales, "desired": desiredLocales})
		return false
	}

	if !compareNullablePermissions(existingPermissions, desiredPermissions) {
		logComparison("context_permissions_mismatch", map[string]any{"existing": existingPermissions, "desired": desiredPermissions})
		return false
	}
	if !compareOptionalBool(existingNSFW, desiredNSFW, false) {
		logComparison("context_nsfw_mismatch", map[string]any{"existing": existingNSFW, "desired": desiredNSFW})
		return false
	}

	if !sameElements(existingIntegrations, desiredIntegrations) {
		logComparison("context_integration_mismatch", map[string]any{"existing": existingIntegrations, "desired": desiredIntegrations})
		return false
	}

	if !sameElements(existingContexts, desiredContexts) {
		logComparison("context_contexts_mismatch", map[string]any{"existing": existingContexts, "desired": desiredContexts})
		return false
	}

	return true
}

func compareOptions(existing []discord.ApplicationCommandOption, desired []discord.ApplicationCommandOption) bool {
	if len(existing) != len(desired) {
		logComparison("option_length_mismatch", map[string]any{"existing": len(existing), "desired": len(desired)})
		return false
	}

	for i := range existing {
		if !compareOption(existing[i], desired[i]) {
			logComparison("option_index_mismatch", map[string]any{"index": i})
			return false
		}
	}

	return true
}

func compareOption(existing discord.ApplicationCommandOption, desired discord.ApplicationCommandOption) bool {
	if existing.Type() != desired.Type() || existing.OptionName() != desired.OptionName() || existing.OptionDescription() != desired.OptionDescription() {
		logComparison("option_metadata_mismatch", map[string]any{
			"existing_type": existing.Type(),
			"desired_type":  desired.Type(),
			"existing_name": existing.OptionName(),
			"desired_name":  desired.OptionName(),
			"existing_desc": existing.OptionDescription(),
			"desired_desc":  desired.OptionDescription(),
		})
		return false
	}

	switch oldOpt := existing.(type) {
	case discord.ApplicationCommandOptionSubCommand:
		newOpt, ok := desired.(discord.ApplicationCommandOptionSubCommand)
		if !ok {
			return false
		}
		if !compareLocaleMap(oldOpt.NameLocalizations, newOpt.NameLocalizations) || !compareLocaleMap(oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations) {
			logComparison("subcommand_locale_mismatch", map[string]any{"existing": oldOpt.NameLocalizations, "desired": newOpt.NameLocalizations})
			return false
		}
		return compareOptions(oldOpt.Options, newOpt.Options)
	case discord.ApplicationCommandOptionSubCommandGroup:
		newOpt, ok := desired.(discord.ApplicationCommandOptionSubCommandGroup)
		if !ok {
			return false
		}
		if !compareLocaleMap(oldOpt.NameLocalizations, newOpt.NameLocalizations) || !compareLocaleMap(oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations) {
			logComparison("subcommand_group_locale_mismatch", map[string]any{"existing": oldOpt.NameLocalizations, "desired": newOpt.NameLocalizations})
			return false
		}
		if len(oldOpt.Options) != len(newOpt.Options) {
			logComparison("subcommand_group_length_mismatch", map[string]any{"existing": len(oldOpt.Options), "desired": len(newOpt.Options)})
			return false
		}
		for i := range oldOpt.Options {
			if !compareSubCommandValue(oldOpt.Options[i], newOpt.Options[i]) {
				logComparison("subcommand_group_option_mismatch", map[string]any{"index": i})
				return false
			}
		}
		return true
	case discord.ApplicationCommandOptionString:
		newOpt, ok := desired.(discord.ApplicationCommandOptionString)
		if !ok {
			return false
		}
		return compareStandardOptionFields(oldOpt.Required, newOpt.Required, oldOpt.NameLocalizations, newOpt.NameLocalizations, oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations) &&
			oldOpt.Autocomplete == newOpt.Autocomplete &&
			compareIntPointer(oldOpt.MinLength, newOpt.MinLength) &&
			compareIntPointer(oldOpt.MaxLength, newOpt.MaxLength) &&
			reflect.DeepEqual(oldOpt.Choices, newOpt.Choices)
	case discord.ApplicationCommandOptionInt:
		newOpt, ok := desired.(discord.ApplicationCommandOptionInt)
		if !ok {
			return false
		}
		return compareStandardOptionFields(oldOpt.Required, newOpt.Required, oldOpt.NameLocalizations, newOpt.NameLocalizations, oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations) &&
			oldOpt.Autocomplete == newOpt.Autocomplete &&
			compareIntPointer(oldOpt.MinValue, newOpt.MinValue) &&
			compareIntPointer(oldOpt.MaxValue, newOpt.MaxValue) &&
			reflect.DeepEqual(oldOpt.Choices, newOpt.Choices)
	case discord.ApplicationCommandOptionBool:
		newOpt, ok := desired.(discord.ApplicationCommandOptionBool)
		if !ok {
			return false
		}
		return compareStandardOptionFields(oldOpt.Required, newOpt.Required, oldOpt.NameLocalizations, newOpt.NameLocalizations, oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations)
	case discord.ApplicationCommandOptionUser:
		newOpt, ok := desired.(discord.ApplicationCommandOptionUser)
		if !ok {
			return false
		}
		return compareStandardOptionFields(oldOpt.Required, newOpt.Required, oldOpt.NameLocalizations, newOpt.NameLocalizations, oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations)
	case discord.ApplicationCommandOptionChannel:
		newOpt, ok := desired.(discord.ApplicationCommandOptionChannel)
		if !ok {
			return false
		}
		return compareStandardOptionFields(oldOpt.Required, newOpt.Required, oldOpt.NameLocalizations, newOpt.NameLocalizations, oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations) &&
			slices.Equal(oldOpt.ChannelTypes, newOpt.ChannelTypes)
	case discord.ApplicationCommandOptionRole:
		newOpt, ok := desired.(discord.ApplicationCommandOptionRole)
		if !ok {
			return false
		}
		return compareStandardOptionFields(oldOpt.Required, newOpt.Required, oldOpt.NameLocalizations, newOpt.NameLocalizations, oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations)
	case discord.ApplicationCommandOptionMentionable:
		newOpt, ok := desired.(discord.ApplicationCommandOptionMentionable)
		if !ok {
			return false
		}
		return compareStandardOptionFields(oldOpt.Required, newOpt.Required, oldOpt.NameLocalizations, newOpt.NameLocalizations, oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations)
	case discord.ApplicationCommandOptionFloat:
		newOpt, ok := desired.(discord.ApplicationCommandOptionFloat)
		if !ok {
			return false
		}
		return compareStandardOptionFields(oldOpt.Required, newOpt.Required, oldOpt.NameLocalizations, newOpt.NameLocalizations, oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations) &&
			oldOpt.Autocomplete == newOpt.Autocomplete &&
			compareFloatPointer(oldOpt.MinValue, newOpt.MinValue) &&
			compareFloatPointer(oldOpt.MaxValue, newOpt.MaxValue) &&
			reflect.DeepEqual(oldOpt.Choices, newOpt.Choices)
	case discord.ApplicationCommandOptionAttachment:
		newOpt, ok := desired.(discord.ApplicationCommandOptionAttachment)
		if !ok {
			return false
		}
		return compareStandardOptionFields(oldOpt.Required, newOpt.Required, oldOpt.NameLocalizations, newOpt.NameLocalizations, oldOpt.DescriptionLocalizations, newOpt.DescriptionLocalizations)
	default:
		logComparison("unsupported_option_type", map[string]any{"type": existing.Type()})
		return false
	}
}

func compareStandardOptionFields(existingRequired bool, desiredRequired bool, existingNameLocales, desiredNameLocales map[discord.Locale]string, existingDescLocales, desiredDescLocales map[discord.Locale]string) bool {
	if existingRequired != desiredRequired {
		logComparison("option_required_mismatch", map[string]any{"existing": existingRequired, "desired": desiredRequired})
		return false
	}

	if !compareLocaleMap(existingNameLocales, desiredNameLocales) {
		logComparison("option_name_locale_mismatch", map[string]any{"existing": existingNameLocales, "desired": desiredNameLocales})
		return false
	}

	if !compareLocaleMap(existingDescLocales, desiredDescLocales) {
		logComparison("option_description_locale_mismatch", map[string]any{"existing": existingDescLocales, "desired": desiredDescLocales})
		return false
	}

	return true
}

func compareSubCommandValue(existing discord.ApplicationCommandOptionSubCommand, desired discord.ApplicationCommandOptionSubCommand) bool {
	if existing.Name != desired.Name || existing.Description != desired.Description {
		logComparison("subcommand_metadata_mismatch", map[string]any{"existing": existing.Name, "desired": desired.Name})
		return false
	}

	if !compareLocaleMap(existing.NameLocalizations, desired.NameLocalizations) {
		logComparison("subcommand_name_locale_mismatch", map[string]any{"existing": existing.NameLocalizations, "desired": desired.NameLocalizations})
		return false
	}

	if !compareLocaleMap(existing.DescriptionLocalizations, desired.DescriptionLocalizations) {
		logComparison("subcommand_description_locale_mismatch", map[string]any{"existing": existing.DescriptionLocalizations, "desired": desired.DescriptionLocalizations})
		return false
	}

	return compareOptions(existing.Options, desired.Options)
}

func compareLocaleMap(existing map[discord.Locale]string, desired map[discord.Locale]string) bool {
	return reflect.DeepEqual(existing, desired)
}

func compareNullablePermissions(existing discord.Permissions, desired omit.Omit[*discord.Permissions]) bool {
	if !desired.OK {
		return existing == 0
	}

	if desired.Value == nil {
		return existing == 0
	}

	return existing == *desired.Value
}

func compareOptionalBool(existing bool, desired *bool, defaultValue bool) bool {
	if desired == nil {
		return existing == defaultValue
	}

	return existing == *desired
}

func compareIntPointer(existing *int, desired *int) bool {
	if existing == nil && desired == nil {
		return true
	}

	if existing == nil || desired == nil {
		logComparison("int_pointer_mismatch", map[string]any{"existing": existing, "desired": desired})
		return false
	}

	return *existing == *desired
}

func compareFloatPointer(existing *float64, desired *float64) bool {
	if existing == nil && desired == nil {
		return true
	}

	if existing == nil || desired == nil {
		logComparison("float_pointer_mismatch", map[string]any{"existing": existing, "desired": desired})
		return false
	}

	return *existing == *desired
}

func sameElements[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	set := make(map[T]int, len(a))
	for _, v := range a {
		set[v]++
	}
	for _, v := range b {
		if set[v] == 0 {
			return false
		}
		set[v]--
	}
	return true
}

func logComparison(reason string, details map[string]any) {
	attrs := []any{"reason", reason}
	for k, v := range details {
		attrs = append(attrs, k, v)
	}

	slog.Debug("command comparator diff", attrs...)
}
