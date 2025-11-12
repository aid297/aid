package str

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

type Rand struct {
	// bufferChan is the buffer for random bytes,
	// every item storing 4 bytes.
	bufferChan chan []byte
}

var (
	Letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"                           // 26
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"                           // 26
	Symbols      = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"                   // 32
	Digits       = "0123456789"                                           // 10
	characters   = Letters + Digits + Symbols                             // 94
	RandApp      Rand
)

var ()

// Buffer size for uint32 random number.
const bufferChanSize = 10000

func (*Rand) New() *Rand {
	ins := &Rand{}
	ins.bufferChan = make(chan []byte, bufferChanSize)

	go ins.asyncProducingRandomBufferBytesLoop()

	return ins
}

// NewRand 实例化：随机字符串
//
//go:fix 推荐使用：New方法
func NewRand() *Rand {
	ins := &Rand{}
	ins.bufferChan = make(chan []byte, bufferChanSize)

	go ins.asyncProducingRandomBufferBytesLoop()

	return ins
}

// asyncProducingRandomBufferBytes is a named goroutine, which uses an asynchronous goroutine
// to produce the random bytes, and a buffer chan to store the random bytes.
// So it has high performance to generate random numbers.
func (my *Rand) asyncProducingRandomBufferBytesLoop() {
	var step int
	for {
		buffer := make([]byte, 1024)
		if n, err := rand.Read(buffer); err != nil {
			panic(err)
		} else {
			// The random buffer from system is very expensive,
			// so fully reuse the random buffer by changing
			// the step with a different number can
			// improve the performance a lot.
			// for _, step = range []int{4, 5, 6, 7} {
			for _, step = range []int{4} {
				for i := 0; i <= n-4; i += step {
					my.bufferChan <- buffer[i : i+4]
				}
			}
		}
	}
}

func (my *Rand) Intn(max int) int {
	if max <= 0 {
		return max
	}
	n := int(binary.LittleEndian.Uint32(<-my.bufferChan)) % max
	if (max > 0 && n < 0) || (max < 0 && n > 0) {
		return -n
	}

	return n
}

// B retrieves and returns random bytes of given length `n`.
func (my *Rand) B(n int) []byte {
	if n <= 0 {
		return nil
	}
	i := 0
	b := make([]byte, n)
	for {
		copy(b[i:], <-my.bufferChan)
		i += 4
		if i >= n {
			break
		}
	}

	return b
}

// N returns a random int between min and max: [min, max].
// The `min` and `max` also support negative numbers.
func (my *Rand) N(min, max int) int {
	if min >= max {
		return min
	}
	if min >= 0 {
		return my.Intn(max-min+1) + min
	}
	// As `Intn` dose not support negative number,
	// so we should first shift the value to right,
	// then call `Intn` to produce the random number,
	// and finally shift the result back to left.
	return my.Intn(max+(0-min)+1) - (0 - min)
}

// S returns a random str which contains digits and letters, and its length is `n`.
// The optional parameter `symbols` specifies whether the result could contain symbols,
// which is false in default.
func (my *Rand) S(n int, symbols ...bool) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = my.B(n)
	)
	for i := range b {
		if len(symbols) > 0 && symbols[0] {
			b[i] = characters[numberBytes[i]%94]
		} else {
			b[i] = characters[numberBytes[i]%62]
		}
	}

	return string(b)
}

// D returns a random time.Duration between min and max: [min, max].
func (my *Rand) D(min, max time.Duration) time.Duration {
	multiple := int64(1)
	if min != 0 {
		for min%10 == 0 {
			multiple *= 10
			min /= 10
			max /= 10
		}
	}
	n := int64(my.N(int(min), int(max)))

	return time.Duration(n * multiple)
}

// GetString randomly picks and returns `n` count of chars from given str `s`.
// It also supports unicode str like Chinese/Russian/Japanese, etc.
func (my *Rand) GetString(s string, n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b     = make([]rune, n)
		runes = []rune(s)
	)
	if len(runes) <= 255 {
		numberBytes := my.B(n)
		for i := range b {
			b[i] = runes[int(numberBytes[i])%len(runes)]
		}
	} else {
		for i := range b {
			b[i] = runes[my.Intn(len(runes))]
		}
	}

	return string(b)
}

// GetDigits returns a random str which contains only digits, and its length is `n`.
func (my *Rand) GetDigits(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = my.B(n)
	)
	for i := range b {
		b[i] = Digits[numberBytes[i]%10]
	}
	return string(b)
}

// GetLetters returns a random str which contains only letters, and its length is `n`.
func (my *Rand) GetLetters(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = my.B(n)
	)
	for i := range b {
		b[i] = Letters[numberBytes[i]%52]
	}

	return string(b)
}

// GetSymbols returns a random str which contains only symbols, and its length is `n`.
func (my *Rand) GetSymbols(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = my.B(n)
	)
	for i := range b {
		b[i] = Symbols[numberBytes[i]%32]
	}

	return string(b)
}

// Perm returns, as a slice of n int numbers, a pseudo-random permutation of the integers [0,n).
// TODO performance improving for large slice producing.
func (my *Rand) Perm(n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := my.Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}

	return m
}

// Meet randomly calculate whether the given probability `num`/`total` is met.
func (my *Rand) Meet(num, total int) bool { return my.Intn(total) < num }

// MeetProb randomly calculate whether the given probability is met.
func (my *Rand) MeetProb(prob float32) bool { return my.Intn(1e7) < int(prob*1e7) }
