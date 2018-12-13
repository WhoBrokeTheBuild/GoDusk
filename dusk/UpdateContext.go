package dusk

// UpdateContext is a context of timing data
type UpdateContext struct {
	FPS         int
	DeltaTime   float32
	ElapsedTime float64
	TotalTime   float64
}
