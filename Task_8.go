package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type Person struct {
	ID          int
	Name        string
	Age         int
	Profession  string
	Hobbies     []string
	Friends     []*Person
	City        string
	Salary      float64
	IsMarried   bool
	Children    int
	Pets        []Pet
	Skills      map[string]int
	Education   []Education
	SocialMedia SocialMedia
	mutex       sync.Mutex
}

type Pet struct {
	Name     string
	Type     string
	Age      int
	Favorite bool
}

type Education struct {
	Institution string
	Degree      string
	Year        int
	GPA         float64
}

type SocialMedia struct {
	Facebook  string
	Instagram string
	Twitter   string
	LinkedIn  string
}

type Network struct {
	People     []*Person
	Statistics NetworkStats
	mutex      sync.RWMutex
}

type NetworkStats struct {
	TotalConnections int
	AverageAge       float64
	PopularHobbies   map[string]int
	PopularCities    map[string]int
}

func (p *Person) AddFriend(friend *Person) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, f := range p.Friends {
		if f.ID == friend.ID {
			return
		}
	}
	p.Friends = append(p.Friends, friend)
}

func (p *Person) RemoveFriend(friendID int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for i, friend := range p.Friends {
		if friend.ID == friendID {
			p.Friends = append(p.Friends[:i], p.Friends[i+1:]...)
			return
		}
	}
}

func (p *Person) Introduce() {
	fmt.Printf("\n🌟 Профиль пользователя %s 🌟\n", p.Name)
	fmt.Printf("ID: %d\n", p.ID)
	fmt.Printf("Имя: %s (%d лет)\n", p.Name, p.Age)
	fmt.Printf("Город: %s\n", p.City)
	fmt.Printf("Профессия: %s (Зарплата: %.2f₽)\n", p.Profession, p.Salary)
	fmt.Printf("Семейное положение: %s\n", map[bool]string{true: "В браке", false: "Не женат/не замужем"}[p.IsMarried])
	fmt.Printf("Количество детей: %d\n", p.Children)

	if len(p.Hobbies) > 0 {
		fmt.Println("\n📚 Хобби:")
		for _, hobby := range p.Hobbies {
			fmt.Printf("- %s\n", hobby)
		}
	}

	if len(p.Pets) > 0 {
		fmt.Println("\n🐾 Питомцы:")
		for _, pet := range p.Pets {
			fmt.Printf("- %s (%s, %d лет)%s\n",
				pet.Name,
				pet.Type,
				pet.Age,
				map[bool]string{true: " ❤️", false: ""}[pet.Favorite])
		}
	}

	fmt.Println("\n🎓 Образование:")
	for _, edu := range p.Education {
		fmt.Printf("- %s, %s (%d) - GPA: %.2f\n",
			edu.Institution,
			edu.Degree,
			edu.Year,
			edu.GPA)
	}

	fmt.Println("\n💪 Навыки:")
	for skill, level := range p.Skills {
		fmt.Printf("- %s: %s\n", skill, strings.Repeat("⭐", level))
	}
}

func (p *Person) ShowFriends() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.Friends) == 0 {
		fmt.Println("\n😢 У меня пока нет друзей")
		return
	}

	fmt.Printf("\n👥 Друзья (%d):\n", len(p.Friends))
	for _, friend := range p.Friends {
		fmt.Printf("- %s (%d лет) - %s из %s\n",
			friend.Name,
			friend.Age,
			friend.Profession,
			friend.City)
	}
}

func (n *Network) CalculateStatistics() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	stats := NetworkStats{
		PopularHobbies: make(map[string]int),
		PopularCities:  make(map[string]int),
	}

	totalAge := 0
	connections := 0

	for _, person := range n.People {
		totalAge += person.Age
		connections += len(person.Friends)

		for _, hobby := range person.Hobbies {
			stats.PopularHobbies[hobby]++
		}
		stats.PopularCities[person.City]++
	}

	stats.TotalConnections = connections
	stats.AverageAge = float64(totalAge) / float64(len(n.People))

	n.Statistics = stats
}

func CreateRandomPerson(id int) *Person {
	names := []string{"Александр", "Мария", "Иван", "Елена", "Дмитрий", "Анна", "Сергей", "Ольга", "Михаил", "Татьяна"}
	surnames := []string{"Иванов", "Петров", "Сидоров", "Смирнов", "Кузнецов", "Попов", "Васильев", "Соколов"}
	professions := []string{"Программист", "Учитель", "Врач", "Инженер", "Художник", "Дизайнер", "Архитектор", "Менеджер"}
	hobbies := []string{"Чтение", "Путешествия", "Спорт", "Музыка", "Кулинария", "Фотография", "Йога", "Рисование"}
	cities := []string{"Москва", "Санкт-Петербург", "Новосибирск", "Екатеринбург", "Казань", "Нижний Новгород"}
	petNames := []string{"Барсик", "Мурка", "Шарик", "Рекс", "Пушок", "Люся"}
	petTypes := []string{"Кот", "Собака", "Хомяк", "Попугай"}
	institutions := []string{"МГУ", "СПБГУ", "МФТИ", "ВШЭ", "МГТУ"}
	degrees := []string{"Бакалавр", "Магистр", "Специалист", "Аспирант"}
	skills := []string{"Программирование", "Аналитика", "Управление проектами", "Коммуникация", "Иностранные языки"}

	rand.Seed(time.Now().UnixNano())

	person := &Person{
		ID:         id,
		Name:       names[rand.Intn(len(names))] + " " + surnames[rand.Intn(len(surnames))],
		Age:        rand.Intn(40) + 20,
		Profession: professions[rand.Intn(len(professions))],
		City:       cities[rand.Intn(len(cities))],
		Salary:     float64(rand.Intn(150000) + 50000),
		IsMarried:  rand.Float32() < 0.5,
		Children:   rand.Intn(3),
		Hobbies:    make([]string, 0),
		Pets:       make([]Pet, 0),
		Skills:     make(map[string]int),
		Education:  make([]Education, 0),
		SocialMedia: SocialMedia{
			Facebook:  fmt.Sprintf("fb.com/user%d", id),
			Instagram: fmt.Sprintf("instagram.com/user%d", id),
			Twitter:   fmt.Sprintf("twitter.com/user%d", id),
			LinkedIn:  fmt.Sprintf("linkedin.com/in/user%d", id),
		},
	}

	numHobbies := rand.Intn(4) + 1
	for i := 0; i < numHobbies; i++ {
		hobby := hobbies[rand.Intn(len(hobbies))]
		if !contains(person.Hobbies, hobby) {
			person.Hobbies = append(person.Hobbies, hobby)
		}
	}

	numPets := rand.Intn(3)
	for i := 0; i < numPets; i++ {
		pet := Pet{
			Name:     petNames[rand.Intn(len(petNames))],
			Type:     petTypes[rand.Intn(len(petTypes))],
			Age:      rand.Intn(10) + 1,
			Favorite: i == 0,
		}
		person.Pets = append(person.Pets, pet)
	}

	numEducations := rand.Intn(2) + 1
	for i := 0; i < numEducations; i++ {
		education := Education{
			Institution: institutions[rand.Intn(len(institutions))],
			Degree:      degrees[rand.Intn(len(degrees))],
			Year:        2010 + rand.Intn(13),
			GPA:         4.0 * rand.Float64(),
		}
		person.Education = append(person.Education, education)
	}

	numSkills := rand.Intn(4) + 2
	for i := 0; i < numSkills; i++ {
		skill := skills[rand.Intn(len(skills))]
		if _, exists := person.Skills[skill]; !exists {
			person.Skills[skill] = rand.Intn(5) + 1
		}
	}

	return person
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func SaveNetworkToFile(network *Network, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(network)
}

func LoadNetworkFromFile(filename string) (*Network, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var network Network
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&network)
	return &network, err
}

func main() {
	network := &Network{}

	numPeople := 20
	network.People = make([]*Person, numPeople)

	var wg sync.WaitGroup
	wg.Add(numPeople)

	for i := 0; i < numPeople; i++ {
		go func(index int) {
			defer wg.Done()
			network.People[index] = CreateRandomPerson(index)
		}(i)
	}

	wg.Wait()

	for i := 0; i < len(network.People); i++ {
		numFriends := rand.Intn(5) + 2
		for j := 0; j < numFriends; j++ {
			friendIndex := rand.Intn(len(network.People))
			if friendIndex != i {
				network.People[i].AddFriend(network.People[friendIndex])
			}
		}
	}

	network.CalculateStatistics()

	fmt.Println("🌐 === СОЦИАЛЬНАЯ СЕТЬ ===")
	fmt.Printf("\n📊 Статистика сети:\n")
	fmt.Printf("Всего пользователей: %d\n", len(network.People))
	fmt.Printf("Общее количество связей: %d\n", network.Statistics.TotalConnections)
	fmt.Printf("Средний возраст: %.1f лет\n", network.Statistics.AverageAge)

	fmt.Println("\n🏆 Популярные хобби:")
	for hobby, count := range network.Statistics.PopularHobbies {
		fmt.Printf("- %s: %d человек\n", hobby, count)
	}

	fmt.Println("\n🏢 Распределение по городам:")
	for city, count := range network.Statistics.PopularCities {
		fmt.Printf("- %s: %d человек\n", city, count)
	}

	fmt.Println("\n👥 === ПРОФИЛИ ПОЛЬЗОВАТЕЛЕЙ ===")
	for _, person := range network.People {
		person.Introduce()
		person.ShowFriends()
		fmt.Println("\n" + strings.Repeat("=", 50))
	}

	err := SaveNetworkToFile(network, "social_network.json")
	if err != nil {
		fmt.Printf("Ошибка при сохранении сети: %v\n", err)
	} else {
		fmt.Println("\n💾 Сеть успешно сохранена в файл social_network.json")
	}
}
