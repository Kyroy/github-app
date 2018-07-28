package config

import (
	"github.com/google/go-github/github"
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
)

const name = ".kyroy.yaml"

func Download(client *github.Client, owner, repo, ref string) (*Config, error) {
	f, err := client.Repositories.DownloadContents(context.Background(), owner, repo, name, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	if err != nil {
		return nil, fmt.Errorf("could not download config file: %v", err)
	}
	var cfg Config
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("could not decode config file: %v", err)
	}
	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config file: %v", err)
	}
	for i, version := range cfg.Versions {
		cfg.Versions[i] = fmt.Sprintf("%s:%s", cfg.dockerImage(), version)
	}
	if cfg.GoImportPath == "" {
		cfg.GoImportPath = fmt.Sprintf("github.com/%s/%s", owner, repo)
	}
	return &cfg, nil
}