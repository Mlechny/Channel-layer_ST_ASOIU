package code

import (
	"log"
	"math/rand"
	"time"
)

func Code(data []byte) []int {

	//Преобразуем массив из байтов в массив из битов
	bits := bytesToBits(data)
	log.Println("Исходный массив в виде битов:", bits)
	generator := []int{1, 0, 0, 1, 1}
	var encodedBits []int
	blockSize := 11

	// Разделяем исходный массив битов на блоки по 11 битов
	for i := 0; i < len(bits); i += blockSize {
		end := i + blockSize
		if end > len(bits) {
			end = len(bits)
		}

		// Создаем новый срез для каждого блока
		block := make([]int, end-i)
		copy(block, bits[i:end])
		log.Printf("Блок из информационных битов [%d:%d]: %v\n", i, end, block)

		//Добавляем 0 в последний блок, если кол-во битов в нем < 11
		if len(block) < blockSize {
			padding := make([]int, blockSize-len(block))
			block = append(block, padding...)
			log.Println("Конечный блок с добавленными нулевыми битами:", block)
		}

		// Добавляем 4 проверочных бита в конец каждого блока
		padding := make([]int, 4)
		block = append(block, padding...)
		log.Println("Блок с проверочными битами:", block)

		// Выполняем полиномиальное деление для каждого блока
		remainder := polynomialDivision(block, generator)
		/*log.Println("Остаток:", remainder)*/

		// Заменяем проверочные биты блока на остаток от полиномиального деления
		copy(block[len(block)-4:], remainder)
		/*log.Println("Закодированный блок:", block)*/

		// Организуем закодированные блоки в единый массив
		encodedBits = append(encodedBits, block...)
		/*log.Println("Закодированный массив:", encodedBits)*/
	}

	//Вносим ошибку с определенной вероятностью
	encodedBits = introduceErrors(encodedBits)
	/*log.Println("Закодированный массив с возможной ошибкой:", encodedBits)*/
	return encodedBits
}

// Преобразование байтов в биты
func bytesToBits(data []byte) []int {
	var bits []int
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bits = append(bits, int((b>>i)&1))
		}
	}
	return bits
}

// Внесение ошибки
func introduceErrors(bits []int) []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if rng.Float64() < 0.1 {
		bitPosition := rng.Intn(len(bits))
		log.Println("Ошибка была внесена в закодированный сегмент")
		log.Println("Позиция ошибочного бита", bitPosition)
		//Инвертируем значение бита, тем самым внося ошибку
		bits[bitPosition] = 1 - bits[bitPosition]
		/*fmt.Println("Ошибка (было):", bits[10])
		bits[10] = 1 - bits[10]
		fmt.Println("Ошибка (стало):", bits[10])*/
	} else {
		log.Println("Ошибка не была внесена в закодированный сегмент")
	}
	return bits
}

// Полиномиальное деление
func polynomialDivision(dividend, divisor []int) []int {
	remainder := make([]int, len(divisor)-1)
	temp := make([]int, len(dividend))
	copy(temp, dividend)

	for i := 0; i < len(dividend)-len(divisor)+1; i++ {
		// Проверяем, нужно ли опустить ноль
		if temp[i] == 0 {
			continue
		}
		// Выполняем XOR для вычитания
		for j := 0; j < len(divisor); j++ {
			temp[i+j] ^= divisor[j]
		}
	}

	// Обрезаем лишние элементы
	remainder = temp[len(dividend)-len(divisor)+1:]
	return remainder
}
