package decode

import (
	"log"
	"math"
)

func Decode(encodedBits []int) []byte {
	var decodedBits []int
	blockSize := 15
	generator := []int{1, 0, 0, 1, 1}

	// Разделяем закодированный массив битов на блоки по 15 битов
	for i := 0; i < len(encodedBits); i += blockSize {
		end := i + blockSize
		if end > len(encodedBits) {
			end = len(encodedBits)
		}

		// Создаем новый срез для каждого блока
		block := make([]int, blockSize)
		copy(block, encodedBits[i:end])
		log.Println("Закодированный блок с проверочными битами:", block)

		//Выполняем полиномиальное деление для каждого блока
		remainder := polynomialDivision(block, generator)
		log.Println("Синдром ошибки:", remainder)

		//Сравниваем остаток от полиномиального деления с 0
		if isZero(remainder) {
			log.Println("Блок не содержит ошибку.")
		} else {
			log.Println("Блок содержит ошибку.")
			log.Println("Байт с ошибкой:", bitsToBytes(block))
			// Исправляем внесенную ошибку
			correctError(block, remainder)
			log.Println("Блок с исправленной ошибкой:", block)
			log.Println("Исправленный байт:", bitsToBytes(block))
		}

		//Обрезаем у каждого блока проверочные биты
		infoBits := block[:11]
		/*fmt.Println("Блок из информационных битов:", infoBits)*/

		//Организуем декодированные блоки в единый массив
		decodedBits = append(decodedBits, infoBits...)
		/*fmt.Println("Декодированный массив:", decodedBits)*/
	}

	//Преобразуем массив из битов в массив из байтов
	decodedBytes := bitsToBytes(decodedBits)

	//Случай, когда в конечном блоке число добавленных незначащих битов кратно 8
	if decodedBytes[len(decodedBytes)-1] == 0 {
		// Удаляем последний байт
		decodedBytes = decodedBytes[:len(decodedBytes)-1]
	}
	return decodedBytes
}

// Преобразование битов в байты
func bitsToBytes(bits []int) []byte {
	byteLen := int(math.Ceil(float64(len(bits)) / 8.0))
	bytes := make([]byte, byteLen)

	for i, bit := range bits {
		if bit == 1 {
			bytes[i/8] |= 1 << (7 - uint(i)%8)
		}
	}

	return bytes
}

// Полиномиальное деление
/*func polynomialDivision(data []int, generator []int) []int {
	n := len(data)
	k := len(generator)

	dividend := append([]int{}, data...)

	for i := 0; i <= n-k; i++ {
		if dividend[i] == 1 {
			for j := 0; j < k; j++ {
				dividend[i+j] ^= generator[j]
			}
		}
	}

	return dividend[n-k+1:]
}*/
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

// Сравнение остатка с 0
func isZero(remainder []int) bool {
	for _, bit := range remainder {
		if bit != 0 {
			return false
		}
	}
	return true
}

// Исправление ошибки
func correctError(block, remainder []int) []int {

	//Массив синдромов ошибок для циклического кода [15, 11]
	syndromes := [][]int{
		{0, 0, 0, 1},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{1, 0, 0, 0},
		{0, 0, 1, 1},
		{0, 1, 1, 0},
		{1, 1, 0, 0},
		{1, 0, 1, 1},
		{0, 1, 0, 1},
		{1, 0, 1, 0},
		{0, 1, 1, 1},
		{1, 1, 1, 0},
		{1, 1, 1, 1},
		{1, 1, 0, 1},
		{1, 0, 0, 1},
	}

	var errorPosition int
	for i, s := range syndromes {
		// Сравниваем текущий элемент с remainder
		if compareArrays(s, remainder) {
			errorPosition = i + 1
			break
		}
	}

	//Инвертируем значение ошибочного бита при отсчете от конца массива
	block[len(block)-errorPosition] = 1 - block[len(block)-errorPosition]
	return block
}

// Сравнение массивов на основе сравнения всех их элементов
func compareArrays(arr1, arr2 []int) bool {
	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}
