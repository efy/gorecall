package webinfo

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrMissingTitle = fmt.Errorf("opengraph: missing required attribute: title")
	ErrMissingType  = fmt.Errorf("opengraph: missing required attribute: type")
	ErrMissingURL   = fmt.Errorf("opengraph: missing required attribute: url")
	ErrMissingImage = fmt.Errorf("opengraph: missing required attribute: image")
)

// Type open graph contains fields that map to the opengraph attributes
// described at http://ogp.me/
type OG struct {
	// Basic metdata
	Title string   `json:"title"`
	Type  string   `json:"type"`
	URL   string   `json:"url"`
	Image *OGImage `json:"image"`

	// Optional
	Description     string   `json:"description"`
	Determiner      string   `json:"determiner"`
	Locale          string   `json:"locale"`
	LocaleAlternate []string `json:"locale_alternate"`
	SiteName        string   `json:"site_name"`
	Audio           *OGAudio `json:"audio"`
	Video           *OGVideo `json:"video"`
}

type OGImage struct {
	Source string `json:"source"`

	// Optional
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"`
	Width     string `json:"width"`
	Height    string `json:"height"`
	Alt       string `json:"alt"`
}

type OGAudio struct {
	Source string `json:"source"`

	// Optional
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"`
}

type OGVideo struct {
	Source string `json:"source"`

	// Optional
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"`
	Width     string `json:"width"`
	Height    string `json:"height"`
}

func (og *OG) Validate() (bool, []error) {
	errs := make([]error, 0)
	valid := true

	if og.Title == "" {
		errs = append(errs, ErrMissingTitle)
		valid = false
	}

	if og.Type == "" {
		errs = append(errs, ErrMissingType)
		valid = false
	}

	if og.URL == "" {
		errs = append(errs, ErrMissingURL)
		valid = false
	}

	if og.Image == nil {
		errs = append(errs, ErrMissingImage)
		valid = false
	}

	return valid, errs
}

func parseOpenGraph(doc *goquery.Document) (*OG, error) {
	og := &OG{}

	ogtitle, exists := doc.Find(`meta[property="og:title"]`).First().Attr("content")
	if exists {
		og.Title = ogtitle
	}

	ogtype, exists := doc.Find(`meta[property="og:type"]`).First().Attr("content")
	if exists {
		og.Type = ogtype
	}

	ogurl, exists := doc.Find(`meta[property="og:url"]`).First().Attr("content")
	if exists {
		og.URL = ogurl
	}

	ogdesc, exists := doc.Find(`meta[property="og:description"]`).First().Attr("content")
	if exists {
		og.Description = ogdesc
	}

	oglocale, exists := doc.Find(`meta[property="og:locale"]`).First().Attr("content")
	if exists {
		og.Locale = oglocale
	}

	ogdeterminer, exists := doc.Find(`meta[property="og:determiner"]`).First().Attr("content")
	if exists {
		og.Determiner = ogdeterminer
	}

	ogname, exists := doc.Find(`meta[property="og:site_name"]`).First().Attr("content")
	if exists {
		og.SiteName = ogname
	}

	ogimgsource, exists := doc.Find(`meta[property="og:image"]`).First().Attr("content")
	if exists {
		ogimg := &OGImage{
			Source: ogimgsource,
		}

		imgtype, exists := doc.Find(`meta[property="og:image:type"]`).First().Attr("content")
		if exists {
			ogimg.Type = imgtype
		}

		width, exists := doc.Find(`meta[property="og:image:width"]`).First().Attr("content")
		if exists {
			ogimg.Width = width
		}

		height, exists := doc.Find(`meta[property="og:image:height"]`).First().Attr("content")
		if exists {
			ogimg.Height = height
		}

		alt, exists := doc.Find(`meta[property="og:image:alt"]`).First().Attr("content")
		if exists {
			ogimg.Alt = alt
		}

		secureurl, exists := doc.Find(`meta[property="og:image:secure_url"]`).First().Attr("content")
		if exists {
			ogimg.SecureURL = secureurl
		}

		og.Image = ogimg
	}

	ogvideosource, exists := doc.Find(`meta[property="og:video"]`).First().Attr("content")
	if exists {
		ogvideo := &OGVideo{
			Source: ogvideosource,
		}

		vtype, exists := doc.Find(`meta[property="og:video:type"]`).First().Attr("content")
		if exists {
			ogvideo.Type = vtype
		}

		width, exists := doc.Find(`meta[property="og:video:width"]`).First().Attr("content")
		if exists {
			ogvideo.Width = width
		}

		height, exists := doc.Find(`meta[property="og:video:height"]`).First().Attr("content")
		if exists {
			ogvideo.Height = height
		}

		secureurl, exists := doc.Find(`meta[property="og:video:secure_url"]`).First().Attr("content")
		if exists {
			ogvideo.SecureURL = secureurl
		}

		og.Video = ogvideo
	}

	ogaudiosource, exists := doc.Find(`meta[property="og:audio"]`).First().Attr("content")
	if exists {
		ogaudio := &OGAudio{
			Source: ogaudiosource,
		}

		atype, exists := doc.Find(`meta[property="og:audio:type"]`).First().Attr("content")
		if exists {
			ogaudio.Type = atype
		}

		secureurl, exists := doc.Find(`meta[property="og:audio:secure_url"]`).First().Attr("content")
		if exists {
			ogaudio.SecureURL = secureurl
		}

		og.Audio = ogaudio
	}

	valid, _ := og.Validate()

	if !valid {
		return og, fmt.Errorf("opengraph: invalid object")
	}

	return og, nil
}
