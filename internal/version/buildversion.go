package version

import (
	"bytes"
	"runtime"
	"text/template"
)

var (
	GitCommit string = "library-import"
	Version   string = "library-import"
	BuildTime string = "library-import"
)

const versionTemplate = `Client:
 Version:      {{.Version}}
 Go version:   {{.GoVersion}}
 Git commit:   {{.GitCommit}}
 Built:        {{.BuildTime}}
 OS/Arch:      {{.Os}}/{{.Arch}}`

type VersionOptions struct {
	GitCommit string
	Version   string
	BuildTime string
	GoVersion string
	Os        string
	Arch      string
}

func String() string {
	var doc bytes.Buffer
	vo := VersionOptions{
		GitCommit: GitCommit,
		Version:   Version,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
	tmpl, _ := template.New("version").Parse(versionTemplate)
	tmpl.Execute(&doc, vo)
	return doc.String()
}
