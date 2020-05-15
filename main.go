package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

const (
	coberturaDTDDecl = "<!DOCTYPE coverage SYSTEM \"http://cobertura.sourceforge.net/xml/coverage-04.dtd\">\n"
	formatHTML       = "html"
	formatCobertura  = "cobertura"
	formatClover     = "clover"
	formatJacoco     = "jacoco"
)

var allowedOutputFormats = allowedFormats{
	formatCobertura,
}

var allowedInputFormats = allowedFormats{
	formatCobertura,
}

const timeFormat = "2006-01-02-15-04-05" //"Jan 2, 2006 at 3:04pm (MST)"
var now = time.Now()

func main() {
	fileName := now.Format(timeFormat)

	var reports multiString
	flag.Var(&reports, "r", "Path to input report. (Required)")

	var outputFormats multiString
	flag.Var(&outputFormats, "f", fmt.Sprintf("Output formats [%s] (default \"%s\")", allowedOutputFormats.String(), allowedOutputFormats[0]))

	inputFormat := flag.String("i", allowedInputFormats[0], fmt.Sprintf("Input format [%s]", allowedInputFormats.String()))
	toStdout := flag.Bool("o", false, "Output to stdout (only usable if single output format is selected)")

	flag.Parse()

	if len(os.Args) < 2 {
		errorExit("")
	}

	if !allowedInputFormats.Exists(*inputFormat) {
		errorExit(fmt.Sprintf("Input format \"%s\" not supported!", *inputFormat))
	}

	if outputFormats.Size() == 0 {
		outputFormats.Set(allowedOutputFormats[0])
	} else {
		for _, f := range outputFormats {
			if !allowedOutputFormats.Exists(f) {
				errorExit(fmt.Sprintf("Output format \"%s\" not supported!", f))
			}
		}

		if outputFormats.Size() > 1 && *toStdout {
			errorExit("-o option makes sense only when you want single output format")
		}
	}

	var files []reportFile
	if reports.Size() == 0 {
		errorExit("No reports are provided!")
	} else {
		for _, r := range reports {
			file := reportFile{
				Path:   r,
				Format: *inputFormat,
			}
			if !file.Exists() {
				errorExit(fmt.Sprintf("File \"%s\" does not exist!", file.Path))
			}
			files = append(files, file)
		}
	}

	var linesValid int64
	var branchesValid int64
	var allPackages []*coberturaPackage
	for _, f := range files {
		p, err := readCobertura(f.Path)
		if err != nil {
			errorExit(fmt.Sprintf("Error processing: %s", err))
		}

		linesValid += p.LinesValid
		branchesValid += p.BranchesValid

		for _, pkg := range p.Packages {
			allPackages = append(allPackages, pkg)
		}
	}

	merged := coberturaCoverage{
		Packages:      allPackages,
		LinesValid:    linesValid,
		BranchesValid: branchesValid,
		Timestamp:     time.Now().UnixNano() / int64(time.Millisecond),
		Version:       "0.1",
	}

	merged.LinesCovered = merged.NumLinesWithHits()
	merged.LineRate = merged.HitRate()
	merged.BranchesCovered = merged.NumBranchesCovered()
	merged.BranchRate = merged.CondRate()

	var out io.Writer
	if *toStdout {
		out = os.Stdout
	} else {
		out, _ = os.Create(fmt.Sprintf("cobertura-%s.xml", fileName))
	}

	fmt.Fprintf(out, xml.Header)
	fmt.Fprintf(out, coberturaDTDDecl)

	encoder := xml.NewEncoder(out)
	encoder.Indent("", "    ")
	err := encoder.Encode(merged)
	if err != nil {
		errorExit(fmt.Sprintf("Error processing: %s", err))
	}

	fmt.Fprintln(out)
}

func readCobertura(file string) (coberturaCoverage, error) {
	xmlFile, err := os.Open(file)
	defer xmlFile.Close()
	if err != nil {
		return coberturaCoverage{}, err
	}

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return coberturaCoverage{}, err
	}

	var coverage coberturaCoverage
	xml.Unmarshal(byteValue, &coverage)
	coverage.Origin = file
	return coverage, nil
}
