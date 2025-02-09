# CLI Google Translator written on Golang

A CLI application that translates text files using Google Translate API.
It supports both Basic and Advanced Google Translate APIs and various language options.
The Basic API is simpler but has fewer features, while the Advanced API offers more control but requires a Google Cloud Project ID.

Example usage:

1. Using Basic API (default):

```bash
./gootrago -i input.txt -o output.txt -t es
```

2. Using Advanced API:

```bash
./gootrago -i input.txt -o output.txt -t uk -p your-project-id -a
```

3. Full options with Advanced API:

```bash
./gootrago --input input.txt --output output.txt --source en --target es --project your-project-id --credentials /path/to/creds.json --advanced
```

Configuration file (`.gootrago.yaml`) can now include API preference:

```yaml
input: default-input.txt
output: default-output.txt
advanced: true
project: your-project-id
credentials: /path/to/credentials.json
```

The program will:

1. Use Basic API by default
2. Switch to Advanced API when `-a` or `--advanced` flag is used
3. Require project ID only when using Advanced API
4. Support all other features (source language detection, custom credentials, etc.) in both modes
5. Indicate which API version was used in the output
