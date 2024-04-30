package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {

	filePath := flag.String("f", "./testdata/measurements.txt", "path for measurements file")
	flag.Parse()
	if *filePath == "" {
		fmt.Println("here")
		flag.Usage()
		return
	}

	err := run(*filePath, os.Stdout)
	if err != nil {
		fmt.Println("File not found at ", filePath)
		panic(err)
	}
}

type stats struct {
	min, max, sum float64
	count         int64
}

func run(filePath string, output io.Writer) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stationStats := make(map[string]stats)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		vals := strings.Split(line, ";")
		station := vals[0]
		temp, err := strconv.ParseFloat(vals[1], 64)
		if err != nil {
			return err
		}

		stat, ok := stationStats[station]
		if !ok {
			stat.min = temp
			stat.max = temp
			stat.sum = temp
			stat.count = 1
		} else {
			stat.min = math.Min(stat.min, temp)
			stat.max = math.Max(stat.max, temp)
			stat.sum += temp
			stat.count++
		}
		stationStats[station] = stat
	}

	printOutput(stationStats, output)
	return nil
}

func printOutput(stationStats map[string]stats, output io.Writer) {
	stationList := make([]string, 0, len(stationStats))
	for station := range stationStats {
		stationList = append(stationList, station)
	}
	sort.Strings(stationList)
	fmt.Fprint(output, "{")
	for i, station := range stationList {
		if i > 0 {
			fmt.Fprint(output, ", ")
		}
		s := stationStats[station]
		mean := s.sum / float64(s.count)
		mean = math.Ceil(mean*10) / 10
		fmt.Fprintf(output, "%s=%.1f/%.1f/%.1f", station, s.min, mean, s.max)
	}
	fmt.Fprint(output, "}\n")
}
