package v1_test

import (
	v12 "com.teamdev/news-aggregator/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe("FeedStatus", func() {
	var fs *v12.FeedStatus
	var condition v12.Condition

	BeforeEach(func() {
		fs = &v12.FeedStatus{}
		condition = v12.Condition{
			Type:           v12.ConditionAdded,
			Status:         true,
			Reason:         "Added successfully",
			Message:        "Feed was added",
			LastUpdateTime: metav1.Time{Time: time.Now()},
		}
	})

	Describe("AddCondition", func() {
		It("should add a condition to the FeedStatus", func() {
			fs.AddCondition(condition)
			Expect(fs.Conditions).To(HaveLen(1))
			Expect(fs.Conditions[0].Type).To(Equal(condition.Type))
			Expect(fs.Conditions[0].Status).To(Equal(condition.Status))
			Expect(fs.Conditions[0].Reason).To(Equal(condition.Reason))
			Expect(fs.Conditions[0].Message).To(Equal(condition.Message))
		})
	})

	Describe("Contains", func() {
		Context("when checking for an existing condition", func() {
			BeforeEach(func() {
				fs.Conditions = []v12.Condition{
					condition,
				}
			})

			It("should return true for existing condition with matching status", func() {
				exists := fs.Contains(v12.ConditionAdded, true)
				Expect(exists).To(BeTrue())
			})

			It("should return false for non-existing condition or mismatched status", func() {
				exists := fs.Contains(v12.ConditionUpdated, false)
				Expect(exists).To(BeFalse())
			})
		})
	})
})
