package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/test-repo-golang-support/models"
)

// OrganizationService handles organization-related operations
type OrganizationService struct {
	orgs        map[string]*models.Organization
	memberships map[string]*models.Membership // key: "userID:orgID"
	mu          sync.RWMutex
}

// NewOrganizationService creates a new OrganizationService instance
func NewOrganizationService() *OrganizationService {
	return &OrganizationService{
		orgs:        make(map[string]*models.Organization),
		memberships: make(map[string]*models.Membership),
	}
}

// =====================================
// Value Receiver Methods on OrganizationService
// =====================================

// Count returns the number of organizations (value receiver)
func (s OrganizationService) Count() int {
	return len(s.orgs)
}

// HasOrgs checks if there are any organizations (value receiver)
func (s OrganizationService) HasOrgs() bool {
	return len(s.orgs) > 0
}

// IsEmpty checks if the service has no organizations (value receiver)
func (s OrganizationService) IsEmpty() bool {
	return len(s.orgs) == 0
}

// MembershipCount returns total membership count (value receiver)
func (s OrganizationService) MembershipCount() int {
	return len(s.memberships)
}

// =====================================
// Pointer Receiver Methods - OrgReader Implementation
// =====================================

// ReadOrg retrieves an organization by ID (pointer receiver)
func (s *OrganizationService) ReadOrg(ctx context.Context, id string) (*models.Organization, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	org, exists := s.orgs[id]
	if !exists {
		return nil, errors.New("organization not found")
	}
	return org, nil
}

// ReadAllOrgs retrieves all organizations (pointer receiver)
func (s *OrganizationService) ReadAllOrgs(ctx context.Context) (models.OrgList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orgs := make(models.OrgList, 0, len(s.orgs))
	for _, org := range s.orgs {
		orgs = append(orgs, *org)
	}
	return orgs, nil
}

// ReadOrgsByOwner retrieves organizations by owner ID (pointer receiver)
func (s *OrganizationService) ReadOrgsByOwner(ctx context.Context, ownerID string) (models.OrgList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orgs := make(models.OrgList, 0)
	for _, org := range s.orgs {
		if org.OwnerID == ownerID {
			orgs = append(orgs, *org)
		}
	}
	return orgs, nil
}

// =====================================
// Pointer Receiver Methods - OrgWriter Implementation
// =====================================

// WriteOrg creates or updates an organization (pointer receiver)
func (s *OrganizationService) WriteOrg(ctx context.Context, org *models.Organization) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if org.ID == "" {
		return errors.New("organization ID is required")
	}
	s.orgs[org.ID] = org
	return nil
}

// DeleteOrg removes an organization (pointer receiver)
func (s *OrganizationService) DeleteOrg(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orgs[id]; !exists {
		return errors.New("organization not found")
	}
	delete(s.orgs, id)

	// Also remove all memberships for this org
	for key, m := range s.memberships {
		if m.OrgID == id {
			delete(s.memberships, key)
		}
	}
	return nil
}

// =====================================
// Pointer Receiver Methods - OrgRepository Implementation
// =====================================

// CountOrgs returns total organization count (pointer receiver)
func (s *OrganizationService) CountOrgs(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.orgs), nil
}

// OrgExists checks if an organization exists (pointer receiver)
func (s *OrganizationService) OrgExists(ctx context.Context, id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.orgs[id]
	return exists, nil
}

// =====================================
// Pointer Receiver Methods - MembershipManager Implementation
// =====================================

// membershipKey generates a unique key for user-org membership
func membershipKey(userID, orgID string) string {
	return fmt.Sprintf("%s:%s", userID, orgID)
}

// AddMember adds a member to an organization (pointer receiver)
func (s *OrganizationService) AddMember(ctx context.Context, membership *models.Membership) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify organization exists
	if _, exists := s.orgs[membership.OrgID]; !exists {
		return errors.New("organization not found")
	}

	key := membershipKey(membership.UserID, membership.OrgID)
	if _, exists := s.memberships[key]; exists {
		return errors.New("membership already exists")
	}

	s.memberships[key] = membership
	return nil
}

// RemoveMember removes a member from an organization (pointer receiver)
func (s *OrganizationService) RemoveMember(ctx context.Context, userID, orgID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := membershipKey(userID, orgID)
	if _, exists := s.memberships[key]; !exists {
		return errors.New("membership not found")
	}

	delete(s.memberships, key)
	return nil
}

// GetMembers retrieves all members of an organization (pointer receiver)
func (s *OrganizationService) GetMembers(ctx context.Context, orgID string) ([]*models.Membership, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	members := make([]*models.Membership, 0)
	for _, m := range s.memberships {
		if m.OrgID == orgID {
			members = append(members, m)
		}
	}
	return members, nil
}

// GetMembership retrieves a specific membership (pointer receiver)
func (s *OrganizationService) GetMembership(ctx context.Context, userID, orgID string) (*models.Membership, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := membershipKey(userID, orgID)
	membership, exists := s.memberships[key]
	if !exists {
		return nil, errors.New("membership not found")
	}
	return membership, nil
}

// UpdateMemberRole updates a member's role (pointer receiver)
func (s *OrganizationService) UpdateMemberRole(ctx context.Context, userID, orgID string, role models.MemberRole) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := membershipKey(userID, orgID)
	membership, exists := s.memberships[key]
	if !exists {
		return errors.New("membership not found")
	}

	membership.ChangeRole(role)
	return nil
}

// =====================================
// Additional Pointer Receiver Methods
// =====================================

// FindOrgByName finds an organization by name (pointer receiver)
func (s *OrganizationService) FindOrgByName(ctx context.Context, name string) (*models.Organization, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, org := range s.orgs {
		if org.Name == name {
			return org, nil
		}
	}
	return nil, errors.New("organization not found")
}

// FindOrgsByIndustry finds organizations by industry (pointer receiver)
func (s *OrganizationService) FindOrgsByIndustry(ctx context.Context, industry string) (models.OrgList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orgs := make(models.OrgList, 0)
	for _, org := range s.orgs {
		if org.Industry == industry {
			orgs = append(orgs, *org)
		}
	}
	return orgs, nil
}

// GetUserOrganizations gets all organizations a user belongs to (pointer receiver)
func (s *OrganizationService) GetUserOrganizations(ctx context.Context, userID string) (models.OrgList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orgIDs := make(map[string]bool)
	for _, m := range s.memberships {
		if m.UserID == userID {
			orgIDs[m.OrgID] = true
		}
	}

	orgs := make(models.OrgList, 0, len(orgIDs))
	for orgID := range orgIDs {
		if org, exists := s.orgs[orgID]; exists {
			orgs = append(orgs, *org)
		}
	}
	return orgs, nil
}

// =====================================
// Standalone Functions for Organization
// =====================================

// CreateOrganization is a standalone function that creates a new organization
func CreateOrganization(id, name, ownerID string) *models.Organization {
	return models.NewOrganization(id, name, ownerID)
}

// GenerateOrgID generates a unique organization ID (standalone function)
func GenerateOrgID() string {
	return fmt.Sprintf("org_%d", time.Now().UnixNano())
}

// GenerateMembershipID generates a unique membership ID (standalone function)
func GenerateMembershipID() string {
	return fmt.Sprintf("mem_%d", time.Now().UnixNano())
}

// CreateMembership is a standalone function that creates a new membership
func CreateMembership(userID, orgID string, role models.MemberRole) *models.Membership {
	return models.NewMembership(GenerateMembershipID(), userID, orgID, role)
}

