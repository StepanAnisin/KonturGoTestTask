package farm

import (
	"log"
	"sync"

	"github.com/StepanAnisin/chickenfarm/pkg/config"
)

type Ranch struct {
	eggsInFridge int
	mutex        sync.Mutex
}

func InitRanch(ranch *Ranch) {
	cfg, err := config.LoadConfig("../config/app.env")
	if err != nil {
		log.Fatal("can not load config:", err)
	}
	var wg sync.WaitGroup
	wg.Add(cfg.ChikensCount + 1)
	defer wg.Wait()
	for i := 0; i < cfg.ChikensCount; i++ {
		go CarryEggs(i, &ranch.mutex, cfg.EggsMinSpawnCount, cfg.EggsMaxSpawnCount,
			cfg.EggsSpawnMinDelay, cfg.EggsSpawnMaxDelay, &ranch.eggsInFridge)
	}
	go FarmerComes(cfg.FarmerCheckMinDelay, cfg.FarmerCheckMaxDelay, cfg.FarmerMaxNeededQuantity,
		cfg.FarmerMinNeededQuantity, &ranch.mutex, &ranch.eggsInFridge)
}

func GetEggsCount(ranch *Ranch) int {
	ranch.mutex.Lock()
	eggsCount := ranch.eggsInFridge
	ranch.mutex.Unlock()
	return eggsCount
}
