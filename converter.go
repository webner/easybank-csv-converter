package main
import ("fmt"
	"os"
	"encoding/csv"
	"log"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"regexp"
	"strings"
	"strconv"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Usage: easybank-csv-converter filename.csv")
		os.Exit(0);
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Can't open file %v\nError: %v\n", os.Args[1], err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(transform.NewReader(file, charmap.Windows1252.NewDecoder()))
	reader.Comma = ';'
	record, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Can't read csv file %v\nError: %v\n", os.Args[1], err.Error())
	}

	writer := csv.NewWriter(os.Stdout)

	descriptionRegex, err := regexp.Compile("(.*[^\\W]*)\\W*([A-Z]{2}/[0-9]{9})\\W*(.*)")
	if err != nil {
		log.Fatal("Can't parse regex %v", err.Error())
	}

	for _, row := range record {
		account := row[0]
		description := row[1]
		date := row[2]
		amountGer := row[4]
		amount, err := strconv.ParseFloat(strings.Replace(strings.Replace(amountGer, ".", "", -1), ",", ".", -1), 64)
		if err != nil {
			log.Fatalf("Unable to parse float %v", row[4])
		}
		currency := row[5]

		groups := descriptionRegex.FindStringSubmatch(description)
		if len(groups) == 0 {
			log.Fatal("Can't parse description")
		}

		memo := strings.TrimSpace(groups[1])
		nr := groups[2]
		payee := groups[3]

		writer.Write([]string {account, nr, memo, payee, date, strconv.FormatFloat(amount, 'f', 2, 64), currency})
	}

	writer.Flush()
}
