package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/test-repo-golang-support/models"
	"github.com/test-repo-golang-support/services"
)

// OrgHandler wraps the organization service and provides HTTP handlers
type OrgHandler struct {
	service *services.OrganizationService
	logger  *log.Logger
}

// NewOrgHandler creates a new OrgHandler instance
func NewOrgHandler(service *services.OrganizationService, logger *log.Logger) *OrgHandler {
	return &OrgHandler{
		service: service,
		logger:  logger,
	}
}

// =====================================
// Organization HTTP Handlers
// =====================================

// GetOrganizations handles GET /organizations - returns all organizations
func (h *OrgHandler) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgs, err := h.service.ReadAllOrgs(ctx)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch organizations")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Organizations retrieved successfully",
		Data:    orgs,
	})
}

// GetOrganization handles GET /organizations/{id} - returns a specific organization
func (h *OrgHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	org, err := h.service.ReadOrg(ctx, id)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Organization not found")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Organization retrieved successfully",
		Data:    org,
	})
}

// CreateOrganization handles POST /organizations - creates a new organization
func (h *OrgHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Industry    string `json:"industry"`
		OwnerID     string `json:"owner_id"`
		Address     struct {
			Street     string `json:"street"`
			City       string `json:"city"`
			State      string `json:"state"`
			Country    string `json:"country"`
			PostalCode string `json:"postal_code"`
		} `json:"address"`
		Contact struct {
			Phone   string `json:"phone"`
			Email   string `json:"email"`
			Website string `json:"website"`
		} `json:"contact"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Name == "" {
		h.respondError(w, http.StatusBadRequest, "Organization name is required")
		return
	}

	if input.OwnerID == "" {
		h.respondError(w, http.StatusBadRequest, "Owner ID is required")
		return
	}

	// Create new organization
	orgID := services.GenerateOrgID()
	org := services.CreateOrganization(orgID, input.Name, input.OwnerID)
	org.UpdateDescription(input.Description)
	org.SetIndustry(input.Industry)

	// Set address if provided
	if input.Address.City != "" {
		org.UpdateAddress(models.Address{
			Street:     input.Address.Street,
			City:       input.Address.City,
			State:      input.Address.State,
			Country:    input.Address.Country,
			PostalCode: input.Address.PostalCode,
		})
	}

	// Set contact info if provided
	if input.Contact.Email != "" {
		org.UpdateContact(models.ContactInfo{
			Phone:   input.Contact.Phone,
			Email:   input.Contact.Email,
			Website: input.Contact.Website,
		})
	}

	// Save organization
	if err := h.service.WriteOrg(ctx, org); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create organization")
		return
	}

	// Create owner membership
	membership := services.CreateMembership(input.OwnerID, orgID, models.MemberRoleOwner)
	if err := h.service.AddMember(ctx, membership); err != nil {
		h.logger.Printf("Warning: Failed to create owner membership: %v", err)
	}

	h.logger.Printf("Created organization: %s (%s)", org.DisplayName(), org.ID)

	h.respondJSON(w, http.StatusCreated, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Organization created successfully",
		Data:    org,
	})
}

// UpdateOrganization handles PUT /organizations/{id} - updates an existing organization
func (h *OrgHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing organization
	org, err := h.service.ReadOrg(ctx, id)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Organization not found")
		return
	}

	var input struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Industry    string          `json:"industry"`
		Size        models.OrgSize  `json:"size"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update organization using pointer receiver methods
	if input.Name != "" {
		org.UpdateName(input.Name)
	}
	if input.Description != "" {
		org.UpdateDescription(input.Description)
	}
	if input.Industry != "" {
		org.SetIndustry(input.Industry)
	}
	if input.Size != "" {
		org.SetSize(input.Size)
	}

	// Save updated organization
	if err := h.service.WriteOrg(ctx, org); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to update organization")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Organization updated successfully",
		Data:    org,
	})
}

// DeleteOrganization handles DELETE /organizations/{id} - deletes an organization
func (h *OrgHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if organization exists
	exists, _ := h.service.OrgExists(ctx, id)
	if !exists {
		h.respondError(w, http.StatusNotFound, "Organization not found")
		return
	}

	if err := h.service.DeleteOrg(ctx, id); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete organization")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Organization deleted successfully",
	})
}

// =====================================
// Membership HTTP Handlers
// =====================================

// GetOrgMembers handles GET /organizations/{id}/members - returns all members
func (h *OrgHandler) GetOrgMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars["id"]

	members, err := h.service.GetMembers(ctx, orgID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch members")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Members retrieved successfully",
		Data:    members,
	})
}

// AddOrgMember handles POST /organizations/{id}/members - adds a new member
func (h *OrgHandler) AddOrgMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars["id"]

	var input struct {
		UserID string            `json:"user_id"`
		Role   models.MemberRole `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.UserID == "" {
		h.respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	if input.Role == "" {
		input.Role = models.MemberRoleMember
	}

	membership := services.CreateMembership(input.UserID, orgID, input.Role)
	if err := h.service.AddMember(ctx, membership); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Member added successfully",
		Data:    membership,
	})
}

// RemoveOrgMember handles DELETE /organizations/{id}/members/{userId} - removes a member
func (h *OrgHandler) RemoveOrgMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars["id"]
	userID := vars["userId"]

	if err := h.service.RemoveMember(ctx, userID, orgID); err != nil {
		h.respondError(w, http.StatusNotFound, "Membership not found")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Member removed successfully",
	})
}

// UpdateMemberRole handles PUT /organizations/{id}/members/{userId} - updates member role
func (h *OrgHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars["id"]
	userID := vars["userId"]

	var input struct {
		Role models.MemberRole `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateMemberRole(ctx, userID, orgID, input.Role); err != nil {
		h.respondError(w, http.StatusNotFound, "Membership not found")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Member role updated successfully",
	})
}

// GetUserOrganizations handles GET /users/{id}/organizations - returns user's organizations
func (h *OrgHandler) GetUserOrganizations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := vars["id"]

	orgs, err := h.service.GetUserOrganizations(ctx, userID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch organizations")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "User organizations retrieved successfully",
		Data:    orgs,
	})
}

// =====================================
// Helper Methods
// =====================================

// respondJSON sends a JSON response (pointer receiver)
func (h *OrgHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("Error encoding response: %v", err)
	}
}

// respondError sends an error response (pointer receiver)
func (h *OrgHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, models.APIResponse{
		Code:    models.ResponseError,
		Message: message,
	})
}

// =====================================
// Route Setup for Organizations
// =====================================

// SetupOrgRoutes configures organization routes
func SetupOrgRoutes(router *mux.Router, h *OrgHandler) {
	// Organization routes
	router.HandleFunc("/organizations", h.GetOrganizations).Methods("GET")
	router.HandleFunc("/organizations/{id}", h.GetOrganization).Methods("GET")
	router.HandleFunc("/organizations", h.CreateOrganization).Methods("POST")
	router.HandleFunc("/organizations/{id}", h.UpdateOrganization).Methods("PUT")
	router.HandleFunc("/organizations/{id}", h.DeleteOrganization).Methods("DELETE")

	// Membership routes
	router.HandleFunc("/organizations/{id}/members", h.GetOrgMembers).Methods("GET")
	router.HandleFunc("/organizations/{id}/members", h.AddOrgMember).Methods("POST")
	router.HandleFunc("/organizations/{id}/members/{userId}", h.RemoveOrgMember).Methods("DELETE")
	router.HandleFunc("/organizations/{id}/members/{userId}", h.UpdateMemberRole).Methods("PUT")

	// User organizations route
	router.HandleFunc("/users/{id}/organizations", h.GetUserOrganizations).Methods("GET")
}

