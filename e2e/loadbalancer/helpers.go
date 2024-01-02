package loadbalancer

import (
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func GetLbPort(mode string) string {
	if mode == "HTTP" {
		return "80"
	}
	return "443"
}

func ExpandBackends(config []interface{}) ([]models.Backend, error) {
	backends := make([]models.Backend, 0, len(config))

	for _, backend := range config {
		detail := backend.(map[string]interface{})

		servers, err := ExpandServers(detail["servers"].([]interface{}))
		if err != nil {
			return nil, err
		}
		r := models.Backend{
			Balance:        detail["balance"].(string),
			CheckboxEnable: detail["checkbox_enable"].(bool),
			DomainName:     detail["domain_name"].(string),
			CheckUrl:       detail["check_url"].(string),
			Servers:        servers,
			HttpCheck:      detail["http_check"].(bool),
		}
		backends = append(backends, r)
	}
	return backends, nil
}

func ExpandServers(config []interface{}) ([]models.Server, error) {
	servers := make([]models.Server, 0, len(config))

	for _, server := range config {
		detail := server.(map[string]interface{})

		r := models.Server{
			BackendName: detail["backend_name"].(string),
			BackendIp:   detail["backend_ip"].(string),
			BackendPort: detail["backend_port"].(string),
		}

		servers = append(servers, r)
	}
	return servers, nil
}

func ExpandAclList(config []interface{}) ([]models.AclListInfo, error) {
	aclList := make([]models.AclListInfo, 0, len(config))
	return aclList, nil
}

func ExpandAclMap(config []interface{}) ([]models.AclMapInfo, error) {
	aclMap := make([]models.AclMapInfo, 0, len(config))
	return aclMap, nil
}

func ExpandVpcList(config []interface{}) ([]models.VpcDetail, error) {
	vpcList := make([]models.VpcDetail, 0, len(config))

	for _, vpc := range config {
		detail := vpc.(map[string]interface{})

		r := models.VpcDetail{
			Network_id: detail["network_id"].(float64),
			VpcName:    detail["vpc_name"].(string),
			Ipv4_cidr:  detail["ipv4_cidr"].(string),
		}

		vpcList = append(vpcList, r)
	}
	return vpcList, nil
}

func ExpandEnableEosLogger(config []interface{}) (models.EosDetail, error) {
	eosDetail := models.EosDetail{}

	for _, eosBucketInfo := range config {
		detail := eosBucketInfo.(map[string]interface{})
		eosDetail.ApplianceId = detail["appliance_id"].(int)
		eosDetail.AccessKey = detail["access_key"].(string)
		eosDetail.Secretkey = detail["secret_key"].(string)
		eosDetail.Bucket = detail["bucket"].(string)
	}
	return eosDetail, nil
}

func ExpandTcpBackend(config []interface{}) ([]models.TcpBackendDetail, error) {
	tcpBackends := make([]models.TcpBackendDetail, 0, len(config))

	for _, tcpBackend := range config {
		detail := tcpBackend.(map[string]interface{})

		servers, err := ExpandServers(detail["servers"].([]interface{}))
		if err != nil {
			return nil, err
		}
		r := models.TcpBackendDetail{
			BackendName: detail["backend_name"].(string),
			Port:        detail["port"].(string),
			Balance:     detail["balance"].(string),
			Servers:     servers,
		}

		tcpBackends = append(tcpBackends, r)
	}
	return tcpBackends, nil
}
