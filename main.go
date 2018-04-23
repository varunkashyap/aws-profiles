package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
)

var (
	credentialsFile string
	verbose         bool
	help            bool
)

func init() {
	flag.StringVar(&credentialsFile, "f", "~/.aws/credentials", "Path to AWS credentials file")
	flag.BoolVar(&verbose, "v", false, "Verbose output for debugging")
	flag.BoolVar(&help, "h", false, "Print command usage")
	flag.Parse()

	if help == true {
		usage := `aws-profiles is a tool to manage multiple AWS profiles using the credentials file

Usage: 

aws-profiles [-f filepath] [-v] [-h] profile-name

  -f	Override the default credentials file location (~/.aws/credentials)
  -v	Turn on verbose logging for debugging
  -h	Print this message
`
		fmt.Println(usage)
		os.Exit(0)
	}

	// disable logging if verbose is false
	if verbose == false {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	requestedProfile := os.Args[1]
	log.Printf("Credentialsfile: %s", expandTildeToUserHome(credentialsFile))
	cfg, err := ini.Load(expandTildeToUserHome(credentialsFile))

	if err != nil {
		fmt.Printf("Unable to read file: %v\n", err)
	}

	for _, section := range cfg.Sections() {
		if section.Name() == requestedProfile {
			keyHash := section.KeysHash()
			for k, v := range keyHash {
				canonicalForm, err := getCanonicalEnvVarName(k)
				if err != nil {
					os.Exit(1)
				}
				fmt.Printf("export %s=%s\n", canonicalForm, v)
			}
		}
	}
}

func expandTildeToUserHome(filePath string) string {
	if strings.HasPrefix(filePath, "~/") {
		return filepath.Join(os.Getenv("HOME"), filePath[2:])
	}
	return filePath
}

func getCanonicalEnvVarName(credFileVar string) (string, error) {
	switch credFileVar {
	case "region":
		return "AWS_DEFAULT_REGION", nil
	case "output":
		return "AWS_DEFAULT_OUTPUT", nil
	case "aws_access_key_id":
		return "AWS_ACCESS_KEY_ID", nil
	case "aws_secret_access_key":
		return "AWS_SECRET_ACCESS_KEY", nil
	case "aws_sts_token":
		return "AWS_STS_TOKEN", nil
	default:
		return "", fmt.Errorf("Unknown credentials file variable: %s", credFileVar)
	}
}
