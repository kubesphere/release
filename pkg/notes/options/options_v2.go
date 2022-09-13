package options

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/release-sdk/github"
)

func (o *Options) ValidateAndFinishV2() (err error) {
	// Add appropriate log filtering

	if o.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if o.ReplayDir != "" && o.RecordDir != "" {
		return errors.New("please do not use record and replay together")
	}

	// Recover for replay if needed
	if o.ReplayDir != "" {
		logrus.Info("Using replay mode")
		return nil
	}

	// The GitHub Token is required if replay is not specified
	// token, ok := os.LookupEnv(github.TokenEnvKey)
	// if ok {
	// 	o.GithubToken = token
	// } else if o.ReplayDir == "" {
	// 	return errors.Errorf(
	// 		"neither environment variable `%s` nor `replay` option is set",
	// 		github.TokenEnvKey,
	// 	)
	// }

	// Check if we want to automatically discover the revisions
	if o.DiscoverMode != RevisionDiscoveryModeNONE {
		if err := o.resolveDiscoverMode(); err != nil {
			return err
		}
	}

	// The start SHA or rev is required.
	if o.StartSHA == "" && o.StartRev == "" {
		return errors.New("the starting commit hash must be set via --start-sha, $START_SHA, --start-rev or $START_REV")
	}

	// The end SHA or rev is required.
	if o.EndSHA == "" && o.EndRev == "" {
		return errors.New("the ending commit hash must be set via --end-sha, $END_SHA, --end-rev or $END_REV")
	}

	// Check if we have to parse a revision
	if (o.StartRev != "" && o.StartSHA == "") || (o.EndRev != "" && o.EndSHA == "") {
		repo, err := o.repo()

		if err != nil {
			fmt.Println("error")
			return err
		}
		if o.StartRev != "" && o.StartSHA == "" {
			sha, err := repo.RevParseTag(o.StartRev)
			if err != nil {
				return errors.Wrapf(err, "resolving %s", o.StartRev)
			}
			logrus.Infof("Using found start SHA: %s", sha)
			o.StartSHA = sha
		}
		if o.EndRev != "" && o.EndSHA == "" {
			sha, err := repo.RevParseTag(o.EndRev)
			if err != nil {
				return errors.Wrapf(err, "resolving %s", o.EndRev)
			}
			logrus.Infof("Using found end SHA: %s", sha)
			o.EndSHA = sha
		}
	}

	// Create the record dir
	if o.RecordDir != "" {
		logrus.Info("Using record mode")
		if err := os.MkdirAll(o.RecordDir, os.FileMode(0o755)); err != nil {
			return err
		}
	}

	// Set GithubBaseURL to https://github.com if it is unset.
	if o.GithubBaseURL == "" {
		o.GithubBaseURL = github.GitHubURL
	}

	if err := o.checkFormatOptions(); err != nil {
		return errors.Wrap(err, "while checking format flags")
	}
	return nil
}
