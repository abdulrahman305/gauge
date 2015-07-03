// Copyright 2015 ThoughtWorks, Inc.

// This file is part of Gauge.

// Gauge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Gauge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Gauge.  If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"github.com/getgauge/common"
)

func ParseSpecFiles(specFiles []string, conceptDictionary *ConceptDictionary) ([]*Specification, []*ParseResult) {
	parseResultsChan := make(chan *ParseResult, len(specFiles))
	specsChan := make(chan *Specification, len(specFiles))
	parseResults := make([]*ParseResult, 0)
	specs := make([]*Specification, 0)

	for _, specFile := range specFiles {
		go parseSpec(specFile, conceptDictionary, specsChan, parseResultsChan)
	}
	for _, _ = range specFiles {
		parseResults = append(parseResults, <-parseResultsChan)
		spec := <-specsChan
		if spec != nil {
			specs = append(specs, spec)
		}
	}
	return specs, parseResults
}

func parseSpec(specFile string, conceptDictionary *ConceptDictionary, specChannel chan *Specification, parseResultChan chan *ParseResult) {
	specFileContent, err := common.ReadFileContents(specFile)
	if err != nil {
		specChannel <- nil
		parseResultChan <- &ParseResult{ParseError: &ParseError{Message: err.Error()}, Ok: false, FileName: specFile}
		return
	}
	spec, parseResult := new(SpecParser).parse(specFileContent, conceptDictionary)
	parseResult.FileName = specFile
	if !parseResult.Ok {
		specChannel <- nil
		parseResultChan <- parseResult
		return
	}
	spec.FileName = specFile
	specChannel <- spec
	parseResultChan <- parseResult
}
