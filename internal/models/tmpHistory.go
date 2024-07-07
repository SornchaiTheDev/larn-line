package models

type TmpHistory struct {
	Message interface{}
}

type TextMessage struct {
	Text string
}

func (t *TextMessage) GetText() string {
	return t.Text
}

type ImageMessage struct {
	Preview  string
	Original string
}

func (t *ImageMessage) GetPreview() string {
	return t.Preview
}

func (t *ImageMessage) GetOriginal() string {
	return t.Original
}
