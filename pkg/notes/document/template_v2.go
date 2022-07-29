package document

const defaultReleaseNotesTemplateV2 = `
{{- $CurrentRevision := .CurrentRevision -}}
{{- $PreviousRevision := .PreviousRevision -}}

{{if .FileDownloads}}
## Downloads for {{$CurrentRevision}}

{{- with .FileDownloads.Source }}

### Source Code

filename | sha512 hash
-------- | -----------
{{range .}}[{{.Name}}]({{.URL}}) | {{.Checksum}}{{println}}{{end}}
{{end}}

{{- with .FileDownloads.Client -}}
### Client Binaries

filename | sha512 hash
-------- | -----------
{{range .}}[{{.Name}}]({{.URL}}) | {{.Checksum}}{{println}}{{end}}
{{end}}

{{- with .FileDownloads.Server -}}
### Server Binaries

filename | sha512 hash
-------- | -----------
{{range .}}[{{.Name}}]({{.URL}}) | {{.Checksum}}{{println}}{{end}}
{{end}}

{{- with .FileDownloads.Node -}}
### Node Binaries

filename | sha512 hash
-------- | -----------
{{range .}}[{{.Name}}]({{.URL}}) | {{.Checksum}}{{println}}{{end}}
{{end -}}
{{- end -}}

{{with .CVEList -}}
## Important Security Information

This release contains changes that address the following vulnerabilities:
{{range .}}
### {{.ID}}: {{.Title}}

{{.Description}}

**CVSS Rating:** {{.CVSSRating}} ({{.CVSSScore}}) [{{.CVSSVector}}]({{.CalcLink}})
{{- if .TrackingIssue -}}
<br>
**Tracking Issue:** {{.TrackingIssue}}
{{- end }}

{{ end }}
{{- end -}}

{{with .NotesWithActionRequired -}}
## Urgent Upgrade Notes 

### (No, really, you MUST read this before you upgrade)

{{range .}}{{println "-" .}} {{end}}
{{end}}

{{- if .NotesV2 -}}

{{ range $area, $NoteCategory  :=  .NotesV2}}
## {{ $area }}

{{ range $NoteCategory}}
### {{.Kind | prettyKind}}

{{range $note := .NoteEntries }}{{println "-" $note}}{{end}}
{{- end -}}

{{- end -}}

{{- end -}}
`
