package drafthouse

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type SeatResponse struct {
	ResponseCode     int            `json:"ResponseCode"`
	SeatLayoutData   SeatLayoutData `json:"SeatLayoutData"`
	ErrorDescription string         `json:"ErrorDescription"`
}

type SeatLayoutData struct {
	AreaCategories []AreaCategory `json:"AreaCategories"`
	Areas          []Area         `json:"Areas"`
	BoundaryLeft   int            `json:"BoundaryLeft"`
	BoundaryRight  int            `json:"BoundaryRight"`
	BoundaryTop    float64        `json:"BoundaryTop"`
	ScreenStart    int            `json:"ScreenStart"`
	ScreenWidth    int            `json:"ScreenWidth"`
}

type AreaCategory struct {
	AreaCategoryCode        string        `json:"AreaCategoryCode"`
	IsInSeatDeliveryEnabled bool          `json:"IsInSeatDeliveryEnabled"`
	SeatsAllocatedCount     int           `json:"SeatsAllocatedCount"`
	SeatsNotAllocatedCount  int           `json:"SeatsNotAllocatedCount"`
	SeatsToAllocate         int           `json:"SeatsToAllocate"`
	SelectedSeats           []interface{} `json:"SelectedSeats"`
}

type Area struct {
	AreaCategoryCode      string  `json:"AreaCategoryCode"`
	ColumnCount           int     `json:"ColumnCount"`
	Description           string  `json:"Description"`
	DescriptionAlt        string  `json:"DescriptionAlt"`
	HasSofaSeatingEnabled bool    `json:"HasSofaSeatingEnabled"`
	Height                float64 `json:"Height"`
	IsAllocatedSeating    bool    `json:"IsAllocatedSeating"`
	Left                  int     `json:"Left"`
	Number                int     `json:"Number"`
	NumberOfSeats         int     `json:"NumberOfSeats"`
	RowCount              int     `json:"RowCount"`
	Rows                  []Row   `json:"Rows"`
	Top                   float64 `json:"Top"`
	Width                 int     `json:"Width"`
}

type Row struct {
	PhysicalName string  `json:"PhysicalName"`
	Seats        []Seats `json:"Seats"`
}

type Seats struct {
	ID             string         `json:"Id"`
	OriginalStatus int            `json:"OriginalStatus"`
	Position       Position       `json:"Position"`
	Priority       int            `json:"Priority"`
	SeatStyle      int            `json:"SeatStyle"`
	SeatsInGroup   []SeatsInGroup `json:"SeatsInGroup"`
	Status         int            `json:"Status"`
}

type SeatsInGroup struct {
	AreaNumber  int `json:"AreaNumber"`
	ColumnIndex int `json:"ColumnIndex"`
	RowIndex    int `json:"RowIndex"`
}

type Position struct {
	AreaNumber  int `json:"AreaNumber"`
	ColumnIndex int `json:"ColumnIndex"`
	RowIndex    int `json:"RowIndex"`
}

type Chart [][]string

const (
	tableEnd        string = "</table>"
	rowStart        string = "<tr>"
	rowEnd          string = "</tr>"
	emptySpot       string = "<td class=\"empty\"></td>"
	seatOpen        string = "<td class=\"open\"></td>"
	seatTaken       string = "<td class=\"taken\"></td>"
	seatConvertable string = "<td class=\"convertable\"></td>"
	seatHandicap    string = "<td class=\"handicap\"></td>"
)

func fillEmptySpots(seatChart Chart) Chart {
	for r := range seatChart {
		for c := range seatChart[r] {
			if seatChart[r][c] == "" {
				seatChart[r][c] = emptySpot
			}
		}
	}
	return seatChart
}

func buildHTMLTable(cinemaName string, chart Chart) template.HTML {
	var seatChart bytes.Buffer
	seatChart.WriteString(fmt.Sprintf("<table class=\"center-block cinema-%s\">", strings.Replace(strings.ToLower(cinemaName), " ", "", -1)))
	for r := range chart {
		seatChart.WriteString(rowStart)
		for c := range chart[r] {
			if chart[r][c] == "" {
				seatChart.WriteString(emptySpot)
			} else {
				seatChart.WriteString(chart[r][c])
			}

		}
		seatChart.WriteString(rowEnd)
	}
	seatChart.WriteString(tableEnd)
	return template.HTML(seatChart.String())
}

func addScreen(seatChart Chart, width int) Chart {
	seatChart[0] = []string{
		fmt.Sprintf("<td colspan=%d class=\"screen\"></td>", width),
	}
	return seatChart
}

func setValue(seatChart Chart, height, width, drafthouseRow, drafthouseColumn int, val string) Chart {
	realRow := height - drafthouseRow - 2
	realCol := width - drafthouseColumn - 1
	seatChart[realRow][realCol] = val
	return seatChart
}

func NewSeatChart(cinemaName string, baseUrl string, res SeatResponse) template.HTML {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic")
			fmt.Println(res.SeatLayoutData)
		}
	}()

	if res.ResponseCode == 67 || res.ResponseCode == 1 {
		log.WithFields(log.Fields{
			"ErrorDescription": res.ErrorDescription,
			"Response code":    res.ResponseCode,
		}).Info("Failed to get seat data")
		return template.HTML(fmt.Sprintf("<img src=\"%s/images/sold_out.svg\">", baseUrl))
	}
	columnCount := res.SeatLayoutData.Areas[0].ColumnCount

	// Add 3 rows. 2 for the front one for the back
	rowCount := res.SeatLayoutData.Areas[0].RowCount + 3
	chart := make([][]string, rowCount)
	for i := range chart {
		chart[i] = make([]string, columnCount)
	}

	rows := res.SeatLayoutData.Areas[0].Rows
	for rowNum, row := range rows {
		for _, seat := range row.Seats {
			colNum := seat.Position.ColumnIndex
			var val string
			switch seat.Status {
			case 0:
				val = seatOpen
			case 1:
				val = seatTaken
			case 7:
				val = seatConvertable
			case 3:
				val = seatHandicap
			}
			chart = setValue(chart, rowCount, columnCount, rowNum, colNum, val)
		}
	}
	chart = addScreen(chart, columnCount)

	return buildHTMLTable(cinemaName, chart)
}

func loadFilmSeats(filmSessions []FilmSession, baseUrl string) {
	var wg sync.WaitGroup

	for i := range filmSessions {
		wg.Add(1)
		filmSession := &filmSessions[i]
		go func(filmSession *FilmSession) {
			seatResponse := getFilmSeats(*filmSession)
			filmSession.SeatChart = NewSeatChart(filmSession.CinemaName, baseUrl, seatResponse)
			wg.Done()
		}(filmSession)

	}
	wg.Wait()
}

func sortFilmSessions(filmSessions []FilmSession) map[string][]FilmSession {
	cinemas := make(map[string][]FilmSession)
	for _, session := range filmSessions {
		if filmSlice, ok := cinemas[session.CinemaName]; ok {
			cinemas[session.CinemaName] = append(filmSlice, session)
		} else {
			cinemas[session.CinemaName] = []FilmSession{session}
		}
	}
	return cinemas
}
