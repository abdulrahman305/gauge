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
	"github.com/getgauge/gauge/gauge_messages"
	"github.com/getgauge/gauge/result"
	"github.com/golang/protobuf/proto"
)

func convertToProtoItem(item Item) *gauge_messages.ProtoItem {
	switch item.kind() {
	case scenarioKind:
		return convertToProtoScenarioItem(item.(*Scenario))
	case stepKind:
		return convertToProtoStepItem(item.(*Step))
	case commentKind:
		return convertToProtoCommentItem(item.(*Comment))
	case tableKind:
		return convertToProtoTableItem(item.(*Table))
	case tagKind:
		return convertToProtoTagItem(item.(*Tags))
	}
	return nil
}

func convertToProtoTagItem(tags *Tags) *gauge_messages.ProtoItem {
	return &gauge_messages.ProtoItem{ItemType: gauge_messages.ProtoItem_Tags.Enum(), Tags: convertToProtoTags(tags)}
}

func convertToProtoStepItem(step *Step) *gauge_messages.ProtoItem {
	if step.isConcept {
		return convertToProtoConcept(step)
	}
	return &gauge_messages.ProtoItem{ItemType: gauge_messages.ProtoItem_Step.Enum(), Step: convertToProtoStep(step)}
}

func convertToProtoStepItems(steps []*Step) []*gauge_messages.ProtoItem {
	protoItems := make([]*gauge_messages.ProtoItem, 0)
	for _, step := range steps {
		protoItems = append(protoItems, convertToProtoStepItem(step))
	}
	return protoItems
}

func convertToProtoScenarioItem(scenario *Scenario) *gauge_messages.ProtoItem {
	scenarioItems := make([]*gauge_messages.ProtoItem, 0)
	for _, item := range scenario.items {
		scenarioItems = append(scenarioItems, convertToProtoItem(item))
	}
	protoScenario := newProtoScenario(scenario)
	return &gauge_messages.ProtoItem{ItemType: gauge_messages.ProtoItem_Scenario.Enum(), Scenario: protoScenario}
}

func convertToProtoConcept(concept *Step) *gauge_messages.ProtoItem {
	protoConcept := &gauge_messages.ProtoConcept{ConceptStep: convertToProtoStep(concept), Steps: convertToProtoStepItems(concept.conceptSteps)}
	protoConceptItem := &gauge_messages.ProtoItem{ItemType: gauge_messages.ProtoItem_Concept.Enum(), Concept: protoConcept}
	return protoConceptItem
}

func convertToProtoStep(step *Step) *gauge_messages.ProtoStep {
	return &gauge_messages.ProtoStep{ActualText: proto.String(step.lineText), ParsedText: proto.String(step.value), Fragments: makeFragmentsCopy(step.fragments)}
}

func convertToProtoTags(tags *Tags) *gauge_messages.ProtoTags {
	return &gauge_messages.ProtoTags{Tags: getAllTags(tags)}

}

func getAllTags(tags *Tags) []string {
	allTags := make([]string, 0)
	for _, tag := range tags.values {
		allTags = append(allTags, *proto.String(tag))
	}
	return allTags
}

func makeFragmentsCopy(fragments []*gauge_messages.Fragment) []*gauge_messages.Fragment {
	copiedFragments := make([]*gauge_messages.Fragment, 0)
	for _, fragment := range fragments {
		copiedFragments = append(copiedFragments, makeFragmentCopy(fragment))
	}
	return copiedFragments
}

func makeFragmentCopy(fragment *gauge_messages.Fragment) *gauge_messages.Fragment {
	if fragment.GetFragmentType() == gauge_messages.Fragment_Text {
		return &gauge_messages.Fragment{FragmentType: gauge_messages.Fragment_Text.Enum(), Text: proto.String(fragment.GetText())}
	} else {
		return &gauge_messages.Fragment{FragmentType: gauge_messages.Fragment_Parameter.Enum(), Parameter: makeParameterCopy(fragment.Parameter)}
	}
}

func makeParameterCopy(parameter *gauge_messages.Parameter) *gauge_messages.Parameter {
	switch parameter.GetParameterType() {
	case gauge_messages.Parameter_Static:
		return &gauge_messages.Parameter{ParameterType: gauge_messages.Parameter_Static.Enum(), Value: proto.String(parameter.GetValue()), Name: proto.String(parameter.GetName())}
	case gauge_messages.Parameter_Dynamic:
		return &gauge_messages.Parameter{ParameterType: gauge_messages.Parameter_Dynamic.Enum(), Value: proto.String(parameter.GetValue()), Name: proto.String(parameter.GetName())}
	case gauge_messages.Parameter_Table:
		return &gauge_messages.Parameter{ParameterType: gauge_messages.Parameter_Table.Enum(), Table: makeTableCopy(parameter.GetTable()), Name: proto.String(parameter.GetName())}
	case gauge_messages.Parameter_Special_String:
		return &gauge_messages.Parameter{ParameterType: gauge_messages.Parameter_Special_String.Enum(), Value: proto.String(parameter.GetValue()), Name: proto.String(parameter.GetName())}
	case gauge_messages.Parameter_Special_Table:
		return &gauge_messages.Parameter{ParameterType: gauge_messages.Parameter_Special_Table.Enum(), Table: makeTableCopy(parameter.GetTable()), Name: proto.String(parameter.GetName())}
	}
	return parameter
}

func makeTableCopy(table *gauge_messages.ProtoTable) *gauge_messages.ProtoTable {
	copiedTable := &gauge_messages.ProtoTable{}
	copiedTable.Headers = makeProtoTableRowCopy(table.GetHeaders())

	copiedRows := make([]*gauge_messages.ProtoTableRow, 0)
	for _, tableRow := range table.GetRows() {
		copiedRows = append(copiedRows, makeProtoTableRowCopy(tableRow))
	}
	copiedTable.Rows = copiedRows
	return copiedTable
}

func makeProtoTableRowCopy(tableRow *gauge_messages.ProtoTableRow) *gauge_messages.ProtoTableRow {
	copiedCells := make([]string, 0)
	return &gauge_messages.ProtoTableRow{Cells: append(copiedCells, tableRow.GetCells()...)}
}

func convertToProtoSteps(steps []*Step) []*gauge_messages.ProtoStep {
	protoSteps := make([]*gauge_messages.ProtoStep, 0)
	for _, step := range steps {
		protoSteps = append(protoSteps, convertToProtoStep(step))
	}
	return protoSteps
}

func convertToProtoCommentItem(comment *Comment) *gauge_messages.ProtoItem {
	return &gauge_messages.ProtoItem{ItemType: gauge_messages.ProtoItem_Comment.Enum(), Comment: &gauge_messages.ProtoComment{Text: proto.String(comment.value)}}
}

func convertToProtoTableItem(table *Table) *gauge_messages.ProtoItem {
	return &gauge_messages.ProtoItem{ItemType: gauge_messages.ProtoItem_Table.Enum(), Table: convertToProtoTableParam(table)}
}

func convertToProtoParameters(args []*StepArg) []*gauge_messages.Parameter {
	params := make([]*gauge_messages.Parameter, 0)
	for _, arg := range args {
		params = append(params, convertToProtoParameter(arg))
	}
	return params
}

func convertToProtoParameter(arg *StepArg) *gauge_messages.Parameter {
	switch arg.ArgType {
	case Static:
		return &gauge_messages.Parameter{ParameterType: gauge_messages.Parameter_Static.Enum(), Value: proto.String(arg.Value), Name: proto.String(arg.Name)}
	case Dynamic:
		return &gauge_messages.Parameter{ParameterType: gauge_messages.Parameter_Dynamic.Enum(), Value: proto.String(arg.Value), Name: proto.String(arg.Name)}
	case TableArg:
		return &gauge_messages.Parameter{ParameterType: gauge_messages.Parameter_Table.Enum(), Table: convertToProtoTableParam(&arg.Table), Name: proto.String(arg.Name)}
	case SpecialString:
		return &gauge_messages.Parameter{ParameterType: gauge_messages.Parameter_Special_String.Enum(), Value: proto.String(arg.Value), Name: proto.String(arg.Name)}
	case SpecialTable:
		return &gauge_messages.Parameter{ParameterType: gauge_messages.Parameter_Special_Table.Enum(), Table: convertToProtoTableParam(&arg.Table), Name: proto.String(arg.Name)}
	}
	return nil
}

func convertToProtoTableParam(table *Table) *gauge_messages.ProtoTable {
	protoTableParam := &gauge_messages.ProtoTable{Rows: make([]*gauge_messages.ProtoTableRow, 0)}
	protoTableParam.Headers = &gauge_messages.ProtoTableRow{Cells: table.Headers}
	for _, row := range table.Rows() {
		protoTableParam.Rows = append(protoTableParam.Rows, &gauge_messages.ProtoTableRow{Cells: row})
	}
	return protoTableParam
}

func addExecutionResult(protoItem *gauge_messages.ProtoItem, protoStepExecutionResult *gauge_messages.ProtoStepExecutionResult) {
	if protoStepExecutionResult != nil {
		protoItem.Step.StepExecutionResult = protoStepExecutionResult
	}
}

func convertToProtoSuiteResult(suiteResult *result.SuiteResult) *gauge_messages.ProtoSuiteResult {
	protoSuiteResult := &gauge_messages.ProtoSuiteResult{
		PreHookFailure:   suiteResult.PreSuite,
		PostHookFailure:  suiteResult.PostSuite,
		Failed:           proto.Bool(suiteResult.IsFailed),
		SpecsFailedCount: proto.Int32(int32(suiteResult.SpecsFailedCount)),
		ExecutionTime:    proto.Int64(suiteResult.ExecutionTime),
		SpecResults:      convertToProtoSpecResult(suiteResult.SpecResults),
		SuccessRate:      proto.Float32(getSuccessRate(len(suiteResult.SpecResults), suiteResult.SpecsFailedCount)),
		Environment:      proto.String(suiteResult.Environment),
		Tags:             proto.String(suiteResult.Tags),
		ProjectName:      proto.String(suiteResult.ProjectName),
		Timestamp:        proto.String(suiteResult.Timestamp),
	}
	return protoSuiteResult
}

func getSuccessRate(totalSpecs int, failedSpecs int) float32 {
	if totalSpecs == 0 {
		return 0
	}
	return (float32)(100.0 * (totalSpecs - failedSpecs) / totalSpecs)
}

func convertToProtoSpecResult(specResults []*result.SpecResult) []*gauge_messages.ProtoSpecResult {
	protoSpecResults := make([]*gauge_messages.ProtoSpecResult, 0)
	for _, specResult := range specResults {
		protoSpecResult := &gauge_messages.ProtoSpecResult{
			ProtoSpec:           specResult.ProtoSpec,
			ScenarioCount:       proto.Int32(int32(specResult.ScenarioCount)),
			ScenarioFailedCount: proto.Int32(int32(specResult.ScenarioFailedCount)),
			Failed:              proto.Bool(specResult.IsFailed),
			FailedDataTableRows: specResult.FailedDataTableRows,
			ExecutionTime:       proto.Int64(specResult.ExecutionTime),
		}
		protoSpecResults = append(protoSpecResults, protoSpecResult)
	}
	return protoSpecResults
}

func convertToProtoSpec(spec *Specification) *gauge_messages.ProtoSpec {
	protoSpec := newProtoSpec(spec)
	protoItems := make([]*gauge_messages.ProtoItem, 0)
	for _, item := range spec.items {
		protoItems = append(protoItems, convertToProtoItem(item))
	}
	protoSpec.Items = protoItems
	return protoSpec
}

func convertToProtoStepValue(stepValue *StepValue) *gauge_messages.ProtoStepValue {
	return &gauge_messages.ProtoStepValue{
		StepValue:              proto.String(stepValue.stepValue),
		ParameterizedStepValue: proto.String(stepValue.parameterizedStepValue),
		Parameters:             stepValue.args,
	}
}

func newProtoSpec(specification *Specification) *gauge_messages.ProtoSpec {
	return &gauge_messages.ProtoSpec{
		Items:         make([]*gauge_messages.ProtoItem, 0),
		SpecHeading:   proto.String(specification.heading.value),
		IsTableDriven: proto.Bool(false),
		FileName:      proto.String(specification.FileName),
		Tags:          getTags(specification.tags),
	}

}

func newSpecResult(specification *Specification) *result.SpecResult {
	return &result.SpecResult{
		ProtoSpec:           newProtoSpec(specification),
		FailedDataTableRows: make([]int32, 0),
	}
}

func newProtoScenario(scenario *Scenario) *gauge_messages.ProtoScenario {
	return &gauge_messages.ProtoScenario{
		ScenarioHeading: proto.String(scenario.heading.value),
		Failed:          proto.Bool(false),
		Tags:            getTags(scenario.tags),
		Contexts:        make([]*gauge_messages.ProtoItem, 0),
		ExecutionTime:   proto.Int64(0),
	}
}

func getTags(tags *Tags) []string {
	if tags != nil {
		return tags.values
	} else {
		return make([]string, 0)
	}
}
