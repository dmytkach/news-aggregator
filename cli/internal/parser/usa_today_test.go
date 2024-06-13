package parser

import (
	"news-aggregator/cli/internal/entity"
	"reflect"
	"testing"
	"time"
)

func TestUsaToday_Parse(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    []entity.News
		wantErr bool
	}{{
		name: "should collect news from the specified sources.",
		file: "../testdata/news.html",
		want: []entity.News{
			{
				Title:       "Astronomers discover an enormous planet made of something as light as cotton candy",
				Description: "Oftentimes when astronomers discover new exoplanets they’re molten hot hellscapes or solid frozen spheres of ice. However, WASP-193b which resides 1,232 light-years is a little different.",
				Link:        "https://www.usatoday.com/videos/news/world/2024/05/15/astronomers-discover-an-enormous-planet-made-of-something-as-light-as-cotton-candy/73697406007/",
				Date:        time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			},
			{
				Title:       "Ukraine's Zelenskyy cancels all foreign trips as Russian offensive intensifies",
				Description: "Zelenskyy's cancellation of all foreign trips may indicate that Ukraine is struggling to contain Russia's offensive. The US said it is 'rushing' more weapons.",
				Link:        "https://www.usatoday.com/story/news/world/2024/05/15/ukraine-zelenskyy-cancels-foreign-trips-russian-offensive-blinken-visit/73697239007/",
				Date:        time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			},
			{
				Title:       "King Charles unveils first official portrait",
				Description: "The portrait painted by artist Jonathan Yeo depicts King Charles III wearing the uniform of the Welsh Guards military unit, against a red background.",
				Link:        "https://www.usatoday.com/videos/news/world/2024/05/14/king-charles-iii-first-portrait-since-his-coronation-unveiled/73689220007/",
				Date:        time.Date(2024, 5, 14, 0, 0, 0, 0, time.UTC),
			},
		},
		wantErr: false,
	},
		{
			name:    "should fail if HTML file does not contain expected elements",
			file:    "../testdata/invalid_news.html",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test parsing missing Html file",
			file:    "../testdata/nonexistent_file.html",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usaTodayParser := &UsaToday{
				FilePath: entity.PathToFile(tt.file),
			}
			got, err := usaTodayParser.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
