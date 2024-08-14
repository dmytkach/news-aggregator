package v1

import (
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestValidateHotNewsKeywords(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).Build()

	tests := []struct {
		name      string
		hotNews   *HotNews
		expectErr bool
	}{
		{
			name: "valid keywords",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords: []string{"news", "update"},
				},
			},
			expectErr: false,
		},
		{
			name: "empty keywords",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords: []string{},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.hotNews.validateHotNews()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateHotNewsFeeds(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	existingFeed := &Feed{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "valid-feed",
			Namespace: "default",
		},
	}

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingFeed).Build()

	tests := []struct {
		name      string
		hotNews   *HotNews
		expectErr bool
	}{
		{
			name: "valid feeds",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords: []string{"feed"},
					Feeds:    []string{"valid-feed"},
				},
			},
			expectErr: false,
		},
		{
			name: "non-existing feed",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords: []string{"feed"},
					Feeds:    []string{"non-existing-feed"},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.hotNews.validateHotNews()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateHotNewsDates(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).Build()

	tests := []struct {
		name      string
		hotNews   *HotNews
		expectErr bool
	}{
		{
			name: "valid dates",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"date"},
					DateStart: "2023-01-01",
					DateEnd:   "2023-01-10",
				},
			},
			expectErr: false,
		},
		{
			name: "end date before start date",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"date"},
					DateStart: "2023-01-10",
					DateEnd:   "2023-01-01",
				},
			},
			expectErr: true,
		},
		{
			name: "invalid date format",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"date"},
					DateStart: "01-01-2023",
					DateEnd:   "01-10-2023",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.hotNews.validateHotNews()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
