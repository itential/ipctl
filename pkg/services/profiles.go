// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package services provides client services for interacting with Itential Automation Platform (IAP)
// resources. This file contains the profile service implementation for managing IAP profiles.
//
// IAP profiles are configuration templates that define the behavior and settings for an IAP instance,
// including application properties, UI settings, authentication configuration, and various service
// properties like logging, Redis, and adapter configurations.
package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

// ApplicationProperties defines the application-level configuration properties for an IAP profile.
// These properties control basic application behavior and file system paths.
type ApplicationProperties struct {
	// Directory specifies the base directory path for the application, typically "./node_modules/"
	Directory string `json:"directory"`
}

// UIProperties defines the user interface configuration properties for an IAP profile.
// These properties control the web UI behavior, page layouts, and visual elements.
type UIProperties struct {
	// Description provides a human-readable description of the UI configuration
	Description    string `json:"description"`
	// Layout specifies the path to the main layout template file
	Layout         string `json:"layout"`
	// Home specifies the path to the home page HTML file
	Home           string `json:"home"`
	// Login specifies the path to the login page HTML file
	Login          string `json:"login"`
	// Profile specifies the path to the profile page HTML file
	Profile        string `json:"profile"`
	// UserConfig specifies the path to the user configuration template
	UserConfig     string `json:"user_config"`
	// GroupConfig specifies the path to the group configuration template
	GroupConfig    string `json:"group_config"`
	// NewUser specifies the path to the new user dialog template
	NewUser        string `json:"new_user"`
	// EditUser specifies the path to the edit user dialog template
	EditUser       string `json:"edit_user"`
	// NewGroup specifies the path to the new group dialog template
	NewGroup       string `json:"new_group"`
	// FavIcon specifies the path to the favicon image file
	FavIcon        string `json:"fav_icon"`
	// AppleTouchIcon specifies the path to the Apple touch icon image file
	AppleTouchIcon string `json:"apple_touch_icon"`
}

// ProfileMetadata contains metadata information about a profile's state and synchronization status.
type ProfileMetadata struct {
	// ActiveSync indicates whether the profile is actively synchronized
	ActiveSync bool `json:"activeSync"`
	// IsActive indicates whether this profile is currently the active profile
	IsActive   bool `json:"isActive"`
}

// Profile represents an IAP (Itential Automation Platform) profile configuration.
// A profile defines all the settings and properties needed to configure an IAP instance,
// including application behavior, UI settings, authentication, logging, and service configurations.
type Profile struct {
	// Id is the unique identifier for the profile
	Id              string                   `json:"id"`
	// Description provides a human-readable description of the profile's purpose
	Description     string                   `json:"description"`
	// LaunchDelay specifies the delay in milliseconds before launching services
	LaunchDelay     int                      `json:"launchDelay"`
	// LaunchTimeout specifies the timeout in milliseconds for service launch operations
	LaunchTimeout   int                      `json:"launchTimeout"`
	// Application contains application-level configuration properties
	Application     ApplicationProperties    `json:"applicationProps"`
	// UI contains user interface configuration properties
	UI              UIProperties             `json:"uiProps"`
	// Authentication contains authentication and security configuration
	Authentication  map[string]interface{}   `json:"authenticationProps,omitempty"`
	// Express contains Express.js web server configuration
	Express         map[string]interface{}   `json:"expressProps,omitempty"`
	// Logger contains logging configuration properties
	Logger          map[string]interface{}   `json:"loggerProps,omitempty"`
	// Alarm contains SNMP alarm and notification configuration
	Alarm           []map[string]interface{} `json:"alarmProps,omitempty"`
	// Audit contains audit logging and compliance configuration
	Audit           map[string]interface{}   `json:"auditProps,omitempty"`
	// Path contains file system path configuration
	Path            map[string]interface{}   `json:"pathProps,omitempty"`
	// TaskWorker contains background task worker configuration
	TaskWorker      map[string]interface{}   `json:"taskWorkerProps,omitempty"`
	// Redis contains Redis cache and session store configuration
	Redis           map[string]interface{}   `json:"redisProps,omitempty"`
	// Services lists the enabled services for this profile
	Services        []string                 `json:"services,omitempty"`
	// System contains system-level configuration properties
	System          map[string]interface{}   `json:"systemProps,omitempty"`
	// Integration contains external integration configuration
	Integration     map[string]interface{}   `json:"integrationProps,omitempty"`
	// Prebuilt contains pre-built automation configuration
	Prebuilt        map[string]interface{}   `json:"prebuiltProps,omitempty"`
	// RetryStrategy contains retry logic and timeout configuration
	RetryStrategy   map[string]interface{}   `json:"retryStrategy,omitempty"`
	// AdapterStrategy contains adapter management and routing configuration
	AdapterStrategy map[string]interface{}   `json:"adapterStrategy,omitempty"`
	// Adapter contains adapter-specific configuration properties
	Adapter         map[string]interface{}   `json:"adapterProps,omitempty"`
}

// ProfileService provides methods for managing IAP profiles through the REST API.
// It handles profile creation, retrieval, modification, import/export, and activation operations.
type ProfileService struct {
	// client is the underlying HTTP client used for API communication
	client *ServiceClient
}

// NewProfileService creates a new ProfileService instance with the provided client.
// The client is used for all HTTP communication with the IAP REST API.
//
// Parameters:
//   - c: The HTTP client to use for API requests
//
// Returns:
//   - *ProfileService: A new ProfileService instance ready for use
func NewProfileService(c client.Client) *ProfileService {
	return &ProfileService{client: NewServiceClient(c)}
}

// NewProfile creates a new Profile instance with the specified name and description.
// This is a convenience constructor that initializes a basic profile structure.
// Additional properties can be set on the returned Profile as needed.
//
// Parameters:
//   - name: The unique identifier for the profile
//   - desc: A human-readable description of the profile
//
// Returns:
//   - Profile: A new Profile instance with Id and Description set
func NewProfile(name, desc string) Profile {
	logger.Trace()

	return Profile{
		Id:          name,
		Description: desc,
	}
}

// GetAll retrieves all profiles from the IAP server.
// This method fetches the complete list of profiles available on the server,
// including both active and inactive profiles.
//
// Returns:
//   - []Profile: A slice containing all profiles found on the server
//   - error: An error if the request fails or the response cannot be parsed
func (svc *ProfileService) GetAll() ([]Profile, error) {
	logger.Trace()

	type Results struct {
		Metadata ProfileMetadata `json:"metadata"`
		Profile  Profile         `json:"profile"`
	}

	type Response struct {
		Results []Results `json:"results"`
		Total   int       `json:"total"`
	}

	var res Response

	if err := svc.client.Get("/profiles", &res); err != nil {
		return nil, err
	}

	var profiles []Profile
	for _, ele := range res.Results {
		profiles = append(profiles, ele.Profile)
	}

	return profiles, nil
}

// Get retrieves a specific profile by its unique identifier.
// This method fetches detailed information about a single profile from the IAP server.
//
// Parameters:
//   - id: The unique identifier of the profile to retrieve
//
// Returns:
//   - *Profile: A pointer to the requested profile, or nil if not found
//   - error: An error if the request fails, the profile doesn't exist, or the response cannot be parsed
func (svc *ProfileService) Get(id string) (*Profile, error) {
	logger.Trace()

	type Response struct {
		Metadata struct {
			IsActive   bool `json:"isActive"`
			ActiveSync bool `json:"activeSync"`
		} `json:"metadata"`
		Profile Profile `json:"profile"`
	}

	var res Response
	var uri = fmt.Sprintf("/profiles/%s", id)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return &res.Profile, nil
}

// GetActiveProfile retrieves the currently active profile from the IAP server.
// This method searches through all profiles to find the one that is currently active.
// An IAP instance can have only one active profile at a time.
//
// Returns:
//   - *Profile: A pointer to the active profile
//   - error: An error if the request fails, no active profile is found, or the response cannot be parsed
func (svc *ProfileService) GetActiveProfile() (*Profile, error) {
	logger.Trace()

	type Results struct {
		Metadata ProfileMetadata `json:"metadata"`
		Profile  Profile         `json:"profile"`
	}

	type Response struct {
		Results []Results `json:"results"`
		Total   int       `json:"total"`
	}

	var res Response

	if err := svc.client.Get("/profiles", &res); err != nil {
		return nil, err
	}

	var active *Profile

	for _, ele := range res.Results {
		if ele.Metadata.IsActive {
			active = &ele.Profile
			break
		}
	}

	if active == nil {
		return nil, errors.New("failed to find the active profile")
	}

	return active, nil
}

// Create creates a new profile on the IAP server.
// This method takes a Profile instance and creates it on the server.
// Only the Id and Description fields are used during creation; other properties
// are set to default values by the server.
//
// Parameters:
//   - in: The Profile instance to create (only Id and Description are used)
//
// Returns:
//   - *Profile: A pointer to the created profile with server-assigned properties
//   - error: An error if the request fails, the profile already exists, or the response cannot be parsed
func (svc *ProfileService) Create(in Profile) (*Profile, error) {
	logger.Trace()

	properties := map[string]interface{}{
		"id":          in.Id,
		"description": in.Description,
	}

	type Response struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    *Profile `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/profiles",
		body:               map[string]interface{}{"properties": properties},
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return res.Data, nil
}

// Delete removes a profile from the IAP server.
// This method permanently deletes the specified profile. This operation cannot be undone.
//
// Parameters:
//   - name: The unique identifier of the profile to delete
//
// Returns:
//   - error: An error if the request fails or the profile doesn't exist
func (svc *ProfileService) Delete(name string) error {
	logger.Trace()
	return svc.client.Delete(fmt.Sprintf("/profiles/%s", name))
}

// Import imports a complete profile configuration to the IAP server.
// This method takes a fully configured Profile instance and imports it,
// including all properties and configurations. Unlike Create, this method
// preserves all profile properties.
//
// Parameters:
//   - in: The complete Profile instance to import
//
// Returns:
//   - *Profile: A pointer to the imported profile as stored on the server
//   - error: An error if the request fails, the profile already exists, or the response cannot be parsed
func (svc *ProfileService) Import(in Profile) (*Profile, error) {
	logger.Trace()

	type Response struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    *Profile `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/profiles/import",
		body:               map[string]interface{}{"properties": in},
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// Export exports a complete profile configuration from the IAP server.
// This method retrieves all configuration properties for the specified profile,
// suitable for backup or migration purposes.
//
// Note: Currently uses the Get profile endpoint due to server-side issues
// with the dedicated export endpoint.
//
// Parameters:
//   - name: The unique identifier of the profile to export
//
// Returns:
//   - *Profile: A pointer to the complete profile configuration
//   - error: An error if the request fails, the profile doesn't exist, or the response cannot be parsed
func (svc *ProfileService) Export(name string) (*Profile, error) {
	logger.Trace()

	type Response struct {
		Metadata Metadata `json:"metadata"`
		Profile  *Profile `json:"profile"`
	}

	var res Response

	// XXX (privateip): The export URI returns an error from the server so
	// using the get profile route instead
	//var uri = fmt.Sprintf("/profiles/%s/export", name)
	var uri = fmt.Sprintf("/profiles/%s", name)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res.Profile, nil
}

// Activate sets the specified profile as the active profile on the IAP server.
// This method activates a profile, making it the current configuration used by the IAP instance.
// Only one profile can be active at a time; activating a profile will deactivate any previously active profile.
//
// Parameters:
//   - name: The unique identifier of the profile to activate
//
// Returns:
//   - error: An error if the request fails, the profile doesn't exist, or the activation fails
func (svc *ProfileService) Activate(name string) error {
	logger.Trace()

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	var res Response
	var uri = fmt.Sprintf("/profiles/%s/active", name)

	if err := svc.client.Put(uri, nil, &res); err != nil {
		return err
	}

	return nil
}
