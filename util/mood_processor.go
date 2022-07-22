package util

import (
	"log"
	"regexp"
	"strings"

	"github.com/pjestin/mood-detector/model/reddit"
)

var NEGATIVE_WORDS = SetFromArray([]string{"bear", "red", "bad", "police", "sell", "sells", "sold", "selling",
	"disbelief", "crisis", "dip", "dips", "dipping", "dipped", "hard", "lose", "loses", "lost", "losing", "plunge",
	"plunges", "plunged", "plunging", "shrink", "shrinks", "shrinking", "shrinked", "avoid", "avoids", "avoiding",
	"avoided", "break", "breaks", "breaking", "broke", "broken", "withdraw", "withdraws", "withdrew", "withdrawn",
	"wary", "risk", "risks", "panic", "panics", "panicking", "panicked", "issue", "issues", "drain", "drains",
	"drained", "draining"})

var POSITIVE_WORDS = SetFromArray([]string{"bull", "green", "good", "recover", "recovers", "recovering",
	"recovered", "amazing", "future", "believe", "believes", "believing", "believed", "unwavering", "confidence",
	"heal", "heals", "healing", "healed", "huge", "adoption", "adopt", "adopts", "adopted", "revival", "launch",
	"buy", "buys", "buying", "bought", "profit", "profits", "hodl", "hodler", "hodlers", "hold", "holds", "holding",
	"held", "safe", "save", "saves", "saving", "sure", "skyrocket", "skyrocketed", "skyrockets"})

func processTextMood(text string) int {
	re := regexp.MustCompile(`\W`)
	var mood int
	for _, word := range re.Split(text, -1) {
		lower := strings.ToLower(word)
		if POSITIVE_WORDS.Contains(lower) {
			log.Println("Positive word:", lower)
			mood++
		}
		if NEGATIVE_WORDS.Contains(lower) {
			log.Println("Negative word:", lower)
			mood--
		}
	}
	return mood
}

func ProcessPostMood(posts []reddit.PostData) int {
	log.Println("Processing post mood; number of posts:", len(posts))
	var mood int
	for _, post := range posts {
		titleMood := processTextMood(post.Title)
		mood += 10 * titleMood
		textMood := processTextMood(post.Selftext)
		mood += textMood
		log.Println("Title mood:", titleMood, "; Text mood:", textMood)
	}
	return mood
}
