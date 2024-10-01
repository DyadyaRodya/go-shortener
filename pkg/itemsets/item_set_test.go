package itemsets

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ItemSet_AddItems(t *testing.T) {
	var itemSet, resSet []string
	tests := []struct {
		name          string
		startItems    []string
		addItems      []string
		expectedItems []string
	}{
		{
			name:          "success",
			startItems:    []string{"1", "2", "3"},
			addItems:      []string{"4", "5", "6"},
			expectedItems: []string{"1", "2", "3", "4", "5", "6"},
		},
		{
			name:          "add to empty",
			startItems:    []string{},
			addItems:      []string{"4", "5", "6"},
			expectedItems: []string{"4", "5", "6"},
		},
		{
			name:          "not updated",
			startItems:    []string{"1", "2", "3"},
			addItems:      []string{"1", "2", "3"},
			expectedItems: []string{"1", "2", "3"},
		},
		{
			name:          "duplicate",
			startItems:    []string{"1", "2", "3", "4", "5"},
			addItems:      []string{"2", "4", "5", "6"},
			expectedItems: []string{"1", "2", "3", "4", "5", "6"},
		},
		{
			name:          "nothing to add",
			startItems:    []string{"1", "2", "3"},
			addItems:      []string{},
			expectedItems: []string{"1", "2", "3"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			itemSet = tc.startItems
			resSet = AddItems(itemSet, tc.addItems)
			require.Equal(t, tc.expectedItems, resSet)
		})
	}
}

func Test_ItemSet_RemoveItems(t *testing.T) {
	var itemSet, resSet []string
	tests := []struct {
		name          string
		startItems    []string
		removeItems   []string
		expectedItems []string
	}{
		{
			name:          "success1",
			startItems:    []string{"1", "2", "3", "4", "5", "6"},
			removeItems:   []string{"4", "5", "6"},
			expectedItems: []string{"1", "2", "3"},
		},
		{
			name:          "success2",
			startItems:    []string{"1", "2", "3", "4", "5", "6"},
			removeItems:   []string{"1", "2", "3"},
			expectedItems: []string{"4", "5", "6"},
		},
		{
			name:          "success3",
			startItems:    []string{"1", "2", "3", "4", "5", "6"},
			removeItems:   []string{"1", "3", "5", "6"},
			expectedItems: []string{"2", "4"},
		},
		{
			name:          "remove from empty",
			startItems:    []string{},
			removeItems:   []string{"4", "5", "6"},
			expectedItems: []string{},
		},
		{
			name:          "not updated",
			startItems:    []string{"1", "2", "3"},
			removeItems:   []string{"4", "5", "6"},
			expectedItems: []string{"1", "2", "3"},
		},
		{
			name:          "nothing to remove",
			startItems:    []string{"1", "2", "3"},
			removeItems:   []string{},
			expectedItems: []string{"1", "2", "3"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			itemSet = tc.startItems
			resSet = RemoveItems(itemSet, tc.removeItems)
			require.Equal(t, tc.expectedItems, resSet)
		})
	}
}

func Test_ItemSet_Intersection(t *testing.T) {
	var resSet []string
	tests := []struct {
		name          string
		itemSet1      []string
		itemSet2      []string
		expectedItems []string
	}{
		{
			name:          "success1",
			itemSet1:      []string{"1", "2", "3", "4", "5", "6"},
			itemSet2:      []string{"4", "5", "6", "7", "8", "9"},
			expectedItems: []string{"4", "5", "6"},
		},
		{
			name:          "success2",
			itemSet1:      []string{"1", "2", "3", "4", "5", "6"},
			itemSet2:      []string{"4", "5", "6"},
			expectedItems: []string{"4", "5", "6"},
		},
		{
			name:          "success3",
			itemSet1:      []string{"4", "5", "6"},
			itemSet2:      []string{"4", "5", "6", "7", "8", "9"},
			expectedItems: []string{"4", "5", "6"},
		},
		{
			name:          "success4",
			itemSet1:      []string{},
			itemSet2:      []string{"4", "5", "6", "7", "8", "9"},
			expectedItems: []string{},
		},
		{
			name:          "success5",
			itemSet1:      []string{"1", "2", "3", "4", "5", "6"},
			itemSet2:      []string{},
			expectedItems: []string{},
		},
		{
			name:          "success6",
			itemSet1:      []string{},
			itemSet2:      []string{},
			expectedItems: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resSet = Intersection(tc.itemSet1, tc.itemSet2)
			require.Equal(t, tc.expectedItems, resSet)
		})
	}
}
