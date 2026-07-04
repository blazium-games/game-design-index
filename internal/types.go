package internal

const SchemaVersion = "1.2"

// MapSchemaVersion is the current gameplay map document version (maps migrate independently).
const MapSchemaVersion = "1.1"

// VariableSchemaVersion is the game variable catalog version.
const VariableSchemaVersion = "1.0"

// UIMenuSchemaVersion is the UI menu catalog version.
const UIMenuSchemaVersion = "1.0"

type Flavor string

const (
	FlavorAction    Flavor = "action"
	FlavorAdventure Flavor = "adventure"
	FlavorStrategy  Flavor = "strategy"
)

type MapType string

const (
	MapTypeGame  MapType = "game"
	MapTypeGenre MapType = "genre"
)

type MechanicRole string

const (
	RoleSignature  MechanicRole = "signature"
	RoleSupporting MechanicRole = "supporting"
	RoleCommon     MechanicRole = "common"
)

type QualityTier string

const (
	QualityCurated  QualityTier = "curated"
	QualityCatalog  QualityTier = "catalog"
	QualityTemplate QualityTier = "template"
	QualityStub     QualityTier = "stub"
)

type MechanicDomain string

const (
	DomainLocomotion  MechanicDomain = "locomotion"
	DomainCombat      MechanicDomain = "combat"
	DomainProgression MechanicDomain = "progression"
	DomainEconomy     MechanicDomain = "economy"
	DomainLevel       MechanicDomain = "level"
	DomainSession     MechanicDomain = "session"
)

type MechanicPhase string

const (
	PhaseEarly    MechanicPhase = "early"
	PhaseMid      MechanicPhase = "mid"
	PhaseLate     MechanicPhase = "late"
	PhaseOptional MechanicPhase = "optional"
)

type ImplementationComplexity string

const (
	ComplexitySmall  ImplementationComplexity = "S"
	ComplexityMedium ImplementationComplexity = "M"
	ComplexityLarge  ImplementationComplexity = "L"
)

type RequirementKind string

const (
	RequirementMechanic RequirementKind = "mechanic"
	RequirementDomain   RequirementKind = "domain"
	RequirementTag      RequirementKind = "tag"
)

type MechanicRequirement struct {
	Kind     RequirementKind `json:"kind"`
	Ref      string          `json:"ref"`
	Optional bool            `json:"optional,omitempty"`
	Note     string          `json:"note,omitempty"`
}

type ParameterKnob struct {
	Name   string `json:"name"`
	Range  string `json:"range,omitempty"`
	Effect string `json:"effect,omitempty"`
}

type SignatureOf struct {
	Games  []string `json:"games,omitempty"`
	Genres []string `json:"genres,omitempty"`
}

type Implementation struct {
	GodotAssetAvailable      bool                      `json:"godot_asset_available"`
	TutorialURL              *string                   `json:"tutorial_url"`
	AssetURL                 *string                   `json:"asset_url"`
	ImplementationComplexity *ImplementationComplexity `json:"implementation_complexity,omitempty"`
}

type DesignGuidance struct {
	WhenToUse     string   `json:"when_to_use,omitempty"`
	WhereToUse    string   `json:"where_to_use,omitempty"`
	WhenToAvoid   []string `json:"when_to_avoid,omitempty"`
	DesignerNotes string   `json:"designer_notes,omitempty"`
}

type SynergyNote struct {
	Slug string `json:"slug"`
	Note string `json:"note,omitempty"`
}

type MechanicExample struct {
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
	MapSlug     string `json:"map_slug,omitempty"`
}

type RelationshipEdge struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Relationship string `json:"relationship,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type RelationshipModel struct {
	Type  string             `json:"type"`
	Edges []RelationshipEdge `json:"edges,omitempty"`
}

type AgentContext struct {
	SummaryForAgents        string   `json:"summary_for_agents,omitempty"`
	GDDPrompt               string   `json:"gdd_prompt,omitempty"`
	ImplementationChecklist []string `json:"implementation_checklist,omitempty"`
}

type MechanicMedia struct {
	SnippetPath string `json:"snippet_path,omitempty"`
}

type MechanicEntry struct {
	SchemaVersion      string                   `json:"schema_version"`
	Slug               string                   `json:"slug"`
	Name               string                   `json:"name"`
	Flavor             Flavor                   `json:"flavor"`
	Summary            string                   `json:"summary"`
	Domain             MechanicDomain           `json:"domain"`
	Tags               []string                 `json:"tags"`
	ProblemSolved      string                   `json:"problem_solved,omitempty"`
	Conflicts          []string                 `json:"conflicts,omitempty"`
	Synergies          []string                 `json:"synergies,omitempty"`
	Prerequisites      []string                 `json:"prerequisites,omitempty"`
	Requirements       []MechanicRequirement    `json:"requirements,omitempty"`
	AntiPatterns       []string                 `json:"anti_patterns,omitempty"`
	ParameterKnobs     []ParameterKnob          `json:"parameter_knobs,omitempty"`
	SignatureOf        *SignatureOf             `json:"signature_of,omitempty"`
	FeaturedIn         []string                 `json:"featured_in,omitempty"`
	CommonIn           []string                 `json:"common_in,omitempty"`
	Implementation     *Implementation          `json:"implementation,omitempty"`
	FlavorRationale    string                   `json:"flavor_rationale,omitempty"`
	DesignGuidance     *DesignGuidance          `json:"design_guidance,omitempty"`
	SynergyNotes       []SynergyNote            `json:"synergy_notes,omitempty"`
	PlayerExperience   string                   `json:"player_experience,omitempty"`
	Complexity         ImplementationComplexity `json:"complexity,omitempty"`
	Examples           []MechanicExample        `json:"examples,omitempty"`
	RelationshipModel  *RelationshipModel       `json:"relationship_model,omitempty"`
	AgentContext       *AgentContext            `json:"agent_context,omitempty"`
	Media              *MechanicMedia           `json:"media,omitempty"`
}

type MechanicTagsCatalog struct {
	SchemaVersion string   `json:"schema_version"`
	Tags          []string `json:"tags"`
}

type MechanicsCatalog struct {
	SchemaVersion string          `json:"schema_version"`
	Mechanics     []MechanicEntry `json:"mechanics"`
}

type Subject struct {
	Kind       string   `json:"kind"`
	Name       string   `json:"name"`
	Genres     []string `json:"genres,omitempty"`
	Influences []string `json:"influences,omitempty"`
}

type Narrative struct {
	Description   string   `json:"description"`
	CoreLoop      string   `json:"core_loop,omitempty"`
	SkillsTested  []string `json:"skills_tested,omitempty"`
	DesignPillars []string `json:"design_pillars,omitempty"`
}

type PlayerCount struct {
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
}

type MapContext struct {
	Year           int          `json:"year,omitempty"`
	Platforms      []string     `json:"platforms,omitempty"`
	Perspective    string       `json:"perspective,omitempty"`
	Dimension      string       `json:"dimension,omitempty"`
	PlayerCount    *PlayerCount `json:"player_count,omitempty"`
	SessionType    string       `json:"session_type,omitempty"`
	WorldStructure string       `json:"world_structure,omitempty"`
}

type MapMetadata struct {
	QualityTier QualityTier `json:"quality_tier,omitempty"`
}

type MapRelationships struct {
	SimilarTo      []string `json:"similar_to,omitempty"`
	GenreMaps      []string `json:"genre_maps,omitempty"`
	InfluenceSlugs []string `json:"influence_slugs,omitempty"`
}

type DesignIntent struct {
	PlayerFantasy string   `json:"player_fantasy,omitempty"`
	EmotionalTone []string `json:"emotional_tone,omitempty"`
	ThemeTags     []string `json:"theme_tags,omitempty"`
	DesignPillars []string `json:"design_pillars,omitempty"`
}

type LoopPhase struct {
	ID        string   `json:"id"`
	Label     string   `json:"label"`
	Mechanics []string `json:"mechanics,omitempty"`
}

type PacingCurve struct {
	Early string `json:"early,omitempty"`
	Mid   string `json:"mid,omitempty"`
	Late  string `json:"late,omitempty"`
}

type MapSystems struct {
	PrimaryLoopPhases []LoopPhase  `json:"primary_loop_phases,omitempty"`
	FailureModel      string       `json:"failure_model,omitempty"`
	EconomyTightness  string       `json:"economy_tightness,omitempty"`
	Metaprogression   *bool        `json:"metaprogression,omitempty"`
	Pacing            *PacingCurve `json:"pacing,omitempty"`
}

type MapComposition struct {
	BaseGenreMaps      []string `json:"base_genre_maps,omitempty"`
	ReferenceGames     []string `json:"reference_games,omitempty"`
	CustomizationNotes string   `json:"customization_notes,omitempty"`
}

type GenreVariant struct {
	ID             string   `json:"id"`
	Label          string   `json:"label"`
	AddSignatures  []string `json:"add_signatures,omitempty"`
	DropSignatures []string `json:"drop_signatures,omitempty"`
	Notes          string   `json:"notes,omitempty"`
}

type MapMechanicBinding struct {
	MechanicSlug string         `json:"mechanic_slug"`
	Role         MechanicRole   `json:"role"`
	MapNotes     string         `json:"map_notes,omitempty"`
	Domain       MechanicDomain `json:"domain"`
	Phase        MechanicPhase  `json:"phase,omitempty"`
	Weight       int            `json:"weight,omitempty"`
	Expression   string         `json:"expression,omitempty"`
	DependsOn    []string       `json:"depends_on,omitempty"`
	Optional     bool           `json:"optional,omitempty"`
}

type ViewFilter struct {
	Role        *MechanicRole `json:"role,omitempty"`
	Flavor      *Flavor       `json:"flavor,omitempty"`
	GenreFilter *string       `json:"genre_filter,omitempty"`
}

type MapView struct {
	ID     string     `json:"id"`
	Label  string     `json:"label"`
	Filter ViewFilter `json:"filter"`
}

type GDDOutline struct {
	Overview         string   `json:"overview,omitempty"`
	CoreLoop         string   `json:"core_loop,omitempty"`
	PlayerGoals      []string `json:"player_goals,omitempty"`
	Constraints      []string `json:"constraints,omitempty"`
	ProgressionNotes string   `json:"progression_notes,omitempty"`
	CombatNotes      string   `json:"combat_notes,omitempty"`
	EconomyNotes     string   `json:"economy_notes,omitempty"`
}

type MechanicRelationship struct {
	FromMechanic string `json:"from_mechanic"`
	ToMechanic   string `json:"to_mechanic"`
	Relationship string `json:"relationship,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type VariableCategory string

const (
	VarCategoryStat     VariableCategory = "stat"
	VarCategoryResource VariableCategory = "resource"
	VarCategoryCurrency VariableCategory = "currency"
	VarCategorySlot     VariableCategory = "slot"
	VarCategoryMeter    VariableCategory = "meter"
	VarCategoryCounter  VariableCategory = "counter"
	VarCategoryFlag     VariableCategory = "flag"
)

type VariableScope string

const (
	ScopePlayer  VariableScope = "player"
	ScopeParty   VariableScope = "party"
	ScopeWorld   VariableScope = "world"
	ScopeSession VariableScope = "session"
	ScopeRun     VariableScope = "run"
)

type ValueKind string

const (
	ValueInteger  ValueKind = "integer"
	ValueFloat    ValueKind = "float"
	ValueBoolean  ValueKind = "boolean"
	ValueEnum     ValueKind = "enum"
	ValueSlotList ValueKind = "slot_list"
	ValueGrid     ValueKind = "grid"
)

type ResetBehavior string

const (
	ResetPerDeath   ResetBehavior = "per_death"
	ResetPerLevel   ResetBehavior = "per_level"
	ResetPerSession ResetBehavior = "per_session"
	ResetPersistent ResetBehavior = "persistent"
)

type EntityRelationship struct {
	Slug         string `json:"slug"`
	Relationship string `json:"relationship,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type GameVariable struct {
	SchemaVersion    string               `json:"schema_version"`
	Slug             string               `json:"slug"`
	Name             string               `json:"name"`
	Summary          string               `json:"summary"`
	Category         VariableCategory     `json:"category"`
	Scope            VariableScope        `json:"scope"`
	ValueKind        ValueKind            `json:"value_kind"`
	SharedRationale  string               `json:"shared_rationale,omitempty"`
	PlayerFocus      string               `json:"player_focus,omitempty"`
	TypicalRange     string               `json:"typical_range,omitempty"`
	ResetBehavior    ResetBehavior        `json:"reset_behavior,omitempty"`
	RelatedMechanics []string             `json:"related_mechanics,omitempty"`
	RelatedVariables []EntityRelationship `json:"related_variables,omitempty"`
	FeaturedIn       []string             `json:"featured_in,omitempty"`
	CommonIn         []string             `json:"common_in,omitempty"`
	Tags             []string             `json:"tags"`
	ParameterKnobs   []ParameterKnob      `json:"parameter_knobs,omitempty"`
	DesignGuidance   *DesignGuidance      `json:"design_guidance,omitempty"`
}

type VariablesCatalog struct {
	SchemaVersion string         `json:"schema_version"`
	Variables     []GameVariable `json:"variables"`
}

type MenuType string

const (
	MenuMain          MenuType = "main"
	MenuPause         MenuType = "pause"
	MenuOptions       MenuType = "options"
	MenuSettings      MenuType = "settings"
	MenuInventory     MenuType = "inventory"
	MenuEquipment     MenuType = "equipment"
	MenuMap           MenuType = "map"
	MenuShop          MenuType = "shop"
	MenuCrafting      MenuType = "crafting"
	MenuDialogue      MenuType = "dialogue"
	MenuHUDOverlay    MenuType = "hud_overlay"
	MenuSaveLoad      MenuType = "save_load"
	MenuCredits       MenuType = "credits"
	MenuStageSelect   MenuType = "stage_select"
	MenuControlsRemap MenuType = "controls_remap"
	MenuLobby         MenuType = "lobby"
	MenuMeeting       MenuType = "meeting"
	MenuVoting        MenuType = "voting"
	MenuHub           MenuType = "hub"
	MenuWeaponSelect  MenuType = "weapon_select"
	MenuMapOverlay    MenuType = "map_overlay"
	MenuBoons         MenuType = "boons"
	MenuMirror        MenuType = "mirror"
)

type MenuLayer string

const (
	LayerMeta          MenuLayer = "meta"
	LayerInGame        MenuLayer = "in_game"
	LayerCombatOverlay MenuLayer = "combat_overlay"
)

type InputContext string

const (
	InputGamepad       InputContext = "gamepad"
	InputKeyboardMouse InputContext = "keyboard_mouse"
	InputTouch         InputContext = "touch"
	InputAny           InputContext = "any"
)

type UIMenu struct {
	SchemaVersion    string               `json:"schema_version"`
	Slug             string               `json:"slug"`
	Name             string               `json:"name"`
	Summary          string               `json:"summary"`
	MenuType         MenuType             `json:"menu_type"`
	Layer            MenuLayer            `json:"layer"`
	InputContext     InputContext         `json:"input_context,omitempty"`
	SharedRationale  string               `json:"shared_rationale,omitempty"`
	TypicalActions   []string             `json:"typical_actions,omitempty"`
	RelatedMechanics []string             `json:"related_mechanics,omitempty"`
	RelatedVariables []string             `json:"related_variables,omitempty"`
	RelatedMenus     []EntityRelationship `json:"related_menus,omitempty"`
	FeaturedIn       []string             `json:"featured_in,omitempty"`
	CommonIn         []string             `json:"common_in,omitempty"`
	Tags             []string             `json:"tags"`
	DesignGuidance   *DesignGuidance      `json:"design_guidance,omitempty"`
}

type UIMenusCatalog struct {
	SchemaVersion string   `json:"schema_version"`
	Menus         []UIMenu `json:"menus"`
}

type VariableRole string

const (
	VarRolePrimary    VariableRole = "primary"
	VarRoleSupporting VariableRole = "supporting"
	VarRoleHidden     VariableRole = "hidden"
)

type MenuRole string

const (
	MenuRoleCore       MenuRole = "core"
	MenuRoleOptional   MenuRole = "optional"
	MenuRoleContextual MenuRole = "contextual"
)

type MapVariableBinding struct {
	VariableSlug     string       `json:"variable_slug"`
	Role             VariableRole `json:"role"`
	MapNotes         string       `json:"map_notes,omitempty"`
	Expression       string       `json:"expression,omitempty"`
	RelatedMechanics []string     `json:"related_mechanics,omitempty"`
}

type MapUIMenuBinding struct {
	MenuSlug          string   `json:"menu_slug"`
	Role              MenuRole `json:"role"`
	MapNotes          string   `json:"map_notes,omitempty"`
	OpensFrom         []string `json:"opens_from,omitempty"`
	DisplaysVariables []string `json:"displays_variables,omitempty"`
	SupportsMechanics []string `json:"supports_mechanics,omitempty"`
}

type VariableRelationship struct {
	FromVariable string `json:"from_variable"`
	ToVariable   string `json:"to_variable"`
	Relationship string `json:"relationship,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type MenuFlowEdge struct {
	FromMenu     string `json:"from_menu"`
	ToMenu       string `json:"to_menu"`
	Relationship string `json:"relationship,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type GameplayMap struct {
	SchemaVersion         string                 `json:"schema_version"`
	Slug                  string                 `json:"slug"`
	Title                 string                 `json:"title"`
	MapType               MapType                `json:"map_type"`
	SourceURL             string                 `json:"source_url"`
	Subject               Subject                `json:"subject"`
	Narrative             Narrative              `json:"narrative"`
	SignatureGameplay     []string               `json:"signature_gameplay"`
	Mechanics             []MapMechanicBinding   `json:"mechanics"`
	Views                 []MapView              `json:"views"`
	Context               *MapContext            `json:"context,omitempty"`
	Metadata              *MapMetadata           `json:"metadata,omitempty"`
	Relationships         *MapRelationships      `json:"relationships,omitempty"`
	DesignIntent          *DesignIntent          `json:"design_intent,omitempty"`
	Systems               *MapSystems            `json:"systems,omitempty"`
	Composition           *MapComposition        `json:"composition,omitempty"`
	Variants              []GenreVariant         `json:"variants,omitempty"`
	GDDOutline            *GDDOutline            `json:"gdd_outline,omitempty"`
	MechanicRelationships []MechanicRelationship `json:"mechanic_relationships,omitempty"`
	Variables             []MapVariableBinding   `json:"variables,omitempty"`
	UIMenus               []MapUIMenuBinding     `json:"ui_menus,omitempty"`
	VariableRelationships []VariableRelationship `json:"variable_relationships,omitempty"`
	MenuFlow              []MenuFlowEdge         `json:"menu_flow,omitempty"`
}

type Bundle struct {
	Root      string
	Mechanics map[string]MechanicEntry
	Variables map[string]GameVariable
	UIMenus   map[string]UIMenu
	Maps      map[string]GameplayMap
}
