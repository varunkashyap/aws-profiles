## AWS Profiles Manager

If you need to interact with more than one AWS accounts but have not yet got around to setting up Bastion accounts and roles yet, you might find your self having to switch credentials multiple times a day. AWS Profiles Manager will help you take some of the pain out of switching while you (hopefully) move away from multiple access keys on your laptop approach!

### Usage:

```bash
eval $(aws-profiles dev)
```

This should pick up a profile named 'dev' from your aws.credentials file. By default, aws-profiles looks for the profiles inside `~/.aws/credentials`. This can be overwritten using the `-f` flag, like so:

```bash
eval $(aws-profiles -f /my/credentials/location dev)
```

### Installation:

For now, you will need go development environment to create a binary. Run the following commands to create a binary for your OS and Architecture:

```bash
go build -o aws-profiles github.com/varunkashyap/aws-profiles
```