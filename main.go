package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
)

var credentialsFile string

func init() {
	flag.StringVar(&credentialsFile, "f", "~/.aws/credentials", "Path to AWS credentials file")
	flag.Parse()
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
