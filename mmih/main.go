package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Funkcja func(input string) (output string, err error)

var UjemnyTenInt error = errors.New("Ujemny ten int")

func GoodInt(input string) (output string, err error) {
	u, err := strconv.Atoi(input)
	if err != nil {
		return
	}
	output = strconv.Itoa(u)
	if u < 0 {
		err = UjemnyTenInt
	}
	return
}

var re = regexp.MustCompile(`\d+(\.\d+)?`)

func AFloat(input string) (before int, after string, err error) {
	bydot := strings.Replace(input, ",", ".", 1)
	bydot = re.FindString(bydot)
	out := strings.Split(bydot, ".")
	before, err = strconv.Atoi(out[0])
	if len(out) > 1 {
		after = out[1]
	} else {
		after = "0"
	}
	if err == nil && len(after) == 0 {
		after = "0"
	}
	if err != nil {
		return
	}
	if before < 0 {
		err = UjemnyTenInt
	}
	afterint, err := strconv.Atoi(after)
	if err == nil && afterint < 0 {
		err = UjemnyTenInt
	}
	return
}

func GoodFloat(input string) (output string, err error) {
	before, after, err := AFloat(input)
	output = strconv.Itoa(before) + "." + after
	return
}

var DuzyTenPercentage error = errors.New("Duży ten percentage")

func Percentage(input string) (output string, err error) {
	before, after, err := AFloat(input)
	var theatoi int
	if err == nil {
		theatoi, err = strconv.Atoi(after)
	}
	if err == nil && before > 99 && theatoi > 0 {
		err = DuzyTenPercentage
	}
	if before == 100 {
		output = "1.000"
		return
	}
	output = fmt.Sprintf("0.%02d%s", before, after)
	return
}

func Year(input string) (output string, err error) {
	if len(input) == 0 {
		return "2017-", nil
	}
	u, err := strconv.Atoi(input)
	if err != nil {
		output = input
		return
	}
	if u < 0 {
		u = -u
	}
	if u < 5 {
		output = strconv.Itoa(2020 + u)
	} else if u < 10 {
		output = strconv.Itoa(2010 + u)
	} else if u < 1000 {
		output = strconv.Itoa(2000 + u)
	} else {
		output = strconv.Itoa(u)
	}
	output = output + "-"
	return
}

func DayMonth(input string) (output string, err error) {
	before, after, err := AFloat(input)
	u, err := strconv.Atoi(after)
	output = fmt.Sprintf("%02d-%02d ", before, u)
	if err == nil && before > 12 {
		err = DuzyTenPercentage
	}
	if err == nil && u > 31 {
		err = DuzyTenPercentage
	}
	return
}

func HourMinute(input string) (output string, err error) {
	before, after, err := AFloat(input)
	u, err := strconv.Atoi(after)
	output = fmt.Sprintf("%02d:%02d", before, u)
	if err == nil && before > 23 {
		err = DuzyTenPercentage
	}
	if err == nil && u > 59 {
		err = DuzyTenPercentage
	}
	return
}

const TABELKA = `

B 2 Obese             Overweight         Thick-set
O
D
Y 1 Lacks exercise    Balanced           Balanced muscular
/
F
A 0 Skinny            Balanced skinny    Skinny muscular
T
     \_0_/              \_1_/              \_2_/
     M   U   S   C   L   E       R   A   T   I   O

`

type Kolumna struct {
	Opis string
	Funkcja
}

var pola = []string{"Date", "Weight", "Points", "Body_fat", "BMI", "Muscle", "Water",
	"Basal metabolism", "Visceral fat", "Bone mass", "Sylwetka body", "Sylwetka muscle"}

var KOLUMNY = []Kolumna{{"Enter `done` instead of year if you want to end\nYear(some digits)", Year},
	{"Month,Day", DayMonth}, {"Hour,Minute", HourMinute},
	{"weight", GoodFloat},
	{"points", GoodInt},
	{"Body Fat %", Percentage}, {"BMI", GoodFloat}, {"Muscle Mass", GoodFloat}, {"Water %", Percentage},
	{"Basal metabolism (Cal)", GoodInt}, {"Visceral Fat Level", GoodInt}, {"Bone Mass", GoodFloat},
	{TABELKA + "Sylwetka Body Fat", GoodInt}, {"Sylwetka Muscle", GoodInt}}

func Odpytaj(reader *bufio.Reader) (output []string, juzkoniec bool) {
	output = make([]string, len(KOLUMNY))
	for i := 0; i < len(KOLUMNY); {
		v := KOLUMNY[i]
		fmt.Print(v.Opis, "  : ")
		input, err := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		if err != nil {
			fmt.Println(err)
			continue
		}
		if input == "back" || input == "--" || input == "++" {
			i--
			fmt.Println(" --- Going back!")
			continue
		}
		if input == "done" {
			juzkoniec = true
			return
		}
		ret, err := v.Funkcja(input)
		output[i] = ret
		if err != nil {
			fmt.Println("ERRONEOUS: ", ret, "***", err)
			continue

		}
		fmt.Println(ret)
		i++
	}
	return
}

func Przerob(input []string) (output []string) {
	output = make([]string, len(pola))
	output[0] = input[0] + input[1] + input[2]
	for i := 1; i < len(pola); i++ {
		output[i] = input[i+2]
	}
	return
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	f, err := os.OpenFile("themmihappendshere.csv", os.O_WRONLY|os.O_APPEND, 0600)
	c := csv.NewWriter(f)
	c.Write(pola)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for {
		odpytanie, koniec := Odpytaj(reader)
		if koniec {
			break
		}
		przerobione := Przerob(odpytanie)
		fmt.Println(przerobione)
		c.Write(przerobione)
	}
	c.Flush()
}
