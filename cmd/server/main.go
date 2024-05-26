package main

import (
	"ST_DataLinkLayer/cmd/code"
	"ST_DataLinkLayer/cmd/decode"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	// Создаем новый роутер Gin
	r := gin.Default()

	// Обрабатываем POST-запросы на /code
	r.POST("/code", func(c *gin.Context) {

		data, err := io.ReadAll(c.Request.Body)
		log.Println("Пришедший сегмент:", data)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.Status(http.StatusOK)
		log.Println("Принят запрос с данными:\n", string(data))

		encodedSegment := code.Code(data)

		if rand.Float64() < 0.02 {
			log.Println("Сегмент был потерян при декодировании")
			return
		} else {
			log.Println("Сегмент не был потерян при декодировании")
		}

		decodedSegment := decode.Decode(encodedSegment)
		log.Println("Декодированный сегмент:", decodedSegment)

		sendURL := "http://<TARGET_IP>:<TARGET_PORT>/send" // Заменить на нужный IP и порт
		resp, err := http.Post(sendURL, "application/json", bytes.NewBuffer(decodedSegment))
		if err != nil {
			log.Printf("Произошла ошибка при отправке POST-запроса: %v", err)
			return
		}
		defer resp.Body.Close()
		/*// Чтение ответа от сервера (лишнее)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return
		}
		log.Printf("Ответ от сервера: %s", body)*/
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}

/*func main() {
	// Пример данных: байты
	data := []byte{123, 10, 32, 32, 32, 47}

	// Кодирование данных
	encodedSegment := code.Code(data)
	fmt.Println("Encoded Segment:", encodedSegment)
	// Декодирование данных
	decodedSegment := decode.Decode(encodedSegment)
	fmt.Println("Decoded Segment:", decodedSegment)
}*/
