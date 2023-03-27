package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/cli/go-gh"
	"github.com/cli/safeexec"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: gh-release-notes <tag-name> [release-title]")
		os.Exit(1)
	}
	tagName := args[0]
	releaseTitle := args[0]
	if len(args) > 1 {
		releaseTitle = args[1]
	}

	body := map[string]interface{}{
		"tag_name": tagName,
	}
	for _, field := range []struct {
		key string
		val interface{}
	}{
		{"target_commitish", targetCommitish},
		{"previous_tag_name", previousTagName},
		{"configuration_file_path", configurationFilePath},
	} {
		if field.val != "" {
			body[field.key] = field.val
		}
	}

	client, err := gh.RESTClient(nil)
	onError(err)
	repo, err := gh.CurrentRepository()
	onError(err)

	endpoint := fmt.Sprintf("repos/%s/%s/releases/generate-notes", repo.Owner(), repo.Name())
	bodyBytes, err := json.Marshal(body)
	onError(err)
	var response response
	err = client.Post(endpoint, bytes.NewReader(bodyBytes), &response)
	onError(err)

	git, err := safeexec.LookPath("git")
	onError(err)
	cmd := exec.Command(
		git, "tag", "--cleanup=verbatim", "-a", tagName, "-m",
		fmt.Sprintf("%s\n\n%s", releaseTitle, response.Body),
	)
	err = cmd.Run()
	onError(err)

	fmt.Printf("Created tag %s!\n", tagName)
}

type response struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

var (
	targetCommitish       string
	previousTagName       string
	configurationFilePath string
)

func init() {
	flag.StringVar(
		&targetCommitish,
		"target",
		"",
		"Specifies the commitish value that will be the target for the release's tag. Required if the supplied tag_name does not reference an existing tag. Ignored if the tag_name already exists.",
	)
	flag.StringVar(
		&previousTagName,
		"previous",
		"",
		"The name of the previous tag to use as the starting point for the release notes. Use to manually specify the range for the set of changes considered as part this release.",
	)
	flag.StringVar(
		&configurationFilePath,
		"config",
		"",
		"Specifies a path to a file in the repository containing configuration settings used for generating the release notes. If unspecified, the configuration file located in the repository at '.github/release.yml' or '.github/release.yaml' will be used. If that is not present, the default configuration will be used.",
	)
}

func onError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
