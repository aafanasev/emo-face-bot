package main

import "fmt"

// Face - emotion
type Face struct {
	FaceRectangle FaceRectanle
	Scores        Scores
}

// FaceRectanle - coordinates
type FaceRectanle struct {
	Left   int
	Top    int
	Width  int
	Height int
}

// Scores - values
type Scores struct {
	Anger     float32
	Contempt  float32
	Disgust   float32
	Fear      float32
	Happiness float32
	Neutral   float32
	Sadness   float32
	Surprise  float32
}

func (f *Face) String() string {
	return fmt.Sprintf("Anger: %.2f\nContempt: %.2f\nDisgust: %.2f\nFear: %.2f\nHappiness: %.2f\nNeutral: %.2f\nSadness: %.2f\nSurprise: %.2f", f.Scores.Anger, f.Scores.Contempt, f.Scores.Disgust, f.Scores.Fear, f.Scores.Happiness, f.Scores.Neutral, f.Scores.Sadness, f.Scores.Surprise)
}
