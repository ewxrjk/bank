package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

func main() {
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "embed.go", "output filename")
	rootCmd.PersistentFlags().StringVarP(&pkg, "package", "p", "main", "output package")
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("embed: %s", err)
	}
}

var pkg, output string

type embedContext struct {
	path     string
	basename string
	varname  string
	mimetype string
}

var rootCmd = cobra.Command{
	Use:           "embed",
	Short:         "Embed type files in Go source",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = embed(output, args); err != nil {
			return
		}
		if err = format(output); err != nil {
			return
		}
		return
	},
}

func embed(output string, args []string) (err error) {
	var w *os.File
	if w, err = os.Create(output); err != nil {
		return
	}
	defer w.Close()
	if _, err = fmt.Fprintf(w, "package %s\n// Automatically generated code, don't edit\n", pkg); err != nil {
		return
	}
	contexts := []embedContext{}
	for _, path := range args {
		context := embedContext{path: path}
		if err = context.process(w); err != nil {
			return
		}
		contexts = append(contexts, context)
	}
	if _, err = fmt.Fprintf(w, "var embedContent=map[string]string{\n"); err != nil {
		return
	}
	for _, context := range contexts {
		if _, err = fmt.Fprintf(w, "\"%s\": %s,\n", context.basename, context.varname); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintf(w, "}\nvar embedType=map[string]string{\n"); err != nil {
		return
	}
	for _, context := range contexts {
		if _, err = fmt.Fprintf(w, "\"%s\": \"%s\",\n", context.basename, context.mimetype); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintf(w, "}\n"); err != nil {
		return
	}
	return
}

var replaceRegexp = regexp.MustCompile(`[^a-zA-Z0-9]`)

func (e *embedContext) process(w io.Writer) (err error) {
	var content []byte
	if content, err = ioutil.ReadFile(e.path); err != nil {
		return
	}
	_, e.basename = path.Split(e.path)
	e.varname = replaceRegexp.ReplaceAllString(e.basename, "")
	e.mimetype = strings.Split(mime.TypeByExtension(path.Ext(e.basename)), ";")[0]
	if _, err = fmt.Fprintf(w, "const %v = (\n\"", e.varname); err != nil {
		return
	}
	col := 0
	for _, r := range []byte(content) {
		switch true {
		case r == '"' || r == '\\':
			if _, err = fmt.Fprintf(w, "\\%c", r); err != nil {
				return
			}
			col += 2
		case r == '\n':
			if _, err = fmt.Fprintf(w, "\\n\"+\n\""); err != nil {
				return
			}
			col = 0
		case r < 32 || r > 127:
			if _, err = fmt.Fprintf(w, "\\x%02x", r); err != nil {
				return
			}
			col += 4
		default:
			if _, err = fmt.Fprintf(w, "%c", r); err != nil {
				return
			}
			col++
		}
		if col >= 64 {
			if _, err = fmt.Fprintf(w, "\"+\n\""); err != nil {
				return
			}
			col = 0
		}
	}
	if _, err = fmt.Fprintf(w, "\")\n"); err != nil {
		return
	}
	return
}

func format(path string) (err error) {
	cmd := exec.Command("gofmt", "-w", path)
	if err = cmd.Run(); err != nil {
		return
	}
	return
}
