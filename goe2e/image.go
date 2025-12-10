package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	imagesPath                 = "images"                          // /images
	imagesOSCategoryPath       = "images/os-category"              // /images/os-category
	imagesImportPath           = "images/import-image"             // /images/import-image
	imagesWindowPermissionPath = "images/window-image-permissions" // /images/window-image-permissions
	imagesSavedPath            = "images/saved-images"             // /images/saved-images
	imagesUpgradeImagePath     = "images/upgradeimage"             // /images/upgradeimage/{templateID}/
)

// ImageService is an interface for interacting with image endpoints
// of the E2E Networks API.
type ImageService interface {
	// OS category and plan operations
	GetOSCategories(context.Context) (*OSCategoryResponse, *Response, error)
	GetImagePlans(context.Context, *ImagePlansRequest) ([]ImagePlan, *Response, error)

	// Custom image operations
	ImportImage(context.Context, *ImportImageRequest) (*ImageImportResult, *Response, error)
	GetWindowsImagePermission(context.Context) (*WindowsPermission, *Response, error)

	// Saved image operations
	GetSavedImages(context.Context) ([]SavedImage, *Response, error)
	GetImage(context.Context, string) (*SavedImage, *Response, error)
	RenameImage(context.Context, string, *RenameImageRequest) (*RenameImageResult, *Response, error)
	DeleteImage(context.Context, string) (*DeleteImageResult, *Response, error)

	// Plan details operations
	GetPlanDetailsFromPlanName(context.Context, int, string) (string, string, *Response, error)
}

// ImageServiceOp handles communication with image related methods of the E2E Networks API.
type ImageServiceOp struct {
	client *Client
}

var _ ImageService = &ImageServiceOp{}

// OSCategory represents an OS category with versions and subcategories
type OSCategory struct {
	OS       string      `json:"OS"`
	Version  []OSVersion `json:"version"`
	Category []string    `json:"category"`
}

// OSVersion represents a version of an operating system
type OSVersion struct {
	NumberOfDomains string `json:"number_of_domains,omitempty"`
	OS              string `json:"os"`
	Version         string `json:"version"`
	SubCategory     string `json:"sub_category"`
	SoftwareVersion string `json:"software_version"`
}

// OSCategoryResponse represents the response from GetOSCategories
type OSCategoryResponse struct {
	CategoryList []OSCategory `json:"category_list"`
}

// ImagePlan represents a plan for a particular OS
type ImagePlan struct {
	Name                     string                 `json:"name"`
	Plan                     string                 `json:"plan"`
	Image                    string                 `json:"image"`
	OS                       OSInfo                 `json:"os"`
	Location                 string                 `json:"location"`
	Specs                    PlanSpecs              `json:"specs"`
	CPUType                  string                 `json:"cpu_type"`
	GPUCardDetails           map[string]interface{} `json:"gpu_card_details,omitempty"`
	NodeDescription          string                 `json:"node_description"`
	InstalledApplications    map[string]interface{} `json:"installed_application_version,omitempty"`
	CanSupportBitNinja       BitNinjaInfo           `json:"can_support_bitninja"`
	BitNinjaDiscountPercent  float64                `json:"bitninja_discount_percentage"`
	AvailableInventoryStatus bool                   `json:"available_inventory_status"`
	Currency                 string                 `json:"currency"`
	IsBlockStorageAttachable bool                   `json:"is_blockstorage_attachable"`
}

// OSInfo represents OS information
type OSInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Image    string `json:"image"`
	Category string `json:"category"`
}

// PlanSpecs represents the specifications of a plan
type PlanSpecs struct {
	ID                   string             `json:"id"`
	SKUName              string             `json:"sku_name"`
	RAM                  string             `json:"ram"`
	CPU                  int                `json:"cpu"`
	DiskSpace            int                `json:"disk_space"`
	PricePerMonth        float64            `json:"price_per_month"`
	PricePerHour         float64            `json:"price_per_hour"`
	Series               string             `json:"series"`
	MinimumBillingAmount float64            `json:"minimum_billing_amount"`
	CommittedSKU         []CommittedSKUInfo `json:"committed_sku,omitempty"`
	Family               string             `json:"family"`
}

// CommittedSKUInfo represents committed SKU information
type CommittedSKUInfo struct {
	CommittedSKUID    int     `json:"committed_sku_id"`
	CommittedSKUName  string  `json:"committed_sku_name"`
	CommittedNodeMsg  string  `json:"committed_node_message"`
	CommittedSKUPrice float64 `json:"committed_sku_price"`
	CommittedUptoDate string  `json:"committed_upto_date"`
	CommittedDays     int     `json:"committed_days"`
}

// BitNinjaInfo represents BitNinja support information
type BitNinjaInfo struct {
	ShowBitNinja bool    `json:"show_bitninja"`
	BitNinjaCost float64 `json:"bitninja_cost"`
}

// ImportImageRequest represents a request to import a custom image
type ImportImageRequest struct {
	ImageName string `json:"image_name"`
	PublicURL string `json:"public_url"`
	Location  string `json:"location"`
	OS        string `json:"os"`
}

// ImageImportResult represents the result of importing an image
type ImageImportResult struct {
	// Empty object typically
}

// WindowsPermission represents Windows image permission response
type WindowsPermission struct {
	IsWindowsAllowed bool `json:"is_windows_allowed"`
}

// SavedImage represents a saved/custom image
type SavedImage struct {
	TemplateID     int           `json:"template_id"`
	IsWindows      bool          `json:"is_windows"`
	VMInfo         []interface{} `json:"vm_info"`
	ImageType      string        `json:"image_type"`
	OSDistribution string        `json:"os_distribution"`
	Name           string        `json:"name"`
	ImageID        string        `json:"image_id"`
	Distro         string        `json:"distro"`
	SKUType        string        `json:"sku_type"`
	ImageState     string        `json:"image_state"`
	RunningVMs     string        `json:"running_vms"`
	CloningOps     string        `json:"cloning_ops"`
	ImageSize      string        `json:"image_size"`
	CreationTime   string        `json:"creation_time"`
}

// RenameImageRequest represents a request to rename an image
type RenameImageRequest struct {
	ActionType string `json:"action_type"`
	Name       string `json:"name"`
}

// RenameImageResult represents the result of renaming an image
type RenameImageResult struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// DeleteImageResult represents the result of deleting an image
type DeleteImageResult struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// ImagePlansRequest represents query parameters for getting image plans
type ImagePlansRequest struct {
	Category        string
	OS              string
	OSVersion       string
	DisplayCategory string
}

// Response wrappers for API calls
type imagePlanRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []ImagePlan `json:"data"`
}

type osListRoot struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    OSCategoryResponse `json:"data"`
}

type savedImageRoot struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    []SavedImage `json:"data"`
}

type imageActionRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type windowsPermissionRoot struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    WindowsPermission `json:"data"`
}

// GetOSCategories retrieves the list of available OS categories and versions.
func (s *ImageServiceOp) GetOSCategories(ctx context.Context) (*OSCategoryResponse, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, imagesOSCategoryPath, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for OS categories: %w", err)
	}

	root := new(osListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve OS categories: %w", err)
	}

	return &root.Data, resp, nil
}

// GetImagePlans retrieves available image plans based on the specified criteria.
func (s *ImageServiceOp) GetImagePlans(ctx context.Context, planReq *ImagePlansRequest) ([]ImagePlan, *Response, error) {
	if planReq == nil {
		return nil, nil, NewArgError("planReq", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, imagesPath, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for image plans: %w", err)
	}

	// Add optional query parameters
	q := req.URL.Query()
	if planReq.Category != "" {
		q.Add("category", planReq.Category)
	}
	if planReq.OS != "" {
		q.Add("os", planReq.OS)
	}
	if planReq.OSVersion != "" {
		q.Add("osversion", planReq.OSVersion)
	}
	if planReq.DisplayCategory != "" {
		q.Add("display_category", planReq.DisplayCategory)
	}
	req.URL.RawQuery = q.Encode()

	root := new(imagePlanRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve image plans: %w", err)
	}

	return root.Data, resp, nil
}

// ImportImage imports a custom image from a public URL.
func (s *ImageServiceOp) ImportImage(ctx context.Context, importReq *ImportImageRequest) (*ImageImportResult, *Response, error) {
	if importReq == nil {
		return nil, nil, NewArgError("importReq", "cannot be nil")
	}
	if importReq.ImageName == "" {
		return nil, nil, NewArgError("ImageName", "cannot be empty")
	}
	if importReq.PublicURL == "" {
		return nil, nil, NewArgError("PublicURL", "cannot be empty")
	}
	if importReq.Location == "" {
		return nil, nil, NewArgError("Location", "cannot be empty")
	}
	if importReq.OS == "" {
		return nil, nil, NewArgError("OS", "cannot be empty")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, imagesImportPath, importReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for importing image (%s): %w", importReq.ImageName, err)
	}

	root := new(imageActionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to import image (%s): %w", importReq.ImageName, err)
	}

	return &ImageImportResult{}, resp, nil
}

// GetWindowsImagePermission checks if the user is allowed to import Windows images.
func (s *ImageServiceOp) GetWindowsImagePermission(ctx context.Context) (*WindowsPermission, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, imagesWindowPermissionPath, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for Windows image permission: %w", err)
	}

	root := new(windowsPermissionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve Windows image permission: %w", err)
	}

	return &root.Data, resp, nil
}

// GetSavedImages retrieves the list of saved/custom images.
func (s *ImageServiceOp) GetSavedImages(ctx context.Context) ([]SavedImage, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, imagesSavedPath, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for saved images: %w", err)
	}

	root := new(savedImageRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve saved images: %w", err)
	}

	return root.Data, resp, nil
}

// GetImage retrieves a single saved image by ID.
func (s *ImageServiceOp) GetImage(ctx context.Context, imageID string) (*SavedImage, *Response, error) {
	if imageID == "" {
		return nil, nil, NewArgError("imageID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", imagesPath, imageID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for getting image (ID: %s): %w", imageID, err)
	}

	type imageRoot struct {
		Code    int        `json:"code"`
		Message string     `json:"message"`
		Data    SavedImage `json:"data"`
	}

	root := new(imageRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get image (ID: %s): %w", imageID, err)
	}

	return &root.Data, resp, nil
}

// RenameImage renames an existing image.
func (s *ImageServiceOp) RenameImage(ctx context.Context, imageID string, renameReq *RenameImageRequest) (*RenameImageResult, *Response, error) {
	if imageID == "" {
		return nil, nil, NewArgError("imageID", "cannot be empty")
	}
	if renameReq == nil {
		return nil, nil, NewArgError("renameReq", "cannot be nil")
	}
	if renameReq.ActionType == "" {
		return nil, nil, NewArgError("ActionType", "cannot be empty")
	}
	if renameReq.Name == "" {
		return nil, nil, NewArgError("Name", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", imagesPath, imageID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, renameReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for renaming image (ID: %s): %w", imageID, err)
	}

	root := new(imageActionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to rename image (ID: %s): %w", imageID, err)
	}

	// Extract the result from the response data
	result := &RenameImageResult{}
	if data, ok := root.Data.(map[string]interface{}); ok {
		if status, ok := data["status"].(bool); ok {
			result.Status = status
		}
		if message, ok := data["message"].(string); ok {
			result.Message = message
		}
	}

	return result, resp, nil
}

// DeleteImage deletes an existing image.
func (s *ImageServiceOp) DeleteImage(ctx context.Context, imageID string) (*DeleteImageResult, *Response, error) {
	if imageID == "" {
		return nil, nil, NewArgError("imageID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", imagesPath, imageID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for deleting image (ID: %s): %w", imageID, err)
	}

	root := new(imageActionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to delete image (ID: %s): %w", imageID, err)
	}

	// Extract the result from the response data
	result := &DeleteImageResult{}
	if data, ok := root.Data.(map[string]interface{}); ok {
		if status, ok := data["status"].(bool); ok {
			result.Status = status
		}
		if message, ok := data["message"].(string); ok {
			result.Message = message
		}
	}

	return result, resp, nil
}

// PlanDetail represents plan details from upgradeimage endpoint
type PlanDetail struct {
	Name  string `json:"name"`
	Plan  string `json:"plan"`
	Specs struct {
		ID string `json:"id"`
	} `json:"specs"`
}

// upgradeImageRoot represents the response from upgradeimage endpoint
type upgradeImageRoot struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Data []PlanDetail `json:"data"`
	} `json:"data"`
}

// GetPlanDetailsFromPlanName retrieves plan details (planID and slugName) by plan name for a given template ID.
// This uses the images/upgradeimage/{templateID}/ endpoint.
func (s *ImageServiceOp) GetPlanDetailsFromPlanName(ctx context.Context, templateID int, planName string) (string, string, *Response, error) {
	if templateID <= 0 {
		return "", "", nil, NewArgError("templateID", "must be greater than 0")
	}
	if planName == "" {
		return "", "", nil, NewArgError("planName", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%d/", imagesUpgradeImagePath, templateID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create request for plan details: %w", err)
	}

	root := new(upgradeImageRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return "", "", resp, fmt.Errorf("failed to retrieve plan details: %w", err)
	}

	for _, item := range root.Data.Data {
		if item.Name == planName {
			return item.Specs.ID, item.Plan, resp, nil
		}
	}

	return "", "", resp, fmt.Errorf("plan name %s not found in template %d", planName, templateID)
}
