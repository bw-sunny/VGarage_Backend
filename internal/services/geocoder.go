package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const GeocoderAPIURL = "https://geocode-maps.yandex.ru/1.x/"
const APIKey = "e39c0614-61a2-47c1-af0c-8d5785077992" // Вставь сюда ключ символ-в-символ

// GeocodeResponse описывает структуру ответа Яндекса (JSON)
type GeocodeResponse struct {
	Response struct {
		GeoObjectCollection struct {
			FeatureMember []struct {
				GeoObject struct {
					Point struct {
						Pos string `json:"pos"` // Строка формата "долгота широта" (например, "37.617644 55.755819")
					} `json:"Point"`
				} `json:"GeoObject"`
			} `json:"featureMember"`
		} `json:"GeoObjectCollection"`
	} `json:"response"`
}

// FetchCoordinates отправляет запрос к Яндексу и возвращает координаты
func FetchCoordinates(address string) (string, error) {
	// Кодируем адрес (пробелы превратятся в %20, запятые в %2C и т.д.)
	escapedAddress := url.QueryEscape(address)

	// Формируем URL. Для Геокодера обязательно указываем формат json
	fullURL := fmt.Sprintf("%s?apikey=%s&geocode=%s&format=json", GeocoderAPIURL, APIKey, escapedAddress)

	// Делаем прямой GET-запрос от имени сервера
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("ошибка сетевого запроса: %v", err)
	}
	defer resp.Body.Close()

	// Если ключ не подошел или заблокирован, мы сразу увидим код ошибки (например, 403)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("яндекс вернул статус-код: %d", resp.StatusCode)
	}

	var result GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ошибка декодирования JSON: %v", err)
	}

	// Проверяем, нашел ли Яндекс что-нибудь по этому адресу
	if len(result.Response.GeoObjectCollection.FeatureMember) == 0 {
		return "", fmt.Errorf("адрес не найден")
	}

	// Забираем координаты первого (самого точного) совпадения
	return result.Response.GeoObjectCollection.FeatureMember[0].GeoObject.Point.Pos, nil
}
