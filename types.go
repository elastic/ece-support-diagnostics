package ecediag

type v0APIresponse struct {
	Ok           bool   `json:"ok"`
	Message      string `json:"message"`
	EulaAccepted bool   `json:"eula_accepted"`
	Hrefs        struct {
		Regions       string `json:"regions"`
		Elasticsearch string `json:"elasticsearch"`
		Logs          string `json:"logs"`
		DatabaseUsers string `json:"database/users"`
	} `json:"hrefs"`
}

// type Clusters struct {
// 	ReturnCount           int `json:"return_count"`
// 	ElasticsearchClusters []struct {
// 		ClusterName           string        `json:"cluster_name"`
// 		AssociatedApmClusters []interface{} `json:"associated_apm_clusters"`
// 		PlanInfo              struct {
// 			Healthy bool          `json:"healthy"`
// 			History []interface{} `json:"history"`
// 		} `json:"plan_info"`
// 		Snapshots struct {
// 			Healthy       bool `json:"healthy"`
// 			Count         int  `json:"count"`
// 			RecentSuccess bool `json:"recent_success"`
// 		} `json:"snapshots"`
// 		AssociatedKibanaClusters []struct {
// 			KibanaID string `json:"kibana_id"`
// 			Enabled  bool   `json:"enabled"`
// 			Links    struct {
// 			} `json:"links"`
// 		} `json:"associated_kibana_clusters"`
// 		Elasticsearch struct {
// 			Healthy   bool `json:"healthy"`
// 			ShardInfo struct {
// 				Healthy         bool `json:"healthy"`
// 				AvailableShards []struct {
// 					InstanceName string `json:"instance_name"`
// 					ShardCount   int    `json:"shard_count"`
// 				} `json:"available_shards"`
// 				UnavailableShards []struct {
// 					InstanceName string `json:"instance_name"`
// 					ShardCount   int    `json:"shard_count"`
// 				} `json:"unavailable_shards"`
// 				UnavailableReplicas []struct {
// 					InstanceName string `json:"instance_name"`
// 					ReplicaCount int    `json:"replica_count"`
// 				} `json:"unavailable_replicas"`
// 			} `json:"shard_info"`
// 			MasterInfo struct {
// 				Healthy bool `json:"healthy"`
// 				Masters []struct {
// 					MasterNodeID       string   `json:"master_node_id"`
// 					MasterInstanceName string   `json:"master_instance_name"`
// 					Instances          []string `json:"instances"`
// 				} `json:"masters"`
// 				InstancesWithNoMaster []interface{} `json:"instances_with_no_master"`
// 			} `json:"master_info"`
// 			BlockingIssues struct {
// 				Healthy      bool          `json:"healthy"`
// 				ClusterLevel []interface{} `json:"cluster_level"`
// 				IndexLevel   []interface{} `json:"index_level"`
// 			} `json:"blocking_issues"`
// 		} `json:"elasticsearch"`
// 		Links struct {
// 		} `json:"links"`
// 		Healthy  bool   `json:"healthy"`
// 		Status   string `json:"status"`
// 		Topology struct {
// 			Healthy   bool `json:"healthy"`
// 			Instances []struct {
// 				Disk struct {
// 					DiskSpaceAvailable int64 `json:"disk_space_available"`
// 					DiskSpaceUsed      int64 `json:"disk_space_used"`
// 				} `json:"disk"`
// 				MaintenanceMode       bool     `json:"maintenance_mode"`
// 				ServiceRunning        bool     `json:"service_running"`
// 				Healthy               bool     `json:"healthy"`
// 				InstanceName          string   `json:"instance_name"`
// 				ServiceVersion        string   `json:"service_version"`
// 				ServiceRoles          []string `json:"service_roles"`
// 				AllocatorID           string   `json:"allocator_id"`
// 				ServiceID             string   `json:"service_id"`
// 				Zone                  string   `json:"zone"`
// 				InstanceConfiguration struct {
// 					ID       string `json:"id"`
// 					Name     string `json:"name"`
// 					Resource string `json:"resource"`
// 				} `json:"instance_configuration"`
// 				ContainerStarted bool `json:"container_started"`
// 				Memory           struct {
// 					InstanceCapacity int64 `json:"instance_capacity"`
// 				} `json:"memory"`
// 			} `json:"instances"`
// 		} `json:"topology"`
// 		Metadata struct {
// 			Ports struct {
// 				HTTP  int `json:"http"`
// 				HTTPS int `json:"https"`
// 			} `json:"ports"`
// 			LastModified string `json:"last_modified"`
// 			Version      int    `json:"version"`
// 			Endpoint     string `json:"endpoint"`
// 			CloudID      string `json:"cloud_id"`
// 		} `json:"metadata"`
// 		ExternalLinks []interface{} `json:"external_links"`
// 		ClusterID     string        `json:"cluster_id"`
// 	} `json:"elasticsearch_clusters"`
// }
