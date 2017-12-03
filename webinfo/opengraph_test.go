package webinfo

import (
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	messagediff "gopkg.in/d4l3k/messagediff.v1"
)

func TestOpenGraphParse(t *testing.T) {
	tt := []struct {
		name string
		html string
		err  error
		out  *OG
	}{
		{
			"og basic metadata",
			`
			<meta property="og:title" content="The Rock" />
			<meta property="og:type" content="video.movie" />
			<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />
			<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />
			`,
			nil,
			&OG{
				Title: "The Rock",
				Type:  "video.movie",
				URL:   "http://www.imdb.com/title/tt0117500/",
				Image: &OGImage{
					Source: "http://ia.media-imdb.com/images/rock.jpg",
				},
			},
		},
		{
			"og image metadata",
			`
			<meta property="og:title" content="The Rock" />
			<meta property="og:type" content="video.movie" />
			<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />

			<meta property="og:image" content="http://example.com/ogp.jpg" />
			<meta property="og:image:secure_url" content="https://secure.example.com/ogp.jpg" />
			<meta property="og:image:type" content="image/jpeg" />
			<meta property="og:image:width" content="400" />
			<meta property="og:image:height" content="300" />
			<meta property="og:image:alt" content="A shiny red apple with a bite taken out" />
			`,
			nil,
			&OG{
				Title: "The Rock",
				Type:  "video.movie",
				URL:   "http://www.imdb.com/title/tt0117500/",
				Image: &OGImage{
					Source:    "http://example.com/ogp.jpg",
					SecureURL: "https://secure.example.com/ogp.jpg",
					Type:      "image/jpeg",
					Width:     "400",
					Height:    "300",
					Alt:       "A shiny red apple with a bite taken out",
				},
			},
		},
		{
			"og video metadata",
			`
			<meta property="og:title" content="The Rock" />
			<meta property="og:type" content="video.movie" />
			<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />
			<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />

			<meta property="og:video" content="http://example.com/movie.swf" />
			<meta property="og:video:secure_url" content="https://secure.example.com/movie.swf" />
			<meta property="og:video:type" content="application/x-shockwave-flash" />
			<meta property="og:video:width" content="400" />
			<meta property="og:video:height" content="300" />
			`,
			nil,
			&OG{
				Title: "The Rock",
				Type:  "video.movie",
				URL:   "http://www.imdb.com/title/tt0117500/",
				Image: &OGImage{
					Source: "http://ia.media-imdb.com/images/rock.jpg",
				},
				Video: &OGVideo{
					Source:    "http://example.com/movie.swf",
					SecureURL: "https://secure.example.com/movie.swf",
					Type:      "application/x-shockwave-flash",
					Width:     "400",
					Height:    "300",
				},
			},
		},
		{
			"og audio metadata",
			`
			<meta property="og:title" content="The Rock" />
			<meta property="og:type" content="video.movie" />
			<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />
			<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />

			<meta property="og:audio" content="http://example.com/sound.mp3" />
			<meta property="og:audio:secure_url" content="https://secure.example.com/sound.mp3" />
			<meta property="og:audio:type" content="audio/mpeg" />
			`,
			nil,
			&OG{
				Title: "The Rock",
				Type:  "video.movie",
				URL:   "http://www.imdb.com/title/tt0117500/",
				Image: &OGImage{
					Source: "http://ia.media-imdb.com/images/rock.jpg",
				},
				Audio: &OGAudio{
					Source:    "http://example.com/sound.mp3",
					SecureURL: "https://secure.example.com/sound.mp3",
					Type:      "audio/mpeg",
				},
			},
		},
	}

	for _, tr := range tt {
		t.Run(tr.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tr.html))
			if err != nil {
				t.Fatal(err)
			}

			og, err := parseOpenGraph(doc)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(og, tr.out) {
				t.Error("expected", tr.out)
				t.Error("got     ", og)

				diff, _ := messagediff.PrettyDiff(og, tr.out)
				t.Error(diff)
			}
		})
	}
}

func TestOpenGraphValidation(t *testing.T) {
	tt := []struct {
		name  string
		html  string
		valid bool
		errs  []error
	}{
		{
			"missing everything",
			``,
			false,
			[]error{
				ErrMissingTitle,
				ErrMissingType,
				ErrMissingURL,
				ErrMissingImage,
			},
		},
		{
			"missing title",
			`
			<meta property="og:type" content="video.movie" />
			<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />
			<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />
			`,
			false,
			[]error{
				ErrMissingTitle,
			},
		},
		{
			"missing image",
			`
			<meta property="og:title" content="The Rock" />
			<meta property="og:type" content="video.movie" />
			<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />
			`,
			false,
			[]error{
				ErrMissingImage,
			},
		},
		{
			"valid",
			`
			<meta property="og:title" content="The Rock" />
			<meta property="og:type" content="video.movie" />
			<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />
			<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />
			`,
			true,
			[]error{},
		},
	}

	for _, tr := range tt {
		t.Run(tr.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tr.html))
			if err != nil {
				t.Fatal(err)
			}
			og, _ := parseOpenGraph(doc)
			valid, errs := og.Validate()

			if tr.valid != valid {
				t.Error("expected", tr.valid)
				t.Error("got     ", valid)
			}
			if !reflect.DeepEqual(tr.errs, errs) {
				t.Error("expected", tr.errs)
				t.Error("got     ", errs)
			}
		})
	}
}
