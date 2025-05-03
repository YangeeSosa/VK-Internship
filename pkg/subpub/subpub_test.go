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

	err = sp.Publish("тест", "тестовое сообщение")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}

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

	err = sp.Publish("тест", "тестовое сообщение")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}

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

	sub.Unsubscribe()

	err = sp.Publish("тест", "тестовое сообщение")
	if err != nil {
		t.Fatalf("Ошибка публикации: %v", err)
	}

	select {
	case msg := <-received:
		t.Errorf("Получено сообщение после отписки: %v", msg)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestClose(t *testing.T) {
	sp, err := NewSubPub()
	if err != nil {
		t.Fatalf("Ошибка создания SubPub: %v", err)
	}
	ctx := context.Background()

	_, err = sp.Subscribe("тест", func(msg interface{}) {})
	if err != nil {
		t.Fatalf("Ошибка подписки: %v", err)
	}

	err = sp.Close(ctx)
	if err != nil {
		t.Fatalf("Ошибка закрытия: %v", err)
	}

	err = sp.Publish("тест", "тестовое сообщение")
	if err != context.Canceled {
		t.Errorf("Ожидалась ошибка context.Canceled, получено: %v", err)
	}

	_, err = sp.Subscribe("тест", func(msg interface{}) {})
	if err != context.Canceled {
		t.Errorf("Ожидалась ошибка context.Canceled, получено: %v", err)
	}
}
