package property

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/text/width"
)

type ImportRequest struct {
	FileHandle io.ReadCloser
}

type ImportResponse struct {
	Message string `json:"message"`
}

type Importer interface {
	Import(context.Context, *ImportRequest) (*ImportResponse, error)
}

func NewImporter(repo PropertyRepository) *ImporterService {
	return &ImporterService{
		repo: repo,
	}
}

type ImporterService struct {
	repo PropertyRepository
}

func (h *ImporterService) Import(ctx context.Context, req *ImportRequest) (*ImportResponse, error) {
	// cleanup file handle object from memory after request is finished,
	// so to avoid mem leak.
	defer req.FileHandle.Close()

	csvReader := csv.NewReader(req.FileHandle)

	i := 0
	for {
		// read line by line
		row, errRow := csvReader.Read()
		if errRow != nil {
			if !errors.Is(errRow, io.EOF) {
				// if an unexpected err occurred when reading the line,
				// log this err:
				log.Println(errRow)
			}
			// finished reading till the end of file, stopping.
			break
		}

		// skip the first line, which is a CSV header:
		if i == 0 {
			i++
			continue
		}

		// we only process rows with complete columns.
		if len(row) != 11 {
			continue
		}

		propertyData, errParseRow := parsePropertyDataFromRow(row)
		if errParseRow != nil {
			log.Printf("Error parsing row %s due to %s, this row is skipped.\n", strings.Join(row, ","), errParseRow.Error())
			continue
		}

		// insert to datastore.
		// possible future improvement:
		// - batch insert, instead of 1-1, so to minimize network roundtrips.
		// - write in background process using some sort of message bus,
		//   so that the customer who accessed the endpoint shall not wait for slow
		//   synchronous process which might fail if the internet connection is bad.
		if errInsert := h.repo.Insert(ctx, propertyData); errInsert != nil {
			// due to time constraint, we will only log the failure.
			log.Printf("Failed to insert %s due to %s.\n", FullAddress(propertyData), errInsert.Error())
			// continue to next row, hoping the next row will successfully inserted.
			// we might apply strategies to level-up the resiliency here,
			// for example: adding retry with backoff
			continue
		}

		log.Printf("Successfully inserted %s.\n", FullAddress(propertyData))

		i++
	}

	msg := fmt.Sprintf("Processed %d properties.", i-1)

	return &ImportResponse{
		Message: msg,
	}, nil
}

// atoi converts string to int32.
// On normal occassion, this function is basically just
// a thin layer for [strconv.Atoi], the only difference
// is that this function is able to parse string with full-width
// characters that might be prevalent in Japanese input.
func atoi(str string) (int32, error) {
	intVal, errConvert := strconv.Atoi(str)
	if errConvert == nil {
		return int32(intVal), nil
	}

	// we fail to parse using narrow-width.
	// now it's time to convert str to narrow-width chars.
	narrowGoStr := width.Narrow.String(str)
	intValNarrow, errConvertNarrow := strconv.Atoi(narrowGoStr)
	if errConvertNarrow != nil {
		// we still unable to parse, nothing we can do
		// other than rejecting the input:
		return 0, errConvertNarrow
	}

	return int32(intValNarrow), nil
}

// parsePropertyDataFromRow will parse each CSV row using best-effort strategy.
// Returns error when any of the pre-set conditions is not met.
//
// Expected column index with their respective data:
// 0 = prefecture
// 1 = city
// 2 = town
// 3 = chome
// 4 = banchi
// 5 = go
// 6 = building name
// 7 = price
// 8 = nearest station
// 9 = property type
// 10 = land area
func parsePropertyDataFromRow(row []string) (PropertyData, error) {
	chome, errChome := atoi(row[3])
	if errChome != nil {
		return PropertyData{}, fmt.Errorf("failed to parse chome: %w", errChome)
	}

	banchi, errBanchi := atoi(row[4])
	if errBanchi != nil {
		return PropertyData{}, fmt.Errorf("failed to parse banchi: %w", errBanchi)
	}

	goVal, errGo := atoi(row[5])
	if errGo != nil {
		return PropertyData{}, fmt.Errorf("failed to parse go: %w", errGo)
	}

	price, errPrice := strconv.ParseInt(row[7], 10, 0)
	if errPrice != nil {
		return PropertyData{}, fmt.Errorf("failed to parse price: %w", errPrice)
	}

	return PropertyData{
		Prefecture:     row[0],
		City:           row[1],
		Town:           row[2],
		Chome:          chome,
		Banchi:         banchi,
		Go:             goVal,
		Building:       row[6],
		Price:          price,
		NearestStation: row[8],
		PropertyType:   row[9],
		LandArea:       row[10],
	}, nil
}

// PropertyData represents a single entity of Property.
//
// For the sake of time constraint given,
// I am using this struct end to end from representing Request,
// inserting into Datastore, reading from the Datastore, and
// also to represent the endpoint Response.
//
// In an ideal situation, we must decouple data structure
// according to where it was used, so to prevent the abstraction
// from leaking everywhere.
type PropertyData struct {
	ID             int32  `db:"property_id" json:"id"`
	FullAddress    string `db:"-" json:"full_address"`
	Prefecture     string `db:"prefecture" json:"prefecture"`
	City           string `db:"city" json:"city"`
	Town           string `db:"town" json:"town"`
	Chome          int32  `db:"chome" json:"chome"`   // 丁目
	Banchi         int32  `db:"banchi" json:"banchi"` // 番地
	Go             int32  `db:"go" json:"go"`         // 号
	Building       string `db:"building" json:"building"`
	Price          int64  `db:"price" json:"price"`
	NearestStation string `db:"nearest_station" json:"nearest_station"`
	PropertyType   string `db:"property_type" json:"property_type"`
	LandArea       string `db:"land_area" json:"land_area"`
}

// FullAddress formats a PropertyData into a contiguous string.
// Pardon me that I don't know how the actual format for Full Address in Japan,
// so this is my best take on this for today.
func FullAddress(pd PropertyData) string {
	return fmt.Sprintf("%s %s %s, %d %d %d, %s", pd.Prefecture, pd.City, pd.Town, pd.Chome, pd.Banchi, pd.Go, pd.Building)
}
