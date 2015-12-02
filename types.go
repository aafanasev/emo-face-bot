package main

import (
	"bytes"
	"fmt"
)

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

func (face *Face) String() string {
	var buffer bytes.Buffer

	if anger := round(face.Scores.Anger * 100); anger > 0 {
		buffer.WriteString(fmt.Sprintf("Anger: %d%%\n", anger))
	}
	if contempt := round(face.Scores.Contempt * 100); contempt > 0 {
		buffer.WriteString(fmt.Sprintf("Contempt: %d%%\n", contempt))
	}
	if disgust := round(face.Scores.Disgust * 100); disgust > 0 {
		buffer.WriteString(fmt.Sprintf("Disgust: %d%%\n", disgust))
	}
	if fear := round(face.Scores.Fear * 100); fear > 0 {
		buffer.WriteString(fmt.Sprintf("Fear: %d%%\n", fear))
	}
	if happiness := round(face.Scores.Happiness * 100); happiness > 0 {
		buffer.WriteString(fmt.Sprintf("Happiness: %d%%\n", happiness))
	}
	if neutral := round(face.Scores.Neutral * 100); neutral > 0 {
		buffer.WriteString(fmt.Sprintf("Neutral: %d%%\n", neutral))
	}
	if sadness := round(face.Scores.Sadness * 100); sadness > 0 {
		buffer.WriteString(fmt.Sprintf("Sadness: %d%%\n", sadness))
	}
	if surprise := round(face.Scores.Surprise * 100); surprise > 0 {
		buffer.WriteString(fmt.Sprintf("Surprise: %d%%\n", surprise))
	}

	return buffer.String()
}
