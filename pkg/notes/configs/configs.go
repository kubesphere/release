package configs

import (
	"os"
	"path/filepath"

	"k8s.io/release/pkg/notes/options"
	"sigs.k8s.io/release-sdk/git"
)

type Config struct {
	Repos              []*options.Options `yaml:"repos"`
	GithubBaseURL      string             `yaml:"github-base-url"`
	GithubUploadURL    string             `yaml:"github-upload-url"`
	Output             string             `yaml:"output"`
	Format             string             `yaml:"format"`
	MarkdownLinks      bool               `yaml:"markdown-links"`
	Debug              bool               `yaml:"debug"`
	ReleaseBucket      string             `yaml:"release-bucket"`
	ReleaseTars        string             `yaml:"release-tars"`
	Toc                bool               `yaml:"toc"`
	ListReleaseNotesV2 bool               `yaml:"list-v2"`
	GithubToken        string             `yaml:"github-token"`
	GoTemplate         string             `yaml:"go-template"`
}

func (c *Config) ValidateAndFinish() error {
	//
	for _, repo := range c.Repos {
		repo.GithubBaseURL = c.GithubBaseURL
		repo.GithubToken = c.GithubToken
		repo.Debug = c.Debug
		repo.AddMarkdownLinks = c.MarkdownLinks
		repo.Format = c.Format
		repo.GoTemplate = c.GoTemplate
		repo.GitCloneFn = git.CloneOrOpenGitHubRepo
		if repo.RepoPath == "" {
			repo.RepoPath = filepath.Join(os.TempDir(), repo.GithubRepo)
		}
		// repo.
		// for kubesphere, always set true
		repo.AddRepoName = true
		err := repo.ValidateAndFinishV2()
		if err != nil {
			return err
		}
	}

	return nil
}
