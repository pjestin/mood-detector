package util

import (
	"regexp"
	"strings"

	"github.com/pjestin/mood-detector/model/reddit"
)

var NEGATIVE_WORDS = SetFromArray([]string{"bear", "red", "bad", "police", "sell", "sells", "sold", "selling",
	"disbelief", "crisis", "dip", "dips", "dipping", "dipped", "hard", "lose", "loses", "lost", "losing", "plunge",
	"plunges", "plunged", "plunging", "shrink", "shrinks", "shrinking", "shrinked", "avoid", "avoids", "avoiding",
	"avoided", "break", "breaks", "breaking", "broke", "broken", "withdraw", "withdraws", "withdrew", "withdrawn",
	"wary", "risk", "risks", "panic", "panics", "panicking", "panicked", "issue", "issues", "drain", "drains",
	"drained", "draining", "bankruptcy", "bankruptcies", "nonsense", "shenanigan", "shenanigans", "scandal",
	"scandals", "dump", "dumped", "worse", "worst", "terrible"})

var POSITIVE_WORDS = SetFromArray([]string{"bull", "green", "good", "recover", "recovers", "recovering",
	"recovered", "amazing", "future", "believe", "believes", "believing", "believed", "unwavering", "confidence",
	"heal", "heals", "healing", "healed", "huge", "adoption", "adopt", "adopts", "adopted", "revival", "launch",
	"buy", "buys", "buying", "bought", "profit", "profits", "hodl", "hodler", "hodlers", "hold", "holds", "holding",
	"held", "safe", "save", "saves", "saving", "sure", "skyrocket", "skyrocketed", "skyrockets", "pump", "pumped",
	"euphoria", "better", "best", "great"})

func processTextMood(text string) int {
	re := regexp.MustCompile(`\W`)
	var mood int
	for _, word := range re.Split(text, -1) {
		lower := strings.ToLower(word)
		if POSITIVE_WORDS.Contains(lower) {
			mood++
		}
		if NEGATIVE_WORDS.Contains(lower) {
			mood--
		}
	}
	return mood
}

func ProcessPostMood(posts []reddit.PostData) int {
	var mood int
	for _, post := range posts {
		titleMood := processTextMood(post.Title)
		mood += 10 * titleMood
		textMood := processTextMood(post.Selftext)
		mood += textMood
	}
	return mood
}
