package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	serverAddress := flag.String("server", "127.0.0.1", "Address of the RabbitMQ server")
	initialChannel := flag.String("channel", "general", "Initial channel to join")
	flag.Parse()

	url := fmt.Sprintf("amqp://guest:guest@%s:5672/", *serverAddress)
	conn, err := amqp.Dial(url)
	failOnError(err, fmt.Sprintf("Unable to open connect to RabbitMQ server %s", url))
	defer conn.Close()

	fmt.Println("Connected to RabbitMQ server")

	ch, err := conn.Channel()
	failOnError(err, "Unable to open channel")
	defer ch.Close()

	var currentChannel string
	var msgs <-chan amqp.Delivery

	var ctx context.Context
	var cancel context.CancelFunc

	// Callback для смены канала
	switchChannel := func(channelName string) {
		if msgs != nil {
			fmt.Printf("Disconnect from %s\n", currentChannel)
			cancel()
		}

		// Определяем обменник
		err := ch.ExchangeDeclare(
			channelName, // наименование обменника
			"fanout",    // обменник будет широковещательным (рассылать сообщения всем подписчикам)
			true,        // живет ли обменник после закрытия соединения
			false,       // обменник удаляется, если к нему больше нет привязанных очередей
			false,       // internal
			false,       // сервер не отправляет подтверждение о создании обменника
			nil,         // arguments
		)
		failOnError(err, "Unable to declare exchange")

		// Определяем очередь
		q, err := ch.QueueDeclare(
			"",    // наименование очереди (может быть пустым, и в этом случае сервер сгенерирует уникальное им)
			false, // cохраняет очередь при перезапуске сервера
			true,  // удаляет очередь автоматически, если она больше не нужна
			false, // ограничивает доступ к очереди только текущим соединением
			false, // сервер не отправляет подтверждение о создании очереди
			nil,   // arguments
		)
		failOnError(err, "Unable to declare queue")

		// Привязываем очередь к обменнику
		err = ch.QueueBind(
			q.Name,      // наименование очереди
			"",          // routing key
			channelName, // наименование обменника
			false,       // сервер не отправляет подтверждение о связывании
			nil,
		)
		failOnError(err, "Unable to bind queue")

		// Подписываемся на сообщения очереди
		msgs, err = ch.Consume(
			q.Name, // наименование очереди
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // arguments
		)
		failOnError(err, "Unable to consume")

		// Запуск горутины для асинхронного вывода сообщения
		ctx, cancel = context.WithCancel(context.Background())
		go func() {
			for {
				select {
				case msg := <-msgs:
					fmt.Printf("%s\n", msg.Body)
				case <-ctx.Done():
					return
				}
			}
		}()

		currentChannel = channelName
		fmt.Printf("Switch to %s\n", currentChannel)
	}

	// Подключаемся к начальному каналу
	switchChannel(*initialChannel)

	// Консольный ввод для отправки сообщений и управления
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.HasPrefix(input, "!switch ") {
			newChannel := strings.TrimPrefix(input, "!switch ")
			switchChannel(newChannel)

		} else if input != "" {
			err := ch.Publish(
				currentChannel, // имя обменника
				"",             // routing key
				false,          // mandatory
				false,          // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(input),
				},
			)
			failOnError(err, "Unable to publish message")
		}
	}
}
