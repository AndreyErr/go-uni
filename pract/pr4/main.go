package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Genre  string
	Rating float32
}

var Titles = []string{"Война и мир", "Преступление и наказание", "Гарри Поттер и философский камень", "Властелин колец: Братство кольца", "1984", "Гордость и предубеждение", "Улитка на склоне", "Сто лет одиночества", "Над пропастью во ржи", "Анна Каренина"}
var Authors = []string{"Лев Толстой", "Федор Достоевский", "Джоан Роулинг", "Джон Толкин", "Джордж Оруэлл", "Джейн Остин", "Аркадий и Борис Стругацкие", "Габриэль Гарсиа Маркес", "Джером Сэлинджер", "Лев Толстой"}
var Genres = []string{"Роман", "Детектив", "Фэнтези", "Приключения", "Антиутопия", "Комедия", "Фантастика", "Магический реализм"}

func (r Book) String() string {
	return fmt.Sprintf("ID: %d, | Title: %s, | Author: %s, | Genre: %s, | Rating: %f", r.ID, r.Title, r.Author, r.Genre, r.Rating)
}

func ConnectToDB() *pgxpool.Pool {
	conn, err := pgxpool.Connect(context.Background(), "postgresql://postgres:7852@localhost:5432/postgres")
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}

func CreateTableIfNotExists(conn *pgxpool.Pool) {
	query := `
	DROP TABLE IF EXISTS BOOKS;
	CREATE TABLE IF NOT EXISTS BOOKS (
	id SERIAL PRIMARY KEY,
	title TEXT,
	author TEXT,
	genre TEXT,
	rating FLOAT4
	);
	`
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		log.Fatalln(err)
	}
}

func CreateBooks(conn *pgxpool.Pool) {
	r := rand.New(rand.NewSource(0)) // Одинаковые книги
	for i := 0; i < 50; i++ {
		title := Titles[r.Intn(len(Titles))]
		author := Authors[r.Intn(len(Authors))]
		genre := Genres[r.Intn(len(Genres))]
		rating := r.Float32()*5 + 1
		_, err := conn.Exec(context.Background(), `
		 INSERT INTO books (title, author, genre, rating)
		 VALUES ($1, $2, $3, $4);
		 `, title, author, genre, rating)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BookByID(conn *pgxpool.Pool, id int) Book {
	var book Book
	err := conn.QueryRow(context.Background(), `
	SELECT id, title, author, genre, rating
	FROM BOOKS
	WHERE id = $1;
	`, id).Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Rating)
	if err != nil {
		log.Fatalln(err)
	}
	return book
}

func BooksByGenre(conn *pgxpool.Pool, genre string, books chan<- Book) {
	rows, err := conn.Query(context.Background(), `
	SELECT id, title, author, genre, rating
	FROM BOOKS
	WHERE genre = $1;
	`, genre)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Rating)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(book.String())
		books <- book
	}
	close(books)
}

func ChangeAllBooksRating(conn *pgxpool.Pool, rating float32, id int) {
	_, err := conn.Exec(context.Background(), `
	UPDATE BOOKS SET rating = $1 WHERE id = $2;
	`, rating, id)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	// Инициализация
	conn := ConnectToDB()
	CreateTableIfNotExists(conn)
	CreateBooks(conn)
	ctx, _ := context.WithCancel(context.Background())
	err := rsocket.Receive().
		OnStart(func() {
			log.Println("Server Started")
		}).
		Acceptor(func(_ context.Context, _ payload.SetupPayload, _ rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			return rsocket.NewAbstractSocket(
				// Request-Response
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					// Получить ID - вернуть книгу с таким ID
					id, err := strconv.Atoi(msg.DataUTF8())
					if err != nil {
						log.Fatalln(err)
						return nil
					}
					
					fmt.Println("Запрос клиента - Получение книги с ID:", id)
					
					book := BookByID(conn, id)
					return mono.Just(payload.NewString(book.String(), ""))
				}),
				// Request-Stream
				rsocket.RequestStream(func(request payload.Payload) flux.Flux {
					// Получить жанр...
					genre := request.DataUTF8()
					
					fmt.Println("Запрос клиента - Получение книг с жанром:", genre)
					
					books := make(chan Book)
					go BooksByGenre(conn, genre, books)
					// ...Вернуть книги с таким жанром
					return flux.Create(func(_ context.Context, s flux.Sink) {
						for book := range books {
							s.Next(payload.NewString(book.String(), ""))
						}
						s.Complete()
					})
				}),
				// Request-Channel
				rsocket.RequestChannel(func(payloads flux.Flux) flux.Flux {
					// Создаем канал для ответов сервера
					responses := make(chan string)
					// Запускаем горутину, которая обрабатывает сообщения от клиента
					go func() {
						// Для каждого сообщения от клиента
						payloads.DoOnNext(func(msg payload.Payload) error {
							// Получаем текст сообщения
							text := msg.DataUTF8()
							fmt.Println("Запрос клиента - Сообщение от клиента:", text)
							// Проверяем, что клиент сказал
							switch text {
							case "пока": // Если клиент хочет закончить чат
								// Отправляем прощальное сообщение
								responses <- "Покеда!"
								// Закрываем канал ответов
								close(responses)
							case "рекомендация": // Если клиент просит рекомендацию
								// Выбираем случайную книгу из базы данных
								id := rand.Intn(50) + 1
								book := BookByID(conn, id)
								// Отправляем описание книги
								responses <- fmt.Sprintf("Рекомендуемая книга для прочтения: %s", book.String())
							default: // Если клиент спрашивает о книге
								// Ищем книгу по названию, автору или жанру
								rows, err := conn.Query(context.Background(), `
									SELECT id, title, author, genre, rating
									FROM BOOKS
									WHERE title = $1 OR author = $1 OR genre = $1;
									`, text)
								if err != nil {
									log.Fatalln(err)
								}
								defer rows.Close()
								// Если нашли хотя бы одну книгу
								if rows.Next() {
									var book Book
									err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Rating)
									if err != nil {
										log.Fatalln(err)
									}
									log.Println(book.String())
									// Отправляем описание первой найденной книги
									responses <- fmt.Sprintf("Найдена книга для вас: %s", book.String())
									rows.Close()
								} else {
									// Иначе отправляем сообщение, что ничего не нашли
									responses <- fmt.Sprintf("Не найдено книг по вашему запросу: %s", text)
								}
							}
							return nil
						}).Subscribe(context.Background())
					}()
					// Возвращаем поток ответов сервера
					return flux.Create(func(_ context.Context, s flux.Sink) {
						for response := range responses {
							s.Next(payload.NewString(response, ""))
						}
						s.Complete()
					})
				}),
				// Fire-and-Forget
				rsocket.FireAndForget(func(msg payload.Payload) {
					id, err := strconv.Atoi(msg.DataUTF8())
					if err != nil {
						log.Fatalln(err)
					}
					log.Println("Изменение рейтинга книги ID:", id)
					ChangeAllBooksRating(conn, 0, id)
				}),
			), nil
		}).
		Transport(rsocket.TCPServer().SetAddr(":7878").Build()).
		Serve(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
