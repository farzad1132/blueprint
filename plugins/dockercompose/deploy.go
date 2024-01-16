package dockercompose

import (
	"fmt"
	"path/filepath"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/blueprint"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/blueprint/ioutil"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/address"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/plugins/docker"
	"github.com/blueprint-uservices/blueprint/plugins/dockercompose/dockergen"
	"golang.org/x/exp/slog"
)

type (
	/*
		Docker compose is the default docker app deployer.  It simply
		generates a docker-compose.yml file on the local filesystem.
	*/
	dockerComposeDeployer interface {
		ir.ArtifactGenerator
	}

	/*
	   A workspace used when deploying a set of containers as a
	   docker-compose application

	   Implements docker.ContainerWorkspace defined in docker/ir.go

	   This workspace generates a docker-compose file at the root of the
	   output directory.  The docker-compose instantiates containers
	   that are either:
	    (a) pre-built images
	    (b) artifacts built using Dockerfiles in the output directory

	*/
	dockerComposeWorkspace struct {
		ir.VisitTrackerImpl

		info docker.ContainerWorkspaceInfo

		ImageDirs map[string]string // map from image name to directory

		DockerComposeFile *dockergen.DockerComposeFile
	}
)

// Implements ir.ArtifactGenerator
func (node *Deployment) GenerateArtifacts(dir string) error {
	slog.Info(fmt.Sprintf("Collecting container instances for deployment %s in %s", node.Name(), dir))
	workspace := NewDockerComposeWorkspace(node.Name(), dir)
	return node.generateArtifacts(workspace)
}

/*
The basic build process of a docker-compose deployment
*/
func (node *Deployment) generateArtifacts(workspace *dockerComposeWorkspace) error {

	// Add any locally-built container images
	for _, node := range ir.Filter[docker.ProvidesContainerImage](node.Nodes) {
		if err := node.AddContainerArtifacts(workspace); err != nil {
			return err
		}
	}

	// Collect all container instances
	for _, node := range ir.Filter[docker.ProvidesContainerInstance](node.Nodes) {
		if err := node.AddContainerInstance(workspace); err != nil {
			return err
		}
	}

	// Build the docker-compose file
	if err := workspace.Finish(); err != nil {
		return err
	}

	// Reset any port assignments for externally-visible servers, since they will currently
	// be assigned to docker-internal ports
	address.ResetPorts(node.Edges)
	return nil
}

func NewDockerComposeWorkspace(name string, dir string) *dockerComposeWorkspace {
	return &dockerComposeWorkspace{
		info: docker.ContainerWorkspaceInfo{
			Path:   filepath.Clean(dir),
			Target: "docker-compose",
		},
		ImageDirs:         make(map[string]string),
		DockerComposeFile: dockergen.NewDockerComposeFile(name, dir, "docker-compose.yml"),
	}
}

// Implements docker.ContainerWorkspace
func (d *dockerComposeWorkspace) Info() docker.ContainerWorkspaceInfo {
	return d.info
}

// Implements docker.ContainerWorkspace
func (d *dockerComposeWorkspace) CreateImageDir(imageName string) (string, error) {
	// Only alphanumeric and underscores are allowed in an proc name
	imageName = ir.CleanName(imageName)
	imageDir, err := ioutil.CreateNodeDir(d.info.Path, imageName)
	d.ImageDirs[imageName] = imageDir
	return imageDir, err
}

// Implements docker.ContainerWorkspace
func (d *dockerComposeWorkspace) DeclarePrebuiltInstance(instanceName string, image string, args ...ir.IRNode) error {
	// Docker containers should assign all internal server ports (typically using address.AssignPorts) before adding an instance
	if err := address.CheckPorts(args); err != nil {
		return blueprint.Errorf("unable to add docker instance %v due to %v", instanceName, err.Error())
	}

	return d.DockerComposeFile.AddImageInstance(instanceName, image, args...)
}

// Implements docker.ContainerWorkspace
func (d *dockerComposeWorkspace) DeclareLocalImage(instanceName string, imageDir string, args ...ir.IRNode) error {
	// Docker containers should assign all internal server ports (typically using address.AssignPorts) before adding an instance
	if err := address.CheckPorts(args); err != nil {
		return blueprint.Errorf("unable to add docker instance %v due to %v", instanceName, err.Error())
	}

	return d.DockerComposeFile.AddBuildInstance(instanceName, imageDir, args...)
}

// Implements docker.ContainerWorkspace
func (d *dockerComposeWorkspace) SetEnvironmentVariable(instanceName string, key string, val string) error {
	return d.DockerComposeFile.AddEnvVar(instanceName, key, val)
}

// Generates the docker-compose file
func (d *dockerComposeWorkspace) Finish() error {
	// Now that all images and instances have been declared, we can generate the docker-compose file
	return d.DockerComposeFile.Generate()
}

func (d *dockerComposeWorkspace) ImplementsBuildContext()       {}
func (d *dockerComposeWorkspace) ImplementsContainerWorkspace() {}