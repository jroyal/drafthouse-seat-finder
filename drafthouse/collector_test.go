package drafthouse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaDrafthouseCollection(t *testing.T) {
	results, err := getDrafthouseFilmMeta("the-void", "A000012990")
	assert.Nil(t, err)
	assert.Equal(t, "s3.drafthouse.com/images/made/void_ver5_xxlg_poster_240_356_81_s_c1.jpg", results.PosterURL)
	assert.NotEmpty(t, results.Description)

}
