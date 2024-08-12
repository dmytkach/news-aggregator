package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAddCondition(t *testing.T) {
	fs := &FeedStatus{}

	condition := Condition{
		Type:           ConditionAdded,
		Status:         true,
		Reason:         "Added successfully",
		Message:        "Feed was added",
		LastUpdateTime: metav1.Time{Time: time.Now()},
	}

	fs.AddCondition(condition)

	assert.Len(t, fs.Conditions, 1, "Condition list should contain one condition")
	assert.Equal(t, condition.Type, fs.Conditions[0].Type, "Condition type should match")
	assert.Equal(t, condition.Status, fs.Conditions[0].Status, "Condition status should match")
	assert.Equal(t, condition.Reason, fs.Conditions[0].Reason, "Condition reason should match")
	assert.Equal(t, condition.Message, fs.Conditions[0].Message, "Condition message should match")
}

func TestContains(t *testing.T) {
	fs := &FeedStatus{
		Conditions: []Condition{
			{
				Type:           ConditionAdded,
				Status:         true,
				Reason:         "Added successfully",
				Message:        "Feed was added",
				LastUpdateTime: metav1.Time{Time: time.Now()},
			},
		},
	}

	exists := fs.Contains(ConditionAdded, true)
	assert.True(t, exists, "ConditionAdded with status true should exist")

	exists = fs.Contains(ConditionUpdated, false)
	assert.False(t, exists, "ConditionUpdated with status false should not exist")
}

func TestUpdateConditions(t *testing.T) {
	fs := &FeedStatus{
		Conditions: []Condition{
			{
				Type:           ConditionAdded,
				Status:         true,
				Reason:         "Added successfully",
				Message:        "Feed was added",
				LastUpdateTime: metav1.Time{Time: time.Now().Add(-time.Hour)},
			},
		},
	}

	fs.updateConditions()

	assert.True(t, fs.Conditions[0].LastUpdateTime.Time.After(time.Now().Add(-time.Minute)), "LastUpdateTime should be updated to the current time")
}
