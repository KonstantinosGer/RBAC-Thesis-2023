package models

type FrontendPolicy map[string][]PermissionInfo

type PermissionInfo struct {
	PermissionId          int    `json:"id"`
	PermissionDescription string `json:"description"`
	PermissionAction      string `json:"action"`
	PermissionResource    string `json:"resource"`
	HasPermission         bool   `json:"has_permission"`
}

//to json object:
// { "category1": [ ["perm descr", "perm id", "true/false"], ... ],
//   "category2": [ ["perm descr", "perm id", "true/false"], ... ] }
