package subpub

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBasicSubscribePublish(t *testing.T) {
	sp, err := NewSubPub()
	if err != nil {
		t.Fatalf("Ошибка создания SubPub: %v", err)
	}
	received := make(chan interface{}, 1)

	sub, err := sp.Subscribe("тест", func(msg interface{}) {
		received <- msg
	})
	if err != nil {
		t.Fatalf("Ошибка подписки: %v", err)
	}
	defer sub.Unsubscribe()

	// Публикуем сообщение
	err = sp.Publish("тест", "тестовое сообщение")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}

	// Проверяем получение
	select {
	case msg := <-received:
		if msg != "тестовое сообщение" {
			t.Errorf("Получено неверное сообщение: %v", msg)
		}
	case <-time.After(time.Second):
		t.Error("Таймаут ожидания сообщения")
	}
}

func TestMultipleSubscribers(t *testing.T) {
	sp, err := NewSubPub()
	if err != nil {
		t.Fatalf("Ошибка создания SubPub: %v", err)
	}
	var wg sync.WaitGroup
	subscriberCount := 10
	messageCount := 0
	var mu sync.Mutex

	// Создаем подписчиков
	for i := 0; i < subscriberCount; i++ {
		wg.Add(1)
		sub, err := sp.Subscribe("тест", func(msg interface{}) {
			mu.Lock()
			messageCount++
			mu.Unlock()
			wg.Done()
		})
		if err != nil {
			t.Fatalf("Ошибка подписки: %v", err)
		}
		defer sub.Unsubscribe()
	}

	// Публикуем сообщение
	err = sp.Publish("тест", "тестовое сообщение")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}

	// Ждем получения всеми подписчиками
	wg.Wait()

	if messageCount != subscriberCount {
		t.Errorf("Получено %d сообщений, ожидалось %d", messageCount, subscriberCount)
	}
}

func TestUnsubscribe(t *testing.T) {
	sp, err := NewSubPub()
	if err != nil {
		t.Fatalf("Ошибка создания SubPub: %v", err)
	}
	received := make(chan interface{}, 1)

	sub, err := sp.Subscribe("тест", func(msg interface{}) {
		received <- msg
	})
	if err != nil {
		t.Fatalf("Ошибка подписки: %v", err)
	}

	// Отписываемся
	sub.Unsubscribe()

	// Публикуем сообщение
	err = sp.Publish("тест", "тестовое сообщение")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}

	// Проверяем, что сообщение не получено
	select {
	case msg := <-received:
		t.Errorf("Получено сообщение после отписки: %v", msg)
	case <-time.After(100 * time.Millisecond):
		// Это ожидаемое поведение
	}
}

func TestClose(t *testing.T) {
	sp, err := NewSubPub()
	if err != nil {
		t.Fatalf("Ошибка создания SubPub: %v", err)
	}
	ctx := context.Background()

	// Подписываемся
	_, err = sp.Subscribe("тест", func(msg interface{}) {})
	if err != nil {
		t.Fatalf("Ошибка подписки: %v", err)
	}

	// Закрываем систему
	err = sp.Close(ctx)
	if err != nil {
		t.Fatalf("Ошибка закрытия: %v", err)
	}

	// Проверяем, что нельзя публиковать
	err = sp.Publish("тест", "тестовое сообщение")
	if err != context.Canceled {
		t.Errorf("Ожидалась ошибка context.Canceled, получено: %v", err)
	}

	// Проверяем, что нельзя подписаться
	_, err = sp.Subscribe("тест", func(msg interface{}) {})
	if err != context.Canceled {
		t.Errorf("Ожидалась ошибка context.Canceled, получено: %v", err)
	}
}
