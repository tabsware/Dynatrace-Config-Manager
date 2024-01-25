// @license
// Copyright 2023 Dynatrace LLC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package match

import (
	"sort"

	"github.com/Dynatrace/Dynatrace-Config-Manager/one-topology/pkg/match/rules"
)

type IndexMap map[string][]int

type IndexEntry struct {
	indexValue string
	matchedIds []int
}

type ByIndexValue []IndexEntry

func (a ByIndexValue) Len() int           { return len(a) }
func (a ByIndexValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndexValue) Less(i, j int) bool { return a[i].indexValue < a[j].indexValue }

func addUniqueValueToIndex(index *IndexMap, value string, itemId int) {

	if value == "" {
		return
	}

	(*index)[value] = append((*index)[value], itemId)

}

func addValueToIndex(index *IndexMap, value interface{}, itemId int) {

	stringValue, isString := value.(string)

	if isString {
		addUniqueValueToIndex(
			index, stringValue, itemId)
		return
	}

	stringSliceValue, isStringSlice := value.([]string)
	if isStringSlice {
		for _, uniqueValue := range stringSliceValue {
			addUniqueValueToIndex(
				index, uniqueValue, itemId)
		}
		return
	}

	sliceValue, isInterfaceSlice := value.([]interface{})

	if isInterfaceSlice {
		for _, uniqueValue := range sliceValue {
			addUniqueValueToIndex(
				index, uniqueValue.(string), itemId)
		}
		return
	}
}

func GetValueFromPath(item interface{}, path []string) interface{} {

	if len(path) <= 0 {
		return nil
	}

	var current interface{}
	current = item

	for _, field := range path {

		fieldValue, ok := (current.(map[string]interface{}))[field]
		if ok {
			current = fieldValue
		} else {
			current = nil
			break
		}

	}

	if current == nil {
		return nil
	} else {
		return current
	}
}

func GetValueFromList(listItemKey rules.ListItemKey, value interface{}) interface{} {

	if value == nil {
		return nil
	}

	sliceValue, isSlice := value.([]interface{})

	if isSlice {
		// pass
	} else {
		return nil
	}

	values := []string{}

	for _, item := range sliceValue {
		itemMap, isMap := item.(map[string]interface{})

		if isMap {
			// pass
		} else {
			return nil
		}

		keyValue, keyFound := itemMap[listItemKey.KeyKey]

		if keyFound {
			// pass
		} else {
			return nil
		}

		if keyValue.(string) == listItemKey.KeyValue {
			valueValue, valueFound := itemMap[listItemKey.ValueKey]

			if valueFound {
				values = append(values, valueValue.(string))
			} else {
				return nil
			}

		}

	}

	if len(values) == 0 {
		return nil
	}

	return values

}

func flattenSortIndex(index *IndexMap) []IndexEntry {

	flatIndex := make([]IndexEntry, len(*index))
	idx := 0

	for indexValue, matchedIds := range *index {
		flatIndex[idx] = IndexEntry{
			indexValue: indexValue,
			matchedIds: matchedIds,
		}
		idx++
	}

	sort.Sort(ByIndexValue(flatIndex))

	return flatIndex
}

func genSortedItemsIndex(indexRule rules.IndexRule, items *MatchProcessingEnv) []IndexEntry {

	index := IndexMap{}

	for _, itemIdx := range *(items.CurrentRemainingMatch) {

		value := GetValueFromPath((*items.RawMatchList.GetValues())[itemIdx], indexRule.Path)
		if (indexRule.ListItemKey != rules.ListItemKey{} && indexRule.ListItemKey.KeyKey != "") {
			value = GetValueFromList(indexRule.ListItemKey, value)
		}
		if value != nil {
			addValueToIndex(&index, value, itemIdx)
		}

	}

	flatSortedIndex := flattenSortIndex(&index)

	return flatSortedIndex
}
