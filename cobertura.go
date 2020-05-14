package main

import (
	"encoding/xml"
	"regexp"
	"strconv"
)

type coberturaCoverage struct {
	Origin          string              `xml:"-"`
	XMLName         xml.Name            `xml:"coverage"`
	LineRate        float32             `xml:"line-rate,attr"`
	BranchRate      float32             `xml:"branch-rate,attr"`
	Version         string              `xml:"version,attr"`
	Timestamp       int64               `xml:"timestamp,attr"`
	LinesCovered    int64               `xml:"lines-covered,attr"`
	LinesValid      int64               `xml:"lines-valid,attr"`
	BranchesCovered int64               `xml:"branches-covered,attr"`
	BranchesValid   int64               `xml:"branches-valid,attr"`
	Complexity      float32             `xml:"complexity,attr"`
	Sources         []*coberturaSource  `xml:"sources>source"`
	Packages        []*coberturaPackage `xml:"packages>package"`
}

type coberturaSource struct {
	Path string `xml:",chardata"`
}

type coberturaPackage struct {
	Name       string            `xml:"name,attr"`
	LineRate   float32           `xml:"line-rate,attr"`
	BranchRate float32           `xml:"branch-rate,attr"`
	Complexity float32           `xml:"complexity,attr"`
	Classes    []*coberturaClass `xml:"classes>class"`
}

type coberturaClass struct {
	Name       string             `xml:"name,attr"`
	Filename   string             `xml:"filename,attr"`
	LineRate   float32            `xml:"line-rate,attr"`
	BranchRate float32            `xml:"branch-rate,attr"`
	Complexity float32            `xml:"complexity,attr"`
	Methods    []*coberturaMethod `xml:"methods>method"`
	Lines      coberturaLines     `xml:"lines>line"`
}

type coberturaMethod struct {
	Name       string         `xml:"name,attr"`
	Signature  string         `xml:"signature,attr"`
	LineRate   float32        `xml:"line-rate,attr"`
	BranchRate float32        `xml:"branch-rate,attr"`
	Complexity float32        `xml:"complexity,attr"`
	Lines      coberturaLines `xml:"lines>line"`
}

type coberturaLine struct {
	Number            int    `xml:"number,attr"`
	Hits              int64  `xml:"hits,attr"`
	Branch            bool   `xml:"branch,attr"`
	ConditionCoverage string `xml:"condition-coverage,attr,omitempty"`
}

// Lines is a slice of Line pointers, with some convenience methods
type coberturaLines []*coberturaLine

// NumLinesWithHits returns the number of lines with a hit count > 0
func (lines coberturaLines) NumLinesWithHits() (numLinesWithHits int64) {
	for _, line := range lines {
		if line.Hits > 0 {
			numLinesWithHits++
		}
	}
	return numLinesWithHits
}

// NumLinesWithHits returns the number of lines with a hit count > 0
func (class coberturaClass) NumLinesWithHits() (numLinesWithHits int64) {
	return class.Lines.NumLinesWithHits()
}

// NumLinesWithHits returns the number of lines with a hit count > 0
func (pkg coberturaPackage) NumLinesWithHits() (numLinesWithHits int64) {
	for _, class := range pkg.Classes {
		numLinesWithHits += class.NumLinesWithHits()
	}
	return numLinesWithHits
}

// NumLinesWithHits returns the number of lines with a hit count > 0
func (cov coberturaCoverage) NumLinesWithHits() (numLinesWithHits int64) {
	for _, pkg := range cov.Packages {
		numLinesWithHits += pkg.NumLinesWithHits()
	}
	return numLinesWithHits
}

// HitRate returns a float32 from 0.0 to 1.0 representing what fraction of lines
// have hits
func (cov coberturaCoverage) HitRate() float32 {
	return float32(cov.NumLinesWithHits()) / float32(cov.LinesValid)
}

var condCvrgRe = regexp.MustCompile(`^\d+%\s+\((\d+)\/(\d+)\)$`)

// NumBranchesCovered returns the number of branches covered
func (lines coberturaLines) NumBranchesCovered() (numBranchesCovered int64) {
	for _, line := range lines {
		if line.Branch {
			s := condCvrgRe.ReplaceAllString(line.ConditionCoverage, `$1`)
			num, _ := strconv.ParseInt(s, 10, 64)
			numBranchesCovered += num
		}
	}
	return numBranchesCovered
}

// NumBranchesCovered returns the number of branches covered
func (class coberturaClass) NumBranchesCovered() (numBranchesCovered int64) {
	return class.Lines.NumBranchesCovered()
}

// NumBranchesCovered returns the number of branches covered
func (pkg coberturaPackage) NumBranchesCovered() (numBranchesCovered int64) {
	for _, class := range pkg.Classes {
		numBranchesCovered += class.NumBranchesCovered()
	}
	return numBranchesCovered
}

// NumBranchesCovered returns the number of branches covered
func (cov coberturaCoverage) NumBranchesCovered() (numBranchesCovered int64) {
	for _, pkg := range cov.Packages {
		numBranchesCovered += pkg.NumBranchesCovered()
	}
	return numBranchesCovered
}

// CondRate returns a float32 from 0.0 to 1.0 representing what fraction of branches
// have hits
func (cov coberturaCoverage) CondRate() float32 {
	return float32(cov.NumBranchesCovered()) / float32(cov.BranchesValid)
}
