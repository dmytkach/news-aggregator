package v1

import (
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestDefaultValues(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	Client = fake.NewClientBuilder().WithScheme(scheme).Build()

	tests := []struct {
		name           string
		hotNews        *HotNews
		expectedValues HotNewsSpec
	}{
		{
			name: "default values works with empty fields",
			hotNews: &HotNews{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-hotnews",
					Namespace: "default",
				},
				Spec: HotNewsSpec{
					Keywords: []string{},
				},
			},
			expectedValues: HotNewsSpec{
				Keywords:  []string{},
				DateStart: "",
				DateEnd:   "",
				Feeds:     []string{},
				SummaryConfig: SummaryConfig{
					TitlesCount: 10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hotNews.Default()

			assert.Equal(t, tt.expectedValues.Keywords, tt.hotNews.Spec.Keywords, "Keywords should match")
			assert.Equal(t, tt.expectedValues.DateStart, tt.hotNews.Spec.DateStart, "DateStart should match")
			assert.Equal(t, tt.expectedValues.DateEnd, tt.hotNews.Spec.DateEnd, "DateEnd should match")
			assert.Equal(t, tt.expectedValues.Feeds, tt.hotNews.Spec.Feeds, "Feeds should match")
			assert.Equal(t, tt.expectedValues.SummaryConfig.TitlesCount, tt.hotNews.Spec.SummaryConfig.TitlesCount, "TitlesCount should match")
		})
	}
}
func TestValidateHotNewsAllErrors(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	existingFeed := &Feed{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "valid-feed",
			Namespace: "default",
		},
	}

	Client = fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingFeed).Build()

	tests := []struct {
		name      string
		hotNews   *HotNews
		expectErr bool
	}{
		{
			name: "multiple errors",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{},
					DateStart: "01-01-2024",
					DateEnd:   "2024-01-01",
					Feeds:     []string{"non-existing-feed"},
				},
			},
			expectErr: true,
		},
		{
			name: "valid configuration",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"news"},
					DateStart: "2024-01-01",
					DateEnd:   "2024-01-10",
					Feeds:     []string{"valid-feed"},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.hotNews.validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateHotNewsCreationAndUpdate(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	existingFeed := &Feed{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "valid-feed",
			Namespace: "default",
		},
	}

	Client = fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingFeed).Build()

	tests := []struct {
		name      string
		hotNews   *HotNews
		expectErr bool
	}{
		{
			name: "valid creation",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"news"},
					DateStart: "2023-01-01",
					DateEnd:   "2023-01-10",
					Feeds:     []string{"valid-feed"},
				},
			},
			expectErr: false,
		},
		{
			name: "valid update",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"update"},
					DateStart: "2023-01-01",
					DateEnd:   "2023-01-10",
					Feeds:     []string{"valid-feed"},
				},
			},
			expectErr: false,
		},
		{
			name: "invalid creation",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{},
					DateStart: "01-01-2023",
					DateEnd:   "2023-01-01",
					Feeds:     []string{"non-existing-feed"},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid update",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{},
					DateStart: "01-01-2023",
					DateEnd:   "2023-01-01",
					Feeds:     []string{"non-existing-feed"},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.name == "invalid creation" {
				_, err = tt.hotNews.ValidateCreate()
			} else if tt.name == "invalid update" {
				_, err = tt.hotNews.ValidateUpdate(&HotNews{})
			} else {
				_, err = tt.hotNews.ValidateCreate()
			}
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestValidateHotNewsKeywords(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	Client = fake.NewClientBuilder().WithScheme(scheme).Build()

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
			_, err := tt.hotNews.validate()
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

	Client = fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingFeed).Build()

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
			_, err := tt.hotNews.validate()
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

	Client = fake.NewClientBuilder().WithScheme(scheme).Build()

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
					DateStart: "2024-01-01",
					DateEnd:   "2024-01-10",
				},
			},
			expectErr: false,
		},
		{
			name: "end date before start date",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"date"},
					DateStart: "2024-01-10",
					DateEnd:   "2024-01-01",
				},
			},
			expectErr: true,
		},
		{
			name: "incorrect start date format",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"date"},
					DateStart: "01-01-2024",
					DateEnd:   "2024-11-11",
				},
			},
			expectErr: true,
		},
		{
			name: "incorrect end date format",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					Keywords:  []string{"date"},
					DateStart: "2024-01-02",
					DateEnd:   "01-10-2024",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.hotNews.validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
