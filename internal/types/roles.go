package types

// A role is like a hall pass that a user is given, and
// serves to turn on/off specific functionality in the
// backend.
type Role struct {
	Model
	Name        string
	Description string

	Users []User `json:"users" gorm:"many2many:user_roles;"`
}
