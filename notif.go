package main

// CloudBuildNotification represents message from Cloud Build
type CloudBuildNotification struct {
	Status        string           `json:"status,omitempty"`
	LogURL        string           `json:"logUrl,omitempty"`
	Source        CloudBuildSource `json:"source,omitempty"`
	Substitutions CloudBuildConfig `json:"substitutions,omitempty"`
}

// CloudBuildSource represents Cloud Build source portion of notification
type CloudBuildSource struct {
	RepoSource RepoSource `json:"repoSource,omitempty"`
}

// RepoSource represents Cloud Build repo source portion of notification
type RepoSource struct {
	RepoName string `json:"repoName,omitempty"`
	TagName  string `json:"tagName,omitempty"`
}

// CloudBuildConfig represents Cloud Build config
type CloudBuildConfig struct {
	ClusterName string `json:"_CLUSTER_NAME,omitempty"`
	AppName     string `json:"_APP_NAME,omitempty"`
}
