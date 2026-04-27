// utils/cvv_test.go
package utils

import (
	"testing"
)

func TestGenerateCVV(t *testing.T) {
	// Тест 1: проверка формата
	cvv, err := GenerateCVV()
	if err != nil {
		t.Fatalf("GenerateCVV failed: %v", err)
	}

	if len(cvv) != 3 {
		t.Errorf("CVV length should be 3, got %d", len(cvv))
	}

	// Тест 2: проверка что все символы - цифры
	for _, c := range cvv {
		if c < '0' || c > '9' {
			t.Errorf("CVV should contain only digits, got %c", c)
		}
	}

	// Тест 3: проверка диапазона 000-999
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		cvv, err := GenerateCVV()
		if err != nil {
			t.Fatalf("GenerateCVV failed: %v", err)
		}

		// Проверяем что в диапазоне
		if cvv < "000" || cvv > "999" {
			t.Errorf("CVV out of range: %s", cvv)
		}

		seen[cvv] = true
	}

	// Ожидаем хотя бы 500 уникальных из 1000 (вероятность коллизий мала)
	if len(seen) < 500 {
		t.Errorf("CVV generation not random enough: only %d unique out of 1000", len(seen))
	}

	t.Logf("Generated %d unique CVVs out of 1000", len(seen))
}

func TestGenerateCVVUniqueness(t *testing.T) {
	// Генерируем 10000 CVV и проверяем распределение
	counts := make(map[string]int)
	total := 10000

	for i := 0; i < total; i++ {
		cvv, err := GenerateCVV()
		if err != nil {
			t.Fatalf("GenerateCVV failed: %v", err)
		}
		counts[cvv]++
	}

	// Теоретически равномерное распределение - 10000/1000 = 10 на каждое значение
	// Допускаем отклонение в 5 раз
	for cvv, count := range counts {
		if count > 50 {
			t.Errorf("CVV %s appears too frequently: %d times", cvv, count)
		}
	}

	t.Logf("Distribution test passed with %d unique values", len(counts))
}

func BenchmarkGenerateCVV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenerateCVV()
		if err != nil {
			b.Fatalf("GenerateCVV failed: %v", err)
		}
	}
}
