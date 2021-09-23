package talk

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	traficmanager "github.com/myroslav-b/gossip/internal/gossip/traficManager"
)

func Talk(ctx context.Context, tc traficmanager.TraficControler) error {
	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if tc.Speak() {
				time.Sleep(10 * time.Millisecond)
				_ = math.Atan(math.Sin(math.Cos(float64(i))))
				fmt.Print("s")
			}
		}
	}
	fmt.Println()
	log.Printf("Server %v", &ctx)
	return nil
}
