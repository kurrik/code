package main

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Config struct {
	Verbose  bool
	TestWord string
	Count    int
}

var CONFIG Config = Config{}

type WordScore struct {
	Word          string
	Score         float64
	TotalSubwords int
	ValidSubwords int
}

func ReadWords(words map[string]bool) (count int) {
	var word string
	count = len(words)
	for {
		if _, err := fmt.Scanf("%v\n", &word); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(fmt.Sprintf("Error reading: %v", err))
			}
		}
		word = strings.TrimSuffix(word, "%")
		word = strings.TrimSuffix(word, "!")
		if _, exists := words[word]; !exists {
			words[word] = true
			count += 1
		}
	}
	return
}

func GetWordSubsets(word string) (subsets []string) {
	var (
		i int
		j int
	)
	for i = 1; i < len(word); i++ {
		for j = 0; j <= len(word)-i; j++ {
			subsets = append(subsets, word[j:j+i])
		}
	}
	sort.Slice(subsets, func(i, j int) bool {
		return subsets[i] < subsets[j]
	})
	return
}

func GetValidCount(words []string, dict map[string]bool) (count int) {
	count = 0
	for i := 0; i < len(words); i++ {
		if _, exists := dict[words[i]]; exists {
			count++
			if CONFIG.Verbose {
				fmt.Printf("- %40v: valid\n", words[i])
			}
		} else {
			if CONFIG.Verbose {
				fmt.Printf("- %40v: invalid\n", words[i])
			}
		}
	}
	return count
}

func ScoreWord(word string, dict map[string]bool) (score WordScore) {
	var (
		subsets []string
	)
	if CONFIG.Verbose {
		fmt.Printf("# %v:\n", word)
	}
	score.Word = word
	subsets = GetWordSubsets(word)
	score.TotalSubwords = len(subsets)
	score.ValidSubwords = GetValidCount(subsets, dict)
	if score.TotalSubwords == 0 {
		score.Score = 0
	} else {
		score.Score = float64(score.ValidSubwords) / float64(score.TotalSubwords)
	}
	return
}

func PrintScore(score WordScore) {
	fmt.Printf("%v, %v/%v valid (%2.4v%%)\n",
		score.Word,
		score.ValidSubwords,
		score.TotalSubwords,
		score.Score*100.0,
	)
}

func main() {
	var (
		dict  map[string]bool
		count int
	)
	flag.BoolVar(&CONFIG.Verbose, "v", false, "Verbose")
	flag.StringVar(&CONFIG.TestWord, "testword", "", "Test with this word")
	flag.IntVar(&CONFIG.Count, "count", 100, "Number of results to print")
	flag.Parse()
	// Start with some words which aren't present in the given dictionary.
	dict = map[string]bool{
		"a": true,
		"i": true,
	}
	count = ReadWords(dict)
	fmt.Printf("Words read %v\n", count)
	if CONFIG.TestWord != "" {
		PrintScore(ScoreWord(CONFIG.TestWord, dict))
	} else {
		scores := make([]WordScore, count)
		i := 0
		for word, _ := range dict {
			scores[i] = ScoreWord(word, dict)
			i++
		}
		sort.Slice(scores, func(i, j int) bool {
			if scores[i].Score == scores[j].Score {
				if scores[i].ValidSubwords == scores[j].ValidSubwords {
					// Tiebreaker 2: alphabetic order.
					return scores[i].Word < scores[j].Word
				}
				// Tiebreaker 1: Number of valid subwords.
				return scores[i].ValidSubwords > scores[j].ValidSubwords
			}
			// Comparison: Highest score.
			return scores[i].Score > scores[j].Score
		})
		for i = 0; i < CONFIG.Count; i++ {
			fmt.Printf("%2v.) ", i)
			PrintScore(scores[i])
		}
	}
}
