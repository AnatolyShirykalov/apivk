package apivk

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v0"
	"strconv"
)

const root = "https://api.vk.com/method/"

type VkResult []interface{}

func (v *VkResult) Pluck(key string) []interface{} {
	ret := make([]interface{}, 0, len(*v))
	for _, e := range *v {
		if e == nil {
			panic("Pluck panic: nil element")
		}
		em := e.(map[string]interface{})
		if key == "id" {
			ret = append(ret, int(em[key].(float64)))
		} else {
			ret = append(ret, em[key])
		}
	}
	return ret
}

func Init(tkn string) func(string, map[string]string) (*resty.Response, error) {
	return func(method string, pars map[string]string) (*resty.Response, error) {
		if _, ok := pars["access_token"]; !ok {
			pars["access_token"] = tkn
		}
		if _, ok := pars["version"]; !ok {
			pars["v"] = "5.62"
		}
		return resty.R().
			SetQueryParams(pars).
			SetHeader("Accept-Charset", "utf-8").
			Get(fmt.Sprintf("%s%s", root, method))
	}
}

func Woodpecker(tkn string) func(string, map[string]string) VkResult {
	get := Init(tkn)
	return func(method string, params map[string]string) VkResult {
		done := false
		ret := make(VkResult, 0, 20)
		for !done {
			resp, err := get(method, params)
			if err != nil {
				panic(err)
			}
			pre := make(map[string]interface{})
			err = json.Unmarshal(resp.Body(), &pre)
			if err != nil {
				panic(err)
			}
			var data map[string]interface{}
			if p, ok := pre["response"]; ok {
				switch p.(type) {
				case map[string]interface{}:
					data = pre
				case []interface{}:
					data = map[string]interface{}{
						"response": map[string]interface{}{
							"items": p,
							"count": float64(len(p.([]interface{}))),
						},
					}
				default:
					panic("default")
				}
				done = true
				if _, ok1 := data["response"]; !ok1 {
					panic("no data.response")
				}
				res := data["response"].(map[string]interface{})
				count := int(res["count"].(float64))
				_, has_count := params["count"]
				if count > 0 && has_count {
					items := res["items"].([]interface{})
					ret = append(ret, items...)
					//fmt.Println("[Debug]", method, count, len(items))
					if count > len(items) {
						var offset int
						if _, has_offset := params["offset"]; !has_offset {
							offset = 0
						} else {
							v := params["offset"]
							offset, err = strconv.Atoi(v)
						}
						if offset+len(items) < count {
							offset += len(items)
							params["offset"] = fmt.Sprintf("%d", offset)
							done = false
						}
					}
					//fmt.Println(len(items), count)
				} else {
					return res["items"].([]interface{})
				}
			} else {
				vk_err := pre["error"].(map[string]interface{})
				if msg := vk_err["error_msg"]; msg == "Too many requests per second" {
					fmt.Println("error", msg, "Try again")
				} else {
					fmt.Println("Unknown error", msg)
					panic(nil)
				}
			}
		}
		return ret
	}
}

func Run(token string) {
	get := Woodpecker(token)
	targets := []string{
		"Авиамоторная",
		"Автозаводская",
		"Академическая",
		"Александровский сад",
		"Алексеевская",
		"Алма-Атинская",
		"Алтуфьево",
		"Аннино",
		"Арбатская",
		"Аэропорт",
		"Бабушкинская",
		"Багратионовская",
		"Баррикадная",
		"Бауманская",
		"Беговая",
		"Белорусская",
		"Беляево",
		"Бибирево",
		"Библиотека имени Ленина",
		"Борисово",
		"Боровицкая",
		"Ботанический сад",
		"Братиславская",
		"Бульвар адмирала Ушакова",
		"Бульвар Дмитрия Донского",
		"Бульвар Рокоссовского",
		"Бунинская аллея",
		"Варшавская",
		"ВДНХ",
		"Владыкино",
		"Водный стадион",
		"Войковская",
		"Волгоградский проспект",
		"Волжская",
		"Волоколамская",
		"Воробьевы горы",
		"Выставочная",
		"Выхино",
		"Деловой центр",
		"Динамо",
		"Дмитровская",
		"Добрынинская",
		"Домодедовская",
		"Достоевская",
		"Дубровка",
		"Жулебино",
		"Зябликово",
		"Измайловская",
		"Калужская",
		"Кантемировская",
		"Каховская",
		"Каширская",
		"Киевская",
		"Китай-город",
		"Кожуховская",
		"Коломенская",
		"Комсомольская",
		"Коньково",
		"Красногвардейская",
		"Краснопресненская",
		"Красносельская",
		"Красные ворота",
		"Крестьянская застава",
		"Кропоткинская",
		"Крылатское",
		"Кузнецкий мост",
		"Кузьминки",
		"Кунцевская",
		"Курская",
		"Кутузовская",
		"Ленинский проспект",
		"Лермонтовский проспект",
		"Лубянка",
		"Люблино",
		"Марксистская",
		"Марьина роща",
		"Марьино",
		"Маяковская",
		"Медведково",
		"Международная",
		"Менделеевская",
		"Митино",
		"Молодежная",
		"Мякинино",
		"Нагатинская",
		"Нагорная",
		"Нахимовский проспект",
		"Новогиреево",
		"Новокосино",
		"Новокузнецкая",
		"Новослободская",
		"Новоясеневская",
		"Новые Черемушки",
		"Октябрьская",
		"Октябрьское поле",
		"Орехово",
		"Отрадное",
		"Охотныйряд",
		"Павелецкая",
		"Парк культуры",
		"Парк Победы",
		"Партизанская",
		"Первомайская",
		"Перово",
		"Петровско-Разумовская",
		"Печатники",
		"Пионерская",
		"Планерная",
		"Площадь Ильича",
		"Площадь Революции",
		"Полежаевская",
		"Полянка",
		"Пражская",
		"Преображенская площадь",
		"Пролетарская",
		"Проспект Вернадского",
		"Проспект Мира",
		"Профсоюзная",
		"Пушкинская",
		"Пятницкое шоссе",
		"Речной вокзал",
		"Рижская",
		"Римская",
		"Рязанский проспект",
		"Савеловская",
		"Свиблово",
		"Севастопольская",
		"Семеновская",
		"Серпуховская",
		"Славянский бульвар",
		"Сокол",
		"Сокольники",
		"Спартак",
		"Спортивная",
		"Сретенский бульвар",
		"Строгино",
		"Студенческая",
		"Сухаревская",
		"Сходненская",
		"Таганская",
		"Тверская",
		"Театральная",
		"Текстильщики",
		"Теплый стан",
		"Тимирязевская",
		"Третьяковская",
		"Тропарево",
		"Трубная",
		"Тульская",
		"Тургеневская",
		"Тушинская",
		"Улица Академика Янгеля",
		"Улица Горчакова",
		"Улица Скобелевская",
		"Улица Старокачаловская",
		"Улица 1905 года",
		"Университет",
		"Филевский парк",
		"Фили",
		"Фрунзенская",
		"Царицыно",
		"Цветной бульвар",
		"Черкизовская",
		"Чертановская",
		"Чеховская",
		"Чистые пруды",
		"Чкаловская",
		"Шаболовская",
		"Шипиловская",
		"Шоссе Энтузиастов",
		"Щелковская",
		"Щукинская",
		"Электрозаводская",
		"Юго-Западная",
		"Южная",
		"Ясенево",
	}
	for _, station := range targets {
		data := get("groups.search", map[string]string{
			"q":     fmt.Sprintf("подслушано %s", station),
			"count": "20",
		})
		str := ""
		for i, el := range data.Pluck("id") {
			if i != 0 {
				str = fmt.Sprintf("%s%s", str, ",")
			}
			str = fmt.Sprintf("%s%d", str, el)
		}
		if len(str) > 0 {
			good := get("groups.getById", map[string]string{
				"group_ids": str,
				"fields":    "members_count",
			})
			fmt.Println(good.Pluck("members_count"))
		}
	}
}
