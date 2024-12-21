package engine

type RewindDrawOptions struct {
	CurrentFrame int
	TotalFrames  int
}

type DrawOptions struct {
	Rewind *RewindDrawOptions
}

type DrawOptionsFunc func(opts *DrawOptions)

func WithRewind(currentFrame, totalFrames int) DrawOptionsFunc {
	return func(opts *DrawOptions) {
		opts.Rewind = &RewindDrawOptions{CurrentFrame: currentFrame, TotalFrames: totalFrames}
	}
}
