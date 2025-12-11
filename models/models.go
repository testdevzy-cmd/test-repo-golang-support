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
type OrgID = string
type OrgList = []Organization

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

// =====================================
// Organization Model (NEW)
// =====================================

// Address represents a physical address
// Embeddable struct for composition
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

// ContactInfo represents contact information
// Another embeddable struct
type ContactInfo struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

// Organization represents a company or organization
// Demonstrates multiple struct embeddings
type Organization struct {
	BaseEntity              // Embedded struct (composition)
	Timestamps              // Embedded struct for soft delete
	Address                 // Embedded struct for address
	ContactInfo             // Embedded struct for contact
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Industry    string      `json:"industry"`
	Size        OrgSize     `json:"size"`
	Active      bool        `json:"active"`
	OwnerID     UserID      `json:"owner_id"` // Using type alias
}

// OrgSize represents organization size category
type OrgSize string

// Organization size constants
const (
	OrgSizeSmall      OrgSize = "small"
	OrgSizeMedium     OrgSize = "medium"
	OrgSizeLarge      OrgSize = "large"
	OrgSizeEnterprise OrgSize = "enterprise"
)

// Membership represents user membership in an organization
type Membership struct {
	BaseEntity            // Embedded struct
	UserID     UserID     `json:"user_id"`
	OrgID      OrgID      `json:"org_id"`
	Role       MemberRole `json:"role"`
	JoinedAt   time.Time  `json:"joined_at"`
}

// MemberRole represents the role of a member in an organization
type MemberRole string

// Member role constants
const (
	MemberRoleOwner  MemberRole = "owner"
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
	MemberRoleGuest  MemberRole = "guest"
)

// =====================================
// Value Receiver Methods on Organization
// =====================================

// DisplayName returns the organization display name (value receiver)
func (o Organization) DisplayName() string {
	return o.Name
}

// IsActive checks if organization is active (value receiver)
func (o Organization) IsActive() bool {
	return o.Active
}

// HasWebsite checks if organization has a website (value receiver)
func (o Organization) HasWebsite() bool {
	return o.ContactInfo.Website != ""
}

// FullAddress returns the complete address as a string (value receiver)
func (o Organization) FullAddress() string {
	return fmt.Sprintf("%s, %s, %s %s, %s",
		o.Address.Street,
		o.Address.City,
		o.Address.State,
		o.Address.PostalCode,
		o.Address.Country,
	)
}

// String implements Stringer interface (value receiver)
func (o Organization) String() string {
	return fmt.Sprintf("Organization{ID: %s, Name: %s, Industry: %s}", o.ID, o.Name, o.Industry)
}

// =====================================
// Pointer Receiver Methods on Organization
// =====================================

// UpdateName updates the organization name (pointer receiver)
func (o *Organization) UpdateName(name string) {
	o.Name = name
	o.UpdatedAt = time.Now()
}

// UpdateDescription updates the description (pointer receiver)
func (o *Organization) UpdateDescription(desc string) {
	o.Description = desc
	o.UpdatedAt = time.Now()
}

// SetIndustry sets the industry (pointer receiver)
func (o *Organization) SetIndustry(industry string) {
	o.Industry = industry
	o.UpdatedAt = time.Now()
}

// SetSize sets the organization size (pointer receiver)
func (o *Organization) SetSize(size OrgSize) {
	o.Size = size
	o.UpdatedAt = time.Now()
}

// UpdateAddress updates the address (pointer receiver)
func (o *Organization) UpdateAddress(addr Address) {
	o.Address = addr
	o.UpdatedAt = time.Now()
}

// UpdateContact updates contact info (pointer receiver)
func (o *Organization) UpdateContact(contact ContactInfo) {
	o.ContactInfo = contact
	o.UpdatedAt = time.Now()
}

// Deactivate marks organization as inactive (pointer receiver)
func (o *Organization) Deactivate() {
	o.Active = false
	now := time.Now()
	o.DeletedAt = &now
	o.UpdatedAt = now
}

// Activate marks organization as active (pointer receiver)
func (o *Organization) Activate() {
	o.Active = true
	o.DeletedAt = nil
	o.UpdatedAt = time.Now()
}

// Serialize converts organization to JSON (pointer receiver)
func (o *Organization) Serialize() ([]byte, error) {
	return json.Marshal(o)
}

// Deserialize populates organization from JSON (pointer receiver)
func (o *Organization) Deserialize(data []byte) error {
	return json.Unmarshal(data, o)
}

// Validate checks if organization data is valid (pointer receiver)
func (o *Organization) Validate() error {
	if o.ID == "" {
		return errors.New("organization ID is required")
	}
	if o.Name == "" {
		return errors.New("organization name is required")
	}
	if o.OwnerID == "" {
		return errors.New("owner ID is required")
	}
	return nil
}

// =====================================
// Constructor Functions for Organization
// =====================================

// NewOrganization creates a new Organization with initialized fields
func NewOrganization(id, name, ownerID string) *Organization {
	now := time.Now()
	return &Organization{
		BaseEntity: BaseEntity{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:    name,
		OwnerID: ownerID,
		Active:  true,
		Size:    OrgSizeSmall,
	}
}

// NewMembership creates a new Membership
func NewMembership(id string, userID UserID, orgID OrgID, role MemberRole) *Membership {
	now := time.Now()
	return &Membership{
		BaseEntity: BaseEntity{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:   userID,
		OrgID:    orgID,
		Role:     role,
		JoinedAt: now,
	}
}

// =====================================
// Value Receiver Methods on Membership
// =====================================

// IsOwner checks if membership is owner role (value receiver)
func (m Membership) IsOwner() bool {
	return m.Role == MemberRoleOwner
}

// IsAdmin checks if membership has admin privileges (value receiver)
func (m Membership) IsAdmin() bool {
	return m.Role == MemberRoleOwner || m.Role == MemberRoleAdmin
}

// CanManageMembers checks if member can manage other members (value receiver)
func (m Membership) CanManageMembers() bool {
	return m.IsAdmin()
}

// =====================================
// Pointer Receiver Methods on Membership
// =====================================

// ChangeRole changes the membership role (pointer receiver)
func (m *Membership) ChangeRole(role MemberRole) {
	m.Role = role
	m.UpdatedAt = time.Now()
}

// Promote promotes member to admin (pointer receiver)
func (m *Membership) Promote() {
	if m.Role == MemberRoleMember || m.Role == MemberRoleGuest {
		m.Role = MemberRoleAdmin
		m.UpdatedAt = time.Now()
	}
}

// Demote demotes member from admin (pointer receiver)
func (m *Membership) Demote() {
	if m.Role == MemberRoleAdmin {
		m.Role = MemberRoleMember
		m.UpdatedAt = time.Now()
	}
}
