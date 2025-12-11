package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Type aliases (using = syntax)
type UserID = string
type UserList = []User

// Type definition (not an alias, creates new type)
type ResponseCode int

// Response code constants
const (
	ResponseOK    ResponseCode = 200
	ResponseError ResponseCode = 500
)

// BaseEntity represents common fields for all entities
// This will be embedded in other structs for composition
type BaseEntity struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Timestamps is another embeddable struct for audit fields
type Timestamps struct {
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// User represents a user in the system
// Demonstrates struct composition through embedding
type User struct {
	BaseEntity            // Embedded struct (composition)
	Timestamps            // Another embedded struct
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	Active     bool       `json:"active"`
}

// Profile represents user profile information
// Also demonstrates struct embedding
type Profile struct {
	BaseEntity        // Embedded struct
	UserID     UserID `json:"user_id"` // Using type alias
	Bio        string `json:"bio"`
	AvatarURL  string `json:"avatar_url"`
	Website    string `json:"website"`
}

// APIResponse is a generic response wrapper
type APIResponse struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
}

// =====================================
// Value Receiver Methods on User
// =====================================

// FullName returns the full name of a user (value receiver)
// Value receivers work on a copy of the struct
func (u User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// IsActive checks if user is active (value receiver)
func (u User) IsActive() bool {
	return u.Active
}

// GetAge returns how long the user has existed (value receiver)
func (u User) GetAge() time.Duration {
	return time.Since(u.CreatedAt)
}

// String implements the Stringer interface (value receiver)
func (u User) String() string {
	return fmt.Sprintf("User{ID: %s, Name: %s, Email: %s}", u.ID, u.FullName(), u.Email)
}

// =====================================
// Pointer Receiver Methods on User
// =====================================

// UpdateEmail updates the user's email (pointer receiver)
// Pointer receivers can modify the original struct
func (u *User) UpdateEmail(email string) {
	u.Email = email
	u.UpdatedAt = time.Now()
}

// UpdateName updates the user's name (pointer receiver)
func (u *User) UpdateName(firstName, lastName string) {
	u.FirstName = firstName
	u.LastName = lastName
	u.UpdatedAt = time.Now()
}

// SetRole sets the user's role (pointer receiver)
func (u *User) SetRole(role string) {
	u.Role = role
	u.UpdatedAt = time.Now()
}

// Deactivate marks the user as inactive (pointer receiver)
func (u *User) Deactivate() {
	u.Active = false
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
}

// Activate marks the user as active (pointer receiver)
func (u *User) Activate() {
	u.Active = true
	u.DeletedAt = nil
	u.UpdatedAt = time.Now()
}

// Serialize converts user to JSON (pointer receiver - implements Serializer)
func (u *User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

// Deserialize populates user from JSON (pointer receiver - implements Serializer)
func (u *User) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}

// Validate checks if user data is valid (pointer receiver - implements Validator)
func (u *User) Validate() error {
	if u.ID == "" {
		return errors.New("user ID is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.FirstName == "" {
		return errors.New("first name is required")
	}
	return nil
}

// =====================================
// Constructor Functions
// =====================================

// NewUser creates a new User with initialized BaseEntity
func NewUser(id, firstName, lastName, email string) *User {
	now := time.Now()
	return &User{
		BaseEntity: BaseEntity{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Active:    true,
	}
}

// NewProfile creates a new Profile for a user
func NewProfile(id string, userID UserID) *Profile {
	now := time.Now()
	return &Profile{
		BaseEntity: BaseEntity{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID: userID,
	}
}

// =====================================
// Value Receiver Methods on BaseEntity
// =====================================

// GetID returns the entity ID (value receiver)
func (b BaseEntity) GetID() string {
	return b.ID
}

// Age returns the age of the entity (value receiver)
func (b BaseEntity) Age() time.Duration {
	return time.Since(b.CreatedAt)
}

// =====================================
// Pointer Receiver Methods on BaseEntity
// =====================================

// Touch updates the UpdatedAt timestamp (pointer receiver)
func (b *BaseEntity) Touch() {
	b.UpdatedAt = time.Now()
}

// =====================================
// Value Receiver Methods on Profile
// =====================================

// HasAvatar checks if profile has an avatar (value receiver)
func (p Profile) HasAvatar() bool {
	return p.AvatarURL != ""
}

// HasWebsite checks if profile has a website (value receiver)
func (p Profile) HasWebsite() bool {
	return p.Website != ""
}

// =====================================
// Pointer Receiver Methods on Profile
// =====================================

// SetBio updates the profile bio (pointer receiver)
func (p *Profile) SetBio(bio string) {
	p.Bio = bio
	p.UpdatedAt = time.Now()
}

// SetAvatarURL updates the avatar URL (pointer receiver)
func (p *Profile) SetAvatarURL(url string) {
	p.AvatarURL = url
	p.UpdatedAt = time.Now()
}
