package bench

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkZapTextPositive(b *testing.B) {
	stream := &blackholeStream{}
	syncStream := zapcore.AddSync(stream)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})

	core := zapcore.NewCore(consoleEncoder, syncStream, highPriority)
	logger := zap.New(core)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("The quick brown fox jumps over the lazy dog")
		}
	})

	if stream.WriteCount() != uint64(b.N) {
		b.Fatalf("Log write count: %d != %d", stream.WriteCount(), uint64(b.N))
	}
}

func BenchmarkZapTextNegative(b *testing.B) {
	stream := &blackholeStream{}
	syncStream := zapcore.AddSync(stream)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return false
	})

	core := zapcore.NewCore(consoleEncoder, syncStream, highPriority)
	logger := zap.New(core)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("The quick brown fox jumps over the lazy dog")
		}
	})

	if stream.WriteCount() != uint64(0) {
		b.Fatalf("Log write count: %d != %d", stream.WriteCount(), uint64(0))
	}
}

func BenchmarkZapJSONNegative(b *testing.B) {
	stream := &blackholeStream{}
	syncStream := zapcore.AddSync(stream)

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return false
	})

	core := zapcore.NewCore(encoder, syncStream, levelEnabler)
	logger := zap.New(core).Sugar()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Infow("The quick brown fox jumps over the lazy dog",
				"rate", "15",
				"low", 16,
				"high", 123.2,
			)
		}
	})

	if stream.WriteCount() != uint64(0) {
		b.Fatalf("Log write count")
	}
}

func BenchmarkZapJSONPositive(b *testing.B) {
	stream := &blackholeStream{}
	syncStream := zapcore.AddSync(stream)

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})

	core := zapcore.NewCore(encoder, syncStream, levelEnabler)
	logger := zap.New(core).Sugar()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Infow("The quick brown fox jumps over the lazy dog",
				"rate", "15",
				"low", 16,
				"high", 123.2,
			)
		}
	})

	if stream.WriteCount() != uint64(b.N) {
		b.Fatalf("Log write count")
	}
}
