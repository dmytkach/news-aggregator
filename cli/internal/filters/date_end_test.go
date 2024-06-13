package filters

import (
	"news-aggregator/cli/internal/entity"
	"reflect"
	"testing"
	"time"
)

func TestDateEnd_Filter(t *testing.T) {

	tests := []struct {
		name    string
		endDate time.Time
		news    []entity.News
		want    []entity.News
	}{
		{
			name:    "should filter news to a given date.",
			endDate: time.Date(2024, time.May, 18, 23, 0, 0, 0, time.UTC),
			news: []entity.News{
				{
					Title:       "Container ship that struck Baltimore bridge will be removed from the site 'within days,' Maryland governor says",
					Description: "Ship that struck Francis Scott Key Bridge in Baltimore will be removed \"within days,\" Maryland Gov. Wes Moore says",
					Link:        "https://www.nbcnews.com/politics/politics-news/francis-scott-key-bridge-ship-removal-wes-moore-baltimore-rcna152955",
					Date:        time.Date(2024, 5, 19, 14, 6, 47, 0, time.UTC),
				},
				{
					Title:       "Harris says more Indian American representation is needed in government",
					Description: "Addressing a crowd of Indian Americans this week, Vice President Kamala Harris asserted the importance of voting and running. But Biden and Harris approval among the group has fallen.",
					Link:        "https://www.nbcnews.com/news/asian-america/kamala-harris-more-indian-american-representation-needed-government-rcna152761",
					Date:        time.Date(2024, 5, 17, 19, 48, 19, 0, time.UTC),
				},
				{
					Title:       "Atlanta officer accused of killing Lyft driver allegedly said victim was ‘gay fraternity’ recruiter",
					Description: "An Atlanta police officer accused of murdering a Lyft driver allegedly said the victim was in a gay fraternity trying to recruit him.",
					Link:        "https://www.nbcnews.com/nbc-out/out-news/atlanta-officer-accused-killing-lyft-driver-allegedly-said-victim-was-rcna152751",
					Date:        time.Date(2024, 5, 17, 14, 29, 43, 0, time.UTC),
				},
			},
			want: []entity.News{
				{
					Title:       "Harris says more Indian American representation is needed in government",
					Description: "Addressing a crowd of Indian Americans this week, Vice President Kamala Harris asserted the importance of voting and running. But Biden and Harris approval among the group has fallen.",
					Link:        "https://www.nbcnews.com/news/asian-america/kamala-harris-more-indian-american-representation-needed-government-rcna152761",
					Date:        time.Date(2024, 5, 17, 19, 48, 19, 0, time.UTC),
				},
				{
					Title:       "Atlanta officer accused of killing Lyft driver allegedly said victim was ‘gay fraternity’ recruiter",
					Description: "An Atlanta police officer accused of murdering a Lyft driver allegedly said the victim was in a gay fraternity trying to recruit him.",
					Link:        "https://www.nbcnews.com/nbc-out/out-news/atlanta-officer-accused-killing-lyft-driver-allegedly-said-victim-was-rcna152751",
					Date:        time.Date(2024, 5, 17, 14, 29, 43, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def := &DateEnd{
				EndDate: tt.endDate,
			}
			if got := def.Filter(tt.news); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
