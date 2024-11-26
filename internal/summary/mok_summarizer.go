package summary

type MockSummarizer struct{}

func NewMockSummarizer() *MockSummarizer {
	return &MockSummarizer{}
}

func (m *MockSummarizer) Summarize(text string) (string, error) {
	// Просто возвращаем предсказуемый результат для тестирования
	return "Это мок-суммаризация текста", nil
}
