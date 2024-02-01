package game

type Emotion int

const (
	EmotionNeutral Emotion = iota

	EmotionPlacidity

	EmotionHappiness
	EmotionJoy
	EmotionElation

	EmotionExcitement

	EmotionSadness
	EmotionTrauma
	EmotionFear

	EmotionAnger
	EmotionHatred
)

type Emotioner interface {
	AddEmotion(emotion Emotion, amount float64)
	SetEmotion(emotion Emotion, amount float64)
	GetEmotion(emotion Emotion) float64
	CloneEmotions() Emotions
}

type Emotions struct {
	emotions EmotionState
}

type EmotionState map[Emotion]float64

func (e *Emotions) AddEmotion(emotion Emotion, amount float64) {
	e.emotions[emotion] += amount
	if e.emotions[emotion] < 0 {
		e.emotions[emotion] = 0
	}
}

func (e *Emotions) SetEmotion(emotion Emotion, amount float64) {
	e.emotions[emotion] = amount
	if e.emotions[emotion] < 0 {
		e.emotions[emotion] = 0
	}
}

func (e *Emotions) GetEmotion(emotion Emotion) float64 {
	return e.emotions[emotion]
}

func (e *Emotions) CloneEmotions() Emotions {
	clone := Emotions{}
	for k, v := range e.emotions {
		clone.emotions[k] = v
	}
	return clone
}
