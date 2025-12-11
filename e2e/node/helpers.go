package node

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func convertLabelToSshKey(m interface{}, ssh_keys []interface{}, project_id string, location string) ([]interface{}, diag.Diagnostics) {
	cfg := m.(*config.Config)
	var goe2eClient *goe2e.Client
	goe2eClient, err := cfg.Goe2eClientForProject(project_id, location)
	if err != nil {
		return nil, diag.Errorf("Error creating goe2e client: %s", err)
	}

	log.Printf("[INFO] Helper Function ssh_keys = %+v", ssh_keys)
	if ssh_keys != nil || len(ssh_keys) > 0 {
		var new_SSH_keys []interface{}
		ctx := context.Background()
		for _, v := range ssh_keys {
			sshKey, _, err := goe2eClient.SSHKeys.GetSSHKeyByLabel(ctx, v.(string))
			log.Printf("[INFO] Helper Function sshKey = %+v", sshKey)
			if err != nil {
				return nil, diag.Errorf("Error fetching SSH key by label (%s): %s", v.(string), err)
			}
			new_SSH_keys = append(new_SSH_keys, sshKey.SSHKey)
		}
		return new_SSH_keys, nil
	}
	return nil, nil
}

// convertIDsToSshKeyContent converts ssh_key_ids (resource IDs) to actual SSH key content strings
// This is required because the API expects the actual key strings, not resource IDs
func convertIDsToSshKeyContent(m interface{}, ssh_key_ids []interface{}, project_id string, location string) ([]interface{}, diag.Diagnostics) {
	cfg := m.(*config.Config)
	var goe2eClient *goe2e.Client
	goe2eClient, err := cfg.Goe2eClientForProject(project_id, location)
	if err != nil {
		return nil, diag.Errorf("Error creating goe2e client: %s", err)
	}

	log.Printf("[INFO] Converting ssh_key_ids to key content: %+v", ssh_key_ids)
	if len(ssh_key_ids) == 0 {
		return nil, nil
	}

	var sshKeyContents []interface{}
	ctx := context.Background()
	for _, idInterface := range ssh_key_ids {
		keyID := idInterface.(string)

		// Fetch SSH key by ID
		sshKey, _, err := goe2eClient.SSHKeys.GetSSHKey(ctx, keyID)
		if err != nil {
			return nil, diag.Errorf("Error fetching SSH key (ID: %s) in project (%s), region (%s): %s", keyID, project_id, location, err)
		}

		// Extract the actual SSH key content
		if sshKey.SSHKey == "" {
			return nil, diag.Errorf("SSH key (ID: %s) has empty content", keyID)
		}

		log.Printf("[INFO] Fetched SSH key (ID: %s): %s...", keyID, sshKey.SSHKey[:min(50, len(sshKey.SSHKey))])
		sshKeyContents = append(sshKeyContents, sshKey.SSHKey)
	}

	return sshKeyContents, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func checkBlockStorage(m interface{}, image_id, project_id string, location string) diag.Diagnostics {
	cfg := m.(*config.Config)
	var goe2eClient *goe2e.Client
	goe2eClient, err := cfg.Goe2eClientForProject(project_id, location)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	ctx := context.Background()
	blockStorage, _, err := goe2eClient.BlockStorage.GetBlockStorage(ctx, image_id)
	if err != nil {
		return diag.Errorf("error finding Block Storage with ID %v: %s", image_id, err.Error())
	}

	if blockStorage.Status != goe2econstants.BlockStorageStatusAvailable {
		return diag.Errorf("Block Storage is in %s state, it must be in %s state", blockStorage.Status, goe2econstants.BlockStorageStatusAvailable)
	}
	return nil
}

func GetStartScripts(start_script string) []interface{} {
	startScripts := make([]interface{}, 0)
	if start_script != "" {
		startScripts = append(startScripts, start_script)
	}
	return startScripts
}

// ============================================
// V3 Helper Functions for Structured Blocks
// ============================================

// expandNetworkInterface parses the network_interface block from Terraform config
// Returns: vpc_id, assign_public_ip, enable_ipv6, security_group_ids
func expandNetworkInterface(niList []interface{}) (string, bool, bool, []int) {
	if len(niList) == 0 {
		return "", false, false, nil
	}

	ni := niList[0].(map[string]interface{})

	vpcID := ""
	if v, ok := ni[tfconstants.AttrVPCID]; ok {
		vpcID = v.(string)
	}

	assignPublicIP := false
	if v, ok := ni[tfconstants.AttrAssignPublicIP]; ok {
		assignPublicIP = v.(bool)
	}

	enableIPv6 := false
	if v, ok := ni[tfconstants.AttrEnableIPv6]; ok {
		enableIPv6 = v.(bool)
	}

	var securityGroupIDs []int
	if v, ok := ni[tfconstants.AttrSecurityGroupIDs]; ok && v != nil {
		sgList := v.([]interface{})
		securityGroupIDs = make([]int, len(sgList))
		for i, sg := range sgList {
			securityGroupIDs[i] = sg.(int)
		}
	}

	return vpcID, assignPublicIP, enableIPv6, securityGroupIDs
}

// flattenNetworkInterface converts API response to network_interface block structure
func flattenNetworkInterface(vpcID string, publicIP string, ipv6Address string, securityGroupIDs []int) []interface{} {
	ni := make(map[string]interface{})

	if vpcID != "" {
		ni[tfconstants.AttrVPCID] = vpcID
	}

	// Infer assign_public_ip from presence of public IP
	ni[tfconstants.AttrAssignPublicIP] = publicIP != ""

	// Infer enable_ipv6 from presence of IPv6 address
	ni[tfconstants.AttrEnableIPv6] = ipv6Address != ""

	if len(securityGroupIDs) > 0 {
		sgList := make([]interface{}, len(securityGroupIDs))
		for i, sg := range securityGroupIDs {
			sgList[i] = sg
		}
		ni[tfconstants.AttrSecurityGroupIDs] = sgList
	}

	return []interface{}{ni}
}

// expandRootDisk parses the root_disk block from Terraform config
// Returns: size_gb, disk_type
func expandRootDisk(rdList []interface{}) (int, string) {
	if len(rdList) == 0 {
		return 0, "standard"
	}

	rd := rdList[0].(map[string]interface{})

	sizeGB := 0
	if v, ok := rd[tfconstants.AttrSizeGB]; ok {
		sizeGB = v.(int)
	}

	diskType := "standard"
	if v, ok := rd[tfconstants.AttrDiskType]; ok && v != nil {
		diskType = v.(string)
	}

	return sizeGB, diskType
}

// flattenRootDisk converts API response to root_disk block structure
func flattenRootDisk(diskInfo string, diskType string) []interface{} {
	rd := make(map[string]interface{})

	// Parse disk info to extract size (format may vary: "100 GB", "100GB", etc.)
	// For now, store as computed field - actual size extraction needs API format
	// TODO: Parse actual size from disk info string when API format is known
	rd[tfconstants.AttrSizeGB] = 0 // Computed from API (default to 0 until parsed)

	rd[tfconstants.AttrDiskType] = diskType

	return []interface{}{rd}
}

// BackupConfig mirrors goe2e.BackupConfig for local use
type BackupConfig struct {
	PlanID               int
	Type                 string
	ExcludeFileFolder    string
	BackupNow            bool
	CompressionType      string
	CompressionLevel     string
	IsEncryptionRequired bool
	EncryptionPassphrase string
	HoursOfDay           string
	StartingMinute       int
	DBEnabled            bool
	DBUsername           string
	DBPassword           string
}

// BackupStatus mirrors goe2e.BackupStatus for local use
type BackupStatus struct {
	Status            string
	Detail            string
	LastRecoveryPoint string
}
