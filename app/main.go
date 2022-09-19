package main

import (
	"config"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var eggsInFridge int = 0 //  общий ресурс

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	var wg sync.WaitGroup
	wg.Add(6)

	//Запуск процесса спауна яиц
	ch := make(chan bool) // канал
	var mutex sync.Mutex  // определяем мьютекс
	for i := 0; i < config.ChikensCount; i++ {
		go carryEggs(i, ch, &mutex, config.EggsMinSpawnCount, config.EggsMaxSpawnCount,
			config.EggsSpawnMinDelay, config.EggsSpawnMaxDelay)
	}
	go farmerComes(ch, config.FarmerCheckMinDelay, config.FarmerCheckMaxDelay, config.FarmerMaxNeededQuantity,
		config.FarmerMinNeededQuantity, &mutex)

	//Запуск Хэндлера
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		fmt.Print(eggsInFridge)
		mutex.Unlock()
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
	wg.Wait()
}

// Генерация рандомного числа в диапазоне
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// TODO
// Функция спауна яиц курицей
// func carryEggs(){
// 1. Генерируем время через сколько заспаунилось яйцо (eggSpawnDelay)
// 2. Генерируем сколько яиц заспаунилось (eggsMin, eggsMax)
// 3. Складываем в любой свободный ресурс
// 4. Приходит фермер, всё забирает. Видимо, обнуляем этот некторый счетчик
// }
func carryEggs(number int, ch chan bool, mutex *sync.Mutex, eggsMinSpawnCount int, eggsMaxSpawnCount int,
	eggsSpawnMinDelay int, eggsSpawnMaxDelay int) {
	for {
		eggsSpawnDelay := random(eggsSpawnMinDelay, eggsSpawnMaxDelay)
		// Calling Sleep method
		time.Sleep(time.Duration(eggsSpawnDelay) * time.Second)
		eggsSpawnCount := random(eggsMinSpawnCount, eggsMaxSpawnCount)
		mutex.Lock()
		eggsInFridge += eggsSpawnCount
		mutex.Unlock()
		log.Print("Курица ", number, " снесла ", eggsSpawnCount, " яиц")
	}
}

// Фермер приходит и забирает яйца
func farmerComes(ch chan bool, FarmerCheckMinDelay int, FarmerCheckMaxDelay int, FarmerMaxNeededQuantity int,
	FarmerMinNeededQuantity int, mutex *sync.Mutex) {
	for {
		farmerCheckDelay := random(FarmerCheckMinDelay, FarmerCheckMaxDelay)
		// Calling Sleep method
		time.Sleep(time.Duration(farmerCheckDelay) * time.Second)
		eggsQuantityNeeded := random(FarmerMinNeededQuantity, FarmerMaxNeededQuantity)
		mutex.Lock()
		if eggsInFridge <= eggsQuantityNeeded {
			eggsInFridge = 0
		}
		if eggsInFridge > eggsQuantityNeeded {
			eggsInFridge -= eggsQuantityNeeded
		}
		mutex.Unlock()
		log.Print("Фермер взял ", eggsQuantityNeeded, " яиц ")
	}
	ch <- true
}
