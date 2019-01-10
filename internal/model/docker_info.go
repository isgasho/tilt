package model

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/docker/distribution/reference"
)

type ImageTarget struct {
	cachePaths   []string
	Ref          reference.Named
	BuildDetails BuildDetails
}

func (di ImageTarget) ID() TargetID {
	return TargetID{
		Type: TargetTypeImage,
		Name: TargetName(di.Ref.String()),
	}
}

func (i ImageTarget) Validate() error {
	if i.Ref == nil {
		return fmt.Errorf("[Validate] Image target missing image ref: %+v", i.BuildDetails)
	}

	switch bd := i.BuildDetails.(type) {
	case StaticBuild:
		if bd.BuildPath == "" {
			return fmt.Errorf("[Validate] Image %q missing build path", i.Ref)
		}
	case FastBuild:
		if bd.BaseDockerfile == "" {
			return fmt.Errorf("[Validate] Image %q missing base dockerfile", i.Ref)
		}

		for _, mnt := range bd.Mounts {
			if !filepath.IsAbs(mnt.LocalPath) {
				return fmt.Errorf(
					"[Validate] Image %q: mount must be an absolute path (got: %s)", i.Ref, mnt.LocalPath)
			}
		}

	default:
		return fmt.Errorf("[Validate] Image %q has neither StaticBuildInfo nor FastBuildInfo", i.Ref)
	}

	return nil
}

type BuildDetails interface {
	buildDetails()
}

func (di ImageTarget) WithBuildDetails(details BuildDetails) ImageTarget {
	di.BuildDetails = details
	return di
}

func (di ImageTarget) WithCachePaths(paths []string) ImageTarget {
	di.cachePaths = append(append([]string{}, di.cachePaths...), paths...)
	sort.Strings(di.cachePaths)
	return di
}

func (di ImageTarget) CachePaths() []string {
	return di.cachePaths
}

type StaticBuild struct {
	Dockerfile string
	BuildPath  string // the absolute path to the files
	BuildArgs  DockerBuildArgs
}

func (StaticBuild) buildDetails()  {}
func (sb StaticBuild) Empty() bool { return reflect.DeepEqual(sb, StaticBuild{}) }

type FastBuild struct {
	BaseDockerfile string
	Mounts         []Mount
	Steps          []Step
	Entrypoint     Cmd
}

func (FastBuild) buildDetails()  {}
func (fb FastBuild) Empty() bool { return reflect.DeepEqual(fb, FastBuild{}) }

var _ TargetSpec = ImageTarget{}
