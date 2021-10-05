package main

import (
	"strings"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

var progNameColor = color.New(color.FgGreen)
var progArgColor = color.New(color.FgMagenta)

var flagNameColor = color.New(color.FgCyan)
var flagValueColor = color.New(color.FgYellow)

var commentColor = color.New(color.FgHiBlack)

type encodingFieldHelp struct {
	short  string
	long   string
	substr string
}

func writeRow(sb *strings.Builder, example string, fields []encodingFieldHelp) {
	const dashes = "--------------------------------------------------------------------------"
	exampleDashes := dashes[:len(example)]
	for _, f := range fields {
		sb.WriteString("  ")
		flagNameColor.Fprint(sb, "-e")
		sb.WriteByte(' ')
		flagValueColor.Fprint(sb, f.short)
		sb.WriteString(", ")
		flagNameColor.Fprint(sb, "-e")
		sb.WriteByte(' ')
		flagValueColor.Fprint(sb, f.long)
		sb.WriteString("              "[len(f.long):])
		i := strings.Index(example, f.substr)
		commentColor.Fprint(sb, exampleDashes[:i])
		sb.WriteString(f.substr)
		commentColor.Fprint(sb, exampleDashes[i+len(f.substr):])
		sb.WriteByte('\n')
	}
}

func encodingsMessage() string {
	var sb strings.Builder
	sb.WriteString("Valid encodings, and their intended usages:\n")
	sb.WriteString("                         ")
	const example1 = "http://user:pass@site.com/index.html?foo=bar#Hello"
	commentColor.Fprint(&sb, example1)
	sb.WriteByte('\n')

	writeRow(&sb, example1, []encodingFieldHelp{
		{short: "s", long: "path-segment", substr: "index.html"},
		{short: "p", long: "path", substr: "/index.html"},
		{short: "q", long: "query", substr: "bar"},
		{short: "h", long: "host", substr: "site.com"},
		{short: "c", long: "cred", substr: "user:pass"},
		{short: "f", long: "frag", substr: "#Hello"},
	})

	sb.WriteString("\n                         ")
	const example2 = "http://[::1%25eth0]/home/index.html"
	commentColor.Fprint(&sb, example2)
	sb.WriteByte('\n')

	writeRow(&sb, example2, []encodingFieldHelp{
		{short: "z", long: "zone", substr: "eth0"},
	})

	return sb.String()
}

func flagsMessage() string {
	var sb strings.Builder
	sb.WriteString("Flags:\n")

	pflag.VisitAll(func (flag *pflag.Flag) {
		if flag.Hidden {
			return
		}
		sb.WriteString("  ")
		var width int
		if flag.Shorthand != "" {
			flagNameColor.Fprintf(&sb, "-%s", flag.Shorthand)
			sb.WriteString(", ")
			width += 3 + len(flag.Shorthand)
		} else {
			sb.WriteString("    ")
			width += 4
		}
		flagNameColor.Fprintf(&sb, "--%s", flag.Name)
		width += 2 + len(flag.Name)
		t := flag.Value.Type()
		if t == "string" {
			if flag.DefValue != "" {
				sb.WriteByte(' ')
				flagValueColor.Fprintf(&sb, `"%s"`, flag.DefValue)
				width += 3 + len(flag.DefValue)
			} else {
				flagValueColor.Fprint(&sb, "string")
				width += 6
			}
		}
		const spaces = "                               "
		sb.WriteString(spaces[width:])
		sb.WriteString(flag.Usage)
		sb.WriteByte('\n')
	})

	return sb.String()
}

func sampleUsageMessage() string {
	var sb strings.Builder
	sb.WriteString(`
Encodes/decodes the input value for HTTP URL by default and prints
the encoded/decoded value to STDOUT.
`)

	sb.WriteString("  ")
	progNameColor.Fprint(&sb, os.Args[0])
	sb.WriteString("             ")
	commentColor.Fprint(&sb, "// read from STDIN")
	sb.WriteString("\n  ")
	progNameColor.Fprint(&sb, os.Args[0])
	sb.WriteByte(' ')
	progArgColor.Fprint(&sb, "myfile.txt")
	sb.WriteString("  ")
	commentColor.Fprint(&sb, "// read from myfile.txt")
	sb.WriteRune('\n')
	//Encodes/decodes the input value for HTTP URL by default and prints
	//the encoded/decoded value to STDOUT.
	//
	//  %s             // read from STDIN
	//  %s myfile.txt  // read from myfile.txt
	return sb.String()
}