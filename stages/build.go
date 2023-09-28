package stages

import (
	"archive/tar"
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os"
)

type DockerBuild struct {
	target    string
	imageName string
	context   string
}

func NewDockerBuild(imageName, context string) *DockerBuild {
	return &DockerBuild{imageName: imageName, context: context}
}

func (d *DockerBuild) Run(cli *client.Client) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	dockerfile := d.context + "/Dockerfile"
	dockerFileReader, err := os.Open(dockerfile)
	readDockerFile, err := io.ReadAll(dockerFileReader)
	if err != nil {
		return err
	}
	tarHeader := &tar.Header{
		Name: dockerfile,
		Size: int64(len(readDockerFile)),
	}

	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return err
	}

	// Writes the dockerfile data to the TAR file
	_, err = tw.Write(readDockerFile)
	if err != nil {
		return err
	}

	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	opts := types.ImageBuildOptions{
		Context: dockerFileTarReader,
		Target:  d.target,
		Tags:    []string{d.imageName},
	}
	build, err := cli.ImageBuild(context.Background(), dockerFileTarReader, opts)
	if err != nil {
		return err
	}
	defer build.Body.Close()
	return nil
}

func (d *DockerBuild) Target(t string) *DockerBuild {
	d.target = t
	return d
}
