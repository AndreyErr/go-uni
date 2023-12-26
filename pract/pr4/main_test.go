package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
)

func connect(t *testing.T) rsocket.Client {
	// Connect to server
	cli, err := rsocket.Connect().
		Transport(rsocket.TCPClient().
			SetHostAndPort("127.0.0.1", 7878).
			Build()).
		Start(context.Background())
	if err != nil {
		t.Error(err)
	}
	return cli
}

func TestRequestResponse(t *testing.T) {
	cli := connect(t)
	defer cli.Close()
	wg := &sync.WaitGroup{}
	// 10 случайных ID
	attempts := 10
	wg.Add(attempts)
	for i := 0; i < attempts; i++ {
		cli.RequestResponse(payload.NewString(fmt.Sprint(rand.Intn(50)+1), "")).
			DoOnSuccess(func(input payload.Payload) error {
				fmt.Println("Ответ сервера - книга:", input.DataUTF8())
				wg.Done()
				return nil
			}).
			Subscribe(context.Background())
	}
	wg.Wait()
}

func TestRequestStream(t *testing.T) {
	cli := connect(t)
	defer cli.Close()
	wg := &sync.WaitGroup{}
	// 3 случайных модели
	attempts := 2
	wg.Add(attempts)
	for i := 0; i < attempts; i++ {
		model := Genres[rand.Intn(len(Genres))]
		stream := cli.RequestStream(payload.NewString(model, ""))
		stream.
			DoOnComplete(func() {
				wg.Done()
			}).
			DoOnNext(func(input payload.Payload) error {
				fmt.Printf("Ответ сервера - Книга (Жанр %s) : %s\n", model, input.DataUTF8())
				return nil
			}).Subscribe(context.Background())
	}
	wg.Wait()
}

func TestRequestChannel(t *testing.T) {
	cli := connect(t)
	defer cli.Close()
	wg := &sync.WaitGroup{}
	// Вывод диалога между клиентом и сервером
	wg.Add(1)
	// Создаем поток сообщений от клиента
	sendFlux := flux.Create(func(_ context.Context, s flux.Sink) {
		// Отправляем разные сообщения
		s.Next(payload.NewString("ааа", ""))
		s.Next(payload.NewString("Фантастика", ""))
		s.Next(payload.NewString("рекомендация", ""))
		s.Next(payload.NewString("Джордж Оруэлл", ""))
		s.Next(payload.NewString("пока", ""))
		s.Complete()
	})
	// Подписываемся на поток ответов от сервера
	cli.RequestChannel(sendFlux).
		DoOnComplete(func() {
			wg.Done()
		}).
		DoOnNext(func(input payload.Payload) error {
			fmt.Println("Ответ сервера:", input.DataUTF8())
			return nil
		}).Subscribe(context.Background())
	wg.Wait()
}


func TestFireAndForget(t *testing.T) {
	cli := connect(t)
	defer cli.Close()
	cli.FireAndForget(payload.NewString(fmt.Sprint(rand.Intn(30)+1), ""))
}