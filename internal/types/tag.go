package types

import "github.com/google/uuid"

// a tag is way to categorize ingredients, and by that very nature,
// products.
//
// * a user may create a tag,	or they may be system specific
// tags.
// * the relationship between tags and users are stored in many2many
// tables.
type Tag struct {
	Model

	// the string representation of a tag, typically
	// one word, and it must be unique to the user.
	Name string `gorm:"type:text;not null;unique;"`

	// This is an optional description of what the Tag
	// represents.
	Description string `gorm:"type:text;"`

	// The list of ingredients that have this tag applied to
	// them.
	Ingredients []Ingredient `gorm:"many2many:ingredient_tags;"`
	Products    []Product    `gorm:"many2many:product_tags;"`
	Users       []User       `gorm:"many2many:user_tags;"`
}

// A tagging rule is our "Auto-Tagging" mechanism. It applies
// tags to new ingredients by performing a regex based search on
// the ingredient labels.
type TaggingRule struct {
	Model

	// When set to true, this states that the text provided
	// should be checked for a "conatins" relationship. If
	// marked as false, we are checking to see if it does
	// not exists, and applying the tag if that's the case.
	Contains bool   `gorm:"type:boolean;default:false"`
	Regex    string `gorm:"type:text;not null;"`

	// This is a foreign key relationship. The TaggingRule
	// requires a "Tag" to apply to the matching ingredients.
	TagID uuid.UUID `gorm:"type:uuid;not null"`
	Tag   Tag       `gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`

	// The tagging set that this rule belongs to, can
	// be left empty.
	TaggingSetID *uuid.UUID `gorm:"type:uuid;"`
	TaggingSet   TaggingSet `gorm:"foreignKey:TaggingSetID;constraint:OnDelete:SET NULL"`

	// Whether this tagging rule is enabled by the user.
	Enabled bool `gorm:"type:boolean;default:true"`

	// The user that created this tagging rule.
	UserId uuid.UUID `gorm:"type:uuid;not null"`
	User   User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

// A tagging set is a collection of tagging rules that can be
// shared (set to Public) by users wishing to share their
// "filters", etc.
type TaggingSet struct {
	Model

	// The name of the tagging set, visible to all users.
	Name string `gorm:"type:text;not null;"`

	// The description of the tagging set, optional.
	Description string `gorm:"type:text;"`

	// Whether this tagging set is a public set that
	// should be available to copy by all users.
	Public bool `gorm:"type:boolean;default:false"`

	// The list of users that have access to the
	// relevant tagging set.
	Users []User `gorm:"many2many:user_tagging_sets;"`
}

// The following is a "tagging job". This is the struct submitted
// to the "Tagging Daemon" so it can process the tagging request.
type TaggingJob struct {
	Model

	// the number of ingredients that have been scanned.
	IngredientScanNum int `gorm:"type:int;default:0"`

	// any errors that came about from the tagging job.
	Error string `gorm:"type:text;"`

	// the tag rule that accompanies it, user informtion can be determined
	// from this struct.
	TagRuleID uuid.UUID   `gorm:"type:uuid;not null;"`
	TagRule   TaggingRule `gorm:"foreignKey:TagRuleID;constraint:OnDelete:CASCADE"`
}
