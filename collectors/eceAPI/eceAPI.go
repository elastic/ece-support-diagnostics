package eceAPI

import "github.com/elastic/ece-support-diagnostics/helpers"

func NewECEversionTask() helpers.Task {
	// platform is collected separately to determine the ECE version
	return helpers.Task{
		Filename: "ece/platform.json",
		Uri:      "/api/v1/platform",
	}
}

func NewECEtasks() *helpers.Tasks {
	return &helpers.Tasks{
		helpers.Task{
			Filename: "ece/allocators.json",
			Uri:      "/api/v1/platform/infrastructure/allocators",
		},
		helpers.Task{
			Filename: "ece/runners.json",
			Uri:      "/api/v1/platform/infrastructure/runners",
		},
		helpers.Task{
			Filename: "ece/proxies.json",
			Uri:      "/api/v1/platform/infrastructure/proxies",
		},
		helpers.Task{
			Filename: "ece/kibana_clusters.json",
			Uri:      "/api/v1/clusters/kibana",
		},
		helpers.Task{
			Filename: "ece/es_clusters.json",
			Uri:      "/api/v1/clusters/elasticsearch",
			Callback: ECEesClusters{},
		},
	}
}

func NewECEdeploymentTasks() *helpers.Tasks {
	// TaskConext will use the attached Metadata interface to template these variables
	return &helpers.Tasks{
		helpers.Task{
			Filename: "elasticsearch/{{ .cluster_id }}/_ece_cluster_info.json",
			Uri:      "/api/v1/clusters/elasticsearch/{{ .cluster_id }}",
		},
		helpers.Task{
			Filename: "elasticsearch/{{ .cluster_id }}/_ece_plan.json",
			Uri:      "/api/v1/clusters/elasticsearch/{{ .cluster_id }}/plan/activity",
		},
	}
}
