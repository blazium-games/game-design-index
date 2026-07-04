package internal

import "fmt"

func validateVariableCategory(label string, c VariableCategory) error {
	switch c {
	case VarCategoryStat, VarCategoryResource, VarCategoryCurrency, VarCategorySlot,
		VarCategoryMeter, VarCategoryCounter, VarCategoryFlag:
		return nil
	default:
		return fmt.Errorf("%s: invalid category %q", label, c)
	}
}

func validateVariableScope(label string, s VariableScope) error {
	switch s {
	case ScopePlayer, ScopeParty, ScopeWorld, ScopeSession, ScopeRun:
		return nil
	default:
		return fmt.Errorf("%s: invalid scope %q", label, s)
	}
}

func validateValueKind(label string, k ValueKind) error {
	switch k {
	case ValueInteger, ValueFloat, ValueBoolean, ValueEnum, ValueSlotList, ValueGrid:
		return nil
	default:
		return fmt.Errorf("%s: invalid value_kind %q", label, k)
	}
}

func validateResetBehavior(label string, r ResetBehavior) error {
	switch r {
	case "", ResetPerDeath, ResetPerLevel, ResetPerSession, ResetPersistent:
		return nil
	default:
		return fmt.Errorf("%s: invalid reset_behavior %q", label, r)
	}
}

func validateVariableRole(label string, r VariableRole) error {
	switch r {
	case VarRolePrimary, VarRoleSupporting, VarRoleHidden:
		return nil
	default:
		return fmt.Errorf("%s: invalid variable role %q", label, r)
	}
}

func validateMenuType(label string, t MenuType) error {
	switch t {
	case MenuMain, MenuPause, MenuOptions, MenuSettings, MenuInventory, MenuEquipment,
		MenuMap, MenuShop, MenuCrafting, MenuDialogue, MenuHUDOverlay, MenuSaveLoad,
		MenuCredits, MenuStageSelect, MenuControlsRemap, MenuLobby, MenuMeeting,
		MenuVoting, MenuHub, MenuWeaponSelect, MenuMapOverlay, MenuBoons, MenuMirror:
		return nil
	default:
		return fmt.Errorf("%s: invalid menu_type %q", label, t)
	}
}

func validateMenuLayer(label string, l MenuLayer) error {
	switch l {
	case LayerMeta, LayerInGame, LayerCombatOverlay:
		return nil
	default:
		return fmt.Errorf("%s: invalid layer %q", label, l)
	}
}

func validateInputContext(label string, c InputContext) error {
	switch c {
	case "", InputGamepad, InputKeyboardMouse, InputTouch, InputAny:
		return nil
	default:
		return fmt.Errorf("%s: invalid input_context %q", label, c)
	}
}

func validateMenuRole(label string, r MenuRole) error {
	switch r {
	case MenuRoleCore, MenuRoleOptional, MenuRoleContextual:
		return nil
	default:
		return fmt.Errorf("%s: invalid menu role %q", label, r)
	}
}

func validateGameVariable(entry GameVariable, tagVocab map[string]struct{}, variableSlugs map[string]struct{}, mechanicSlugs map[string]struct{}) error {
	label := entry.Slug
	if label == "" {
		label = entry.Name
	}
	if entry.SchemaVersion != VariableSchemaVersion {
		return fmt.Errorf("%s: schema_version must be %q", label, VariableSchemaVersion)
	}
	if err := validateSlug(label, entry.Slug); err != nil {
		return err
	}
	if entry.Name == "" {
		return fmt.Errorf("%s: name is required", label)
	}
	if entry.Summary == "" {
		return fmt.Errorf("%s: summary is required", label)
	}
	if err := validateVariableCategory(label, entry.Category); err != nil {
		return err
	}
	if err := validateVariableScope(label, entry.Scope); err != nil {
		return err
	}
	if err := validateValueKind(label, entry.ValueKind); err != nil {
		return err
	}
	if err := validateResetBehavior(label, entry.ResetBehavior); err != nil {
		return err
	}
	if len(entry.Tags) == 0 {
		return fmt.Errorf("%s: at least one tag is required", label)
	}
	for _, t := range entry.Tags {
		if tagVocab != nil {
			if _, ok := tagVocab[t]; !ok {
				return fmt.Errorf("%s: unknown tag %q", label, t)
			}
		}
	}
	for _, m := range entry.RelatedMechanics {
		if err := validateSlug(label+".related_mechanics", m); err != nil {
			return err
		}
		if mechanicSlugs != nil {
			if _, ok := mechanicSlugs[m]; !ok {
				return fmt.Errorf("%s: related_mechanics references unknown mechanic %q", label, m)
			}
		}
	}
	for _, rel := range entry.RelatedVariables {
		if err := validateSlug(label+".related_variables", rel.Slug); err != nil {
			return err
		}
		if variableSlugs != nil {
			if _, ok := variableSlugs[rel.Slug]; !ok {
				return fmt.Errorf("%s: related_variables references unknown variable %q", label, rel.Slug)
			}
		}
	}
	return nil
}

func validateUIMenu(entry UIMenu, tagVocab map[string]struct{}, menuSlugs map[string]struct{}, variableSlugs map[string]struct{}, mechanicSlugs map[string]struct{}) error {
	label := entry.Slug
	if label == "" {
		label = entry.Name
	}
	if entry.SchemaVersion != UIMenuSchemaVersion {
		return fmt.Errorf("%s: schema_version must be %q", label, UIMenuSchemaVersion)
	}
	if err := validateSlug(label, entry.Slug); err != nil {
		return err
	}
	if entry.Name == "" {
		return fmt.Errorf("%s: name is required", label)
	}
	if entry.Summary == "" {
		return fmt.Errorf("%s: summary is required", label)
	}
	if err := validateMenuType(label, entry.MenuType); err != nil {
		return err
	}
	if err := validateMenuLayer(label, entry.Layer); err != nil {
		return err
	}
	if err := validateInputContext(label, entry.InputContext); err != nil {
		return err
	}
	if len(entry.Tags) == 0 {
		return fmt.Errorf("%s: at least one tag is required", label)
	}
	for _, t := range entry.Tags {
		if tagVocab != nil {
			if _, ok := tagVocab[t]; !ok {
				return fmt.Errorf("%s: unknown tag %q", label, t)
			}
		}
	}
	for _, m := range entry.RelatedMechanics {
		if err := validateSlug(label+".related_mechanics", m); err != nil {
			return err
		}
		if mechanicSlugs != nil {
			if _, ok := mechanicSlugs[m]; !ok {
				return fmt.Errorf("%s: related_mechanics references unknown mechanic %q", label, m)
			}
		}
	}
	for _, v := range entry.RelatedVariables {
		if err := validateSlug(label+".related_variables", v); err != nil {
			return err
		}
		if variableSlugs != nil {
			if _, ok := variableSlugs[v]; !ok {
				return fmt.Errorf("%s: related_variables references unknown variable %q", label, v)
			}
		}
	}
	for _, rel := range entry.RelatedMenus {
		if err := validateSlug(label+".related_menus", rel.Slug); err != nil {
			return err
		}
		if menuSlugs != nil {
			if _, ok := menuSlugs[rel.Slug]; !ok {
				return fmt.Errorf("%s: related_menus references unknown menu %q", label, rel.Slug)
			}
		}
	}
	return nil
}

func validateVariablesCatalog(root string, catalog VariablesCatalog, mechanicSlugs map[string]struct{}) error {
	if catalog.SchemaVersion != "" && catalog.SchemaVersion != VariableSchemaVersion {
		return fmt.Errorf("variables catalog schema_version must be %q", VariableSchemaVersion)
	}
	tagVocab, err := LoadVariableTags(root)
	if err != nil {
		return err
	}
	slugs := make(map[string]struct{}, len(catalog.Variables))
	for _, entry := range catalog.Variables {
		slugs[entry.Slug] = struct{}{}
	}
	for _, entry := range catalog.Variables {
		if err := validateGameVariable(entry, tagVocab, slugs, mechanicSlugs); err != nil {
			return err
		}
	}
	return nil
}

func validateUIMenusCatalog(root string, catalog UIMenusCatalog, mechanicSlugs map[string]struct{}, variableSlugs map[string]struct{}) error {
	if catalog.SchemaVersion != "" && catalog.SchemaVersion != UIMenuSchemaVersion {
		return fmt.Errorf("ui-menus catalog schema_version must be %q", UIMenuSchemaVersion)
	}
	tagVocab, err := LoadMenuTags(root)
	if err != nil {
		return err
	}
	slugs := make(map[string]struct{}, len(catalog.Menus))
	for _, entry := range catalog.Menus {
		slugs[entry.Slug] = struct{}{}
	}
	for _, entry := range catalog.Menus {
		if err := validateUIMenu(entry, tagVocab, slugs, variableSlugs, mechanicSlugs); err != nil {
			return err
		}
	}
	return nil
}

func validateMapVariableBindings(label string, bindings []MapVariableBinding) error {
	for _, b := range bindings {
		if err := validateSlug(label+".variables", b.VariableSlug); err != nil {
			return err
		}
		if err := validateVariableRole(label+".variables", b.Role); err != nil {
			return err
		}
		for _, m := range b.RelatedMechanics {
			if err := validateSlug(label+".variables.related_mechanics", m); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateMapUIMenuBindings(label string, bindings []MapUIMenuBinding) error {
	for _, b := range bindings {
		if err := validateSlug(label+".ui_menus", b.MenuSlug); err != nil {
			return err
		}
		if err := validateMenuRole(label+".ui_menus", b.Role); err != nil {
			return err
		}
		for _, from := range b.OpensFrom {
			if err := validateSlug(label+".ui_menus.opens_from", from); err != nil {
				return err
			}
		}
		for _, v := range b.DisplaysVariables {
			if err := validateSlug(label+".ui_menus.displays_variables", v); err != nil {
				return err
			}
		}
		for _, m := range b.SupportsMechanics {
			if err := validateSlug(label+".ui_menus.supports_mechanics", m); err != nil {
				return err
			}
		}
	}
	return nil
}
