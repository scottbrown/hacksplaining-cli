# hacksplaining

A CLI for the [Hacksplaining](https://hacksplaining.com) API. Manage users and track security training progress from the command line.

## Installation

Requires Go 1.21+ and [Task](https://taskfile.dev).

```bash
git clone git@github.com:scottbrown/hacksplaining-cli.git
cd hacksplaining-cli
task build
```

Or install directly to your `$GOPATH/bin`:

```bash
task install
```

## Configuration

The CLI needs a Hacksplaining API key. It checks these locations in order:

1. `HACKSPLAINING_API_KEY` environment variable
2. `~/.hacksplaining` file
3. `/etc/.hacksplaining` file

The config file should contain only the API key:

```bash
echo "your-api-key-here" > ~/.hacksplaining
chmod 600 ~/.hacksplaining
```

## Usage

### List all users

```bash
hacksplaining users
hacksplaining users --json
```

### Get a specific user's progress

```bash
hacksplaining users user@example.com
hacksplaining users user@example.com --json
```

### Add a user

```bash
hacksplaining add user@example.com
hacksplaining add user@example.com --group-id 123
hacksplaining add user@example.com --subject "Welcome" --message "Please complete your training"
```

### Remove a user

```bash
hacksplaining remove user@example.com
```

### Send a training reminder

```bash
hacksplaining remind user@example.com
hacksplaining remind user@example.com --subject "Reminder" --message "Please finish your training"
```

### Print version

```bash
hacksplaining version
```

## Development

```bash
task build    # Build the binary
task test     # Run tests
task lint     # Run linter
task clean    # Remove build artifacts
```

The version is automatically derived from git tags. To set an explicit version at build time:

```bash
go build -ldflags "-X github.com/scottbrown/hacksplaining-cli/cmd.Version=1.0.0" -o hacksplaining .
```

## License

See [LICENSE](LICENSE) for details.
