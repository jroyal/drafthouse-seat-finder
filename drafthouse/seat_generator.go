package drafthouse

import (
	"fmt"
	"html/template"
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

type SeatChart struct {
	width  int
	height int
	Charts [][]template.HTML
}

func (s *SeatChart) fillEmptySpots() {
	for r := range s.Charts {
		for c := range s.Charts[r] {
			if s.Charts[r][c] == "" {
				s.Charts[r][c] = template.HTML("<td class=\"empty\"></td>")
			}
		}
	}
}

func (s *SeatChart) addScreen() {
	s.Charts[0] = []template.HTML{
		template.HTML(fmt.Sprintf("<td colspan=%d class=\"screen\"></td>", s.width)),
	}
}

func setValue(seatChart *SeatChart, drafthouseRow, drafthouseColumn int, val template.HTML) {
	realRow := seatChart.height - drafthouseRow - 2
	realCol := seatChart.width - drafthouseColumn - 1
	seatChart.Charts[realRow][realCol] = val
}

func NewSeatChart(res SeatResponse) SeatChart {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic")
			fmt.Println(res.SeatLayoutData)
		}
	}()
	if res.ResponseCode == 67 {
		log.WithFields(log.Fields{
			"ErrorDescription": res.ErrorDescription,
			"Response code":    res.ResponseCode,
		}).Info("Failed to get seat data")
	}
	columnCount := res.SeatLayoutData.Areas[0].ColumnCount

	// Add 3 rows. 2 for the front one for the back
	rowCount := res.SeatLayoutData.Areas[0].RowCount + 3
	chart := make([][]template.HTML, rowCount)
	for i := range chart {
		chart[i] = make([]template.HTML, columnCount)
	}
	seatChart := SeatChart{
		Charts: chart,
		height: rowCount,
		width:  columnCount,
	}

	rows := res.SeatLayoutData.Areas[0].Rows
	for rowNum, row := range rows {
		for _, seat := range row.Seats {
			colNum := seat.Position.ColumnIndex
			var val string
			switch seat.Status {
			case 0:
				val = "<td class=\"open\"></td>"
			case 1:
				val = "<td class=\"taken\"></td>"
			case 7:
				val = "<td class=\"convertable\"></td>"
			case 3:
				val = "<td class=\"handicap\"></td>"
			}
			setValue(&seatChart, rowNum, colNum, template.HTML(val))
		}
	}

	seatChart.fillEmptySpots()
	seatChart.addScreen()
	return seatChart
}

func loadFilmSeats(filmSessions []FilmSession) []SeatChart {
	var seatCharts []SeatChart
	var wg sync.WaitGroup

	for i := range filmSessions {
		wg.Add(1)
		filmSession := &filmSessions[i]
		go func(filmSession *FilmSession) {
			seatResponse := getFilmSeats(*filmSession)
			filmSession.SeatChart = NewSeatChart(seatResponse)
			wg.Done()
		}(filmSession)

	}
	wg.Wait()
	return seatCharts
}
