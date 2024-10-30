# Example

Shows how the `gosip` client can't be used to list or download files from Sharepoint if they contain special characters in their names, e.g. `#`, `?`, `%`.

## Usage

```bash
go run . -site-url https://contoso.sharepoint.com/sites/contoso -folder "/sites/contoso/Shared Documents/Example"
```
