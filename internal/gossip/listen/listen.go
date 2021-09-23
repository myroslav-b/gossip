package listen

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

func Listen(ctx context.Context) error {
	for i := 0; i < 200; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(10 * time.Millisecond)
			_ = math.Atan(math.Cos(math.Sin(float64(i))))
			fmt.Print("c")
		}
	}
	fmt.Println()
	log.Printf("Client %v", &ctx)
	return errors.New("BAD")
}
