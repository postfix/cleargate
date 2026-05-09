package models

import "time"

type ToolSpec struct {
	APIVersion string         `yaml:"apiVersion" validate:"required"`
	Kind       string         `yaml:"kind" validate:"required"`
	Metadata   Metadata       `yaml:"metadata" validate:"required"`
	Runtime    Runtime        `yaml:"runtime" validate:"required"`
	Sandbox    Sandbox        `yaml:"sandbox"`
	Inputs     []Input        `yaml:"inputs,omitempty" validate:"dive"`
	Flags      []Flag         `yaml:"flags,omitempty" validate:"dive"`
	Positionals []Positional   `yaml:"positionals,omitempty" validate:"dive"`
	Outputs    []Output       `yaml:"outputs,omitempty" validate:"dive"`
	Presets    []Preset       `yaml:"presets,omitempty"`
	Security   SecurityPolicy `yaml:"security,omitempty"`
	Groups     []Group        `yaml:"groups,omitempty"`
}

type Metadata struct {
	Name          string   `yaml:"name" validate:"required"`
	DisplayName   string   `yaml:"displayName"`
	Description   string   `yaml:"description"`
	Version       string   `yaml:"version" validate:"required"`
	Owner         string   `yaml:"owner"`
	Tags          []string `yaml:"tags,omitempty"`
	Homepage      string   `yaml:"homepage,omitempty"`
	Documentation []string `yaml:"documentation,omitempty"`
}

type Runtime struct {
	Executable       string            `yaml:"executable" validate:"required"`
	Argv0            string            `yaml:"argv0"`
	ContainerImage   string            `yaml:"containerImage,omitempty"`
	WorkingDirectory string            `yaml:"workingDirectory,omitempty"`
	TimeoutSeconds   int               `yaml:"timeoutSeconds,omitempty" validate:"min=1"`
	MaxMemoryMB      int               `yaml:"maxMemoryMB,omitempty"`
	MaxCPUPercent    int               `yaml:"maxCPUPercent,omitempty"`
	MaxStdoutBytes   int               `yaml:"maxStdoutBytes,omitempty"`
	MaxStderrBytes   int               `yaml:"maxStderrBytes,omitempty"`
	Environment      EnvironmentPolicy `yaml:"environment,omitempty"`
	Network          NetworkPolicy     `yaml:"network,omitempty"`
}

type EnvironmentPolicy struct {
	Allowlist []string          `yaml:"allowlist,omitempty"`
	Fixed     map[string]string `yaml:"fixed,omitempty"`
}

type NetworkPolicy struct {
	Enabled      bool     `yaml:"enabled"`
	AllowedCIDRs []string `yaml:"allowedCIDRs,omitempty"`
}

type Sandbox struct {
	Profile                string `yaml:"profile,omitempty"`
	User                   string `yaml:"user,omitempty"`
	ReadonlyRootFilesystem bool   `yaml:"readonlyRootFilesystem,omitempty"`
	SeccompProfile         string `yaml:"seccompProfile,omitempty"`
	MaxProcesses           int    `yaml:"maxProcesses,omitempty"`
}

type Input struct {
	ID                string   `yaml:"id" validate:"required"`
	Type              string   `yaml:"type" validate:"required"`
	Required          bool     `yaml:"required"`
	Destination       string   `yaml:"destination"`
	MaxSizeMB         int      `yaml:"maxSizeMB,omitempty"`
	AllowedExtensions []string `yaml:"allowedExtensions,omitempty"`
	Unpack            *Unpack  `yaml:"unpack,omitempty"`
}

type Unpack struct {
	Enabled           bool `yaml:"enabled"`
	MaxFiles          int  `yaml:"maxFiles,omitempty"`
	MaxExpandedSizeMB int  `yaml:"maxExpandedSizeMB,omitempty"`
	DenySymlinks      bool `yaml:"denySymlinks,omitempty"`
}

type Flag struct {
	ID         string      `yaml:"id" validate:"required"`
	Type       string      `yaml:"type" validate:"required"`
	FlagString string      `yaml:"flag,omitempty"`
	Required   bool        `yaml:"required,omitempty"`
	Default    interface{} `yaml:"default,omitempty"`
	Values     []string    `yaml:"values,omitempty"`
	Cli      CliRender   `yaml:"cli,omitempty"`
	Ui       UiRender    `yaml:"ui,omitempty"`
}

type CliRender struct {
	Render RenderRules `yaml:"render"`
}

type RenderRules struct {
	WhenTrue   []string       `yaml:"whenTrue,omitempty"`
	Sequence   []string       `yaml:"sequence,omitempty"`
	KeyValue   string         `yaml:"keyValue,omitempty"`
	RepeatEach *RepeatRules   `yaml:"repeatEach,omitempty"`
}

type RepeatRules struct {
	Sequence []string `yaml:"sequence,omitempty"`
}

type UiRender struct {
	Label    string `yaml:"label"`
	Category string `yaml:"category,omitempty"`
	Widget   string `yaml:"widget,omitempty"`
}

type Positional struct {
	ID     string `yaml:"id"`
	Source string `yaml:"source"`
	Order  int    `yaml:"order"`
}

type Output struct {
	ID           string `yaml:"id"`
	Type         string `yaml:"type"`
	Path         string `yaml:"path"`
	Downloadable bool   `yaml:"downloadable,omitempty"`
	ContentType  string `yaml:"contentType,omitempty"`
}

type Preset struct {
	ID          string                 `yaml:"id" json:"id"`
	ToolID      string                 `yaml:"tool_id" json:"tool_id"`
	Name        string                 `yaml:"name" json:"name"`
	Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Visibility  string                 `yaml:"visibility,omitempty" json:"visibility,omitempty"`
	Locked      bool                   `yaml:"locked,omitempty" json:"locked,omitempty"`
	Values      map[string]interface{} `yaml:"values" json:"values"`
}

type SecurityPolicy struct {
	RiskLevel             string     `yaml:"riskLevel,omitempty"`
	RequireApprovalForRun bool       `yaml:"requireApprovalForRun,omitempty"`
	DenyUnknownFlags      bool       `yaml:"denyUnknownFlags,omitempty"`
	AllowRawArgs          bool       `yaml:"allowRawArgs,omitempty"`
	PathPolicy            PathPolicy `yaml:"pathPolicy,omitempty"`
	DeniedFlags           []string   `yaml:"deniedFlags,omitempty"`
	Network               NetworkPolicy `yaml:"network,omitempty"`
}

type PathPolicy struct {
	AllowAbsolutePaths   bool `yaml:"allowAbsolutePaths,omitempty"`
	AllowParentTraversal bool `yaml:"allowParentTraversal,omitempty"`
	AllowSymlinks        bool `yaml:"allowSymlinks,omitempty"`
}

type Group struct {
	ID      string   `yaml:"id"`
	Type    string   `yaml:"type"`
	Members []string `yaml:"members"`
}

type AuditLog struct {
	JobID     string    `json:"job_id"`
	ToolID    string    `json:"tool_id"`
	ExitCode  int       `json:"exit_code"`
	CreatedAt time.Time `json:"created_at"`
}
