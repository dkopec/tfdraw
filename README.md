# tfdraw
A simple cli tool to convert the output of terraform graph to mermaidjs markdown.

## Usage

Clone repo, build and install then you can pipe the output of terraform show --json to tfdraw which then will output the markdown with the mermaid diagram.

```powershell
go build
go install
terraform show --json | tfdraw > diagram.md
```
