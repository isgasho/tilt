package model

import (
	"fmt"

	"github.com/windmilleng/tilt/internal/yaml"
)

type DockerComposeTarget struct {
	Name       TargetName
	ConfigPath string
	Mounts     []Mount
	YAMLRaw    []byte // for diff'ing when config files change
	DfRaw      []byte // for diff'ing when config files change
}

func (t DockerComposeTarget) ID() TargetID {
	return TargetID{
		Type: TargetTypeDockerCompose,
		Name: t.Name,
	}
}

func (dc DockerComposeTarget) Validate() error {
	if dc.ID().Empty() {
		return fmt.Errorf("[Validate] DockerCompose resources missing name:\n%s", dc.YAMLRaw)
	}

	return nil
}

type K8sTarget struct {
	Name         TargetName
	YAML         string
	PortForwards []PortForward
}

func (k8s K8sTarget) Validate() error {
	if k8s.ID().Empty() {
		return fmt.Errorf("[Validate] K8s resources missing name:\n%s", k8s.YAML)
	}

	if k8s.YAML == "" {
		return fmt.Errorf("[Validate] K8s resources %q missing YAML", k8s.Name)
	}

	return nil
}

func (k8s K8sTarget) ID() TargetID {
	return TargetID{
		Type: TargetTypeK8s,
		Name: k8s.Name,
	}
}

func (k8s K8sTarget) AppendYAML(y string) K8sTarget {
	if k8s.YAML == "" {
		k8s.YAML = y
	} else {
		k8s.YAML = yaml.ConcatYAML(k8s.YAML, y)
	}
	return k8s
}

var _ TargetSpec = K8sTarget{}
var _ TargetSpec = DockerComposeTarget{}
